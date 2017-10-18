package manifold

import (
	"context"
	"fmt"
	"net/http"

	"github.com/manifoldco/go-base64"
)

// DefaultURLPattern is the default pattern used for connecting to Manifold's
// API hosts.
const DefaultURLPattern = "https://api.%s.manifold.co/v1"

// Client is the Manifold API client.
type Client struct {
	client http.Client
	IdentityClient
	CatalogClient
	MarketplaceClient
}

// New returns a new API client with the default configuration
func New(cfgs ...ConfigFunc) *Client {
	c := &Client{
		http.Client{Transport: http.DefaultTransport},
		*NewIdentity(),
		*NewCatalog(),
		*NewMarketplace(),
	}

	c.IdentityClient.common.backend.(*defaultBackend).client = &c.client
	c.CatalogClient.common.backend.(*defaultBackend).client = &c.client
	c.MarketplaceClient.common.backend.(*defaultBackend).client = &c.client

	ForURLPattern(DefaultURLPattern)(c)
	for _, cfg := range cfgs {
		cfg(c)
	}

	return c
}

// ConfigFunc is a func that configures the client during New
type ConfigFunc func(*Client)

// ForURLPattern returns a configuration func to set the URL pattern for all
// endpoints.
func ForURLPattern(pattern string) ConfigFunc {
	return func(c *Client) {
		c.IdentityClient.common.backend.(*defaultBackend).base = fmt.Sprintf(pattern, "identity")
		c.CatalogClient.common.backend.(*defaultBackend).base = fmt.Sprintf(pattern, "catalog")
		c.MarketplaceClient.common.backend.(*defaultBackend).base = fmt.Sprintf(pattern, "marketplace")
	}
}

// WithAPIKey returns a configuration func to set the API key to use for
// authentication.
func WithAPIKey(key string) ConfigFunc {
	return func(c *Client) {
		ot := c.client.Transport
		c.client.Transport = rtFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Authorization", "Bearer "+key)
			return ot.RoundTrip(r)
		})
	}
}

type rtFunc func(*http.Request) (*http.Response, error)

func (rt rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return rt(r) }

// Login logs a user in to Manifold using the provided email and password. It
// returns the user's JWT auth token on success. This token is not stored on the
// API client; you must instantiate a new one to use it.
func (c *IdentityClient) Login(ctx context.Context, email, password string) (string, error) {
	lt, err := c.Tokens.CreateLogin(ctx, &LoginTokenRequest{Email: email})
	if err != nil {
		return "", err
	}

	salt, err := base64.NewFromString(lt.Salt)
	if err != nil {
		return "", err
	}

	_, pk, err := deriveKeypair(password, salt)
	if err != nil {
		return "", err
	}

	sig := sign(pk, lt.Token).String()
	at, err := c.Tokens.CreateAuth(ctx, "Bearer "+lt.Token, &AuthTokenRequest{
		Type:          "auth",
		LoginTokenSig: sig,
	})
	if err != nil {
		return "", err
	}

	return at.Body.Token, nil
}