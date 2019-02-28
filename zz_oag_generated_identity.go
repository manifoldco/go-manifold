package manifold

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// This file is automatically generated by oag (https://github.com/jbowes/oag)
// DO NOT EDIT

const baseIdentityURL = "https://api.identity.manifold.co/v1"

// APIToken is a data type for API communication.
type APIToken struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`

	Body *struct {
		FirstFour   string  `json:"first_four"`
		LastFour    string  `json:"last_four"`
		Role        string  `json:"role"`
		UserID      ID      `json:"user_id"`
		TeamID      ID      `json:"team_id"`
		Token       *string `json:"token"` // Optional
		Description string  `json:"description"`
	} `json:"body"` // Optional
}

// APITokenRequest is a data type for API communication.
type APITokenRequest struct {
	Description string `json:"description"`
	Role        string `json:"role"`
	UserID      *ID    `json:"user_id"` // Optional
	TeamID      *ID    `json:"team_id"` // Optional
}

// AcceptInvite is a data type for API communication.
type AcceptInvite struct {
	Token string `json:"token"`
}

// AnalyticsEvent is a data type for API communication.
type AnalyticsEvent struct {
	EventName  string                    `json:"event_name"`
	UserID     ID                        `json:"user_id"`
	Properties *AnalyticsEventProperties `json:"properties"` // Optional
}

// AnalyticsEventProperties is a data type for API communication.
type AnalyticsEventProperties struct {
	Platform *string `json:"platform"` // Optional
}

// AuthToken is a data type for API communication.
type AuthToken struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`

	Body struct {
		Token     string `json:"token"`
		UserID    ID     `json:"user_id"`
		Mechanism string `json:"mechanism"`
	} `json:"body"`
}

// AuthTokenRequest is a data type for API communication.
type AuthTokenRequest struct {
	Type          string `json:"type"`
	LoginTokenSig string `json:"login_token_sig"`
}

// CreateInvite is a data type for API communication.
type CreateInvite struct {
	Body CreateInviteBody `json:"body"`
}

// CreateInviteBody is a data type for API communication.
type CreateInviteBody struct {
	Name   string  `json:"name"`
	Role   *string `json:"role"` // Optional
	Email  string  `json:"email"`
	TeamID ID      `json:"team_id"`
}

// CreateTeam is a data type for API communication.
type CreateTeam struct {
	Body CreateTeamBody `json:"body"`
}

// CreateTeamBody is a data type for API communication.
type CreateTeamBody struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// CreateUser is a data type for API communication.
type CreateUser struct {
	Body CreateUserBody `json:"body"`
}

// CreateUserBody is a data type for API communication.
type CreateUserBody struct {
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	PublicKey LoginPublicKey `json:"public_key"`
}

// ForgotPassword is a data type for API communication.
type ForgotPassword struct {
	Email     string         `json:"email"`
	Token     string         `json:"token"`
	PublicKey LoginPublicKey `json:"public_key"`
}

// ForgotPasswordCreate is a data type for API communication.
type ForgotPasswordCreate struct {
	Email string `json:"email"`
}

// Invite is a data type for API communication.
type Invite struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		TeamID    ID     `json:"team_id"`
		Role      string `json:"role"`
		InvitedBy *ID    `json:"invited_by"` // Optional
		Token     string `json:"token"`
	} `json:"body"`
}

// LoginPublicKey is a data type for API communication.
type LoginPublicKey struct {
	Salt  string `json:"salt"`
	Value string `json:"value"`
	Alg   string `json:"alg"`
}

// LoginTokenRequest is a data type for API communication.
type LoginTokenRequest struct {
	Email string `json:"email"`
}

// LoginTokenResponse is a data type for API communication.
type LoginTokenResponse struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

