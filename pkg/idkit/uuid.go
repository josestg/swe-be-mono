package idkit

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var _staticUUID = uuid.MustParse("f8635b1a-3524-4441-8c22-10edf1c45407")

// StaticUUID returns a static UUID that is used for testing.
func StaticUUID() uuid.UUID { return _staticUUID }

var (
	// ErrTeapot is a special error that is used for testing. The UUIDv255 provider will return
	// this error when Request or FromStr is called.
	ErrTeapot = errors.New("i'm a teapot")
)

// UUIDProvider provides an API to generate and parse UUID.
type UUIDProvider interface {
	// Request requests a new UUID based on the provider.
	Request(ctx context.Context) (uuid.UUID, error)

	// FromStr converts a string to UUID based on the provider.
	// If the string is not a valid UUID or the UUID version is not the same as the provider,
	// it will return an error.
	FromStr(ctx context.Context, s string) (uuid.UUID, error)
}

// uuidProvider is type that implements UUIDProvider, the type is private to avoid
// direct usage of the type.
type uuidProvider uuid.Version

const (

	// UUIDv1 is a UUID provider that generates and parses UUIDv1.
	UUIDv1 = uuidProvider(0x01)

	// UUIDv4 is a UUID provider that generates and parses UUIDv4.
	UUIDv4 = uuidProvider(0x04)

	// UUIDv253 is a special UUID provider that always returns uuid.Nil for Request and FromStr.
	// This provider is useful for testing.
	UUIDv253 = uuidProvider(0xfd)

	// UUIDv254 is a special UUID provider that always returns a non-nil StaticUUID() for Request and FromStr.
	// This provider is useful for testing.
	UUIDv254 = uuidProvider(0xfe)

	// UUIDv255 is a special UUID provider that always returns error for Request and FromStr.
	// This provider is useful for testing.
	UUIDv255 = uuidProvider(0xff)
)

func (u uuidProvider) Request(_ context.Context) (uuid.UUID, error) {
	switch u {
	case UUIDv1:
		return uuid.NewUUID()
	case UUIDv4:
		return uuid.NewRandom()
	case UUIDv253:
		return uuid.Nil, nil
	case UUIDv254:
		return StaticUUID(), nil
	case UUIDv255:
		return uuid.Nil, ErrTeapot
	default:
		return uuid.Nil, fmt.Errorf("unknown uuid provider for version: %d", u)
	}
}

func (u uuidProvider) FromStr(_ context.Context, s string) (uuid.UUID, error) {
	switch u {
	case UUIDv253:
		return uuid.Nil, nil
	case UUIDv254:
		return StaticUUID(), nil
	case UUIDv255:
		return uuid.Nil, ErrTeapot

	}

	uid, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("from string: %w", err)
	}

	if uid.Version() != uuid.Version(u) {
		return uuid.Nil, fmt.Errorf("invalid uuid version: %s", uid.Version())
	}

	return uid, nil
}
