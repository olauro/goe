package goe

import (
	"errors"
	"reflect"
	"slices"

	"github.com/olauro/goe/enum"
	"github.com/olauro/goe/query"
)

var ErrInvalidWhere = errors.New("goe: invalid where operation. try sending a pointer as parameter")
var ErrNoMatchesTables = errors.New("don't have any relationship")
var ErrNotManyToMany = errors.New("don't have a many to many relationship")

type builder struct {
	query        Query
	pkFieldId    int //insert
	inserts      []field
	fields       []field
	fieldsSelect []fieldSelect
	fieldIds     []int    //insert and update
	joins        []string //select
	joinsArgs    []field  //select
	tables       []int
	brs          []query.Operation
}

func createBuilder(typeQuery enum.QueryType) *builder {
	return &builder{
		query: Query{Type: typeQuery},
	}
}

func (b *builder) buildSelect() {
	b.query.Attributes = make([]Attribute, 0, len(b.fieldsSelect))

	len := len(b.fieldsSelect)
	if len == 0 {
		return
	}

	for i := range b.fieldsSelect[:len-1] {
		b.fieldsSelect[i].buildAttributeSelect(b)
	}

	b.fieldsSelect[len-1].buildAttributeSelect(b)
}

func (b *builder) buildSelectJoins(join string, fields []field) {
	j := len(b.joinsArgs)
	b.joinsArgs = append(b.joinsArgs, make([]field, 2)...)
	b.tables = append(b.tables, make([]int, 1)...)
	b.joins = append(b.joins, join)
	b.joinsArgs[j] = fields[0]
	b.joinsArgs[j+1] = fields[1]
}

func (b *builder) buildSqlSelect() (err error) {
	b.buildSelect()
	err = b.buildTables()
	if err != nil {
		return err
	}
	return b.buildWhere()
}

func (b *builder) buildSqlInsert(v reflect.Value) (pkFieldId int) {
	b.buildInsert()
	pkFieldId = b.buildValues(v)
	return pkFieldId
}

func (b *builder) buildSqlInsertBatch(v reflect.Value) (pkFieldId int) {
	b.buildInsert()
	pkFieldId = b.buildBatchValues(v)
	return pkFieldId
}

func (b *builder) buildSqlUpdate(v reflect.Value) (err error) {
	b.buildUpdate()
	b.buildSet(v)
	err = b.buildWhere()
	return err
}

func (b *builder) buildSqlDelete() (err error) {
	b.query.Tables = make([]string, 1)
	b.query.Tables[0] = b.fields[0].table()
	err = b.buildWhere()
	return err
}

func (b *builder) buildWhere() error {
	if len(b.brs) == 0 {
		return nil
	}
	b.query.WhereOperations = make([]Where, 0, len(b.brs))

	argsCount := len(b.query.Arguments) + 1
	b.query.WhereIndex = len(b.query.Arguments) + 1
	for _, v := range b.brs {
		switch v.Type {
		case enum.OperationWhere:
			b.query.Arguments = append(b.query.Arguments, v.Value.GetValue())

			b.query.WhereOperations = append(b.query.WhereOperations, Where{
				Attribute: Attribute{
					Name:         v.Attribute,
					Table:        v.Table,
					FunctionType: v.Function,
				},
				Operator: v.Operator,
				Type:     v.Type,
			})
			argsCount++
		case enum.OperationAttributeWhere:
			b.query.WhereOperations = append(b.query.WhereOperations, Where{
				Attribute: Attribute{
					Name:  v.Attribute,
					Table: v.Table,
				},
				Operator:       v.Operator,
				AttributeValue: Attribute{Name: v.AttributeValue, Table: v.AttributeValueTable},
				Type:           v.Type,
			})

		case enum.OperationIsWhere:
			b.query.WhereOperations = append(b.query.WhereOperations, Where{
				Attribute: Attribute{
					Name:  v.Attribute,
					Table: v.Table,
				},
				Operator: v.Operator,
				Type:     v.Type,
			})

		case enum.LogicalWhere:
			b.query.WhereOperations = append(b.query.WhereOperations, Where{
				Operator: v.Operator,
				Type:     v.Type,
			})

		}
	}
	return nil
}

