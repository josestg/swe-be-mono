package sqlxkit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

//go:generate mockgen -destination=mock$GOPACKAGE/$GOFILE -package=mock$GOPACKAGE . Conn,DB,Tx

// Config holds database configuration.
type Config struct {
	MaxOpenConnections int
	MaxIdleConnections int
	StructTagName      string // default: sql
}

// Option is function to customize Config.
type Option func(*Config)

// apply applies Option to Config.
func (f Option) apply(c *Config) { f(c) }

// DefaultOption sets Config with default values.
func DefaultOption() Option {
	return func(cfg *Config) {
		cfg.MaxOpenConnections = 0 // unlimited.
		cfg.MaxIdleConnections = 2 // default from sqlx.
		cfg.StructTagName = "sql"
	}
}

// Reader is a subset of sqlx.DB that only has read-only methods.
type Reader interface {
	// QueryxContext queries the database and returns an *sqlx.Rows.
	// Any placeholder parameters are replaced with supplied args.
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)

	// QueryRowxContext queries the database and returns an *sqlx.Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
}

// Writer is a subset of sqlx.DB that only has write-only methods.
type Writer interface {
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// NamedExecContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

// Binder is a subset of sqlx.DB that only has bind methods.
type Binder interface {
	// Rebind transforms a query from sqlx.QUESTION to the DB driver's bindvar type.
	Rebind(query string) string

	// BindNamed binds a query using the DB driver's bindvar type.
	// Any named placeholder parameters are replaced with fields from arg.
	BindNamed(query string, arg any) (string, []any, error)
}

// Preparer is a subset of sqlx.DB that only has prepared statement methods.
type Preparer interface {
	// PreparexContext returns a sqlx.Stmt instead of a sql.Stmt
	// The provided context is used until the returned statement is closed.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// PrepareNamedContext returns a sqlx.NamedStmt instead of a sql.NamedStmt
	// The provided context is used until the returned statement is closed.
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

// DB is the interface for working with database.
type DB interface {
	Reader
	Writer
	Binder
	Preparer
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// Tx is the interface for working with database transaction.
type Tx interface {
	Reader
	Writer
	Binder
	Preparer
}

// Conn is the DB with Connection management.
type Conn interface {
	DB

	// Driver returns the underlying sql driver.
	Driver() driver.Driver

	// Close closes the underlying sql.DB connection.
	Close() error

	// Conn returns the underlying driver connection.
	Conn(ctx context.Context) (*sql.Conn, error)

	// PingContext verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error
}

// Open opens a database connection with given driver and dsn.
// Options can be used to override default config.
func Open(driver, dsn string, options ...Option) (Conn, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	ApplyConfig(db, options...)
	return db, nil
}

// ApplyConfig applies given options to db.
// This function is useful when you want to apply options to existing db.
// For example, using mock in test but want the same config as production.
func ApplyConfig(db *sqlx.DB, options ...Option) *sqlx.DB {
	var cfg Config
	DefaultOption().apply(&cfg)
	// override default config.
	for _, opt := range options {
		opt.apply(&cfg)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.Mapper = reflectx.NewMapperFunc(cfg.StructTagName, strings.ToLower)
	// ... other options in the future.
	return db
}

// Atomic is a function that executes transaction.
type Atomic func(ctx context.Context, tx Tx) (context.Context, error)

// Exec executes the transaction.
// The Exec is a syntactic sugar for applying Atomic function.
func (f Atomic) Exec(ctx context.Context, tx Tx) (context.Context, error) { return f(ctx, tx) }

// NoopTransaction is a transaction that does nothing.
func NoopTransaction(ctx context.Context, _ Tx) (context.Context, error) { return ctx, nil }

// ExecTransaction executes the given transactions in a single transaction.
// The order of the transactions is the same as the order of the arguments, and the context is passed to the next
// transaction to provide some metadata to the next transaction calls. It can be useful for tracing and passing
// data between transactions if needed.
//
// If one of the transaction cause error, the next transactions will not be executed and all transactions will be
// rolling back. Otherwise, all transactions will be committed.
func ExecTransaction(ctx context.Context, db DB, transactions ...Atomic) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	for i := 0; i < len(transactions); i++ {
		// don't use `:=`, because we need to replace ctx with the returned ctx to next calls.
		ctx, err = transactions[i].Exec(ctx, tx)
		// if one of transaction cause error, it should be rollback.
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(
					err, // the root cause error.
					fmt.Errorf("rolling back transactions[%d] failed with error: %w", i, rollbackErr),
				)
			}
			return fmt.Errorf("evaluating transactions[%d]: %w", i, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// ErrUnexpectedAffectedRows is an error that is returned when the affected rows
// is not equal to the expected.
var ErrUnexpectedAffectedRows = errors.New("unexpected affected rows")

type execOption struct {
	verifyAffected bool
	expectAffected int64
	lastInsertedID *int64
	readAffected   bool
	affectedRows   *int64
}

// ExecOption is an option for Exec.
// This option is used to modify the behavior of Exec.
type ExecOption func(*execOption)

// WithVerifyAffectedRows is an option that verifies the affected rows.
func WithVerifyAffectedRows(expected int64) ExecOption {
	return func(opt *execOption) {
		opt.verifyAffected = true
		opt.expectAffected = expected
	}
}

// WithReadAffectedRows is an option that reads the affected rows.
func WithReadAffectedRows(dst *int64) ExecOption {
	return func(opt *execOption) {
		opt.readAffected = true
		opt.affectedRows = dst
	}
}

// WithReadLastInsertedID is an option that reads the last inserted ID.
func WithReadLastInsertedID(dst *int64) ExecOption {
	return func(opt *execOption) {
		opt.lastInsertedID = dst
	}
}

// NamedExec is a helper function to execute named query in transaction or without transaction if db is not
// transactional.
// This helper simply the process for verifying affected rows and reading last inserted ID by using ExecOption.
func NamedExec(query string, arg any, opts ...ExecOption) Atomic {
	var conf execOption
	for _, opt := range opts {
		opt(&conf)
	}
	return func(ctx context.Context, tx Tx) (context.Context, error) {
		ctx, err := doNamedExec(ctx, &conf, tx, query, arg)
		if err != nil {
			return ctx, fmt.Errorf("sqlxkit: NamedExec: %w", err)
		}
		return ctx, nil
	}
}

func doNamedExec(ctx context.Context, conf *execOption, db Tx, query string, arg any) (context.Context, error) {
	res, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return ctx, fmt.Errorf("exec query, error: %w", err)
	}

	if err := doAffectedRowsAction(conf, res); err != nil {
		return ctx, fmt.Errorf("process read affected rows: %w", err)
	}

	if err := doLastInsertIDAction(conf, res); err != nil {
		return ctx, fmt.Errorf("process read last insert id: %w", err)
	}

	return ctx, nil
}

func doAffectedRowsAction(conf *execOption, res sql.Result) error {
	if conf.verifyAffected || conf.readAffected {
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("get affected rows: %w", err)
		}

		if conf.readAffected {
			*conf.affectedRows = n
		}

		if conf.verifyAffected {
			if n != conf.expectAffected {
				return fmt.Errorf("expected=%d, got=%d: %w", conf.expectAffected, n, ErrUnexpectedAffectedRows)
			}
		}
	}
	return nil
}

func doLastInsertIDAction(conf *execOption, res sql.Result) error {
	if conf.lastInsertedID != nil {
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("get last inserted ID: %w", err)
		}
		*conf.lastInsertedID = id
	}
	return nil
}
