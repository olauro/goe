package goe

import (
	"context"
	"errors"
	"reflect"

	"github.com/olauro/goe/enum"
	"github.com/olauro/goe/model"
	"github.com/olauro/goe/query/where"
)

type save[T any] struct {
	table       *T
	tx          Transaction
	errNotFound error
	update      *stateUpdate[T]
}

// Save is a wrapper over [Update] for more simple updates,
// uses the value for create a where matching the primary keys
// and includes for update all non-zero values excluding the primary keys.
// If the record don't exists returns a [ErrNotFound].
//
// Save uses [context.Background] internally;
// to specify the context, use [SaveContext].
//
// # Examples
//
//	// updates animal name on record id 1
//	err = goe.Save(db.Animal).ByValue(Animal{Id: 1, Name: "Cat"})
func Save[T any](table *T) *save[T] {
	return SaveContext(context.Background(), table)
}

// SaveContext is a wrapper over [Update] for more simple updates,
// uses the value for create a where matching the primary keys
// and includes for update all non-zero values excluding the primary keys.
//
// See [Save] for examples.
func SaveContext[T any](ctx context.Context, table *T) *save[T] {
	return &save[T]{update: UpdateContext(ctx, table), table: table, errNotFound: ErrNotFound}
}

func (s *save[T]) OnTransaction(tx Transaction) *save[T] {
	s.update.conn = tx
	s.tx = tx
	return s
}

// Replace the ErrNotFound with err
func (s *save[T]) OnErrNotFound(err error) *save[T] {
	s.errNotFound = err
	return s
}

func (s *save[T]) ByValue(v T) error {
	if s.update.err != nil {
		return s.update.err
	}

	if _, err := Find(s.table).OnErrNotFound(s.errNotFound).OnTransaction(s.tx).ById(v); err != nil {
		return err
	}

	argsSave := getArgsSave(addrMap.mapField, s.table, v)
	if argsSave.err != nil {
		return argsSave.err
	}

	wheres := make([]model.Operation, 0, len(argsSave.argsWhere))
	wheres = append(wheres, where.Equals(&argsSave.argsWhere[0], argsSave.valuesWhere[0]))
	for i := 1; i < len(argsSave.argsWhere); i++ {
		wheres = append(wheres, where.And())
		wheres = append(wheres, where.Equals(&argsSave.argsWhere[i], argsSave.valuesWhere[i]))
	}

	s.update.builder.sets = argsSave.sets
	return s.update.Wheres(wheres...)
}

func (s *save[T]) AndFindByValue(v T) (*T, error) {
	err := s.ByValue(v)
	if err != nil {
		return nil, err
	}
	return Find(s.table).OnErrNotFound(s.errNotFound).OnTransaction(s.tx).ById(v)
}

func (s *save[T]) OrCreateByValue(v T) (*T, error) {
	err := s.ByValue(v)
	if err != nil {
		if errors.Is(err, s.errNotFound) {
			return Create(s.table).OnTransaction(s.tx).ByValue(v)
		}
		return nil, err
	}
	return Find(s.table).OnErrNotFound(s.errNotFound).OnTransaction(s.tx).ById(v)
}

type stateUpdate[T any] struct {
	conn    Connection
	builder builder
	ctx     context.Context
	err     error
}

// Update updates records in the given table
//
// Update uses [context.Background] internally;
// to specify the context, use [UpdateContext].
//
// # Examples
//
//	// update only the attribute IdJobTitle from PersonJobTitle with the value 3
//	err = goe.Update(db.PersonJobTitle).
//	Sets(update.Set(&db.PersonJobTitle.IdJobTitle, 3)).
//	Wheres(
//		where.Equals(&db.PersonJobTitle.PersonId, 2),
//		where.And(),
//		where.Equals(&db.PersonJobTitle.IdJobTitle, 1),
//	)
//
//	// update all animals name to Cat
//	goe.Update(db.Animal).Sets(update.Set(&db.Animal.Name, "Cat")).Wheres()
func Update[T any](table *T) *stateUpdate[T] {
	return UpdateContext(context.Background(), table)
}

// Update updates records in the given table
//
// See [Update] for examples
func UpdateContext[T any](ctx context.Context, table *T) *stateUpdate[T] {
	return createUpdateState[T](ctx)
}

// Sets one or more arguments for update
func (s *stateUpdate[T]) Sets(sets ...model.Set) *stateUpdate[T] {
	if s.err != nil {
		return s
	}

	for i := range sets {
		if field := getArg(sets[i].Attribute, addrMap.mapField, nil); field != nil {
			s.builder.sets = append(s.builder.sets, set{attribute: field, value: sets[i].Value})
		}
	}

	return s
}

func (s *stateUpdate[T]) OnTransaction(tx Transaction) *stateUpdate[T] {
	s.conn = tx
	return s
}

// Wheres receives [model.Operation] as where operations from where sub package
func (s *stateUpdate[T]) Wheres(brs ...model.Operation) error {
	if s.err != nil {
		return s.err
	}
	s.err = helperWhere(&s.builder, addrMap.mapField, brs...)
	if s.err != nil {
		return s.err
	}

	s.err = s.builder.buildUpdate()
	if s.err != nil {
		return s.err
	}

	if s.conn == nil {
		s.conn = s.builder.sets[0].attribute.getDb().driver.NewConnection()
	}

	return handlerValues(s.conn, s.builder.query, s.ctx)
}

type argSave struct {
	sets        []set
	argsWhere   []any
	valuesWhere []any
	err         error
}

func getArgsSave[T any](addrMap map[uintptr]field, table *T, value T) argSave {
	if table == nil {
		return argSave{err: errors.New("goe: invalid argument. try sending a pointer to a database mapped struct as argument")}
	}

	tableOf := reflect.ValueOf(table).Elem()

	if tableOf.Kind() != reflect.Struct {
		return argSave{err: errors.New("goe: invalid argument. try sending a pointer to a database mapped struct as argument")}
	}

	valueOf := reflect.ValueOf(value)

	sets := make([]set, 0)
	pksWhere, valuesWhere := make([]any, 0, valueOf.NumField()), make([]any, 0, valueOf.NumField())

	var addr uintptr
	for i := 0; i < valueOf.NumField(); i++ {
		if !valueOf.Field(i).IsZero() {
			addr = uintptr(tableOf.Field(i).Addr().UnsafePointer())
			if addrMap[addr] != nil {
				if addrMap[addr].isPrimaryKey() {
					pksWhere = append(pksWhere, tableOf.Field(i).Addr().Interface())
					valuesWhere = append(valuesWhere, valueOf.Field(i).Interface())
					continue
				}
				sets = append(sets, set{attribute: addrMap[addr], value: valueOf.Field(i).Interface()})
			}
		}
	}
	if len(pksWhere) == 0 {
		return argSave{err: errors.New("goe: invalid value. pass a value with a primary key filled")}
	}
	return argSave{sets: sets, argsWhere: pksWhere, valuesWhere: valuesWhere}
}

func createUpdateState[T any](ctx context.Context) *stateUpdate[T] {
	return &stateUpdate[T]{builder: createBuilder(enum.UpdateQuery), ctx: ctx}
}
