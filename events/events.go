package events

import (
	"encoding/json"
	"fmt"
	"reflect"
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

	// TypeResourceProjectChanged represents a move operation
	TypeResourceProjectChanged Type = "resource.project.changed"

	// TypeResourceOwnerChanged represents a transfer operation
	TypeResourceOwnerChanged Type = "resource.owner.changed"

	// TypeOperationFailed represents a failed operation
	TypeOperationFailed Type = "operation.failed"

	// TypeResourceMeasuresAdded represents a change on resource usage
	TypeResourceMeasuresAdded = "resource.measures.added"

	// TypeResourceMeasuresFailed represents a failed change to resource usage
	TypeResourceMeasuresFailed = "resource.measures.failed"

	// TypeAccountStatusUpdated represents an account state update
	TypeAccountStatusUpdated = "account.status.updated"
)

// SourceType represents where the request came from.
type SourceType string

// Validate whether source type is valid
func (s SourceType) Validate(interface{}) error {
	switch s {
	case SourceDashboard, SourceCLI, SourceSystem, SourceProvider:
		return nil
	default:
		return manifold.NewError(errors.BadRequestError, fmt.Sprintf("invalid source type %q", s))
	}
}

const (
	// SourceDashboard is a request coming from the dashboard
	SourceDashboard SourceType = "dashboard"

	// SourceCLI is a request coming from the cli
	SourceCLI SourceType = "cli"

	// SourceSystem is an internal request
	SourceSystem SourceType = "system"

	// SourceProvider is a request from provider
	SourceProvider SourceType = "provider"
)

// Event represents meaningful activities performed on the system.
type Event struct {
	ID            manifold.ID `json:"id"`
	StructType    string      `json:"type"`
	StructVersion int         `json:"version"`
	Body          Body        `json:"body"`
}

