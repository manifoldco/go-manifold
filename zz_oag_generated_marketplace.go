package manifold

import (
	context "context"
	fmt "fmt"
	http "net/http"
	url "net/url"
	strconv "strconv"
)

// This file is automatically generated by oag (https://github.com/jbowes/oag)
// DO NOT EDIT

// CreateProject is a data type for API communication.
type CreateProject struct {
	Body CreateProjectBody `json:"body"`
}

// CreateProjectBody is a data type for API communication.
type CreateProjectBody struct {
	UserID      *ID     `json:"user_id"` // Optional
	TeamID      *ID     `json:"team_id"` // Optional
	Name        string  `json:"name"`
	Label       string  `json:"label"`
	Description *string `json:"description"` // Optional
}

// Credential is a data type for API communication.
type Credential struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		ResourceID ID `json:"resource_id"`

		// Map of configuration variable names to values, names
		// must IEEE 1003.1 - 2001 Standard (checked in code).
		Values    map[string]string `json:"values"`
		Source    string            `json:"source"`
		CreatedAt string            `json:"created_at"`
		UpdatedAt string            `json:"updated_at"`
	} `json:"body"`
}

// Project is a data type for API communication.
type Project struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		UserID      *ID     `json:"user_id"` // Optional
		TeamID      *ID     `json:"team_id"` // Optional
		Name        string  `json:"name"`
		Label       string  `json:"label"`
		Description *string `json:"description"` // Optional
	} `json:"body"`
}

// ProjectsListOpts holds optional argument values
type ProjectsListOpts struct {
	Me *bool `json:"me"` // Only list projects with the user as the owner

	// ID of the Team to filter projects by, stored as a
	// base32 encoded 18 byte identifier.
	TeamID *ID     `json:"team_id"`
	Label  *string `json:"label"` // Filter projects by a label, returns one or zero results.
}

// PublicUpdateProject is a data type for API communication.
// Update project patch request
type PublicUpdateProject struct {
	Body PublicUpdateProjectBody `json:"body"`
}

// PublicUpdateProjectBody is a data type for API communication.
type PublicUpdateProjectBody struct {
	Name        *string `json:"name"`        // Optional
	Label       *string `json:"label"`       // Optional
	Description *string `json:"description"` // Optional
}

// PublicUpdateResource is a data type for API communication.
// Shape of request used to update a resource through a public api
type PublicUpdateResource struct {
	Body PublicUpdateResourceBody `json:"body"`
}

// PublicUpdateResourceBody is a data type for API communication.
type PublicUpdateResourceBody struct {
	Name  *string `json:"name"`  // Optional
	Label *string `json:"label"` // Optional
}

// Resource is a data type for API communication.
type Resource struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		Name      string `json:"name"`
		Label     string `json:"label"`
		UserID    *ID    `json:"user_id"`    // Optional
		TeamID    *ID    `json:"team_id"`    // Optional
		ProductID *ID    `json:"product_id"` // Optional
		PlanID    *ID    `json:"plan_id"`    // Optional
		RegionID  *ID    `json:"region_id"`  // Optional
		ProjectID *ID    `json:"project_id"` // Optional
		Source    string `json:"source"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"body"`
}

// ResourcesListOpts holds optional argument values
type ResourcesListOpts struct {
	// ID of the Team to filter Resources by, stored as a
	// base32 encoded 18 byte identifier.
	TeamID *ID `json:"team_id"`

	// ID of the Product to filter Resources by, stored as a
	// base32 encoded 18 byte identifier.
	ProductID *ID `json:"product_id"`

	// ID of the Project to filter Resources by, stored as a
	// base32 encoded 18 byte identifier.
	ProjectID *ID     `json:"project_id"`
	Label     *string `json:"label"` // Filter resources by a label, returns one or zero results.
}

// CredentialIter Iterates over a result set of Credentials.
type CredentialIter struct {
	page []Credential
	i    int

	err   error
	first bool
}

