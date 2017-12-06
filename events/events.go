package events

import (
	json "encoding/json"
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	manifold "github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/marketplace/ptr"
)

// EventType represents the different types of events.
type EventType string

const (
	// EventUserCreated represents a user creation
	EventUserCreated EventType = "user.created"

	// EventResourceProvisioned represents a resource provision
	EventResourceProvisioned EventType = "resource.provisioned"

	// EventResourceCreated represents a resource creation
	EventResourceCreated EventType = "resource.created"
)

// Event represents meaningful activities performed on the system.
type Event struct {
	ID            manifold.ID `json:"id"`
	StructType    string      `json:"type"`
	StructVersion int         `json:"version"`
	Body          EventBody   `json:"body"`
}

// GetID returns the ID associated with this event
func (e *Event) GetID() manifold.ID { return e.ID }

// Version returns the data structure version of this event
func (e *Event) Version() int { return e.StructVersion }

// Type returns the idtype object for this struct type
func (e *Event) Type() idtype.Type { return idtype.ActivityEvent }

// Validate returns whether or not the given Event is valid
func (e *Event) Validate(v interface{}) error {
	if err := e.ID.Validate(v); err != nil {
		return err
	}

	if e.Version() != 1 {
		return manifold.NewError(merrors.BadRequestError, "Expected version to be 1")
	}

	return e.Body.Validate(v)
}

