package gateway

import manifold "github.com/manifoldco/go-manifold"

type endpoint struct {
	backend manifold.Backend
}

// DefaultBackend returns an instance of the default manifold.Backend configuration.
func DefaultBackend() manifold.Backend {
	return manifold.DefaultBackend()
}
