package query

import (
	"fmt"

	"github.com/olauro/goe/enum"
)

type Joins interface {
	FirstArg() any
	Join() enum.JoinType
	SecondArg() any
}

type Aggregate interface {
	Aggregate() enum.AggregateType
}

type FunctionType interface {
	GetType() enum.FunctionType
}

type ValueOperation interface {
	GetValue() any
}

type Operation struct {
	Type                enum.WhereType
	Arg                 any
	Value               ValueOperation
	Operator            enum.OperatorType
	Attribute           string
	Table               string
	Function            enum.FunctionType
	AttributeValue      string
	AttributeValueTable string
}

type Function[T any] struct {
	Field *T
	Type  enum.FunctionType
	Value T
}

func (f *Function[string]) Scan(src any) error {
	v, ok := src.(string)
	if !ok {
		return fmt.Errorf("error scan function")
	}

	f.Value = v
	return nil
}

func (f Function[T]) GetValue() any {
	return f.Value
}

func (f Function[T]) GetType() enum.FunctionType {
	return f.Type
}

type Count struct {
	Field any
	Value int64
}

func (c *Count) Scan(src any) error {
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("error scan aggregate")
	}

	c.Value = v
	return nil
}

func (c Count) Aggregate() enum.AggregateType {
	return enum.CountAggregate
}
