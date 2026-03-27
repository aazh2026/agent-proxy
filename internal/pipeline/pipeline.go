package pipeline

import (
	"context"
	"net/http"
)

type Request struct {
	HTTPRequest *http.Request
	UserID      string
	Model       string
	Provider    string
	Token       string
	Body        []byte
	Stream      bool
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Stream     bool
}

type Stage interface {
	Name() string
	Process(ctx context.Context, req *Request) (*Request, error)
}

type ResponseStage interface {
	Name() string
	ProcessResponse(ctx context.Context, req *Request, resp *Response) (*Response, error)
}

type Pipeline struct {
	stages         []Stage
	responseStages []ResponseStage
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stages:         make([]Stage, 0),
		responseStages: make([]ResponseStage, 0),
	}
}

func (p *Pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *Pipeline) AddResponseStage(stage ResponseStage) {
	p.responseStages = append(p.responseStages, stage)
}

func (p *Pipeline) Execute(ctx context.Context, req *Request) (*Request, error) {
	currentReq := req
	for _, stage := range p.stages {
		var err error
		currentReq, err = stage.Process(ctx, currentReq)
		if err != nil {
			return nil, err
		}
	}
	return currentReq, nil
}

func (p *Pipeline) ProcessResponse(ctx context.Context, req *Request, resp *Response) (*Response, error) {
	currentResp := resp
	for _, stage := range p.responseStages {
		var err error
		currentResp, err = stage.ProcessResponse(ctx, req, currentResp)
		if err != nil {
			return nil, err
		}
	}
	return currentResp, nil
}

type ValidationStage struct{}

func NewValidationStage() *ValidationStage {
	return &ValidationStage{}
}

func (s *ValidationStage) Name() string {
	return "validation"
}

func (s *ValidationStage) Process(ctx context.Context, req *Request) (*Request, error) {
	return req, nil
}

type RoutingStage struct {
	resolver func(model string) (provider string, err error)
}

func NewRoutingStage(resolver func(model string) (provider string, err error)) *RoutingStage {
	return &RoutingStage{resolver: resolver}
}

func (s *RoutingStage) Name() string {
	return "routing"
}

func (s *RoutingStage) Process(ctx context.Context, req *Request) (*Request, error) {
	provider, err := s.resolver(req.Model)
	if err != nil {
		return nil, err
	}
	req.Provider = provider
	return req, nil
}

type TokenStage struct {
	resolver func(userID, provider string) (token string, err error)
}

func NewTokenStage(resolver func(userID, provider string) (token string, err error)) *TokenStage {
	return &TokenStage{resolver: resolver}
}

func (s *TokenStage) Name() string {
	return "token"
}

func (s *TokenStage) Process(ctx context.Context, req *Request) (*Request, error) {
	token, err := s.resolver(req.UserID, req.Provider)
	if err != nil {
		return nil, err
	}
	req.Token = token
	return req, nil
}
