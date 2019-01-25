// Package idtype contains our enumeration of all registered types.
package idtype

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gobuffalo/flect"
)

// Type is the enumerated list of all registered types, global within the
// marketplace. Types are encoded within 12 bits in an ID, allowing for
// 4096 distinct types.
type Type uint16

// All of our public types. They are explicitly given their hex values,
// as these should not change.
//
// ############################################################################
// HOW TO ADD A NEW TYPE:
// Find the group that it fits in with.
// Add it there, with the next available number in sequence for the group.
//
// If there is no group that fits, create a new group section.
// Create your new section allowing for 100 items in the previous section.
// ie 0x64 + the value of the first item in the previous section. This gives
// us room to add more items sequentially to each section.
//
// If we run out of space to add a new 100 group block, then we can start
// filling in 50 group blocks from the start, and so on.
//
// It is a lot of work, but it lets us keep everything neat and tidy!
// ############################################################################
//
var (
	// users, orgs, teams, invites
	User                Type = 0x000 // User object
	ForgotPasswordToken Type = 0x001 // ForgotPasswordToken object
	Team                Type = 0x002
	TeamMembership      Type = 0x003
	Invite              Type = 0x004

	// Authentication
	Token                  Type = 0x064 // Dashboard Auth
	OAuthCredential        Type = 0x065 // OAuth Client ID/Secret Pair
	OAuthAuthorizationCode Type = 0x066 // OAuth Authorization COde
	OAuthAccessToken       Type = 0x067 // OAuth Access Token
	APIToken               Type = 0x068 // API access tokens

	// Things you can buy, and how you buy them
	Provider Type = 0x0C8
	Product  Type = 0x0C9
	Plan     Type = 0x0CA
	Region   Type = 0x0CB

	// Long-running actions for creating or modifying resources
	Operation Type = 0x12C
	Callback  Type = 0x12D

	// Provisioned resources and related
	Resource   Type = 0x190
	Credential Type = 0x191
	Project    Type = 0x192 // Resources can be grouped into projects

	// Billing objects
	BillingProfile Type = 0x1F4

	Trial             Type = 0x1F5
	SubscriptionEvent Type = 0x1F6
	Invoice           Type = 0x1F7 // Charges for users
	InvoiceEvent      Type = 0x1F8 // record of payment or failure to pay

	// A copy of a SubscriptionEvent recorded in the provider's event stream
	ProviderSubscriptionEvent Type = 0x1F9
	Payout                    Type = 0x1FA // Payouts for providers
	PayoutEvent               Type = 0x1FB // record of payout of failure to payout
	PayoutProfile             Type = 0x1FC // Payout profile for providers

	ActivityEventJob Type = 0x25D
	ActivityEvent    Type = 0x25E

	// Partner objects
	Partner              Type = 0x2C1
	PartnerProfile       Type = 0x2C2
	PartnerProfileAccess Type = 0x2C3

	// Values from 0xF00 to 0x1000 are reserved for Manifold private internal
	// only use.
	ManifoldInternalReserved Type = 0xF00

	// 13 filled bits would push the type into our version. If you reach this
	// point you're adding the 4097th type.
	TypeOverflow = 0x1000
)

// Upper returns the upper byte of the type
func (t Type) Upper() byte {
	o := make([]byte, 2)
	binary.BigEndian.PutUint16(o, uint16(t))

	return o[0]
}

// Lower returns the lower byte of the type
func (t Type) Lower() byte {
	o := make([]byte, 2)
	binary.BigEndian.PutUint16(o, uint16(t))

	return o[1]
}

// Mutable returns whether or not those type is mutable
func (t Type) Mutable() bool {
	defn := getDefn(t)
	return defn.mutable
}

// Collection returns the name for a collection of these types
func (t Type) Collection() string {
	defn := getDefn(t)
	return flect.Pluralize(defn.name)
}

// Name returns the name of a single instance of this type
func (t Type) Name() string {
	defn := getDefn(t)
	return defn.name
}

// String returns the string representation of the primitive type
func (t Type) String() string {
	return t.Name()
}

// Decode decodes a Type from a byte pair
func Decode(upper, lower byte) Type {
	return Type(binary.BigEndian.Uint16([]byte{upper, lower}))
}

func getDefn(t Type) definition {
	defn, ok := definitions[t]
	if !ok {
		panic("Type does not have a definition: " + string(t))
	}

	return defn
}

// Register registers a type into the global namespace, allowing lookup of its
// mutable and name properties from only the bytes that make up the type.
//
// We use this when determining a type from an ID.
// Any types used *must* be registered. Ideally, call this from your package's
// init function.
//
// Register panics if it is called twice with the same Type.
func Register(typ Type, mutable bool, name string) {
	if _, ok := definitions[typ]; ok {
		panic("Type already registered: " + string(typ))
	}
	definitions[typ] = definition{mutable, name}
}

// TypeFromString will return the type from a string interpretation of the type.
// If the type is not found, this will panic.
func TypeFromString(str string) Type {
	for t, d := range definitions {
		if d.name == str {
			return t
		}
	}

	panic("Type not registered")
}

var definitions = map[Type]definition{}

type definition struct {
	mutable bool
	name    string
}

func init() {
	inflections, err := json.Marshal(map[string]string{
		"access": "access",
	})
	if err != nil {
		panic(fmt.Sprintf("unable to encode inflections: %s", err))
	}
	err = flect.LoadInflections(bytes.NewBuffer(inflections))
	if err != nil {
		panic(fmt.Sprintf("unable to load inflections: %s", err))
	}

	Register(User, true, "user")
	Register(ForgotPasswordToken, true, "forgot_password_token")
	Register(Team, true, "team")
	Register(TeamMembership, true, "team_membership")
	Register(Invite, true, "invite")

	Register(Token, true, "token")
	Register(OAuthCredential, true, "oauth_credential")
	Register(OAuthAuthorizationCode, true, "oauth_authorization_code")
	Register(OAuthAccessToken, true, "oauth_access_token")
	Register(APIToken, true, "api_token")

	Register(Provider, true, "provider")
	Register(Product, true, "product")
	Register(Plan, true, "plan")
	Register(Region, true, "region")

	Register(Operation, true, "operation")
	Register(Callback, true, "callback")

	Register(Resource, true, "resource")
	Register(Credential, true, "credential")
	Register(Project, true, "project")

	Register(BillingProfile, true, "billing_profile")
	Register(PayoutProfile, true, "payout_profile")
	Register(Trial, false, "trial")
	Register(SubscriptionEvent, false, "subscription_event")
	Register(Invoice, false, "invoice")
	Register(InvoiceEvent, false, "invoice_event")

	Register(ProviderSubscriptionEvent, false, "provider_subscription_event")
	Register(Payout, false, "payout")
	Register(PayoutEvent, false, "payout_event")

	Register(ActivityEventJob, true, "event_job")
	Register(ActivityEvent, true, "event")

	Register(Partner, true, "partner")
	Register(PartnerProfile, true, "partner_profile")
	Register(PartnerProfileAccess, true, "partner_profile_access")
}
