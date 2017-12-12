package events

import (
	json "encoding/json"
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	manifold "github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
)

// Type represents the different types of events.
type Type string

const (
	// TypeUserCreated represents a user creation
	TypeUserCreated Type = "user.created"

	// TypeOperationProvisioned represents a provision operation
	TypeOperationProvisioned Type = "operation.provisioned"

	// TypeResourceCreated represents a resource creation
	TypeResourceCreated Type = "resource.created"
)

// State represents the state of an event. Events usually starts only
// containing raw references to the system and later are processed to expand
// objects information.
type State string

const (
	// StatePending represents the event information is pending expansion. That's
	// accomplished by an asynchronous job. Users don't have access to an event
	// while it is pending.
	StatePending State = "pending"

	// StateTracking represents the event is sending information to a 3rd-party
	// analytics. Users can access an event while is in the tracking state.
	StateTracking State = "tracking"

	// StateDone represents the event information expansion is done. The event
	// is immutable going forward. Users can access a completed event.
	StateDone State = "done"
)

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

// Event represents meaningful activities performed on the system.
type Event struct {
	ID            manifold.ID `json:"id"`
	StructType    string      `json:"type"`
	StructVersion int         `json:"version"`
	State         State       `json:"state"`
	Body          Body        `json:"body"`
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

// StateType returns the event state type as a string
func (e *Event) StateType() string {
	return string(e.Body.Type())
}

// GetState returns the event's state
func (e *Event) GetState() string {
	return string(e.State)
}

// SetState sets the event's state
func (e *Event) SetState(state string) {
	e.State = State(state)
}

// SetUpdatedAt sets the event's updated at time to the current time.
func (e *Event) SetUpdatedAt() {
	e.Body.SetUpdatedAt()
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

	v := BaseBody{}
	err = json.Unmarshal(t.Body, &v)
	if err != nil {
		return err
	}

	var body Body
	switch v.Type() {
	case TypeOperationProvisioned:
		body = &OperationProvisioned{}
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

// Body represents methods all Events must implement.
type Body interface {
	Validate(interface{}) error

	Type() Type
	SetType(string)

	ActorID() manifold.ID
	SetActorID(manifold.ID)

	ScopeID() manifold.ID
	SetScopeID(manifold.ID)

	RefID() manifold.ID
	SetRefID(manifold.ID)

	CreatedAt() *strfmt.DateTime
	SetCreatedAt(*strfmt.DateTime)

	UpdatedAt() *strfmt.DateTime
	SetUpdatedAt()

	Source() *string
	SetSource(*string)

	IPAddress() string
	SetIPAddress(string)
}

// BaseBody contains data associated with all events.
type BaseBody struct {
	EventType       Type        `json:"type"`
	StructActorID   manifold.ID `json:"actor_id"`
	StructScopeID   manifold.ID `json:"scope_id"`
	StructRefID     manifold.ID `json:"ref_id"`
	StructCreatedAt time.Time   `json:"created_at"`
	StructUpdatedAt time.Time   `json:"updated_at"`
	StructSource    SourceType  `json:"source"`
	StructIPAddress string      `json:"ip_address"`
}

// Validate returns an error if the BaseEventBody is not valid
func (b *BaseBody) Validate(v interface{}) error {
	// TODO: luiz - add validation
	return nil
}

// Type returns the body's EventType
func (b *BaseBody) Type() Type {
	return b.EventType
}

// SetType sets the body's EventType
func (b *BaseBody) SetType(s string) {
	b.EventType = Type(s)
}

// ActorID returns the body's ActorID
func (b *BaseBody) ActorID() manifold.ID {
	return b.StructActorID
}

// SetActorID sets the body's ActorID
func (b *BaseBody) SetActorID(id manifold.ID) {
	b.StructActorID = id
}

// ScopeID returns the body's ScopeID
func (b *BaseBody) ScopeID() manifold.ID {
	return b.StructScopeID
}

// SetScopeID sets the body's ScopeID
func (b *BaseBody) SetScopeID(id manifold.ID) {
	b.StructScopeID = id
}

// RefID returns the body's RefID
func (b *BaseBody) RefID() manifold.ID {
	return b.StructRefID
}

// SetRefID sets the body's RefID
func (b *BaseBody) SetRefID(id manifold.ID) {
	b.StructRefID = id
}

// CreatedAt returns the body's CreatedAt
func (b *BaseBody) CreatedAt() *strfmt.DateTime {
	t := strfmt.DateTime(b.StructCreatedAt)
	return &t
}

// SetCreatedAt sets the body's CreatedAt
func (b *BaseBody) SetCreatedAt(t *strfmt.DateTime) {
	if t == nil {
		b.StructCreatedAt = time.Now().UTC()
	} else {
		b.StructCreatedAt = time.Time(*t)
	}
}

// UpdatedAt returns the body's CreatedAt
func (b *BaseBody) UpdatedAt() *strfmt.DateTime {
	t := strfmt.DateTime(b.StructUpdatedAt)
	return &t
}

// SetUpdatedAt sets the body's CreatedAt
func (b *BaseBody) SetUpdatedAt() {
	b.StructUpdatedAt = time.Now().UTC()
}

// Source returns the body's Source
func (b *BaseBody) Source() *string {
	s := string(b.StructSource)
	return &s
}

// SetSource sets the body's Source
func (b *BaseBody) SetSource(s *string) {
	if s == nil {
		b.StructSource = SourceSystem
	} else {
		b.StructSource = SourceType(*s)
	}
}

// IPAddress returns the body's IPAddress
func (b *BaseBody) IPAddress() string {
	return b.StructIPAddress
}

// SetIPAddress sets the body's IPAddress
func (b *BaseBody) SetIPAddress(ip string) {
	b.StructIPAddress = ip
}

// OperationProvisioned represents a provision operation event.
type OperationProvisioned struct {
	BaseBody
	Data OperationProvisionedData `json:"data"`
}

// OperationProvisionedData holds the event specific data.
type OperationProvisionedData struct {
	OperationID manifold.ID `json:"operation_id"`
	Source      string      `json:"source"`

	ResourceID manifold.ID `json:"resource_id"`
	Resource   *Resource   `json:"resource,omitempty"`

	UserID *manifold.ID `json:"user_id,omitempty"`
	User   *User        `json:"user,omitempty"`

	TeamID *manifold.ID `json:"team_id,omitempty"`
	Team   *Team        `json:"team,omitempty"`

	ProjectID *manifold.ID `json:"project_id,omitempty"`
	Project   *Project     `json:"project,omitempty"`

	ProviderID *manifold.ID `json:"provider_id,omitempty"`
	Provider   *Provider    `json:"provider,omitempty"`

	ProductID *manifold.ID `json:"product_id,omitempty"`
	Product   *Product     `json:"product,omitempty"`

	PlanID *manifold.ID `json:"plan_id,omitempty"`
	Plan   *Plan        `json:"plan,omitempty"`

	RegionID *manifold.ID `json:"region_id,omitempty"`
	Region   *Region      `json:"region,omitempty"`
}

// Resource is a simplified version for events data.
type Resource struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// User is a simplified version for events data.
type User struct {
	ID    manifold.ID `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
}

// Team is a simplified version for events data.
type Team struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Project is a simplified version for events data.
type Project struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Provider is a simplified version for events data.
type Provider struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Product is a simplified version for events data.
type Product struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Plan is a simplified version for events data.
type Plan struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
	Cost int         `json:"cost"`
}

// Region is a simplified version for events data.
type Region struct {
	ID       manifold.ID `json:"id"`
	Name     string      `json:"name"`
	Platform string      `json:"platform"`
	Location string      `json:"location"`
	Priority float64     `json:"priority"`
}
