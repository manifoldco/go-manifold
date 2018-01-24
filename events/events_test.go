package events

import (
	"reflect"
	"testing"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
)

func TestAnalytics(t *testing.T) {
	id, _ := manifold.NewID(idtype.Resource)
	tid, _ := manifold.NewID(idtype.Team)

	event := Event{
		Body: &OperationProvisioned{
			BaseBody: BaseBody{
				StructScope: &Scope{
					ID:   tid,
					Name: "manifold",
				},
			},
			Data: &OperationProvisionedData{
				Source: "catalog",
				Resource: &Resource{
					ID:   id,
					Name: "database",
				},
				Project: &Project{
					ID:   id,
					Name: "local",
				},
				ProviderID: &id,
				Provider: &Provider{
					ID:   id,
					Name: "Degraffdb",
				},
				Product: &Product{
					ID:   id,
					Name: "Generator",
				},
				Plan: &Plan{
					ID:   id,
					Name: "Static",
					Cost: 0,
				},
				Region: &Region{
					ID:       id,
					Name:     "US East",
					Location: "US",
					Platform: "AWS",
					Priority: 1,
				},
			},
		},
	}

	expect := map[string]interface{}{
		"type":            "catalog",
		"resource_id":     id,
		"resource_name":   "database",
		"project_id":      id,
		"project_name":    "local",
		"provider_id":     id,
		"provider_name":   "Degraffdb",
		"product_id":      id,
		"product_name":    "Generator",
		"plan_id":         id,
		"plan_name":       "Static",
		"plan_cost":       0,
		"region_id":       id,
		"region_name":     "US East",
		"region_location": "US",
		"region_platform": "AWS",
		"region_priority": float64(1),
		"team_id":         tid,
		"team_name":       "manifold",
	}

	got := event.Analytics()

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Expect properties map to eq %v, got %v", expect, got)
	}
}
