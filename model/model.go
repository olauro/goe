package model

import "github.com/olauro/goe/enum"

type Attribute struct {
	Table         string
	Name          string
	AggregateType enum.AggregateType
	FunctionType  enum.FunctionType
}

type JoinArgument struct {
	Table string
	Name  string
}

type Join struct {
	Table          string
	FirstArgument  JoinArgument
	JoinOperation  enum.JoinType
	SecondArgument JoinArgument
}

type Where struct {
	Type           enum.WhereType
	Attribute      Attribute
	Operator       enum.OperatorType
	AttributeValue Attribute
	SizeIn         uint
	QueryIn        *Query
}

type OrderBy struct {
	Desc      bool
	Attribute Attribute
}

type Query struct {
	Type       enum.QueryType
	Attributes []Attribute
	Tables     []string

	Joins   []Join   //Select
	Limit   uint     //Select
	Offset  uint     //Select
	OrderBy *OrderBy //Select

	WhereOperations []Where //Select, Update and Delete
	WhereIndex      int     //Start of where position arguments $1, $2...
	Arguments       []any

	ReturningId    *Attribute //Insert
	BatchSizeQuery int        //Insert
	SizeArguments  int        //Insert

	RawSql string
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

type Set struct {
	Attribute any
	Value     any
}

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