// New returns a new event without a body.
func New() (*Event, error) {
	id, err := manifold.NewID(idtype.ActivityEvent)
	if err != nil {
		return nil, err
	}

	evt := &Event{
		ID:            id,
		StructType:    "event",
		StructVersion: 1,
	}

	return evt, nil
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

// Analytics returns a property map for analytics consumption.
func (e *Event) Analytics() map[string]interface{} {
	m := analyticsProperties(e.Body)

	scope := e.Body.Scope()

	if scope != nil && scope.ID.Type() == idtype.Team {
		m["team_id"] = scope.ID
		m["team_name"] = scope.Name
	}

	return m
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
	case TypeResourceProjectChanged:
		body = &ResourceProjectChanged{}
	case TypeResourceOwnerChanged:
		body = &ResourceOwnerChanged{}
	case TypeOperationFailed:
		body = &OperationFailed{}
	case TypeResourceMeasuresAdded:
		body = &ResourceMeasuresAdded{}
	case TypeResourceMeasuresFailed:
		body = &ResourceMeasuresFailed{}
	case TypeAccountStatusUpdated:
		body = &AccountStatusUpdated{}
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

	Actor() *Actor
	SetActor(*Actor)

	Scope() *Scope
	SetScope(*Scope)

	RefID() manifold.ID
	SetRefID(manifold.ID)

	CreatedAt() *strfmt.DateTime
	SetCreatedAt(*strfmt.DateTime)

	Source() *string
	SetSource(*string)

	IPAddress() string
	SetIPAddress(string)
}

// BaseBody contains data associated with all events.
type BaseBody struct {
	EventType       Type        `json:"type"`
	StructActor     *Actor      `json:"actor,omitempty"`
	StructScope     *Scope      `json:"scope,omitempty"`
	StructRefID     manifold.ID `json:"ref_id"`
	StructCreatedAt time.Time   `json:"created_at"`
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

// Actor returns the body's Actor
func (b *BaseBody) Actor() *Actor {
	return b.StructActor
}

// SetActor returns the body's Actor
func (b *BaseBody) SetActor(a *Actor) {
	b.StructActor = a
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
	Data *OperationProvisionedData `json:"data"`
}

// OperationProvisionedData holds the event specific data.
type OperationProvisionedData struct {
	Operation Operation `json:"operation"`
	Source    string    `json:"source" analytics:"type"`
	Resource  Resource  `json:"resource"`
	Project   *Project  `json:"project,omitempty"`
	Provider  *Provider `json:"provider,omitempty"`
	Product   *Product  `json:"product,omitempty"`
	Plan      *Plan     `json:"plan,omitempty"`
	Region    *Region   `json:"region,omitempty"`
}

// OperationDeprovisioned represents a deprovision operation event.
type OperationDeprovisioned struct {
	BaseBody
	Data *OperationDeprovisionedData `json:"data"`
}

// OperationDeprovisionedData holds the event specific data.
type OperationDeprovisionedData struct {
	Operation Operation `json:"operation"`
	Source    string    `json:"source" analytics:"type"`
	Resource  Resource  `json:"resource"`
	Project   *Project  `json:"project,omitempty"`
	Provider  *Provider `json:"provider,omitempty"`
	Product   *Product  `json:"product,omitempty"`
	Plan      *Plan     `json:"plan,omitempty"`
	Region    *Region   `json:"region,omitempty"`
}

// OperationResized represents a resize operation event.
type OperationResized struct {
	BaseBody
	Data *OperationResizedData `json:"data"`
}

// OperationResizedData holds the event specific data.
type OperationResizedData struct {
	Operation Operation `json:"operation"`
	Source    string    `json:"source" analytics:"type"`
	Resource  Resource  `json:"resource"`
	Project   *Project  `json:"project,omitempty"`
	Provider  *Provider `json:"provider,omitempty"`
	Product   *Product  `json:"product,omitempty"`
	OldPlan   *Plan     `json:"old_plan,omitempty"`
	NewPlan   *Plan     `json:"new_plan,omitempty"`
	Region    *Region   `json:"region,omitempty"`
}

// ResourceProjectChanged records a move operation event
type ResourceProjectChanged struct {
	BaseBody
	Data *ResourceProjectChangedData `json:"data"`
}

// ResourceProjectChangedData holds the specific move event details
type ResourceProjectChangedData struct {
	Operation  Operation `json:"operation"`
	Source     string    `json:"source" analytics:"type"`
	Resource   Resource  `json:"resource"`
	Provider   *Provider `json:"provider,omitempty"`
	Product    *Product  `json:"product,omitempty"`
	Plan       *Plan     `json:"plan,omitempty"`
	OldProject *Project  `json:"old_project,omitempty"`
	NewProject *Project  `json:"new_project,omitempty"`
	Region     *Region   `json:"region,omitempty"`
}

// ResourceOwnerChanged records a transfer operation event
type ResourceOwnerChanged struct {
	BaseBody
	Data *ResourceOwnerChangedData `json:"data"`
}

// ResourceOwnerChangedData holds the specific transfer event data
type ResourceOwnerChangedData struct {
	Operation Operation `json:"operation"`
	Source    string    `json:"source" analytics:"type"`
	Resource  Resource  `json:"resource"`
	OldUser   *User     `json:"old_user,omitempty"`
	OldTeam   *Team     `json:"old_team,omitempty"`
	NewUser   *User     `json:"new_user,omitempty"`
	NewTeam   *Team     `json:"new_team,omitempty"`
	Plan      *Plan     `json:"plan,omitempty"`
	Provider  *Provider `json:"provider,omitempty"`
	Product   *Product  `json:"product,omitempty"`
	Project   *Project  `json:"project,omitempty"`
	Region    *Region   `json:"region,omitempty"`
}

// OperationFailed represents a resize operation event.
type OperationFailed struct {
	BaseBody
	Data *OperationFailedData `json:"data"`
}

// OperationFailedData holds the event specific data.
type OperationFailedData struct {
	Operation Operation      `json:"operation"`
	Resource  *Resource      `json:"resource,omitempty"`
	Project   *Project       `json:"project,omitempty"`
	Provider  *Provider      `json:"provider,omitempty"`
	Product   *Product       `json:"product,omitempty"`
	Plan      *Plan          `json:"plan,omitempty"`
	Region    *Region        `json:"region,omitempty"`
	Error     OperationError `json:"error"`
}

// ResourceMeasuresAdded represents a change on resource usage.
type ResourceMeasuresAdded struct {
	BaseBody
	Data *ResourceMeasuresAddedData `json:"data"`
}

// ResourceMeasuresAddedData holds the event specific data.
type ResourceMeasuresAddedData struct {
	Resource Resource         `json:"resource"`
	Project  *Project         `json:"project,omitempty"`
	Provider Provider         `json:"provider"`
	Product  Product          `json:"product"`
	Plan     Plan             `json:"plan"`
	Region   Region           `json:"region"`
	Measures map[string]int64 `json:"measures"`
}

// ResourceMeasuresFailed represents a failed change to resource usage.
type ResourceMeasuresFailed struct {
	BaseBody
	Data *ResourceMeasuresFailedData `json:"data"`
}

// ResourceMeasuresFailedData holds the event specific data.
type ResourceMeasuresFailedData struct {
	Resource Resource `json:"resource"`
	Project  *Project `json:"project,omitempty"`
	Provider Provider `json:"provider"`
	Product  Product  `json:"product"`
	Plan     Plan     `json:"plan"`
	Region   Region   `json:"region"`
	Error    Error    `json:"error"`
}

// AccountStatusUpdated represents an account state change
type AccountStatusUpdated struct {
	BaseBody
	Data *AccountStatusUpdatedData `json:"data"`
}

// AccountStatusUpdatedData holds the event specific data.
type AccountStatusUpdatedData struct {
	Reason   string `json:"reason,omitempty"`
	OldState string `json:"old_state"`
	NewState string `json:"new_state"`
}

// Actor represents a simplified version of either a user or a team.
type Actor struct {
	ID    manifold.ID `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email,omitempty"`
}

// Validate returns whether or not the given Actor is valid
func (Actor) Validate(v interface{}) error {
	return nil
}

// Scope represents a simplified version of either a user or a team.
type Scope struct {
	ID    manifold.ID `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email,omitempty"`
}

// Validate returns whether or not the given Scope is valid
func (Scope) Validate(v interface{}) error {
	return nil
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
	ID   manifold.ID `json:"id" analytics:"resource_id"`
	Name string      `json:"name" analytics:"resource_name"`
}

// Validate returns whether or not the given Resource is valid
func (Resource) Validate(v interface{}) error {
	return nil
}

// Project is a simplified version for events data.
type Project struct {
	ID   manifold.ID `json:"id" analytics:"project_id"`
	Name string      `json:"name" analytics:"project_name"`
}

// Validate returns whether or not the given Project is valid
func (Project) Validate(v interface{}) error {
	return nil
}

// Provider is a simplified version for events data.
type Provider struct {
	ID   manifold.ID `json:"id" analytics:"provider_id"`
	Name string      `json:"name" analytics:"provider_name"`
}

// Validate returns whether or not the given Provider is valid
func (Provider) Validate(v interface{}) error {
	return nil
}

// Product is a simplified version for events data.
type Product struct {
	ID   manifold.ID `json:"id" analytics:"product_id"`
	Name string      `json:"name" analytics:"product_name"`
}

// Validate returns whether or not the given Product is valid
func (Product) Validate(v interface{}) error {
	return nil
}

// Plan is a simplified version for events data.
type Plan struct {
	ID   manifold.ID `json:"id" analytics:"plan_id"`
	Name string      `json:"name" analytics:"plan_name"`
	Cost int         `json:"cost" analytics:"plan_cost"`
}

// Validate returns whether or not the given Plan is valid
func (Plan) Validate(v interface{}) error {
	return nil
}

// Region is a simplified version for events data.
type Region struct {
	ID       manifold.ID `json:"id" analytics:"region_id"`
	Name     string      `json:"name" analytics:"region_name"`
	Platform string      `json:"platform" analytics:"region_platform"`
	Location string      `json:"location" analytics:"region_location"`
	Priority float64     `json:"priority" analytics:"region_priority"`
}

// Validate returns whether or not the given Region is valid
func (Region) Validate(v interface{}) error {
	return nil
}

// Error represents an error message.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Validate returns whether or not the given Error is valid
func (e Error) Validate(v interface{}) error {
	if e.Message == "" {
		return fmt.Errorf("invalid error message %q", e.Message)
	}

	if e.Code <= 0 {
		return fmt.Errorf("invalid error code %d", e.Code)
	}

	return nil
}

// analyticsProperties generates property map for analytics consumption by
// finding all struct fields with the analytics tag and copying its value to a
// map
func analyticsProperties(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return m
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() || !field.CanInterface() {
				continue
			}

			mm := analyticsProperties(field.Interface())
			for k, v := range mm {
				m[k] = v
			}
		case reflect.Struct:
			mm := analyticsProperties(field.Addr().Interface())
			for k, v := range mm {
				m[k] = v
			}
		default:
			tag := t.Field(i).Tag.Get("analytics")
			if tag != "" {
				m[tag] = field.Interface()
			}
		}
	}

	return m
}
