package goe

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/olauro/goe/query"
)

var ErrInvalidWhere = errors.New("goe: invalid where operation. try sending a pointer as parameter")
var ErrNoMatchesTables = errors.New("don't have any relationship")
var ErrNotManyToMany = errors.New("don't have a many to many relationship")

type builder struct {
	sql           *strings.Builder
	driver        Driver
	structPkName  string //insert
	returning     []byte //insert
	inserts       []field
	froms         []byte
	fields        []field
	fieldsSelect  []fieldSelect
	argsAny       []any
	structColumns []string //update
	attrNames     []string //insert and update
	orderBy       string
	limit         uint     //select
	offset        uint     //select
	joins         []string //select
	joinsArgs     []field  //select
	tables        []uint
	brs           []query.Operator
}

func createBuilder(d Driver) *builder {
	return &builder{
		sql:    &strings.Builder{},
		driver: d,
	}
}

func (b *builder) buildSelect() {
	b.sql.Write(b.driver.Select())

	len := len(b.fieldsSelect)
	if len == 0 {
		return
	}

	for i := range b.fieldsSelect[:len-1] {
		b.fieldsSelect[i].buildAttributeSelect(b)
		b.sql.WriteByte(',')
	}

	b.fieldsSelect[len-1].buildAttributeSelect(b)
}

func (b *builder) buildSelectJoins(join string, fields []field) {
	j := len(b.joinsArgs)
	b.joinsArgs = append(b.joinsArgs, make([]field, 2)...)
	b.tables = append(b.tables, make([]uint, 1)...)
	b.joins = append(b.joins, join)
	b.joinsArgs[j] = fields[0]
	b.joinsArgs[j+1] = fields[1]
}

func (b *builder) buildPage() {
	if b.limit != 0 {
		b.sql.WriteString(fmt.Sprintf(" LIMIT %v", b.limit))
	}
	if b.offset != 0 {
		b.sql.WriteString(fmt.Sprintf(" OFFSET %v", b.offset))
	}
}

func (b *builder) buildSqlSelect() (err error) {
	err = b.buildTables()
	if err != nil {
		return err
	}
	err = b.buildWhere()
	b.sql.WriteString(b.orderBy)
	b.buildPage()
	b.sql.WriteByte(';')
	return err
}

func (b *builder) buildSqlUpdate() (err error) {
	err = b.buildWhere()
	b.sql.WriteByte(';')
	return err
}

func (b *builder) buildSqlDelete() (err error) {
	err = b.buildWhere()
	b.sql.WriteByte(';')
	return err
}

func (b *builder) buildWhere() error {
	if len(b.brs) == 0 {
		return nil
	}
	b.sql.WriteByte('\n')
	b.sql.WriteString("WHERE ")
	argsCount := len(b.argsAny) + 1
	for _, op := range b.brs {
		switch v := op.(type) {
		case query.Operation:
			v.ValueFlag = fmt.Sprintf("$%v", argsCount)
			b.sql.WriteString(v.Operation())
			b.argsAny = append(b.argsAny, v.Value)
			argsCount++
		default:
			b.sql.WriteString(v.Operation())
		}
	}
	return nil
}

func (b *builder) buildTables() (err error) {
	b.sql.Write(b.driver.From())
	b.sql.Write(b.froms)
	c := 1
	for i := range b.joins {
		err = buildJoins(b.joins[i], b.sql, b.joinsArgs[i+c-1], b.joinsArgs[i+c-1+1], b.tables, i+1)
		if err != nil {
			return err
		}
		c++
	}
	return nil
}

func buildJoins(join string, sql *strings.Builder, f1, f2 field, tables []uint, tableIndice int) error {
	sql.WriteByte('\n')
	if !slices.Contains(tables, f2.getTableId()) {
		sql.WriteString(fmt.Sprintf("%v %v on (%v = %v)", join, string(f2.table()), f1.getSelect(), f2.getSelect()))
		tables[tableIndice] = f2.getTableId()
		return nil
	}
	//TODO: update this to write
	sql.WriteString(fmt.Sprintf("%v %v on (%v = %v)", join, string(f1.table()), f1.getSelect(), f2.getSelect()))
	tables[tableIndice] = f1.getTableId()
	return nil
}

