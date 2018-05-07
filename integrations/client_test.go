package integrations_test

import (
	"context"
	"os"
	"testing"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

var testClient *integrations.Client

func TestGetResource(t *testing.T) {
	ctx := context.Background()

	t.Run("without a valid resource", func(t *testing.T) {
		invalidResource := &primitives.Resource{}

		_, err := testClient.GetResource(ctx, nil, invalidResource)
		expectErrorEqual(t, err, integrations.ErrResourceInvalid)
	})

	t.Run("with a valid resource", func(t *testing.T) {
		t.Run("with a non-existing project", func(t *testing.T) {
			resource := &primitives.Resource{
				Name: "custom-resource1",
			}

			_, err := testClient.GetResource(ctx, strPtr("non-existing"), resource)
			expectErrorEqual(t, err, integrations.ErrProjectNotFound)
		})

		t.Run("with an existing project", func(t *testing.T) {
			project := strPtr("kubernetes-secrets")

			t.Run("with an existing resource", func(t *testing.T) {
				resource := &primitives.Resource{
					Name: "custom-resource1",
				}

				res, err := testClient.GetResource(ctx, project, resource)
				expectNoError(t, err)
				if l := "custom-resource1"; res.Body.Label != l {
					t.Fatalf("Expected label to equal '%s', got '%s'", l, res.Body.Label)
				}
			})

			t.Run("with a non existing resource", func(t *testing.T) {
				resource := &primitives.Resource{
					Name: "non-existing-resource",
				}

				_, err := testClient.GetResource(ctx, project, resource)
				expectErrorEqual(t, err, integrations.ErrResourceNotFound)
			})
		})
	})
}

func TestGetResources(t *testing.T) {
	ctx := context.Background()

	t.Run("with an invalid resource", func(t *testing.T) {
		invalidResource := &primitives.Resource{}

		_, err := testClient.GetResources(ctx, nil, []*primitives.Resource{invalidResource})
		expectErrorEqual(t, err, integrations.ErrResourceInvalid)
	})

	t.Run("with valid resources", func(t *testing.T) {
		resources := []*primitives.Resource{
			{
				Name: "custom-resource1",
			},
			{
				Name: "custom-resource2",
			},
		}

		t.Run("with a non-existing project", func(t *testing.T) {
			_, err := testClient.GetResources(ctx, strPtr("non-existing"), resources)
			expectErrorEqual(t, err, integrations.ErrProjectNotFound)
		})

		t.Run("with an existing project", func(t *testing.T) {
			project := strPtr("kubernetes-secrets")

			t.Run("with one non-existing resource", func(t *testing.T) {
				nonExisting := &primitives.Resource{
					Name: "non-existing",
				}
				nr := append(resources, nonExisting)
				_, err := testClient.GetResources(ctx, project, nr)
				expectErrorEqual(t, err, integrations.ErrResourceNotFound)
			})

			t.Run("with all existing resources", func(t *testing.T) {
				res, err := testClient.GetResources(ctx, project, resources)
				expectNoError(t, err)
				if len(res) != 2 {
					t.Fatalf("Expected '2' resources to be loaded, got '%d'", len(res))
				}
			})
		})
	})
}

func TestGetProjectCredentialValues(t *testing.T) {
	ctx := context.Background()

	t.Run("with an invalid project", func(t *testing.T) {
		invalidProject := &primitives.Project{}

		_, err := testClient.GetProjectCredentialValues(ctx, invalidProject)
		expectErrorEqual(t, err, integrations.ErrProjectInvalid)
	})

	t.Run("with a valid project", func(t *testing.T) {
		project := &primitives.Project{
			Name: "kubernetes-secrets",
		}

		t.Run("without a valid credentials subset", func(t *testing.T) {
			creds, err := testClient.GetProjectCredentialValues(ctx, project)
			expectNoError(t, err)

			if len(creds) != 2 {
				t.Fatalf("Expected '2' CredentialValues for 'kubernetes-secrets', got '%d'", len(creds))
			}

			data, err := integrations.FlattenResourcesCredentialValues(creds)
			expectNoError(t, err)

			expectStringEqual(t, data["TOKEN_SECRET"], "my-secret-token-secret")
			expectStringEqual(t, data["TOKEN_ID"], "my-secret-token-id")
			expectStringEqual(t, data["PASSWORD"], "manifold-secret")
			expectStringEqual(t, data["USERNAME"], "manifold")
		})

		t.Run("with a valid credentials subset", func(t *testing.T) {
			sub := &primitives.Resource{
				Name: "custom-resource1",
				Credentials: []*primitives.Credential{
					{
						Key: "TOKEN_ID",
					},
				},
			}

			subProj := &primitives.Project{
				Name:      project.Name,
				Resources: []*primitives.Resource{sub},
			}

			creds, err := testClient.GetProjectCredentialValues(ctx, subProj)
			expectNoError(t, err)

			if len(creds) != 1 {
				t.Fatalf("Expected '1' CredentialValues in project 'kubernetes-secrets' for sub request. Got '%d'", len(creds))
			}

			if len(creds["custom-resource1"]) != 1 {
				t.Fatalf("Expected '1' CredentailValues in resource 'custom-resource1', got '%d'", len(creds["custom-resource1"]))
			}
			data, err := integrations.FlattenResourcesCredentialValues(creds)
			expectNoError(t, err)
			expectStringEqual(t, data["TOKEN_ID"], "my-secret-token-id")

		})

		t.Run("with a non-existing key", func(t *testing.T) {
			t.Run("with a default value", func(t *testing.T) {
				sub := &primitives.Resource{
					Name: "custom-resource1",
					Credentials: []*primitives.Credential{
						{
							Key:     "NON_EXISTING",
							Default: "my-default-value",
						},
					},
				}

				subProj := &primitives.Project{
					Name:      project.Name,
					Resources: []*primitives.Resource{sub},
				}

				creds, err := testClient.GetProjectCredentialValues(ctx, subProj)
				expectNoError(t, err)

				if len(creds) != 1 {
					t.Fatalf("Expected '1' CustomCredentials value for 'kubere-secrets' sub request. Got '%d'", len(creds))
				}

				data, err := integrations.FlattenResourcesCredentialValues(creds)
				expectNoError(t, err)
				expectStringEqual(t, data["NON_EXISTING"], "my-default-value")
			})
			t.Run("without a defalt value", func(t *testing.T) {
				sub := &primitives.Resource{
					Name: "custom-resource1",
					Credentials: []*primitives.Credential{
						{
							Key: "NON_EXISTING",
						},
					},
				}

				subProj := &primitives.Project{
					Name:      project.Name,
					Resources: []*primitives.Resource{sub},
				}

				_, err := testClient.GetProjectCredentialValues(ctx, subProj)
				expectErrorEqual(t, err, integrations.ErrCredentialDefaultNotSet)
			})
		})

		t.Run("with an invalid credentials subset", func(t *testing.T) {
			sub := &primitives.Resource{
				Name: "custom-resource1",
				Credentials: []*primitives.Credential{
					{
						Name: "invalid",
					},
				},
			}

			subProj := &primitives.Project{
				Name:      project.Name,
				Resources: []*primitives.Resource{sub},
			}

			_, err := testClient.GetProjectCredentialValues(ctx, subProj)
			expectErrorEqual(t, err, integrations.ErrProjectInvalid)
		})
	})
}

func TestGetResourceCredentialValues(t *testing.T) {
	ctx := context.Background()

	t.Run("with an invalid resource", func(t *testing.T) {
		invalidResource := &primitives.Resource{}

		_, err := testClient.GetResourceCredentialValues(ctx, nil, invalidResource)
		expectErrorEqual(t, err, integrations.ErrResourceInvalid)
	})

	t.Run("with a valid resource", func(t *testing.T) {
		res := &primitives.Resource{
			Name: "custom-resource1",
		}

		t.Run("with a non-existing project", func(t *testing.T) {
			_, err := testClient.GetResourceCredentialValues(ctx, strPtr("non-existing"), res)
			expectErrorEqual(t, err, integrations.ErrProjectNotFound)
		})

		t.Run("with an existing project", func(t *testing.T) {
			project := strPtr("kubernetes-secrets")

			t.Run("without credentials subset", func(t *testing.T) {
				creds, err := testClient.GetResourceCredentialValues(ctx, project, res)
				expectNoError(t, err)

				if len(creds) != 2 {
					t.Fatalf("Expected '2' CredentialValues for 'custom-resource1', got '%d'", len(creds))
				}
				data, err := integrations.FlattenResourceCredentialValues(creds)
				expectNoError(t, err)

				expectStringEqual(t, data["TOKEN_ID"], "my-secret-token-id")
				expectStringEqual(t, data["TOKEN_SECRET"], "my-secret-token-secret")
			})

			t.Run("with a valid credential subset", func(t *testing.T) {
				sub := &primitives.Resource{
					Name: "custom-resource1",
					Credentials: []*primitives.Credential{
						{
							Key: "TOKEN_ID",
						},
					},
				}

				creds, err := testClient.GetResourceCredentialValues(ctx, project, sub)
				expectNoError(t, err)
				if len(creds) != 1 {
					t.Fatalf("Expected '1' CredentialValues for 'custom-resource1', got '%d'", len(creds))
				}

				data, err := integrations.FlattenResourceCredentialValues(creds)
				expectNoError(t, err)
				expectStringEqual(t, data["TOKEN_ID"], "my-secret-token-id")
			})

			t.Run("with a non existing key", func(t *testing.T) {
				t.Run("with a default value", func(t *testing.T) {
					sub := &primitives.Resource{
						Name: "custom-resource1",
						Credentials: []*primitives.Credential{
							{
								Key:     "NON_EXISTING",
								Default: "my-default-value",
							},
						},
					}

					creds, err := testClient.GetResourceCredentialValues(ctx, project, sub)
					expectNoError(t, err)
					if len(creds) != 1 {
						t.Fatalf("Expected '1' CredentialValues for 'custom-resource1', got '%d'", len(creds))
					}

					data, err := integrations.FlattenResourceCredentialValues(creds)
					expectNoError(t, err)
					expectStringEqual(t, data["NON_EXISTING"], "my-default-value")
				})

				t.Run("without a default value", func(t *testing.T) {
					sub := &primitives.Resource{
						Name: "custom-resource1",
						Credentials: []*primitives.Credential{
							{
								Key: "NON_EXISTING",
							},
						},
					}

					_, err := testClient.GetResourceCredentialValues(ctx, project, sub)
					expectErrorEqual(t, err, integrations.ErrCredentialDefaultNotSet)
				})
			})

			t.Run("with an invalid credential subset", func(t *testing.T) {
				sub := &primitives.Resource{
					Name: "custom-resource1",
					Credentials: []*primitives.Credential{
						{
							Name: "Invalid",
						},
					},
				}

				_, err := testClient.GetResourceCredentialValues(ctx, project, sub)
				expectErrorEqual(t, err, integrations.ErrResourceInvalid)
			})
		})
	})
}

func init() {
	testClient = newClient()
}

func newClient() *integrations.Client {
	c, err := integrations.NewClient(
		manifold.New(
			manifold.WithAPIToken(os.Getenv("MANIFOLD_API_TOKEN")),
		),
		strPtr(os.Getenv("MANIFOLD_TEAM")),
	)

	if err != nil {
		panic("Could not set up the test client: " + err.Error())
	}

	return c
}

func expectErrorEqual(t *testing.T, act, exp error) {
	if act == nil {
		t.Fatalf("Expected error not to be 'nil' but '%s'", exp.Error())
	}

	if exp != act {
		t.Fatalf("Expected error '%s', to equal '%s'", act.Error(), exp.Error())
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error to have occurred, got '%s'", err)
	}
}

func expectStringEqual(t *testing.T, act, exp string) {
	if act != exp {
		t.Fatalf("Expected '%s' to equal '%s'", act, exp)
	}
}

func strPtr(str string) *string {
	return &str
}
