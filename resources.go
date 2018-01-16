package manifold

import (
	context "context"
	"errors"
)

var (
	// ErrResourceLabelNotFound is the error that is used when the ForLabel
	// function is called but no resource is found.
	ErrResourceLabelNotFound = errors.New("Resource for given label is not found")
)

// ForLabel will find the resource for a specific label.
func (c *ResourcesClient) ForLabel(ctx context.Context, label string, teamID, productID, projectID *ID) (*Resource, error) {
	resourcesList := c.List(ctx, &ResourcesListOpts{
		TeamID:    teamID,
		ProductID: productID,
		ProjectID: projectID,
		Label:     ptrString(label),
	})
	defer resourcesList.Close()

	for resourcesList.Next() {
		resource, err := resourcesList.Current()
		if err != nil {
			return nil, err
		}

		if resource.Body.Label == label {
			return resource, nil
		}
	}

	return nil, ErrResourceLabelNotFound
}
