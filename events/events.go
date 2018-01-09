package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
)

// Type represents the different types of events.
type Type string

const (
	// TypeOperationProvisioned represents a provision operation
	TypeOperationProvisioned Type = "operation.provisioned"

	// TypeOperationDeprovisioned represents a deprovision operation
	TypeOperationDeprovisioned Type = "operation.deprovisioned"

	// TypeOperationResized represents a provision operation
	TypeOperationResized Type = "operation.resized"

	// TypeOperationFailed represents a failed operation
	TypeOperationFailed Type = "operation.failed"
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
		return manifold.NewError(errors.BadRequestError, "Expected version to be 1")
	}

	return e.Body.Validate(v)
}

// StateType returns the event state type as a string
func (e *Event) StateType() string {
	return string(e.Body.Type())
}

// GetState returns the event's state
func (e *Event) GetState() string {
	return string(e.Body.State())
}

// SetState sets the event's state
func (e *Event) SetState(state string) {
	e.Body.SetState(State(state))
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
	o := outEvent{}
	err := json.Unmarshal(b, &o)
	if err != nil {
		return err
	}

	v := BaseBody{}
	err = json.Unmarshal(o.Body, &v)
	if err != nil {
		return err
	}

	var body Body
	switch v.Type() {
	case TypeOperationProvisioned:
		body = &OperationProvisioned{}
	case TypeOperationDeprovisioned:
		body = &OperationDeprovisioned{}
	case TypeOperationResized:
		body = &OperationResized{}
	case TypeOperationFailed:
		body = &OperationFailed{}
	default:
		return fmt.Errorf("Unrecognized Operation Type: %s", v.Type())
	}

	err = json.Unmarshal(o.Body, body)
	if err != nil {
		return err
	}

	e.ID = o.ID
	e.StructVersion = o.StructVersion
	e.StructType = o.StructType
	e.Body = body

	return nil
}