func (b *builder) buildInsert() {
	//TODO: Set a drive type to share stm
	b.sql.WriteString("INSERT ")
	b.sql.WriteString("INTO ")

	b.attrNames = make([]string, 0, len(b.fields))

	f := b.fields[0]
	b.sql.Write(f.table())
	b.sql.WriteString(" (")
	for i := range b.fields {
		b.fields[i].buildAttributeInsert(b)
	}

	b.inserts[0].writeAttributeInsert(b)
	for _, f := range b.inserts[1:] {
		b.sql.WriteByte(',')
		f.writeAttributeInsert(b)
	}

	b.sql.WriteString(") ")
	b.sql.WriteString("VALUES ")
}

func (b *builder) buildValues(value reflect.Value) string {
	b.sql.WriteByte(40)
	b.argsAny = make([]any, 0, len(b.attrNames))

	c := 2
	b.sql.WriteString("$1")
	buildValueField(value.FieldByName(b.attrNames[0]), b)
	a := b.attrNames[1:]
	for i := range a {
		b.sql.WriteByte(',')
		b.sql.WriteString(fmt.Sprintf("$%v", c))
		buildValueField(value.FieldByName(a[i]), b)
		c++
	}
	b.sql.WriteByte(')')
	if b.returning != nil {
		b.sql.Write(b.returning)
	}
	return b.structPkName

}

func (b *builder) buildBatchValues(value reflect.Value) string {
	b.argsAny = make([]any, 0, len(b.attrNames))

	c := 1
	buildBatchValues(value.Index(0), b, &c)
	c++
	for j := 1; j < value.Len(); j++ {
		b.sql.WriteByte(',')
		buildBatchValues(value.Index(j), b, &c)
		c++
	}
	if b.returning != nil {
		b.sql.Write(b.returning)
	}
	return b.structPkName

}

func buildBatchValues(value reflect.Value, b *builder, c *int) {
	b.sql.WriteByte(40)
	b.sql.WriteString(fmt.Sprintf("$%v", *c))
	buildValueField(value.FieldByName(b.attrNames[0]), b)
	a := b.attrNames[1:]
	for i := range a {
		b.sql.WriteByte(',')
		b.sql.WriteString(fmt.Sprintf("$%v", *c+1))
		buildValueField(value.FieldByName(a[i]), b)
		*c++
	}
	b.sql.WriteByte(')')
}

func buildValueField(valueField reflect.Value, b *builder) {
	b.argsAny = append(b.argsAny, valueField.Interface())
}

func (b *builder) buildUpdate() {
	//TODO: Set a drive type to share stm
	b.sql.WriteString("UPDATE ")

	b.structColumns = make([]string, 0, len(b.fields))
	b.attrNames = make([]string, 0, len(b.fields))

	b.sql.Write(b.fields[0].table())
	b.sql.WriteString(" SET ")
	b.fields[0].buildAttributeUpdate(b)

	a := b.fields[1:]
	for i := range a {
		a[i].buildAttributeUpdate(b)
	}
}

func (b *builder) buildSet(value reflect.Value) {
	b.argsAny = make([]any, 0, len(b.attrNames))
	var c uint16 = 1
	buildSetField(value.FieldByName(b.structColumns[0]), b.attrNames[0], b, c)

	a := b.attrNames[1:]
	s := b.structColumns[1:]
	for i := range a {
		b.sql.WriteByte(',')
		c++
		buildSetField(value.FieldByName(s[i]), a[i], b, c)
	}
}

func buildSetField(valueField reflect.Value, FieldName string, b *builder, c uint16) {
	b.sql.WriteString(fmt.Sprintf("%v = $%v", FieldName, c))
	b.argsAny = append(b.argsAny, valueField.Interface())
	c++
}

func (b *builder) buildDelete() {
	//TODO: Set a drive type to share stm
	b.sql.WriteString("DELETE FROM ")
	b.sql.Write(b.fields[0].table())
}
