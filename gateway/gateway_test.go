package gateway_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/gateway"
)

func TestConfig_WithUserAgent(t *testing.T) {
	hct := &headerCheckTransport{}
	http.DefaultTransport = hct
	defaultAgent := fmt.Sprintf("go-manifold-%s", manifold.Version)

	t.Run("without extra configuration", func(t *testing.T) {
		c := gateway.New()

		hct.reset()
		hct.expectHeaderEquals(t, "User-Agent", defaultAgent)

		c.Products.List(context.Background(), nil)
	})

	t.Run("with extra configuration", func(t *testing.T) {
		c := gateway.New(gateway.WithUserAgent("test"))

		hct.reset()
		hct.expectHeaderEquals(t, "User-Agent", fmt.Sprintf("%s (test)", defaultAgent))

		c.Products.List(context.Background(), nil)
	})

	t.Run("with multiple user agents", func(t *testing.T) {
		c := gateway.New(gateway.WithUserAgent("test"), gateway.WithUserAgent("manifold"))

		hct.reset()
		hct.expectHeaderEquals(t, "User-Agent", fmt.Sprintf("%s (test)", defaultAgent))

		c.Products.List(context.Background(), nil)
	})

	t.Run("with multiple calls - [GH 38]", func(t *testing.T) {
		c := gateway.New(gateway.WithUserAgent("test"), gateway.WithUserAgent("manifold"))

		hct.reset()
		hct.expectHeaderEquals(t, "User-Agent", fmt.Sprintf("%s (test)", defaultAgent))

		c.Products.List(context.Background(), nil)
		c.Products.List(context.Background(), nil)
	})
}

func TestConfig_WithAPIToken(t *testing.T) {
	hct := &headerCheckTransport{}
	http.DefaultTransport = hct

	t.Run("without extra configuration", func(t *testing.T) {
		token := os.Getenv("MANIFOLD_API_TOKEN")
		os.Unsetenv("MANIFOLD_API_TOKEN")

		defer func() {
			os.Setenv("MANIFOLD_API_TOKEN", token)
		}()

		c := gateway.New()

		hct.reset()
		hct.expectHeaderEquals(t, "Authorization", "")

		c.Products.List(context.Background(), nil)
	})

	t.Run("with env var", func(t *testing.T) {
		token := os.Getenv("MANIFOLD_API_TOKEN")
		os.Setenv("MANIFOLD_API_TOKEN", "s3cr3t")

		defer func() {
			os.Setenv("MANIFOLD_API_TOKEN", token)
		}()

		c := gateway.New()

		hct.reset()
		hct.expectHeaderEquals(t, "Authorization", "Bearer s3cr3t")

		r := c.Products.List(context.Background(), nil)
		if r != nil {
			for r.Next() {
				p, err := r.Current()
				fmt.Println(p, err)
			}
		}

		t.Error("blah")
	})

	t.Run("with extra configuration", func(t *testing.T) {
		c := gateway.New(gateway.WithAPIToken("test-token"))

		hct.reset()
		hct.expectHeaderEquals(t, "Authorization", "Bearer test-token")

		c.Products.List(context.Background(), nil)
	})
}

type headerCheckTransport struct {
	t      *testing.T
	checks map[string]string
}

func (hct *headerCheckTransport) reset() {
	hct.checks = map[string]string{}
}

func (hct *headerCheckTransport) expectHeaderEquals(t *testing.T, key, value string) {
	hct.t = t
	hct.checks[key] = value
}

func (hct *headerCheckTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	for key, value := range hct.checks {
		if h := r.Header.Get(key); h != value {
			hct.t.Errorf("Expected header '%s' to be '%s', got '%s')", key, value, h)
		}
	}

	// return an error here so the test doesn't trip over nil values
	return nil, errors.New("not successful")
}
