package manifold

import (
	context "context"
	"errors"
)

var (
	// ErrProjectLabelNotFound is the error that is used when the ForLabel
	// function is called but no project is found.
	ErrProjectLabelNotFound = errors.New("Project for given label is not found")
)

// ForLabel will find the project for a specific label.
func (c *ProjectsClient) ForLabel(ctx context.Context, label string, me *bool, teamID *ID) (*Project, error) {
	projectsList := c.List(ctx, &ProjectsListOpts{
		Me:     me,
		TeamID: teamID,
		Label:  ptrString(label),
	})
	defer projectsList.Close()

	for projectsList.Next() {
		project, err := projectsList.Current()
		if err != nil {
			return nil, err
		}

		if project.Body.Label == label {
			return project, nil
		}
	}

	return nil, ErrProjectLabelNotFound
}
