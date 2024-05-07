package prepare

import (
	"context"
	"database/sql"
	"sync"
)

func NewPrepareDB(db *sql.DB) *PreparDB {
	return &PreparDB{db: db}
}

type PreparDB struct {
	db      *sql.DB
	stmtMap sync.Map
}

func (d *PreparDB) Prepare(query string) StmtFace {
	if stmtRef, load := d.stmtMap.Load(query); load {
		return stmtRef.(*stmt)
	}

	ps, err := d.db.Prepare(query)
	if err != nil {
		return &errstmt{err: err}
	}

	ss := &stmt{m: ps}
	if old, loaded := d.stmtMap.LoadOrStore(query, ss); loaded {
		ps.Close()
		return old.(*stmt)
	}
	return ss
}

type StmtFace interface {
	QueryRowContext(ctx context.Context, args ...any) RowFace
	QueryContext(ctx context.Context, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, args ...any) (sql.Result, error)
}

type RowFace interface {
	Scan(dest ...any) error
}

type stmt struct {
	m *sql.Stmt
}

func (s *stmt) QueryRowContext(ctx context.Context, args ...any) RowFace {
	return s.m.QueryRowContext(ctx, args...)
}

func (s *stmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return s.m.QueryContext(ctx, args...)
}

func (s *stmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return s.m.ExecContext(ctx, args...)
}

type errstmt struct {
	err error
}

func (s *errstmt) QueryRowContext(ctx context.Context, args ...any) RowFace {
	return &errrow{err: s.err}
}

func (s *errstmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return nil, s.err
}

func (s *errstmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return nil, s.err
}

type errrow struct {
	err error
}

func (s *errrow) Scan(dest ...any) error {
	return s.err
}