type outEvent struct {
	ID            manifold.ID     `json:"id"`
	StructType    string          `json:"type"`
	StructVersion int             `json:"version"`
	Body          json.RawMessage `json:"body"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for an event
func (e *Event) UnmarshalJSON(b []byte) error {
	t := outEvent{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	v := BaseEventBody{}
	err = json.Unmarshal(t.Body, &v)
	if err != nil {
		return err
	}

	var body EventBody
	switch v.Type() {
	case EventUserCreated:
		body = &UserCreated{}
	case EventResourceProvisioned:
		body = &ResourceProvisioned{}
	case EventResourceCreated:
		body = &ResourceCreated{}
	default:
		return errors.New("Unrecognized Operation Type: " + string(v.Type()))
	}

	err = json.Unmarshal(t.Body, body)
	if err != nil {
		return err
	}

	e.ID = t.ID
	e.StructVersion = t.StructVersion
	e.StructType = t.StructType
	e.Body = body

	return nil
}

// EventBody represents methods all Events must implement.
type EventBody interface {
	Validate(interface{}) error

	Type() EventType
	SetType(string)

	ActorID() manifold.ID
	SetActorID(manifold.ID)

	ScopeID() manifold.ID
	SetScopeID(manifold.ID)

	RefID() manifold.ID
	SetRefID(manifold.ID)

	CreatedAt() *strfmt.DateTime
	SetCreatedAt(*strfmt.DateTime)

	Source() *string
	SetSource(*string)

	IPAddress() string
	SetIPAddress(string)
}

// BaseEventBody contains data associated with all events.
type BaseEventBody struct {
	EventType       EventType   `json:"type"`
	StructActorID   manifold.ID `json:"actor_id"`
	StructScopeID   manifold.ID `json:"scope_id"`
	StructRefID     manifold.ID `json:"ref_id"`
	StructCreatedAt time.Time   `json:"created_at"`
	StructSource    SourceType  `json:"source"`
	StructIPAddress string      `json:"ip_address"`
}

// Validate returns an error if the BaseEventBody is not valid
func (b *BaseEventBody) Validate(v interface{}) error {
	// TODO: luiz - add validation
	return nil
}

// Type returns the body's EventType
func (b *BaseEventBody) Type() EventType {
	return b.EventType
}

// SetType sets the body's EventType
func (b *BaseEventBody) SetType(s string) {
	b.EventType = EventType(s)
}

// ActorID returns the body's ActorID
func (b *BaseEventBody) ActorID() manifold.ID {
	return b.StructActorID
}

// SetActorID sets the body's ActorID
func (b *BaseEventBody) SetActorID(id manifold.ID) {
	b.StructActorID = id
}

// ScopeID returns the body's ScopeID
func (b *BaseEventBody) ScopeID() manifold.ID {
	return b.StructScopeID
}

// SetScopeID sets the body's ScopeID
func (b *BaseEventBody) SetScopeID(id manifold.ID) {
	b.StructScopeID = id
}

// RefID returns the body's RefID
func (b *BaseEventBody) RefID() manifold.ID {
	return b.StructRefID
}

// SetRefID sets the body's RefID
func (b *BaseEventBody) SetRefID(id manifold.ID) {
	b.StructRefID = id
}

// CreatedAt returns the body's CreatedAt
func (b *BaseEventBody) CreatedAt() *strfmt.DateTime {
	t := strfmt.DateTime(b.StructCreatedAt)
	return &t
}

// SetCreatedAt sets the body's CreatedAt
func (b *BaseEventBody) SetCreatedAt(t *strfmt.DateTime) {
	if t == nil {
		b.StructCreatedAt = time.Now().UTC()
	} else {
		b.StructCreatedAt = time.Time(*t)
	}
}

// Source returns the body's Source
func (b *BaseEventBody) Source() *string {
	return ptr.String(string(b.StructSource))
}

// SetSource sets the body's Source
func (b *BaseEventBody) SetSource(s *string) {
	if s == nil {
		b.StructSource = SourceSystem
	} else {
		b.StructSource = SourceType(*s)
	}
}

// IPAddress returns the body's IPAddress
func (b *BaseEventBody) IPAddress() string {
	return b.StructIPAddress
}

// SetIPAddress sets the body's IPAddress
func (b *BaseEventBody) SetIPAddress(ip string) {
	b.StructIPAddress = ip
}

// ResourceProvisioned represents a resource provision event.
type ResourceProvisioned struct {
	BaseEventBody
	Data ResourceProvisionedData `json:"data"`
}

// ResourceProvisionedData holds the event specific data.
type ResourceProvisionedData struct {
	ResourceID   manifold.ID `json:"resource_id"`
	ResourceName string      `json:"_resource_name"`
	Source       string      `json:"source"`

	UserID    *manifold.ID `json:"user_id,omitempty"`
	UserName  string       `json:"_user_name,omitempty"`
	UserEmail string       `json:"_user_email,omitempty"`

	TeamID   *manifold.ID `json:"team_id,omitempty"`
	TeamName string       `json:"_team_name,omitempty"`

	ProjectID   *manifold.ID `json:"project_id,omitempty"`
	ProjectName string       `json:"_project_name,omitempty"`

	Provider     *manifold.ID `json:"provider_id,omitempty"`
	ProviderName string       `json:"_provider_name,omitempty"`

	ProductID   *manifold.ID `json:"product_id,omitempty"`
	ProductName string       `json:"_product_name,omitempty"`

	PlanID   *manifold.ID `json:"plan_id,omitempty"`
	PlanName string       `json:"_plan_name,omitempty"`
	PlanCost int          `json:"_plan_cost,omitempty"`

	RegionID       *manifold.ID `json:"region_id,omitempty"`
	RegionName     string       `json:"_region_name,omitempty"`
	RegionPlatform string       `json:"_region_platform,omitempty"`
	RegionLocation string       `json:"_region_location,omitempty"`
	RegionPriority int          `json:"_region_priority,omitempty"`
}

// UserCreated represents a user signup event.
type UserCreated struct {
	BaseEventBody
	Data UserCreatedData `json:"data"`
}

// UserCreatedData holds the event specific data.
type UserCreatedData struct {
	UserID   manifold.ID `json:"user_id"`
	Email    string      `json:"email"`
	UserName string      `json:"user_name"`
}

// ResourceCreated represents a resource creation event.
type ResourceCreated struct {
	BaseEventBody
	Data ResourceCreatedData `json:"data"`
}

// ResourceCreatedData holds the event specific data.
type ResourceCreatedData struct {
	Name       string       `json:"name"`
	Label      string       `json:"label"`
	ResourceID manifold.ID  `json:"resource_id"`
	OwnerID    manifold.ID  `json:"owner_id"`
	ProductID  *manifold.ID `json:"product_id,omitempty"`
	PlanID     *manifold.ID `json:"plan_id,omitempty"`
	RegionID   *manifold.ID `json:"region_id,omitempty"`
	ProjectID  *manifold.ID `json:"project_id,omitempty"`
	Source     string       `json:"source"`
}

// SourceType represents where the request came from.
type SourceType string

const (
	// SourceDashboard is a request coming from the dashboard
	SourceDashboard SourceType = "dashboard"

	// SourceCLI is a request coming from the cli
	SourceCLI SourceType = "cli"

	// SourceSystem is an internal request
	SourceSystem SourceType = "system"
)