// MemberProfile is a data type for API communication.
type MemberProfile struct {
	UserID       ID     `json:"user_id"`
	MembershipID ID     `json:"membership_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
}

// PublicInvite is a data type for API communication.
type PublicInvite struct {
	Team Team   `json:"team"`
	Role string `json:"role"`

	Invited struct {
		Name  *string `json:"name"`  // Optional
		Email *string `json:"email"` // Optional
	} `json:"invited"`

	InvitedBy struct {
		Name  *string `json:"name"`  // Optional
		Email *string `json:"email"` // Optional
	} `json:"invited_by"`
}

// Team is a data type for API communication.
type Team struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		Name  string `json:"name"`
		Label string `json:"label"`
	} `json:"body"`
}

// TeamMembership is a data type for API communication.
type TeamMembership struct {
	ID      ID     `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	Body struct {
		UserID ID     `json:"user_id"`
		TeamID ID     `json:"team_id"`
		Role   string `json:"role"`
	} `json:"body"`
}

// TokensListOpts holds optional argument values
type TokensListOpts struct {
	Me *bool `json:"me"` // Only list tokens with the user as the owner

	// ID of the Team to filter tokens by, stored as a
	// base32 encoded 18 byte identifier.
	TeamID *ID `json:"team_id"`
}

// UpdateTeam is a data type for API communication.
type UpdateTeam struct {
	Body UpdateTeamBody `json:"body"`
}

// UpdateTeamBody is a data type for API communication.
type UpdateTeamBody struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// UpdateUser is a data type for API communication.
type UpdateUser struct {
	Body *UpdateUserBody `json:"body"` // Optional
}

// UpdateUserBody is a data type for API communication.
type UpdateUserBody struct {
	Name         *string         `json:"name"`           // Optional
	Email        *string         `json:"email"`          // Optional
	PublicKey    *LoginPublicKey `json:"public_key"`     // Optional
	AuthTokenSig *string         `json:"auth_token_sig"` // Optional
}

