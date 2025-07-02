package gobizap

import(
	"context"
)

type DbRows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
	//Values() ([]any, error)
}

type DbRow interface {
	Scan(dest ...any) error
}

type DbConn interface {
	//Prepare(ctx context.Context, queryID string, query string) error
	//Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRow(ctx context.Context, query string, args ...any) DbRow
	//Query(ctx context.Context, query string, args ...any) (DbRows, error)
}