// Close closes the CredentialIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *CredentialIter) Close() {}

// Next advances the CredentialIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *CredentialIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current Credential, and an optional error. Once an error has been returned,
// the CredentialIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *CredentialIter) Current() (*Credential, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// ProjectIter Iterates over a result set of Projects.
type ProjectIter struct {
	page []Project
	i    int

	err   error
	first bool
}

// Close closes the ProjectIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *ProjectIter) Close() {}

// Next advances the ProjectIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *ProjectIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current Project, and an optional error. Once an error has been returned,
// the ProjectIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *ProjectIter) Current() (*Project, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// ResourceIter Iterates over a result set of Resources.
type ResourceIter struct {
	page []Resource
	i    int

	err   error
	first bool
}

// Close closes the ResourceIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *ResourceIter) Close() {}

// Next advances the ResourceIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *ResourceIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current Resource, and an optional error. Once an error has been returned,
// the ResourceIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *ResourceIter) Current() (*Resource, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// CredentialsClient provides access to the /credentials APIs
type CredentialsClient endpoint

// List corresponds to the GET /credentials endpoint.
//
// List credentials
func (c *CredentialsClient) List(ctx context.Context, resourceID []ID) *CredentialIter {
	iter := CredentialIter{
		first: true,
		i:     -1,
	}

	p := "/credentials"

	q := make(url.Values)
	for _, v := range resourceID {
		b, err := v.MarshalText()
		if err != nil {
			return &iter
		}
		q.Add("resource_id", string(b))
	}

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, q, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		switch code {
		case 400, 401, 500:
			return &Error{}
		default:
			return nil
		}
	})
	return &iter
}

// InternalClient provides access to the /internal APIs
type InternalClient endpoint

// DeleteProjects corresponds to the DELETE /internal/projects/:id endpoint.
//
// Internal delete project route
// End-point used by the Provisioning Worker to delete a project.
// provisioning handles the deletion to ensure no concurrent operations are
// using the project.
func (c *InternalClient) DeleteProjects(ctx context.Context, id ID) error {
	idBytes, err := id.MarshalText()
	if err != nil {
		return err
	}

	p := fmt.Sprintf("/internal/projects/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodDelete, p, nil, nil)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		switch code {
		case 400, 401, 404, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return err
	}

	return nil
}

// ProjectsClient provides access to the /projects APIs
type ProjectsClient endpoint

// Create corresponds to the POST /projects endpoint.
//
// Create a new project
func (c *ProjectsClient) Create(ctx context.Context, createProject *CreateProject) (*Project, error) {
	p := "/projects"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, createProject)
	if err != nil {
		return nil, err
	}

	var resp Project
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get corresponds to the GET /projects/:id endpoint.
//
// Retrieve a project by its ID
func (c *ProjectsClient) Get(ctx context.Context, id ID) (*Project, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/projects/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Project
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 401, 403, 404, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List corresponds to the GET /projects endpoint.
//
// List all provisioned projects
func (c *ProjectsClient) List(ctx context.Context, opts *ProjectsListOpts) *ProjectIter {
	iter := ProjectIter{
		first: true,
		i:     -1,
	}

	p := "/projects"

	var q url.Values
	if opts != nil {
		q = make(url.Values)
		if opts.Me != nil {
			q.Set("me", strconv.FormatBool(*opts.Me))
		}

		if opts.TeamID != nil {
			b, err := opts.TeamID.MarshalText()
			if err != nil {
				return &iter
			}
			q.Set("team_id", string(b))
		}

		if opts.Label != nil {
			q.Set("label", *opts.Label)
		}
	}

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, q, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		switch code {
		case 400, 401, 500:
			return &Error{}
		default:
			return nil
		}
	})
	return &iter
}

// Update corresponds to the PATCH /projects/:id endpoint.
//
// Update a project's name or description.
func (c *ProjectsClient) Update(ctx context.Context, id ID, publicUpdateProject *PublicUpdateProject) (*Project, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/projects/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodPatch, p, nil, publicUpdateProject)
	if err != nil {
		return nil, err
	}

	var resp Project
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 400, 401, 403, 409, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ResourcesClient provides access to the /resources APIs
type ResourcesClient endpoint