// User is a data type for API communication.
type User struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`

	Body struct {
		Name             string         `json:"name"`
		Email            string         `json:"email"`
		PublicKey        LoginPublicKey `json:"public_key"`
		VerificationCode *string        `json:"verification_code"` // Optional
		State            string         `json:"state"`
	} `json:"body"`
}

// VerifyEmail is a data type for API communication.
type VerifyEmail struct {
	Body *VerifyEmailBody `json:"body"` // Optional
}

// VerifyEmailBody is a data type for API communication.
type VerifyEmailBody struct {
	VerificationCode string `json:"verification_code"`
}

// APITokenIter Iterates over a result set of APITokens.
type APITokenIter struct {
	page []APIToken
	i    int

	err   error
	first bool
}

// Close closes the APITokenIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *APITokenIter) Close() {}

// Next advances the APITokenIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *APITokenIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current APIToken, and an optional error. Once an error has been returned,
// the APITokenIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *APITokenIter) Current() (*APIToken, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// InviteIter Iterates over a result set of Invites.
type InviteIter struct {
	page []Invite
	i    int

	err   error
	first bool
}

// Close closes the InviteIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *InviteIter) Close() {}

// Next advances the InviteIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *InviteIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current Invite, and an optional error. Once an error has been returned,
// the InviteIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *InviteIter) Current() (*Invite, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// MemberProfileIter Iterates over a result set of MemberProfiles.
type MemberProfileIter struct {
	page []MemberProfile
	i    int

	err   error
	first bool
}

// Close closes the MemberProfileIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *MemberProfileIter) Close() {}

// Next advances the MemberProfileIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *MemberProfileIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current MemberProfile, and an optional error. Once an error has been returned,
// the MemberProfileIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *MemberProfileIter) Current() (*MemberProfile, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// TeamIter Iterates over a result set of Teams.
type TeamIter struct {
	page []Team
	i    int

	err   error
	first bool
}

// Close closes the TeamIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *TeamIter) Close() {}

// Next advances the TeamIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *TeamIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current Team, and an optional error. Once an error has been returned,
// the TeamIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *TeamIter) Current() (*Team, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// TeamMembershipIter Iterates over a result set of TeamMemberships.
type TeamMembershipIter struct {
	page []TeamMembership
	i    int

	err   error
	first bool
}

// Close closes the TeamMembershipIter and releases any associated resources.
// After Close, any calls to Current will return an error.
func (i *TeamMembershipIter) Close() {}

// Next advances the TeamMembershipIter and returns a boolean indicating if the end has been reached.
// Next must be called before the first call to Current.
// Calls to Current after Next returns false will return an error.
func (i *TeamMembershipIter) Next() bool {
	if i.first && i.err != nil {
		i.first = false
		return true
	}
	i.first = false
	i.i++
	return i.i < len(i.page)
}

// Current returns the current TeamMembership, and an optional error. Once an error has been returned,
// the TeamMembershipIter is closed, or the end of iteration is reached, subsequent calls to Current
// will return an error.
func (i *TeamMembershipIter) Current() (*TeamMembership, error) {
	if i.err != nil {
		return nil, i.err
	}
	return &i.page[i.i], nil
}

// AnalyticsClient provides access to the /analytics APIs
type AnalyticsClient endpoint

// Create corresponds to the POST /analytics/ endpoint.
//
// An endpoint used by the cli to push analytics into segment.
func (c *AnalyticsClient) Create(ctx context.Context, analyticsEvent *AnalyticsEvent) error {
	p := "/analytics/"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, analyticsEvent)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		switch code {
		case 400, 401, 500:
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

// InvitesClient provides access to the /invites APIs
type InvitesClient endpoint

// Create corresponds to the POST /invites endpoint.
//
// Create a new invite
func (c *InvitesClient) Create(ctx context.Context, createInvite *CreateInvite) (*Invite, error) {
	p := "/invites"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, createInvite)
	if err != nil {
		return nil, err
	}

	var resp Invite
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateAccept corresponds to the POST /invites/accept endpoint.
//
// Accept an invite
func (c *InvitesClient) CreateAccept(ctx context.Context, acceptInvite *AcceptInvite) error {
	p := "/invites/accept"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, acceptInvite)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete corresponds to the DELETE /invites/:id/ endpoint.
//
// Revoke an existing invite
func (c *InvitesClient) Delete(ctx context.Context, id ID) error {
	idBytes, err := id.MarshalText()
	if err != nil {
		return err
	}

	p := fmt.Sprintf("/invites/%s/", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodDelete, p, nil, nil)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// Get corresponds to the GET /invites/:token/ endpoint.
//
// Retrieve an invite's details by token
func (c *InvitesClient) Get(ctx context.Context, token string) (*PublicInvite, error) {
	p := fmt.Sprintf("/invites/%s/", token)

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp PublicInvite
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List corresponds to the GET /invites endpoint.
//
// List of invites for the user
func (c *InvitesClient) List(ctx context.Context, teamID ID) *InviteIter {
	iter := InviteIter{
		first: true,
		i:     -1,
	}

	p := "/invites"

	q := make(url.Values)
	teamIDBytes, err := teamID.MarshalText()
	if err != nil {
		return &iter
	}
	q.Set("team_id", string(teamIDBytes))

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, q, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		return &Error{}
	})
	return &iter
}

// MembershipsClient provides access to the /memberships APIs
type MembershipsClient endpoint

// Delete corresponds to the DELETE /memberships/:id endpoint.
//
// Remove a member from a team
func (c *MembershipsClient) Delete(ctx context.Context, id ID) error {
	idBytes, err := id.MarshalText()
	if err != nil {
		return err
	}

	p := fmt.Sprintf("/memberships/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodDelete, p, nil, nil)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// List corresponds to the GET /memberships endpoint.
//
// List memberships for the user
func (c *MembershipsClient) List(ctx context.Context) *TeamMembershipIter {
	iter := TeamMembershipIter{
		first: true,
		i:     -1,
	}

	p := "/memberships"

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		return &Error{}
	})
	return &iter
}

// SelfClient provides access to the /self APIs
type SelfClient endpoint

// Get corresponds to the GET /self endpoint.
//
// Retrieve the underlying user represented by a Token
func (c *SelfClient) Get(ctx context.Context) (*User, error) {
	p := "/self"

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp User
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		switch code {
		case 401, 500:
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

// TeamsClient provides access to the /teams APIs
type TeamsClient endpoint

// Create corresponds to the POST /teams endpoint.
//
// Create a new team
func (c *TeamsClient) Create(ctx context.Context, createTeam *CreateTeam) (*Team, error) {
	p := "/teams"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, createTeam)
	if err != nil {
		return nil, err
	}

	var resp Team
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get corresponds to the GET /teams/:id endpoint.
//
// Get a single team's profile
func (c *TeamsClient) Get(ctx context.Context, id ID) (*Team, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/teams/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Team
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List corresponds to the GET /teams endpoint.
//
// List teams for the current authenticated user
func (c *TeamsClient) List(ctx context.Context) *TeamIter {
	iter := TeamIter{
		first: true,
		i:     -1,
	}

	p := "/teams"

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		return &Error{}
	})
	return &iter
}

// ListMembers corresponds to the GET /teams/:id/members endpoint.
//
// Get team members by team id
func (c *TeamsClient) ListMembers(ctx context.Context, id ID) *MemberProfileIter {
	iter := MemberProfileIter{
		first: true,
		i:     -1,
	}

	idBytes, err := id.MarshalText()
	if err != nil {
		return &iter
	}

	p := fmt.Sprintf("/teams/%s/members", string(idBytes))

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, nil, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		return &Error{}
	})
	return &iter
}

// Update corresponds to the PATCH /teams/:id endpoint.
//
// Update team profile
func (c *TeamsClient) Update(ctx context.Context, id ID, updateTeam *UpdateTeam) (*Team, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/teams/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodPatch, p, nil, updateTeam)
	if err != nil {
		return nil, err
	}

	var resp Team
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// TokensClient provides access to the /tokens APIs
type TokensClient endpoint

// Create corresponds to the POST /tokens endpoint.
//
// Create a new api token
func (c *TokensClient) Create(ctx context.Context, apitokenRequest *APITokenRequest) (*APIToken, error) {
	p := "/tokens"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, apitokenRequest)
	if err != nil {
		return nil, err
	}

	var resp APIToken
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateAuth corresponds to the POST /tokens/auth endpoint.
//
// Exchange a login token for a general use auth token
func (c *TokensClient) CreateAuth(ctx context.Context, authorization string, authTokenRequest *AuthTokenRequest) (*AuthToken, error) {
	p := "/tokens/auth"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, authTokenRequest)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", authorization)

	var resp AuthToken
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateLogin corresponds to the POST /tokens/login endpoint.
//
// Create a new login token
func (c *TokensClient) CreateLogin(ctx context.Context, loginTokenRequest *LoginTokenRequest) (*LoginTokenResponse, error) {
	p := "/tokens/login"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, loginTokenRequest)
	if err != nil {
		return nil, err
	}

	var resp LoginTokenResponse
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete corresponds to the DELETE /tokens/:token endpoint.
//
// Revoke an auth token for log out
func (c *TokensClient) Delete(ctx context.Context, token string) error {
	p := fmt.Sprintf("/tokens/%s", token)

	req, err := c.backend.NewRequest(http.MethodDelete, p, nil, nil)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// List corresponds to the GET /tokens endpoint.
//
// List API tokens
func (c *TokensClient) List(ctx context.Context, tokensType string, opts *TokensListOpts) *APITokenIter {
	iter := APITokenIter{
		first: true,
		i:     -1,
	}

	p := "/tokens"

	q := make(url.Values)
	q.Set("type", tokensType)

	if opts != nil {

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
	}

	var req *http.Request
	req, iter.err = c.backend.NewRequest(http.MethodGet, p, q, nil)
	if iter.err != nil {
		return &iter
	}

	_, iter.err = c.backend.Do(ctx, req, &iter.page, func(code int) error {
		return &Error{}
	})
	return &iter
}

// UsersClient provides access to the /users APIs
type UsersClient endpoint

// Create corresponds to the POST /users endpoint.
//
// Create a new user
func (c *UsersClient) Create(ctx context.Context, createUser *CreateUser) (*User, error) {
	p := "/users"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, createUser)
	if err != nil {
		return nil, err
	}

	var resp User
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateForgotPassword corresponds to the POST /users/forgot-password endpoint.
//
// Reset password with an emailed forgot password token
func (c *UsersClient) CreateForgotPassword(ctx context.Context, forgotPassword *ForgotPassword) error {
	p := "/users/forgot-password"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, forgotPassword)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateForgotPasswordToken corresponds to the POST /users/forgot-password/token endpoint.
//
// User has forgotten password, send a reset token
func (c *UsersClient) CreateForgotPasswordToken(ctx context.Context, forgotPasswordCreate *ForgotPasswordCreate) error {
	p := "/users/forgot-password/token"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, forgotPasswordCreate)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateVerify corresponds to the POST /users/verify endpoint.
//
// Verify email address
func (c *UsersClient) CreateVerify(ctx context.Context, verifyEmail *VerifyEmail) error {
	p := "/users/verify"

	req, err := c.backend.NewRequest(http.MethodPost, p, nil, verifyEmail)
	if err != nil {
		return err
	}

	_, err = c.backend.Do(ctx, req, nil, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return err
	}

	return nil
}

// Update corresponds to the PATCH /users/:id endpoint.
//
// Update user profile
func (c *UsersClient) Update(ctx context.Context, id ID, updateUser *UpdateUser) (*User, error) {
	idBytes, err := id.MarshalText()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("/users/%s", string(idBytes))

	req, err := c.backend.NewRequest(http.MethodPatch, p, nil, updateUser)
	if err != nil {
		return nil, err
	}

	var resp User
	_, err = c.backend.Do(ctx, req, &resp, func(code int) error {
		return &Error{}
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Backend defines the low-level interface for communicating with the remote api.
type Backend interface {
	NewRequest(method, path string, query url.Values, body interface{}) (*http.Request, error)
	Do(ctx context.Context, request *http.Request, v interface{}, errFn func(int) error) (*http.Response, error)
}

// DefaultBackend returns an instance of the default Backend configuration.
func DefaultBackend() Backend {
	return &defaultBackend{client: &http.Client{}, base: baseIdentityURL}
}

type defaultBackend struct {
	client *http.Client
	base   string
}

func (b *defaultBackend) NewRequest(method, path string, query url.Values, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	url := b.base
	if path[0] != '/' {
		url += "/"
	}
	url += path
	if q := query.Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (b *defaultBackend) Do(ctx context.Context, request *http.Request, v interface{}, errFn func(int) error) (*http.Response, error) {
	request = request.WithContext(ctx)

	resp, err := b.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		if errFn == nil {
			return nil, nil
		}

		apiErr := errFn(resp.StatusCode)
		if apiErr == nil {
			return nil, nil
		}

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, apiErr
	}

	if v != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(v); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

type endpoint struct {
	backend Backend
}

// IdentityClient is an API client for all endpoints.
type IdentityClient struct {
	common endpoint // Reuse a single struct instead of allocating one for each endpoint on the heap.

	Analytics   *AnalyticsClient
	Invites     *InvitesClient
	Memberships *MembershipsClient
	Self        *SelfClient
	Teams       *TeamsClient
	Tokens      *TokensClient
	Users       *UsersClient
}

// NewIdentity returns a new IdentityClient with the default configuration.
func NewIdentity() *IdentityClient {
	c := &IdentityClient{}
	c.common.backend = DefaultBackend()

	c.Analytics = (*AnalyticsClient)(&c.common)
	c.Invites = (*InvitesClient)(&c.common)
	c.Memberships = (*MembershipsClient)(&c.common)
	c.Self = (*SelfClient)(&c.common)
	c.Teams = (*TeamsClient)(&c.common)
	c.Tokens = (*TokensClient)(&c.common)
	c.Users = (*UsersClient)(&c.common)

	return c
}