// Body represents methods all Events must implement.
type Body interface {
	Validate(interface{}) error

	Type() Type
	SetType(string)

	State() State
	SetState(State)

	ActorID() manifold.ID
	SetActorID(manifold.ID)

	Actor() *Actor
	SetActor(*Actor)

	ScopeID() manifold.ID
	SetScopeID(manifold.ID)

	Scope() *Scope
	SetScope(*Scope)

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
	StructState     State       `json:"state"`
	StructActorID   manifold.ID `json:"actor_id"`
	StructActor     *Actor      `json:"actor,omitempty"`
	StructScopeID   manifold.ID `json:"scope_id"`
	StructScope     *Scope      `json:"scope,omitempty"`
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

// State returns the body's State
func (b *BaseBody) State() State {
	return b.StructState
}

// SetState sets the body's State
func (b *BaseBody) SetState(s State) {
	b.StructState = s
}

// ActorID returns the body's ActorID
func (b *BaseBody) ActorID() manifold.ID {
	return b.StructActorID
}

// SetActorID sets the body's ActorID
func (b *BaseBody) SetActorID(id manifold.ID) {
	b.StructActorID = id
}

// Actor returns the body's Actor
func (b *BaseBody) Actor() *Actor {
	return b.StructActor
}

// SetActor returns the body's Actor
func (b *BaseBody) SetActor(a *Actor) {
	b.StructActor = a
}

// ScopeID returns the body's ScopeID
func (b *BaseBody) ScopeID() manifold.ID {
	return b.StructScopeID
}

// SetScopeID sets the body's ScopeID
func (b *BaseBody) SetScopeID(id manifold.ID) {
	b.StructScopeID = id
}

// Scope returns the body's Scope
func (b *BaseBody) Scope() *Scope {
	return b.StructScope
}

// SetScope returns the body's Scope
func (b *BaseBody) SetScope(s *Scope) {
	b.StructScope = s
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

// OperationDeprovisioned represents a deprovision operation event.
type OperationDeprovisioned struct {
	BaseBody
	Data OperationDeprovisionedData `json:"data"`
}

// OperationDeprovisionedData holds the event specific data.
type OperationDeprovisionedData struct {
	OperationID manifold.ID `json:"operation_id"`

	UserID *manifold.ID `json:"user_id,omitempty"`
	User   *User        `json:"user,omitempty"`

	TeamID *manifold.ID `json:"team_id,omitempty"`
	Team   *Team        `json:"team,omitempty"`
}

// OperationResized represents a resize operation event.
type OperationResized struct {
	BaseBody
	Data OperationResizedData `json:"data"`
}

// OperationResizedData holds the event specific data.
type OperationResizedData struct {
	OperationID manifold.ID `json:"operation_id"`

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

	OldPlanID manifold.ID `json:"old_plan_id"`
	OldPlan   *Plan       `json:"old_plan,omitempty"`

	NewPlanID manifold.ID `json:"new_plan_id"`
	NewPlan   *Plan       `json:"new_plan,omitempty"`

	RegionID *manifold.ID `json:"region_id,omitempty"`
	Region   *Region      `json:"region,omitempty"`
}

// OperationFailed represents a resize operation event.
type OperationFailed struct {
	BaseBody
	Data OperationFailedData `json:"data"`
}

// OperationFailedData holds the event specific data.
type OperationFailedData struct {
	OperationID manifold.ID `json:"operation_id"`
	Operation   *Operation  `json:"operation"`

	ResourceID *manifold.ID `json:"resource_id,omitempty"`
	Resource   *Resource    `json:"resource,omitempty"`

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

	Error OperationError `json:"error"`
}

// Actor represents a simplified version either a user or a team.
type Actor struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Scope represents a simplified version either a user or a team.
type Scope struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// User is a simplified version for events data.
type User struct {
	ID    manifold.ID `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
}

// Validate returns whether or not the given User is valid
func (User) Validate(v interface{}) error {
	return nil
}

// Team is a simplified version for events data.
type Team struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Validate returns whether or not the given Team is valid
func (Team) Validate(v interface{}) error {
	return nil
}

// Operation is a simplified version for events data.
type Operation struct {
	ID        manifold.ID `json:"id"`
	Type      string      `json:"type"`
	State     string      `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Validate returns whether or not the given Operation is valid
func (Operation) Validate(v interface{}) error {
	return nil
}

// OperationError represents an error message.
type OperationError struct {
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Attempt   int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate returns whether or not the given OperationError is valid
func (OperationError) Validate(v interface{}) error {
	return nil
}

// Resource is a simplified version for events data.
type Resource struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Validate returns whether or not the given Resource is valid
func (Resource) Validate(v interface{}) error {
	return nil
}

// Project is a simplified version for events data.
type Project struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Validate returns whether or not the given Project is valid
func (Project) Validate(v interface{}) error {
	return nil
}

// Provider is a simplified version for events data.
type Provider struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Validate returns whether or not the given Provider is valid
func (Provider) Validate(v interface{}) error {
	return nil
}

// Product is a simplified version for events data.
type Product struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
}

// Validate returns whether or not the given Product is valid
func (Product) Validate(v interface{}) error {
	return nil
}

// Plan is a simplified version for events data.
type Plan struct {
	ID   manifold.ID `json:"id"`
	Name string      `json:"name"`
	Cost int         `json:"cost"`
}

// Validate returns whether or not the given Plan is valid
func (Plan) Validate(v interface{}) error {
	return nil
}

// Region is a simplified version for events data.
type Region struct {
	ID       manifold.ID `json:"id"`
	Name     string      `json:"name"`
	Platform string      `json:"platform"`
	Location string      `json:"location"`
	Priority float64     `json:"priority"`
}

// Validate returns whether or not the given Region is valid
func (Region) Validate(v interface{}) error {
	return nil
}