// Get corresponds to the GET /resources/:id endpoint.
//
// Retrieve a resource
func (c *ResourcesClient) Get(ctx context.Context, id ID) (*Resource, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/resources/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Resource
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 401, 404, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetConfig corresponds to the GET /resources/:id/config endpoint.
//
// Get the custom config for a resource
// Get the custom config for a resource. Only resources with a `source` of
// `custom` have custom configs. Requesting the custom config for a resource
// of any other source will result in a `400` response.
func (c *ResourcesClient) GetConfig(ctx context.Context, id ID) (*map[string]string, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/resources/%s/config", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp map[string]string
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 400, 401, 404, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List corresponds to the GET /resources/ endpoint.
//
// List all provisioned resources
func (c *ResourcesClient) List(ctx context.Context, opts *ResourcesListOpts) *ResourceIter {
	iter := ResourceIter{
		first: true,
		i:     -1,
	}

	p := "/resources/"

	var q url.Values
	if opts != nil {
		q = make(url.Values)
		if opts.TeamID != nil {
			b, err := opts.TeamID.MarshalText()
			if err != nil {
				return &iter
			}
			q.Set("team_id", string(b))
		}

		if opts.ProductID != nil {
			b, err := opts.ProductID.MarshalText()
			if err != nil {
				return &iter
			}
			q.Set("product_id", string(b))
		}

		if opts.ProjectID != nil {
			b, err := opts.ProjectID.MarshalText()
			if err != nil {
				return &iter
			}
			q.Set("project_id", string(b))
		}

		if opts.Label != nil {
			q.Set("label", *opts.Label)
		}
	}

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, q, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		switch code {
		case 401, 500:
			return &Error{}
		default:
			return nil
		}
	})
	return &iter
}

// Update corresponds to the PATCH /resources/:id endpoint.
//
// Update a resource name or other property that doesn't require
// communication with the provider.
func (c *ResourcesClient) Update(ctx context.Context, id ID, publicUpdateResource *PublicUpdateResource) (*Resource, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/resources/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodPatch, p, nil, publicUpdateResource)
	if err != nil {
		return nil, err
	}

	var resp Resource
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 400, 401, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateConfig corresponds to the PATCH /resources/:id/config endpoint.
//
// Modify the custom config for a resource
// Modify the custom config for a resource. Only resources with a `source`
// of `custom` have custom configs. Attempting to modify the custom config
// for a resource of any other source will result in a `400` response.
//
// Custom config is modified via JSON merge patching. To change or modify a
// key, include that key with the new value. To remove a key, include that
// key in the request with a null value. Any keys not included in the
// request are ignored and left unchanged.
func (c *ResourcesClient) UpdateConfig(ctx context.Context, id ID, configPatch *map[string]interface{}) (*map[string]string, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/resources/%s/config", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodPatch, p, nil, configPatch)
	if err != nil {
		return nil, err
	}

	var resp map[string]string
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 400, 401, 404, 409, 500:
			return &Error{}
		default:
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// MarketplaceClient is an API client for all endpoints.
type MarketplaceClient struct {
	common endpoint // Reuse a single struct instead of allocating one for each endpoint on the heap.

	Credentials *CredentialsClient
	Internal    *InternalClient
	Projects    *ProjectsClient
	Resources   *ResourcesClient
}

// NewMarketplace returns a new MarketplaceClient with the default configuration.
func NewMarketplace() *MarketplaceClient {
	c := &MarketplaceClient{}
	c.common.backend = DefaultBackend()

	c.Credentials = (*CredentialsClient)(&c.common)
	c.Internal = (*InternalClient)(&c.common)
	c.Projects = (*ProjectsClient)(&c.common)
	c.Resources = (*ResourcesClient)(&c.common)

	return c
}