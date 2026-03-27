package pipeline

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StreamingProxy struct {
	forwardingStage *ForwardingStage
}

func NewStreamingProxy(forwardingStage *ForwardingStage) *StreamingProxy {
	return &StreamingProxy{
		forwardingStage: forwardingStage,
	}
}

func (sp *StreamingProxy) ProxyStream(ctx context.Context, w http.ResponseWriter, req *Request) error {
	body, headers, statusCode, err := sp.forwardingStage.ForwardStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to forward stream: %w", err)
	}
	defer body.Close()

	if statusCode != http.StatusOK {
		return fmt.Errorf("provider returned status %d", statusCode)
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming not supported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for key, values := range headers {
		if strings.HasPrefix(strings.ToLower(key), "x-") {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return fmt.Errorf("failed to write to client: %w", err)
		}
		flusher.Flush()

		if line == "data: [DONE]" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func (sp *StreamingProxy) ProxyStreamWithContext(ctx context.Context, w http.ResponseWriter, req *Request) error {
	done := make(chan error, 1)

	go func() {
		done <- sp.ProxyStream(ctx, w, req)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
