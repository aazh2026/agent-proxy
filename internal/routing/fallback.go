package routing

import (
	"github.com/openclaw/agent-proxy/internal/logging"
)

type FallbackChain struct {
	primary   string
	fallbacks []string
}

func NewFallbackChain(primary string, fallbacks []string) *FallbackChain {
	return &FallbackChain{
		primary:   primary,
		fallbacks: fallbacks,
	}
}

func (c *FallbackChain) GetChain() []string {
	chain := []string{c.primary}
	chain = append(chain, c.fallbacks...)
	return chain
}

func (c *FallbackChain) Primary() string {
	return c.primary
}

func (c *FallbackChain) HasFallback() bool {
	return len(c.fallbacks) > 0
}

type FallbackRouter struct {
	chains map[string]*FallbackChain
}

func NewFallbackRouter() *FallbackRouter {
	return &FallbackRouter{
		chains: make(map[string]*FallbackChain),
	}
}

func (r *FallbackRouter) AddChain(model string, primary string, fallbacks []string) {
	r.chains[model] = NewFallbackChain(primary, fallbacks)
}

func (r *FallbackRouter) GetChain(model string) *FallbackChain {
	if chain, ok := r.chains[model]; ok {
		return chain
	}
	return nil
}

func (r *FallbackRouter) ExecuteWithFallback(model string, fn func(provider string) error) error {
	chain := r.GetChain(model)
	if chain == nil {
		return fn("")
	}

	for _, provider := range chain.GetChain() {
		logging.Debug("Trying provider %s for model %s", provider, model)
		err := fn(provider)
		if err == nil {
			return nil
		}
		logging.Warn("Provider %s failed: %v", provider, err)
	}

	return ErrAllTokensFailed
}
