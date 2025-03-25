package where

import (
	"reflect"

	"github.com/olauro/goe/enum"
	"github.com/olauro/goe/model"
)

type valueOperation struct {
	value any
}

func (vo valueOperation) GetValue() any {
	if result, ok := vo.value.(model.ValueOperation); ok {
		return result.GetValue()
	}
	return vo.value
}

// # Example
//
//	// delete Food with Id fc1865b4-6f2d-4cc6-b766-49c2634bf5c4
//	Wheres(where.Equals(&db.Food.Id, "fc1865b4-6f2d-4cc6-b766-49c2634bf5c4"))
//
//	// generate: WHERE "animals"."idhabitat" IS NULL
//	Wheres(where.Equals(&db.Animal.IdHabitat, nil))
func Equals[T any, A *T | **T](a A, v T) model.Operation {
	if reflect.ValueOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil() {
		return model.Operation{Arg: a, Operator: enum.Is, Type: enum.OperationIsWhere}
	}
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Equals, Type: enum.OperationWhere}
}

// # Example
//
//	// get all foods that name are not Cookie
//	Wheres(where.NotEquals(&db.Food.Name, "Cookie"))
//
//	// generate: WHERE "animals"."idhabitat" IS NOT NULL
//	Wheres(where.NotEquals(&db.Animal.IdHabitat, nil))
func NotEquals[T any, A *T | **T](a A, v T) model.Operation {
	if reflect.ValueOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil() {
		return model.Operation{Arg: a, Operator: enum.IsNot, Type: enum.OperationIsWhere}
	}
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.NotEquals, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that was created after 09 of october 2024 at 11:50AM
//	Wheres(where.Greater(&db.Animal.CreateAt, time.Date(2024, time.October, 9, 11, 50, 00, 00, time.Local)))
func Greater[T any, A *T | **T](a A, v T) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Greater, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that was created in or after 09 of october 2024 at 11:50AM
//	Wheres(where.GreaterEquals(&db.Animal.CreateAt, time.Date(2024, time.October, 9, 11, 50, 00, 00, time.Local)))
func GreaterEquals[T any, A *T | **T](a A, v T) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.GreaterEquals, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that was updated before 09 of october 2024 at 11:50AM
//	Wheres(where.Less(&db.Animal.UpdateAt, time.Date(2024, time.October, 9, 11, 50, 00, 00, time.Local)))
func Less[T any, A *T | **T](a A, v T) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Less, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that was updated in or before 09 of october 2024 at 11:50AM
//	Wheres(where.LessEquals(&db.Animal.UpdateAt, time.Date(2024, time.October, 9, 11, 50, 00, 00, time.Local)))
func LessEquals[T any, A *T | **T](a A, v T) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.LessEquals, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that has a "at" in his name
//	Wheres(where.Like(&db.Animal.Name, "%at%"))
func Like[T any](a *T, v string) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Like, Type: enum.OperationWhere}
}

// # Example
//
//	// get all animals that has a "at" in his name
//	Wheres(where.Like(&db.Animal.Name, "%at%"))
func NotLike[T any](a *T, v string) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.NotLike, Type: enum.OperationWhere}
}

// # Example
//
//	// where in using a slice
//	Wheres(where.In(&db.Animal.Name, []string{"Cat", "Dog"}))
//
//	// AsQuery for get the query result from a select query
//	querySelect, err := goe.Select(&struct{ Name *string }{Name: &db.Animal.Name}).From(db.Animal).AsQuery()
//
//	// Use querySelect on in
//	rows, err := goe.Select(db.Animal).From(db.Animal).Wheres(where.In(&db.Animal.Name, querySelect).AsSlice()
func In[T any, V []T | *model.Query](a *T, mq V) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: mq}, Operator: enum.In, Type: enum.OperationInWhere}
}

// # Example
//
//	// where not in using a slice
//	Wheres(where.NotIn(&db.Animal.Name, []string{"Cat", "Dog"}))
//
//	// AsQuery for get the query result from a select query
//	querySelect, err := goe.Select(&struct{ Name *string }{Name: &db.Animal.Name}).From(db.Animal).AsQuery()
//
//	// Use querySelect on not in
//	rows, err := goe.Select(db.Animal).From(db.Animal).Wheres(where.NotIn(&db.Animal.Name, querySelect).AsSlice()
func NotIn[T any, V []T | *model.Query](a *T, mq V) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: mq}, Operator: enum.NotIn, Type: enum.OperationInWhere}
}

// # Example
//
//	Wheres(
//		where.Equals(&db.Animal.Status, "Eating"),
//		where.And(),
//		where.Like(&db.Animal.Name, "%Cat%"))
func And() model.Operation {
	return model.Operation{Operator: enum.And, Type: enum.LogicalWhere}
}

// # Example
//
//	Wheres(
//		where.Equals(&db.Animal.Status, "Eating"),
//		where.Or(),
//		where.Like(&db.Animal.Name, "%Cat%"))
func Or() model.Operation {
	return model.Operation{Operator: enum.Or, Type: enum.LogicalWhere}
}

// # Example
//
//	// implicit join using EqualsArg
//	goe.Select(db.Animal).
//	From(db.Animal, db.AnimalFood, db.Food).
//	Wheres(
//		where.EqualsArg[int](&db.Animal.Id, &db.AnimalFood.IdAnimal),
//		where.And(),
//		where.EqualsArg[uuid.UUID](&db.Food.Id, &db.AnimalFood.IdFood)).
//	AsSlice()
func EqualsArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Equals, Type: enum.OperationAttributeWhere}
}

// # Example
//
//	Wheres(where.NotEqualsArg(&db.Job.Id, &db.Person.Id))
func NotEqualsArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.NotEquals, Type: enum.OperationAttributeWhere}
}

// # Example
//
//	Wheres(where.GreaterArg(&db.Stock.Minimum, &db.Drinks.Stock))
func GreaterArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Greater, Type: enum.OperationAttributeWhere}
}

// # Example
//
//	Wheres(where.GreaterEqualsArg(&db.Drinks.Reorder, &db.Drinks.Stock))
func GreaterEqualsArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.GreaterEquals, Type: enum.OperationAttributeWhere}
}

// # Example
//
//	Wheres(where.LessArg(&db.Exam.Score, &db.Data.Minimum))
func LessArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.Less, Type: enum.OperationAttributeWhere}
}

// # Example
//
//	Wheres(where.LessEqualsArg(&db.Exam.Score, &db.Data.Minimum))
func LessEqualsArg[T any, A *T | **T](a A, v A) model.Operation {
	return model.Operation{Arg: a, Value: valueOperation{value: v}, Operator: enum.LessEquals, Type: enum.OperationAttributeWhere}
}
