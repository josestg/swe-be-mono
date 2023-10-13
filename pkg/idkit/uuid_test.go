package idkit

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDProvider_Request(t *testing.T) {
	uv1, err := UUIDv1.Request(context.Background())
	expectNoError(t, err)
	expectTrue(t, uv1.Version() == 0x01)

	uv4, err := UUIDv4.Request(context.Background())
	expectNoError(t, err)
	expectTrue(t, uv4.Version() == 0x04)

	uv0, err := UUIDv253.Request(context.Background())
	expectNoError(t, err)
	expectTrue(t, uv0 == uuid.Nil)

	uv254, err := UUIDv254.Request(context.Background())
	expectNoError(t, err)
	expectTrue(t, uv254 == StaticUUID())

	_, err = UUIDv255.Request(context.Background())
	expectTrue(t, errors.Is(err, ErrTeapot))

	_, err = uuidProvider(0x00).Request(context.Background())
	expectTrue(t, err != nil)
}

func TestUUIDProvider_FromStr(t *testing.T) {
	const v1str = "ffd42014-69c5-11ee-8c99-0242ac120002"
	const v4str = "a4c670b4-0dd8-4958-908c-55865b7ce52f"

	uv1, err := UUIDv1.FromStr(context.Background(), v1str)
	expectNoError(t, err)
	expectTrue(t, uv1 == uuid.MustParse(v1str))

	uv4, err := UUIDv4.FromStr(context.Background(), v4str)
	expectNoError(t, err)
	expectTrue(t, uv4 == uuid.MustParse(v4str))

	uv0, err := UUIDv253.FromStr(context.Background(), v1str)
	expectNoError(t, err)
	expectTrue(t, uv0 == uuid.Nil)

	uv254, err := UUIDv254.FromStr(context.Background(), v1str)
	expectNoError(t, err)
	expectTrue(t, uv254 == StaticUUID())

	_, err = UUIDv255.FromStr(context.Background(), v1str)
	expectTrue(t, errors.Is(err, ErrTeapot))

	_, err = uuidProvider(0x00).FromStr(context.Background(), v1str)
	expectTrue(t, err != nil)

	_, err = UUIDv1.FromStr(context.Background(), "invalid-uuid")
	expectTrue(t, err != nil)
}

func expectTrue(t *testing.T, ok bool) {
	t.Helper()
	if !ok {
		t.Errorf("expect true; got false")
	}
}

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expect no error; got an error: %v", err)
	}
}
