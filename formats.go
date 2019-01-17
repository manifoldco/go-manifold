package manifold

import (
	"regexp"

	"github.com/asaskevich/govalidator"

	"github.com/manifoldco/go-manifold/errors"
)

var (
	labelRegex             = regexp.MustCompile("^[a-z0-9][a-z0-9-_]{1,128}$")
	nameRegex              = regexp.MustCompile(`^[a-zA-Z0-9][a-z0-9A-Z\. \-_]{2,128}$`)
	codeRegex              = regexp.MustCompile("^([0-9abcdefghjkmnpqrtuvwxyz]{16}|[0-9]{6})$")
	featureValueLabelRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-_\.]{1,128}$`)
	annotationKeyRegex     = regexp.MustCompile(`^(?:[a-z0-9][a-z0-9-\.\/]{0,62}[a-z0-9]|[a-z0-9])$`)
	annotationValueRegex   = regexp.MustCompile(`^(?:[a-zA-Z0-9][a-zA-Z0-9-\.\/]{0,252}[a-zA-Z0-9]|[a-zA-Z0-9])$`)
	credentialKeyRegex     = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,999}$`)
	maxCredentialBodySize  = 32 * 1024
)

var (
	errInvalidLabel = NewError(errors.BadRequestError, "Invalid label")
	errInvalidName  = NewError(errors.BadRequestError, "Invalid name")
	errInvalidEmail = NewError(errors.BadRequestError, "Invalid Email Provided")
	errInvalidCode  = NewError(errors.BadRequestError,
		"Invalid email verification code provided")
	errInvalidFeatureValueLabel = NewError(errors.BadRequestError,
		"Invalid feature value label provided")
	errInvalidAnnotationKey = NewError(errors.BadRequestError,
		"Invalid annotation key provided")
	errInvalidAnnotationValue = NewError(errors.BadRequestError,
		"Invalid annotation value provided")
	errInvalidCredentialKey = NewError(errors.BadRequestError,
		"Invalid credential key provided")
	errInvalidCredentialValue = NewError(errors.BadRequestError,
		"Invalid credential body provided")
)

// Label represents any object's label field
type Label string

// Validate ensures the label value is valid
func (lbl Label) Validate(_ interface{}) error {
	if labelRegex.Match([]byte(lbl)) {
		return nil
	}

	return errInvalidLabel
}

// Name represents any object's name field
type Name string

// Validate ensures the name value is valid
func (n Name) Validate(_ interface{}) error {
	if nameRegex.Match([]byte(n)) {
		return nil
	}

	return errInvalidName
}

// Email represents any email field
type Email string

// Validate ensures that the email is valid
func (e Email) Validate(_ interface{}) error {
	if govalidator.IsEmail(string(e)) {
		return nil
	}

	return errInvalidEmail
}

// Code represents a manifold verification code ( E-Mail Verification )
type Code string

// Validate ensures the name value is valid
func (c Code) Validate(_ interface{}) error {
	if codeRegex.Match([]byte(c)) {
		return nil
	}

	return errInvalidCode
}

// FeatureValueLabel represents any object's label field
type FeatureValueLabel string

// Validate ensures the label value is valid
func (lbl FeatureValueLabel) Validate(_ interface{}) error {
	if featureValueLabelRegex.Match([]byte(lbl)) {
		return nil
	}

	return errInvalidFeatureValueLabel
}

// AnnotationKey represents any annotation map's key
type AnnotationKey string

// Validate ensures the annotation key is valid
func (key AnnotationKey) Validate(_ interface{}) error {
	if annotationKeyRegex.Match([]byte(key)) {
		return nil
	}

	return errInvalidAnnotationKey
}

// AnnotationValue represents any annotation map's value in the array of values
type AnnotationValue string

// Validate ensures the annotation value is valid
func (val AnnotationValue) Validate(_ interface{}) error {
	if annotationValueRegex.Match([]byte(val)) {
		return nil
	}

	return errInvalidAnnotationValue
}

// CredentialKey represents a key for a secret credential associated with a resource
type CredentialKey string

// Validate ensures the credential key is valid
func (key CredentialKey) Validate(_ interface{}) error {
	if credentialKeyRegex.Match([]byte(key)) {
		return nil
	}

	return errInvalidCredentialKey
}

// CredentialBody represents the body for a secret credential associated with a resource
type CredentialBody string

// Validate ensures the credential body is valid
func (body CredentialBody) Validate(_ interface{}) error {
	if len(body) > maxCredentialBodySize {
		return errInvalidCredentialValue
	}

	return nil
}
