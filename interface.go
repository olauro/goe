package goe

import (
	"context"
	"database/sql"
)

type field interface {
	fieldSelect
	fieldDb
	isPrimaryKey() bool
	getTableId() int
	getSelect() string
	getAttributeName() []byte
	table() []byte
	buildAttributeInsert(*builder)
	writeAttributeInsert(*builder)
	buildAttributeUpdate(*builder)
}

type fieldDb interface {
	getDb() *DB
}

type fieldSelect interface {
	fieldDb
	buildAttributeSelect(*builder)
}

type Driver interface {
	Name() string
	MigrateContext(context.Context, *Migrator, Connection) (string, error)
	DropTable(string, Connection) (string, error)
	DropColumn(table, column string, conn Connection) (string, error)
	RenameColumn(table, oldColumn, newColumn string, conn Connection) (string, error)
	Init(*DB)
	KeywordHandler(string) string
	Sql
}

type Sql interface {
	Select() []byte
	From() []byte
	Where() []byte
	Insert() []byte
	Values() []byte
	Returning([]byte) []byte
	Update() []byte
	Set() []byte
	Delete() []byte
}

type Connection interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
