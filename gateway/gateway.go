package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	manifold "github.com/manifoldco/go-manifold"
)

// Client is the Manifold API client.
type Client struct {
	client http.Client
	APIClient
}

type defaultBackend struct {
	client *http.Client
	base   string
}

const baseGatewayURL = "https://api.manifold.co/v1"

// DefaultBackend returns an instance of the default Backend configuration.
func DefaultBackend() manifold.Backend {
	return &defaultBackend{client: &http.Client{}, base: baseGatewayURL}
}

// New returns a new API client with the default configuration
func New(cfgs ...ConfigFunc) *Client {
	c := &Client{
		http.Client{Transport: http.DefaultTransport},
		*NewAPI(),
	}

	c.APIClient.common.backend.(*defaultBackend).client = &c.client

	for _, cfg := range cfgs {
		cfg(c)
	}

	ForURLPattern(manifold.DefaultURLPattern)(c)

	// We need to do this after we've set the configuration. In case someone
	// provided a UserAgent, it will get loaded and overwrite our defaults since
	// we re-assign the previous transport after this.
	WithUserAgent("")(c)
	WithAPIToken(os.Getenv("MANIFOLD_API_TOKEN"))(c)

	return c
}

func (b *defaultBackend) NewRequest(method, path string, query url.Values, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	url := b.base
	if path[0] != '/' {
		url += "/"
	}
	url += path
	if q := query.Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (b *defaultBackend) Do(ctx context.Context, request *http.Request, v interface{}, errFn func(int) error) (*http.Response, error) {
	request = request.WithContext(ctx)

	resp, err := b.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		if errFn == nil {
			return nil, nil
		}

		apiErr := errFn(resp.StatusCode)
		if apiErr == nil {
			return nil, nil
		}

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, apiErr
	}

	if v != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(v); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

type endpoint struct {
	backend manifold.Backend
}

// ConfigFunc is a func that configures the client during New
type ConfigFunc func(*Client)

// ForURLPattern returns a configuration func to set the URL pattern for all
// endpoints.
func ForURLPattern(pattern string) ConfigFunc {
	return func(c *Client) {
		c.APIClient.common.backend.(*defaultBackend).base = baseGatewayURL
	}
}

// WithAPIToken returns a configuration func to set the API key to use for
// authentication.
func WithAPIToken(token string) ConfigFunc {
	return func(c *Client) {
		ot := c.client.Transport
		c.client.Transport = rtFunc(func(r *http.Request) (*http.Response, error) {
			if token != "" {
				r.Header.Set("Authorization", "Bearer "+token)
			}
			return ot.RoundTrip(r)
		})
	}
}

// WithUserAgent sets a specific user agent on the client. This will overwrite
// any 'User-Agent' header that has been set before. We will prepend the
// specified agent with `go-manifold-$version`.
func WithUserAgent(agent string) ConfigFunc {
	return func(c *Client) {
		ot := c.client.Transport
		c.client.Transport = rtFunc(func(r *http.Request) (*http.Response, error) {
			nagent := agent
			if agent != "" {
				nagent = fmt.Sprintf(" (%s)", nagent)
			}

			r.Header.Set("User-Agent", fmt.Sprintf("go-manifold-%s%s", manifold.Version, nagent))
			return ot.RoundTrip(r)
		})
	}
}

type rtFunc func(*http.Request) (*http.Response, error)

func (rt rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return rt(r) }
