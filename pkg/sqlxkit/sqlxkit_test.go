package sqlxkit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

const (
	simpleMock = "sqlxkit_test_simple_mock"
)

func init() {
	sql.Register(simpleMock, &simpleMockDriver{})
}

func TestOpen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, err := Open(simpleMock, "foo")
		expectNoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		_, ok := db.Driver().(*simpleMockDriver)
		expectTrue(t, ok)
	})

	t.Run("ok with options", func(t *testing.T) {
		db, err := Open(simpleMock, "foo", func(c *Config) { c.StructTagName = "json" })
		expectNoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		_, ok := db.Driver().(*simpleMockDriver)
		expectTrue(t, ok)
	})

	t.Run("failed", func(t *testing.T) {
		db, err := Open(simpleMock+"unregistered", "foo", func(c *Config) { c.StructTagName = "json" })
		expectTrue(t, db == nil)
		expectTrue(t, err != nil)
	})

}

var errExample = errors.New("an error")

func TestExecTransaction(t *testing.T) {
	t.Run("begin transaction failed", func(t *testing.T) {
		db, mock, teardown := Setup(t)
		t.Cleanup(teardown)

		mock.ExpectBegin().WillReturnError(errExample)

		err := ExecTransaction(context.Background(), db, NoopTransaction)
		expectTrue(t, errors.Is(err, errExample))
	})

	t.Run("commit failed", func(t *testing.T) {
		db, mock, teardown := Setup(t)
		t.Cleanup(teardown)

		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errExample)

		err := ExecTransaction(context.Background(), db, NoopTransaction, NoopTransaction)
		expectTrue(t, errors.Is(err, errExample))
	})

	t.Run("rollback failed", func(t *testing.T) {
		db, mock, teardown := Setup(t)
		t.Cleanup(teardown)

		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(errExample)

		fail := func(ctx context.Context, tx Tx) (context.Context, error) {
			return ctx, errExample
		}

		err := ExecTransaction(context.Background(), db, NoopTransaction, fail, func(ctx context.Context, tx Tx) (context.Context, error) {
			t.Fatalf("should not be called")
			return ctx, nil
		})

		expectTrue(t, errors.Is(err, errExample))
	})

	t.Run("committed", func(t *testing.T) {
		db, mock, teardown := Setup(t)
		t.Cleanup(teardown)

		mock.ExpectBegin()
		mock.ExpectCommit()

		err := ExecTransaction(context.Background(), db, NoopTransaction, NoopTransaction)
		expectNoError(t, err)
	})

	t.Run("rollback", func(t *testing.T) {
		db, mock, teardown := Setup(t)
		t.Cleanup(teardown)

		mock.ExpectBegin()
		mock.ExpectRollback()

		fail := func(ctx context.Context, tx Tx) (context.Context, error) {
			return ctx, errExample
		}

		err := ExecTransaction(context.Background(), db, NoopTransaction, fail, func(ctx context.Context, tx Tx) (context.Context, error) {
			t.Fatalf("should not be called")
			return ctx, nil
		})

		expectTrue(t, errors.Is(err, errExample))
	})
}

func TestNamedExec(t *testing.T) {
	namedArg := map[string]any{"bar": "baz"}
	//goland:noinspection ALL
	nameQuery := "INSERT foo (bar) VALUES(:bar);"

	query, anyArg, err := sqlx.Named(nameQuery, namedArg)
	expectNoError(t, err)

	valArg := make([]driver.Value, 0, len(anyArg))
	for _, arg := range anyArg {
		valArg = append(valArg, driver.Value(arg))
	}

	t.Run("success", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := NamedExec(nameQuery, namedArg, WithVerifyAffectedRows(1)).
			Exec(context.Background(), db)
		expectNoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnError(errExample)

		_, err := NamedExec(nameQuery, namedArg, WithVerifyAffectedRows(1)).
			Exec(context.Background(), db)

		expectTrue(t, errors.Is(err, errExample))
	})

	t.Run("rows affected failed", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewErrorResult(errExample))

		_, err := NamedExec(nameQuery, namedArg, WithVerifyAffectedRows(1)).
			Exec(context.Background(), db)

		expectTrue(t, errors.Is(err, errExample))
	})

	t.Run("rows affected not matched", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := NamedExec(nameQuery, namedArg, WithVerifyAffectedRows(1)).
			Exec(context.Background(), db)

		expectTrue(t, errors.Is(err, ErrUnexpectedAffectedRows))
	})

	t.Run("get last id failed", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(&lastInsertedIDFail{})

		var id int64
		_, err := NamedExec(nameQuery, namedArg, WithReadLastInsertedID(&id)).Exec(context.Background(), db)
		expectTrue(t, errors.Is(err, errExample))
		expectTrue(t, id == 0)
	})

	t.Run("read affected rows enabled", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewResult(123, 1))

		var id, affected int64
		_, err := NamedExec(nameQuery, namedArg, WithReadAffectedRows(&affected), WithReadLastInsertedID(&id)).
			Exec(context.Background(), db)

		expectNoError(t, err)
		expectTrue(t, id == 123)
		expectTrue(t, affected == 1)
	})

	t.Run("read affected rows but failed", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewErrorResult(errExample))

		var id, affected int64
		_, err := NamedExec(nameQuery, namedArg, WithReadAffectedRows(&affected), WithReadLastInsertedID(&id)).
			Exec(context.Background(), db)

		expectTrue(t, errors.Is(err, errExample))
		expectTrue(t, id == 0)
		expectTrue(t, affected == 0)
	})

	t.Run("all options enabled", func(t *testing.T) {
		db, mock, teardown := Setup(t, sqlmock.QueryMatcherEqual)
		t.Cleanup(teardown)

		mock.ExpectExec(query).
			WithArgs(valArg...).
			WillReturnResult(sqlmock.NewResult(123, 1))

		var id, affected int64
		_, err := NamedExec(nameQuery, namedArg,
			WithReadAffectedRows(&affected),
			WithReadLastInsertedID(&id),
			WithVerifyAffectedRows(1),
		).Exec(context.Background(), db)

		expectNoError(t, err)
		expectTrue(t, id == 123)
		expectTrue(t, affected == 1)
	})
}

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func expectTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("expected true, got false")
	}
}

type simpleMockDriver struct{ driver.Conn }

func (d *simpleMockDriver) Open(string) (driver.Conn, error) { return d, nil }

type lastInsertedIDFail struct{}

func (f *lastInsertedIDFail) RowsAffected() (int64, error) { return 0, nil }
func (f *lastInsertedIDFail) LastInsertId() (int64, error) { return 0, errExample }

func Setup(t *testing.T, matchers ...sqlmock.QueryMatcher) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	defaultMatcher := sqlmock.QueryMatcherEqual
	for _, m := range matchers {
		defaultMatcher = m
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(defaultMatcher))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sql-mock")

	teardown := func() {
		defer func() { _ = dbx.Close() }()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
	return dbx, mock, teardown
}
