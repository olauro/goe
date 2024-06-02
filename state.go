package goe

import (
	"fmt"
	"reflect"
)

type stateSelect struct {
	conn    conn
	addrMap map[string]any
	builder *builder
}

func createSelectState(conn conn, qt int8) *stateSelect {
	return &stateSelect{conn: conn, builder: createBuilder(qt)}
}

func (s *stateSelect) Where(brs ...operator) SelectWhere {
	where(s.builder, brs...)
	return s
}

func (s *stateSelect) Join(tables ...any) Select {
	s.builder.args = append(s.builder.args, getArgs(tables...)...)
	s.builder.buildSelectJoins(s.addrMap)
	return s
}

func (s *stateSelect) querySelect(args []string) *stateSelect {
	s.builder.args = args
	s.builder.buildSelect(s.addrMap)
	return s
}

func (s *stateSelect) Result(target any) {
	value := reflect.ValueOf(target)

	if value.Kind() != reflect.Ptr {
		fmt.Printf("%v: target result needs to be a pointer try &animals\n", pkg)
		return
	}

	//generate query
	s.builder.buildSql()

	fmt.Println(s.builder.sql)
	handlerResult(s.conn, s.builder.sql.String(), value.Elem(), s.builder.argsAny)
}

/*
State Insert
*/
type stateInsert struct {
	conn    conn
	builder *builder
}

func createInsertState(conn conn, qt int8) *stateInsert {
	return &stateInsert{conn: conn, builder: createBuilder(qt)}
}

func (s *stateInsert) queryInsert(args []string, addrMap map[string]any) Insert {
	s.builder.args = args
	s.builder.buildInsert(addrMap)
	return s
}

func (s *stateInsert) Value(target any) {
	value := reflect.ValueOf(target)

	if value.Kind() != reflect.Ptr {
		fmt.Printf("%v: target result needs to be a pointer try &animals\n", pkg)
		return
	}

	value = value.Elem()

	idName := s.builder.buildValues(value)

	//generate query
	s.builder.buildSql()

	fmt.Println(s.builder.sql)
	handlerValuesReturning(s.conn, s.builder.sql.String(), value, s.builder.argsAny, idName)
}

func (s *stateInsert) queryInsertBetwent(args []string, addrMap map[string]any) InsertBetwent {
	s.builder.args = args
	s.builder.buildInsertManyToMany(addrMap)
	return s
}

func (s *stateInsert) Values(v1 any, v2 any) {
	s.builder.argsAny = append(s.builder.argsAny, v1)
	s.builder.argsAny = append(s.builder.argsAny, v2)

	s.builder.buildValuesManyToMany()

	s.builder.buildSql()

	fmt.Println(s.builder.sql)
	handlerValues(s.conn, s.builder.sql.String(), s.builder.argsAny)
}

/*
State Update
*/
type stateUpdate struct {
	conn    conn
	builder *builder
}

func createUpdateState(conn conn, qt int8) *stateUpdate {
	return &stateUpdate{conn: conn, builder: createBuilder(qt)}
}

func (s *stateUpdate) Where(brs ...operator) UpdateWhere {
	where(s.builder, brs...)
	return s
}

func (s *stateUpdate) queryUpdate(args []string, addrMap map[string]any) Update {
	s.builder.args = args
	s.builder.buildUpdate(addrMap)
	return s
}

func (s *stateUpdate) Value(target any) {
	value := reflect.ValueOf(target)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		fmt.Printf("%v: value for update needs to be a struct\n", pkg)
		return
	}

	s.builder.buildSet(value)

	//generate query
	s.builder.buildSql()

	fmt.Println(s.builder.sql)
	handlerValues(s.conn, s.builder.sql.String(), s.builder.argsAny)

}

type stateUpdateBetwent struct {
	conn    conn
	builder *builder
}

func createUpdateBetwentState(conn conn, qt int8) *stateUpdateBetwent {
	return &stateUpdateBetwent{conn: conn, builder: createBuilder(qt)}
}

func (s *stateUpdateBetwent) Where(brs ...operator) UpdateWhere {
	where(s.builder, brs...)
	return s
}

func (s *stateUpdateBetwent) queryUpdateBetwent(args []string, addrMap map[string]any) Update {
	s.builder.args = args
	s.builder.buildUpdateBetwent(addrMap)
	return s
}

func (s *stateUpdateBetwent) Value(value any) {
	s.builder.argsAny = append(s.builder.argsAny, value)

	s.builder.buildSetBetwent()

	s.builder.buildeSqlUpdateBetwent()

	fmt.Println(s.builder.sql)
	handlerValues(s.conn, s.builder.sql.String(), s.builder.argsAny)
}

type stateDelete struct {
	conn    conn
	builder *builder
}

func createDeleteState(conn conn, qt int8) *stateDelete {
	return &stateDelete{conn: conn, builder: createBuilder(qt)}
}

func (s *stateDelete) queryDelete(args []string, addrMap map[string]any) Delete {
	s.builder.args = args
	s.builder.buildDelete(addrMap)
	return s
}

func (s *stateDelete) Where(brs ...operator) {
	where(s.builder, brs...)

	s.builder.buildSqlDelete()

	fmt.Println(s.builder.sql)
	handlerValues(s.conn, s.builder.sql.String(), s.builder.argsAny)
}

type stateDeleteIn struct {
	conn    conn
	builder *builder
}

func createDeleteInState(conn conn, qt int8) *stateDeleteIn {
	return &stateDeleteIn{conn: conn, builder: createBuilder(qt)}
}

func (s *stateDeleteIn) queryDeleteIn(args []string, addrMap map[string]any) DeleteIn {
	s.builder.args = args
	s.builder.buildDeleteIn(addrMap)
	return s
}

func (s *stateDeleteIn) Where(values ...any) {
	s.builder.argsAny = append(s.builder.argsAny, values...)

	s.builder.buildSqlDeleteIn()

	fmt.Println(s.builder.sql)
	handlerValues(s.conn, s.builder.sql.String(), s.builder.argsAny)
}

func where(builder *builder, brs ...operator) {
	builder.brs = brs
	for _, br := range builder.brs {
		if op, ok := br.(complexOperator); ok {
			builder.tables.add(createStatement(op.pk.table, writeTABLE))
			builder.pks.add(op.pk)
		}
	}
}