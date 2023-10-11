package env

import (
	"reflect"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	const (
		key     = "TESTING_ENV_STRING"
		initial = "STRING INITIAL VALUE"
		val     = "STRING VALUE"
	)

	if got := String(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, val)
	if got := String(key, initial); got != val {
		t.Errorf("expected using the env value")
	}
}

func TestInt(t *testing.T) {
	const (
		key       = "TESTING_ENV_INT"
		initial   = 1
		val       = 2
		valString = "2"
	)

	if got := Int(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, valString)
	if got := Int(key, initial); got != val {
		t.Errorf("expected using the env value")
	}

	// Test: expected panic if the env value not a number.
	assertPanic(t, func() {
		t.Setenv(key, "INVALID NUMBER FORMAT")
		_ = Int(key, initial) // this should be panic.
	})
}

func TestInt64(t *testing.T) {
	const (
		key       = "TESTING_ENV_INT64"
		initial   = int64(1)
		val       = int64(2)
		valString = "2"
	)

	if got := Int64(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, valString)
	if got := Int64(key, initial); got != val {
		t.Errorf("expected using the env value")
	}

	// Test: expected panic if the env value not a number.
	assertPanic(t, func() {
		t.Setenv(key, "INVALID NUMBER FORMAT")
		_ = Int64(key, initial) // this should be panic.
	})
}

func TestFloat64(t *testing.T) {
	const (
		key       = "TESTING_ENV_FLOAT64"
		initial   = float64(1)
		val       = float64(2)
		valString = "2"
	)

	if got := Float64(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, valString)
	if got := Float64(key, initial); got != val {
		t.Errorf("expected using the env value")
	}

	// Test: expected panic if the env value not a number.
	assertPanic(t, func() {
		t.Setenv(key, "INVALID NUMBER FORMAT")
		_ = Float64(key, initial) // this should be panic.
	})
}

func TestDuration(t *testing.T) {
	const (
		key       = "TESTING_ENV_DURATION"
		initial   = 1 * time.Second
		val       = 2 * time.Second
		valString = "2s"
	)

	if got := Duration(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, valString)
	if got := Duration(key, initial); got != val {
		t.Errorf("expected using the env value")
	}

	// Test: expected panic if the env value not a duration.
	assertPanic(t, func() {
		t.Setenv(key, "INVALID DURATION FORMAT")
		_ = Duration(key, initial) // this should be panic.
	})
}

func TestBool(t *testing.T) {
	const (
		key       = "TESTING_ENV_BOOL"
		initial   = true
		val       = false
		valString = "false"
	)

	if got := Bool(key, initial); got != initial {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, valString)
	if got := Bool(key, initial); got != val {
		t.Errorf("expected using the env value")
	}

	// Test: expected panic if the env value not a boolean.
	assertPanic(t, func() {
		t.Setenv(key, "INVALID BOOL FORMAT")
		_ = Bool(key, initial) // this should be panic.
	})
}

func TestList_String(t *testing.T) {
	var (
		key     = "TESTING_ENV_STRING"
		initial = []string{"A", "B", "C"}
		val     = "A, B, C"
	)

	got := StringList(key, initial)
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, val)
	got = StringList(key, []string{})
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the env value")
	}
}

func TestList_Int64(t *testing.T) {
	var (
		key     = "TESTING_ENV_STRING"
		initial = []int64{1, 2, 3}
		val     = "1, 2, 3"
	)

	got := List(key, Parsers.Int64(), initial)
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, val)
	got = List(key, Parsers.Int64(), []int64{})
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the env value")
	}
}

func TestList_ParserFailed(t *testing.T) {
	var (
		key     = "TESTING_ENV_STRING"
		initial = []int64{1, 2, 3}
		val     = "1, X, 3"
	)

	got := List(key, Parsers.Int64(), initial)
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the initial value")
	}

	t.Setenv(key, val)
	got = List(key, Parsers.Int64(), initial)
	if !reflect.DeepEqual(got, initial) {
		t.Errorf("expected using the env value")
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	f()
}