func (b *builder) buildTables() (err error) {
	if len(b.joins) != 0 {
		b.query.Joins = make([]Join, 0, len(b.joins))
	}
	c := 1
	for i := range b.joins {
		buildJoins(b, b.joins[i], b.joinsArgs[i+c-1], b.joinsArgs[i+c-1+1], b.tables, i+1)
		c++
	}
	return nil
}

func buildJoins(b *builder, join string, f1, f2 field, tables []int, tableIndice int) {
	if slices.Contains(tables, f1.getTableId()) {
		b.query.Joins = append(b.query.Joins, Join{
			Table:          f2.table(),
			FirstArgument:  JoinArgument{Table: f1.table(), Name: f1.getAttributeName()},
			JoinOperation:  join,
			SecondArgument: JoinArgument{Table: f2.table(), Name: f2.getAttributeName()}})

		tables[tableIndice] = f2.getTableId()
		return
	}
	b.query.Joins = append(b.query.Joins, Join{
		Table:          f1.table(),
		FirstArgument:  JoinArgument{Table: f1.table(), Name: f1.getAttributeName()},
		JoinOperation:  join,
		SecondArgument: JoinArgument{Table: f2.table(), Name: f2.getAttributeName()}})

	tables[tableIndice] = f1.getTableId()
}

func (b *builder) buildInsert() {

	b.fieldIds = make([]int, 0, len(b.fields))
	b.query.Attributes = make([]Attribute, 0, len(b.fields))

	f := b.fields[0]
	b.query.Tables = make([]string, 1)
	b.query.Tables[0] = f.table()
	for i := range b.fields {
		b.fields[i].buildAttributeInsert(b)
	}

	b.inserts[0].writeAttributeInsert(b)
	for _, f := range b.inserts[1:] {
		f.writeAttributeInsert(b)
	}

}

func (b *builder) buildValues(value reflect.Value) int {
	//update to index
	b.query.Arguments = make([]any, 0, len(b.fieldIds))

	c := 2
	b.query.Arguments = append(b.query.Arguments, value.Field(b.fieldIds[0]).Interface())

	a := b.fieldIds[1:]
	for i := range a {
		b.query.Arguments = append(b.query.Arguments, value.Field(a[i]).Interface())
		c++
	}
	b.query.SizeArguments = len(b.fieldIds)
	return b.pkFieldId

}

func (b *builder) buildBatchValues(value reflect.Value) int {
	b.query.Arguments = make([]any, 0, len(b.fieldIds))

	c := 1
	buildBatchValues(value.Index(0), b, &c)
	c++
	for j := 1; j < value.Len(); j++ {
		buildBatchValues(value.Index(j), b, &c)
		c++
	}
	b.query.BatchSizeQuery = value.Len()
	b.query.SizeArguments = len(b.fieldIds)
	return b.pkFieldId

}

func buildBatchValues(value reflect.Value, b *builder, c *int) {
	b.query.Arguments = append(b.query.Arguments, value.Field(b.fieldIds[0]).Interface())

	a := b.fieldIds[1:]
	for i := range a {
		b.query.Arguments = append(b.query.Arguments, value.Field(a[i]).Interface())
		*c++
	}
}

func (b *builder) buildUpdate() {

	b.fieldIds = make([]int, 0, len(b.fields))
	b.query.Attributes = make([]Attribute, 0, len(b.fields))
	b.query.Tables = make([]string, 1)
	b.query.Tables[0] = b.fields[0].table()

	b.fields[0].buildAttributeUpdate(b)

	a := b.fields[1:]
	for i := range a {
		a[i].buildAttributeUpdate(b)
	}
}

func (b *builder) buildSet(value reflect.Value) {
	b.query.Arguments = make([]any, 0, len(b.fieldIds))
	var c uint16 = 1
	buildSetField(value.Field(b.fieldIds[0]), b.fields[0].getAttributeName(), b, c)

	for i := 1; i < len(b.fieldIds); i++ {
		c++
		buildSetField(value.Field(b.fieldIds[i]), b.fields[i].getAttributeName(), b, c)
	}
}

func buildSetField(valueField reflect.Value, attributeName string, b *builder, c uint16) {
	b.query.Attributes = append(b.query.Attributes, Attribute{Name: attributeName})
	b.query.Arguments = append(b.query.Arguments, valueField.Interface())
	c++
}
