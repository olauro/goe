package goe

import (
	"context"
	"log"
	"reflect"
)

type stateInsert[T any] struct {
	config  *Config
	conn    Connection
	builder *builder
	ctx     context.Context
	err     error
}

// Insert uses [context.Background] internally;
// to specify the context, use [query.InsertContext].
//
// # Example
func Insert[T any](table *T, tx ...*Tx) *stateInsert[T] {
	return InsertContext[T](context.Background(), table, tx...)
}

// InsertContext creates a insert state for table
func InsertContext[T any](ctx context.Context, table *T, tx ...*Tx) *stateInsert[T] {
	fields, err := getArgsTable(addrMap, table)

	var state *stateInsert[T]
	if err != nil {
		state = new(stateInsert[T])
		state.err = err
		return state
	}
	db := fields[0].getDb()

	if tx != nil {
		state = createInsertState[T](tx[0].SqlTx, db.Config, ctx, db.Driver, err)
	} else {
		state = createInsertState[T](db.SqlDB, db.Config, ctx, db.Driver, err)
	}
	state.builder.fields = fields
	return state
}

func (s *stateInsert[T]) One(value *T) error {
	if s.err != nil {
		return s.err
	}

	if value == nil {
		return ErrInvalidInsertValue
	}

	v := reflect.ValueOf(value).Elem()

	pkFieldId := s.builder.buildSqlInsert(v)

	sql := s.builder.sql.String()
	if s.config.LogQuery {
		log.Println("\n" + sql)
	}
	if s.builder.returning != nil {
		return handlerValuesReturning(s.conn, sql, v, s.builder.argsAny, pkFieldId, s.ctx)
	}
	return handlerValues(s.conn, sql, s.builder.argsAny, s.ctx)
}

func (s *stateInsert[T]) All(value []T) error {
	if len(value) == 0 {
		return ErrEmptyBatchValue
	}

	valueOf := reflect.ValueOf(value)

	pkFieldId := s.builder.buildSqlInsertBatch(valueOf)

	Sql := s.builder.sql.String()
	if s.config.LogQuery {
		log.Println("\n" + Sql)
	}
	return handlerValuesReturningBatch(s.conn, Sql, valueOf, s.builder.argsAny, pkFieldId, s.ctx)
}

func createInsertState[T any](conn Connection, c *Config, ctx context.Context, d Driver, e error) *stateInsert[T] {
	return &stateInsert[T]{conn: conn, builder: createBuilder(d), config: c, ctx: ctx, err: e}
}

func getArgsTable[T any](AddrMap map[uintptr]field, table *T) ([]field, error) {
	if table == nil {
		return nil, ErrInvalidArg
	}
	fields := make([]field, 0)

	valueOf := reflect.ValueOf(table).Elem()
	if valueOf.Kind() != reflect.Struct {
		return nil, ErrInvalidArg
	}

	var fieldOf reflect.Value
	for i := 0; i < valueOf.NumField(); i++ {
		fieldOf = valueOf.Field(i)
		if fieldOf.Kind() == reflect.Slice && fieldOf.Type().Elem().Kind() == reflect.Struct {
			continue
		}
		addr := uintptr(fieldOf.Addr().UnsafePointer())
		if AddrMap[addr] != nil {
			fields = append(fields, AddrMap[addr])
		}
	}

	if len(fields) == 0 {
		return nil, ErrInvalidArg
	}
	return fields, nil
}
