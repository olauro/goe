package tests_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/olauro/goe"
)

func TestPostgresSelect(t *testing.T) {
	db, err := SetupPostgres()
	if err != nil {
		t.Fatalf("Expected database, got error: %v", err)
	}

	err = db.Delete(db.AnimalFood).Where()
	if err != nil {
		t.Fatalf("Expected delete AnimalFood, got error: %v", err)
	}
	err = db.Delete(db.Flag).Where()
	if err != nil {
		t.Fatalf("Expected delete flags, got error: %v", err)
	}
	err = db.Delete(db.Animal).Where()
	if err != nil {
		t.Fatalf("Expected delete animals, got error: %v", err)
	}
	err = db.Delete(db.Food).Where()
	if err != nil {
		t.Fatalf("Expected delete foods, got error: %v", err)
	}
	err = db.Delete(db.Habitat).Where()
	if err != nil {
		t.Fatalf("Expected delete habitats, got error: %v", err)
	}
	err = db.Delete(db.Info).Where()
	if err != nil {
		t.Fatalf("Expected delete infos, got error: %v", err)
	}
	err = db.Delete(db.Status).Where()
	if err != nil {
		t.Fatalf("Expected delete status, got error: %v", err)
	}
	err = db.Delete(db.UserRole).Where()
	if err != nil {
		t.Fatalf("Expected delete user roles, got error: %v", err)
	}
	err = db.Delete(db.User).Where()
	if err != nil {
		t.Fatalf("Expected delete users, got error: %v", err)
	}
	err = db.Delete(db.Role).Where()
	if err != nil {
		t.Fatalf("Expected delete roles, got error: %v", err)
	}
	err = db.Delete(db.PersonJob).Where()
	if err != nil {
		t.Fatalf("Expected delete personJobs, got error: %v", err)
	}
	err = db.Delete(db.Job).Where()
	if err != nil {
		t.Fatalf("Expected delete jobs, got error: %v", err)
	}
	err = db.Delete(db.Person).Where()
	if err != nil {
		t.Fatalf("Expected delete persons, got error: %v", err)
	}

	weathers := []Weather{
		{Name: "Hot"},
		{Name: "Cold"},
		{Name: "Wind"},
		{Name: "Nice"},
		{Name: "Ocean"},
	}
	err = db.Insert(db.Weather).Value(&weathers)
	if err != nil {
		t.Fatalf("Expected insert weathers, got error: %v", err)
	}

	habitats := []Habitat{
		{Id: uuid.New(), Name: "City", IdWeather: weathers[0].Id, NameWeather: "Test"},
		{Id: uuid.New(), Name: "Jungle", IdWeather: weathers[3].Id},
		{Id: uuid.New(), Name: "Savannah", IdWeather: weathers[0].Id},
		{Id: uuid.New(), Name: "Ocean", IdWeather: weathers[2].Id},
	}
	err = db.Insert(db.Habitat).Value(&habitats)
	if err != nil {
		t.Fatalf("Expected insert habitats, got error: %v", err)
	}

	status := []Status{
		{Name: "Cat Alive"},
		{Name: "Dog Alive"},
		{Name: "Big Dog Alive"},
	}

	err = db.Insert(db.Status).Value(&status)
	if err != nil {
		t.Fatalf("Expected insert habitats, got error: %v", err)
	}

	infos := []Info{
		{Id: uuid.New().NodeID(), Name: "Little Cat", IdStatus: status[0].Id, NameStatus: "Test"},
		{Id: uuid.New().NodeID(), Name: "Big Dog", IdStatus: status[2].Id},
	}
	err = db.Insert(db.Info).Value(&infos)
	if err != nil {
		t.Fatalf("Expected insert infos, got error: %v", err)
	}

	animals := []Animal{
		{Name: "Cat", IdHabitat: &habitats[0].Id, IdInfo: &infos[0].Id},
		{Name: "Dog", IdHabitat: &habitats[0].Id, IdInfo: &infos[1].Id},
		{Name: "Forest Cat", IdHabitat: &habitats[1].Id},
		{Name: "Bear", IdHabitat: &habitats[1].Id},
		{Name: "Lion", IdHabitat: &habitats[2].Id},
		{Name: "Puma", IdHabitat: &habitats[1].Id},
		{Name: "Snake", IdHabitat: &habitats[1].Id},
		{Name: "Whale"},
	}
	err = db.Insert(db.Animal).Value(&animals)
	if err != nil {
		t.Fatalf("Expected insert animals, got error: %v", err)
	}

	foods := []Food{{Id: uuid.New(), Name: "Meat"}, {Id: uuid.New(), Name: "Grass"}}
	err = db.Insert(db.Food).Value(&foods)
	if err != nil {
		t.Fatalf("Expected insert foods, got error: %v", err)
	}

	animalFoods := []AnimalFood{
		{IdFood: foods[0].Id, IdAnimal: animals[0].Id},
		{IdFood: foods[0].Id, IdAnimal: animals[1].Id}}
	err = db.Insert(db.AnimalFood).Value(&animalFoods)
	if err != nil {
		t.Fatalf("Expected insert animalFoods, got error: %v", err)
	}

	users := []User{
		{Name: "Lauro Santana", Email: "lauro@email.com"},
		{Name: "John Constantine", Email: "hunter@email.com"},
		{Name: "Harry Potter", Email: "harry@email.com"},
	}
	err = db.Insert(db.User).Value(&users)
	if err != nil {
		t.Fatalf("Expected insert users, got error: %v", err)
	}

	roles := []Role{
		{Name: "Administrator"},
		{Name: "User"},
		{Name: "Mid-Level"},
	}
	err = db.Insert(db.Role).Value(&roles)
	if err != nil {
		t.Fatalf("Expected insert roles, got error: %v", err)
	}

	tt := time.Now().AddDate(0, 0, 10)
	userRoles := []UserRole{
		{IdUser: users[0].Id, IdRole: roles[0].Id, EndDate: &tt},
		{IdUser: users[1].Id, IdRole: roles[2].Id},
	}
	err = db.Insert(db.UserRole).Value(&userRoles)
	if err != nil {
		t.Fatalf("Expected insert user roles, got error: %v", err)
	}

	persons := []Person{
		{Name: "Jhon"},
		{Name: "Laura"},
		{Name: "Luana"},
	}
	err = db.Insert(db.Person).Value(&persons)
	if err != nil {
		t.Fatalf("Expected insert persons, got error: %v", err)
	}

	jobs := []Job{
		{Name: "Developer"},
		{Name: "Designer"},
	}
	err = db.Insert(db.Job).Value(&jobs)
	if err != nil {
		t.Fatalf("Expected insert jobs, got error: %v", err)
	}

	personJobs := []PersonJob{
		{IdPerson: persons[0].Id, IdJob: jobs[0].Id, CreatedAt: time.Now()},
		{IdPerson: persons[1].Id, IdJob: jobs[0].Id, CreatedAt: time.Now()},
		{IdPerson: persons[2].Id, IdJob: jobs[1].Id, CreatedAt: time.Now()},
	}
	err = db.Insert(db.PersonJob).Value(&personJobs)
	if err != nil {
		t.Fatalf("Expected insert personJobs, got error: %v", err)
	}

	testCases := []struct {
		desc     string
		testCase func(t *testing.T)
	}{
		{
			desc: "Select",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(&db.Animal.Id, &db.Animal.Name).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != len(animals) {
					t.Errorf("Expected %v animals, got %v", len(animals), len(a))
				}
			},
		},
		{
			desc: "Select_One_Field",
			testCase: func(t *testing.T) {
				var a []int
				err = db.Select(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != len(animals) {
					t.Errorf("Expected %v animals, got %v", len(animals), len(a))
				}
				if a[0] == a[1] {
					t.Errorf("Expected a select, got same values: %v and %v", a[0], a[1])
				}
			},
		},
		{
			desc: "Select_Where_Equals",
			testCase: func(t *testing.T) {
				var a Animal
				err = db.Select(db.Animal).Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a.Name != animals[0].Name {
					t.Errorf("Expected a %v, got %v", animals[0].Name, a.Name)
				}
			},
		},
		{
			desc: "Select_Slice_Not_Found_One_Field",
			testCase: func(t *testing.T) {
				var a []int
				err = db.Select(&db.Animal.Id).Where(db.Equals(&db.Animal.Id, 0)).Scan(&a)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Not_Found_One_Field",
			testCase: func(t *testing.T) {
				var a int
				err = db.Select(&db.Animal.Id).Where(db.Equals(&db.Animal.Id, 0)).Scan(&a)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Slice_Not_Found",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Where(db.Equals(&db.Animal.Id, 0)).Scan(&a)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Not_Found",
			testCase: func(t *testing.T) {
				var a Animal
				err = db.Select(db.Animal).Where(db.Equals(&db.Animal.Id, 0)).Scan(&a)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Pointer_As_Scan",
			testCase: func(t *testing.T) {
				var a *Animal
				err = db.Select(db.Animal).Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a.Id != animals[0].Id {
					t.Errorf("Expected a %v, got : %v", animals[0].Id, a.Id)
				}
			},
		},
		{
			desc: "Select_Pointer_As_Scan_One_Field",
			testCase: func(t *testing.T) {
				var a *int
				err = db.Select(&db.Animal.Id).Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if *a != animals[0].Id {
					t.Errorf("Expected a %v, got : %v", animals[0].Id, a)
				}
			},
		},
		{
			desc: "Select_Where_Like",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Where(db.Like(&db.Animal.Name, "%Cat%")).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Errorf("Expected %v animals, got %v", 2, len(a))
				}
			},
		},
		{
			desc: "Select_Order_By_Asc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a[0].Id > a[1].Id {
					t.Errorf("Expected animals order by asc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Order_By_Desc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).OrderByDesc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a[0].Id < a[1].Id {
					t.Errorf("Expected animals order by desc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Page",
			testCase: func(t *testing.T) {
				var a []Animal
				var pageSize uint = 5
				err = db.Select(db.Animal).Page(1, pageSize).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != int(pageSize) {
					t.Errorf("Expected %v animals, got %v", pageSize, len(a))
				}
			},
		},
		{
			desc: "Select_Join",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != len(animalFoods) {
					t.Errorf("Expected 1 animal, got %v", len(a))
				}
				if a[0].Name != animals[0].Name {
					t.Errorf("Expected %v, got %v", animals[0].Name, a[0].Name)
				}
			},
		},
		{
			desc: "Select_Join_Where",
			testCase: func(t *testing.T) {
				var f []Food
				err = db.Select(db.Food).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Where(
						db.Equals(&db.Animal.Name, animals[0].Name)).Scan(&f)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(f) != 1 {
					t.Fatalf("Expected 1 food, got %v", len(f))
				}
				if f[0].Name != foods[0].Name {
					t.Errorf("Expected %v, got %v", foods[0].Name, f[0].Name)
				}
			},
		},
		{
			desc: "Select_Join_Where_And_Equals_Find_0",
			testCase: func(t *testing.T) {
				var f []Food
				err = db.Select(db.Food).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Where(
						db.Equals(&db.Animal.Name, animals[0].Name),
						db.And(),
						db.Equals(&db.Food.Id, foods[1].Id),
					).Scan(&f)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(f) != 0 {
					t.Errorf("Expected 0 food, got %v", len(f))
				}
			},
		},
		{
			desc: "Select_Inverted_Join_Where_And_Equals_Find_0",
			testCase: func(t *testing.T) {
				var f []Food
				err = db.Select(db.Food).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Where(
						db.Equals(&db.Animal.Name, animals[0].Name),
						db.And(),
						db.Equals(&db.Food.Id, foods[1].Id),
					).Scan(&f)
				if !errors.Is(err, goe.ErrNotFound) {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(f) != 0 {
					t.Errorf("Expected 0 food, got %v", len(f))
				}
			},
		},
		{
			desc: "Select_Join_Where_And_Equals_Find_1",
			testCase: func(t *testing.T) {
				var f []Food
				err = db.Select(db.Food).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Where(
						db.Equals(&db.Animal.Name, animals[0].Name),
						db.And(),
						db.Equals(&db.Food.Id, foods[0].Id),
					).Scan(&f)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(f) != 1 {
					t.Errorf("Expected 1 food, got %v", len(f))
				}
			},
		},
		{
			desc: "Select_Join_Order_By_Asc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a[0].Id > a[1].Id {
					t.Errorf("Expected animals order by asc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Join_Order_By_Desc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					OrderByDesc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a[0].Id < a[1].Id {
					t.Errorf("Expected animals order by desc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Join_Where_Order_By_Asc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Where(
						db.Equals(&db.Food.Id, foods[0].Id),
					).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Fatalf("Expected 2 animals, got %v", len(a))
				}
				if a[0].Id > a[1].Id {
					t.Errorf("Expected animals order by asc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Join_Where_Order_By_Desc",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Where(
						db.Equals(&db.Food.Id, foods[0].Id),
					).OrderByDesc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Fatalf("Expected 2 animals, got %v", len(a))
				}
				if a[0].Id < a[1].Id {
					t.Errorf("Expected animals order by desc, got %v", a)
				}
			},
		},
		{
			desc: "Select_Join_Many_To_One",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Join(&db.Animal.IdHabitat, &db.Habitat.Id).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				for i := range a {
					if a[i].Id != animals[i].Id {
						t.Errorf("Expected %v, got %v", a[0].Id, animals[0].Id)
					}
				}
			},
		},
		{
			desc: "Select_Inverted_Join_Many_To_One",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Join(&db.Habitat.Id, &db.Animal.IdHabitat).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				for i := range a {
					if a[i].Id != animals[i].Id {
						t.Errorf("Expected %v, got %v", a[0].Id, animals[0].Id)
					}
				}
			},
		},
		{
			desc: "Select_Left_Join_Many_To_One",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).LeftJoin(&db.Habitat.Id, &db.Animal.IdHabitat).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != len(animals) {
					t.Errorf("Expected %v, got %v", len(animals), len(a))
				}
				if a[len(a)-1].IdHabitat != nil {
					t.Errorf("Expected nil, got value")
				}
			},
		},
		{
			desc: "Select_Join_Many_To_Many_And_Many_To_One",
			testCase: func(t *testing.T) {
				var f []Food
				err = db.Select(db.Food).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Animal.IdHabitat, &db.Habitat.Id).Where(db.Equals(&db.Habitat.Id, habitats[0].Id)).
					Scan(&f)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(f) != 2 {
					t.Errorf("Expected 2, got : %v", len(f))
				}
			},
		},
		{
			desc: "Select_Join_One_To_One",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Join(&db.Animal.IdInfo, &db.Info.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Errorf("Expected 2, got : %v", len(a))
				}
			},
		},
		{
			desc: "Select_Inverted_Join_One_To_One",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Join(&db.Info.Id, &db.Animal.IdInfo).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Errorf("Expected 2, got : %v", len(a))
				}
			},
		},
		{
			desc: "Select_Animal_Join_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Join(&db.Animal.IdInfo, &db.Info.Id).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 2 {
					t.Errorf("Expected 2, got : %v", len(a))
				}
			},
		},
		{
			desc: "Select_Info_Join_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var i []Info
				err = db.Select(db.Info).Join(&db.Animal.IdInfo, &db.Info.Id).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&i)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(i) != 2 {
					t.Errorf("Expected 2, got : %v", len(i))
				}
			},
		},
		{
			desc: "Select_Info_Inverted_Join_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var i []Info
				err = db.Select(db.Info).Join(&db.Animal.IdInfo, &db.Info.Id).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).
					Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&i)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(i) != 2 {
					t.Errorf("Expected 2, got : %v", len(i))
				}
			},
		},
		{
			desc: "Select_Info_Join_Status_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var s []Info
				err = db.Select(db.Info).Join(&db.Status.Id, &db.Info.IdStatus).Join(&db.Animal.IdInfo, &db.Info.Id).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&s)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(s) != 2 {
					t.Errorf("Expected 2, got : %v", len(s))
				}
			},
		},
		{
			desc: "Select_Status_Inverted_Join_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var s []Status
				err = db.Select(db.Status).Join(&db.Info.IdStatus, &db.Status.Id).Join(&db.Info.Id, &db.Animal.IdInfo).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&s)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(s) != 2 {
					t.Errorf("Expected 2, got : %v", len(s))
				}
			},
		},
		{
			desc: "Select_Status_Join_One_To_One_And_Many_To_Many",
			testCase: func(t *testing.T) {
				var s []Status
				err = db.Select(db.Status).Join(&db.Status.Id, &db.Info.IdStatus).Join(&db.Animal.IdInfo, &db.Info.Id).
					Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).Where(db.Equals(&db.Food.Id, foods[0].Id)).Scan(&s)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(s) != 2 {
					t.Errorf("Expected 2, got : %v", len(s))
				}
			},
		},
		{
			desc: "Select_Animal_By_Weather_Join_One_To_Many",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).
					Join(&db.Animal.IdHabitat, &db.Habitat.Id).
					Join(&db.Habitat.IdWeather, &db.Weather.Id).
					Where(db.Equals(&db.Weather.Id, weathers[3].Id)).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(a) != 4 {
					t.Errorf("Expected 4, got : %v", len(a))
				}
			},
		},
		{
			desc: "Select_Weather_By_Animal_Join_One_To_Many",
			testCase: func(t *testing.T) {
				var w []Weather
				err = db.Select(db.Weather).
					Join(&db.Weather.Id, &db.Habitat.IdWeather).
					Join(&db.Habitat.Id, &db.Animal.IdHabitat).
					Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&w)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(w) != 1 {
					t.Errorf("Expected 1, got : %v", len(w))
				}
			},
		},
		{
			desc: "Select_Join_Page",
			testCase: func(t *testing.T) {
				var a []Animal
				var pageSize uint = 2
				err = db.Select(db.Animal).Join(&db.Animal.Id, &db.AnimalFood.IdAnimal).
					Join(&db.Food.Id, &db.AnimalFood.IdFood).Page(1, pageSize).Scan(&a)
				if err != nil {
					t.Errorf("Expected a page select, got error: %v", err)
				}
				if len(a) != int(pageSize) {
					t.Errorf("Expected %v animals, got %v", pageSize, len(a))
				}
			},
		},
		{
			desc: "Select_Join_Name",
			testCase: func(t *testing.T) {
				var h Habitat
				err = db.Select(db.Habitat).
					Join(&db.Habitat.Name, &db.Weather.Name).Scan(&h)
				if err != nil {
					t.Fatalf("Expected a select, got error: %v", err)
				}
				if h.Name != "Ocean" {
					t.Errorf("Expected Ocean, got : %v", h.Name)
				}
			},
		},
		{
			desc: "Select_User_And_Roles",
			testCase: func(t *testing.T) {
				var q []struct {
					User    string
					Role    *string
					EndTime *time.Time
				}
				err = db.Select(&db.User.Name, &db.Role.Name, &db.UserRole.EndDate).
					LeftJoin(&db.User.Id, &db.UserRole.IdUser).
					LeftJoin(&db.UserRole.IdRole, &db.Role.Id).OrderByAsc(&db.User.Id).Scan(&q)
				if err != nil {
					t.Fatalf("Expected a select, got error: %v", err)
				}
				if len(q) != len(users) {
					t.Errorf("Expected %v, got : %v", len(users), len(q))
				}
				if q[0].EndTime == nil {
					t.Errorf("Expected a value, got : %v", q[0].EndTime)
				}
			},
		},
		{
			desc: "Select_Persons_And_Jobs",
			testCase: func(t *testing.T) {
				pj := []struct {
					Job    string
					Person string
				}{}
				err = db.Select(&db.Person.Name, &db.Job.Name).
					Join(&db.Person.Id, &db.PersonJob.IdPerson).
					Join(&db.PersonJob.IdJob, &db.Job.Id).Scan(&pj)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if len(pj) != len(personJobs) {
					t.Errorf("Expected %v, got : %v", len(personJobs), len(pj))
				}
			},
		},
		{
			desc: "Select_User_And_Roles_Inverted",
			testCase: func(t *testing.T) {
				var q []struct {
					User    string
					Role    *string
					EndTime *time.Time
				}
				err = db.Select(&db.User.Name, &db.Role.Name, &db.UserRole.EndDate).
					LeftJoin(&db.UserRole.IdUser, &db.User.Id).
					LeftJoin(&db.UserRole.IdRole, &db.Role.Id).OrderByAsc(&db.User.Id).Scan(&q)
				if err != nil {
					t.Fatalf("Expected a select, got error: %v", err)
				}
				if len(q) != len(users) {
					t.Fatalf("Expected %v, got : %v", len(users), len(q))
				}
				if q[0].EndTime == nil {
					t.Errorf("Expected a value, got : %v", q[0].EndTime)
				}
			},
		},
		{
			desc: "Select_Aggregate_Count",
			testCase: func(t *testing.T) {
				var c int
				err = db.Select(db.Count(&db.User.Name)).Scan(&c)
				if err != nil {
					t.Fatalf("Expected a select, got error: %v", err)
				}
				if c != len(users) {
					t.Errorf("Expected %v, got : %v", len(users), c)
				}
			},
		},
		{
			desc: "Select_Join_Aggregate_Count",
			testCase: func(t *testing.T) {
				var c int
				err = db.Select(db.Count(&db.User.Name)).
					LeftJoin(&db.UserRole.IdUser, &db.User.Id).
					LeftJoin(&db.UserRole.IdRole, &db.Role.Id).
					Scan(&c)
				if err != nil {
					t.Fatalf("Expected a select, got error: %v", err)
				}
				if c != len(users) {
					t.Errorf("Expected %v, got : %v", len(users), c)
				}
			},
		},
		{
			desc: "Select_Anonymous_Struct",
			testCase: func(t *testing.T) {
				var a struct {
					Id1 int
					Id2 string
				}
				err = db.Select(&db.Animal.Id, &db.Animal.Name).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a.Id1 != animals[0].Id {
					t.Errorf("Expected %v, got : %v", animals[0].Id, a.Id1)
				}
				if a.Id2 != animals[0].Name {
					t.Errorf("Expected %v, got : %v", animals[0].Name, a.Id2)
				}
			},
		},
		{
			desc: "Select_Anonymous_Struct_2",
			testCase: func(t *testing.T) {
				var a struct {
					Id  int
					Id2 string
				}
				err = db.Select(&db.Animal.Id, &db.Animal.Name).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				if a.Id != animals[0].Id {
					t.Errorf("Expected %v, got : %v", animals[0].Id, a.Id)
				}
				if a.Id2 != animals[0].Name {
					t.Errorf("Expected %v, got : %v", animals[0].Name, a.Id2)
				}
			},
		},
		{
			desc: "Select_Anonymous_Struct_Slice_3",
			testCase: func(t *testing.T) {
				var a []struct {
					AnimalId        int
					AnimalName      string
					AnimalIdHabitat uuid.UUID
					AnimalIdInfo    []byte
					HabitatId       uuid.UUID
					HabitatName     string
					IdWeather       int
					NameWeather     string
				}
				err = db.Select(db.Animal, db.Habitat).Join(&db.Animal.IdHabitat, &db.Habitat.Id).OrderByAsc(&db.Animal.Id).Scan(&a)
				if err != nil {
					t.Errorf("Expected a select, got error: %v", err)
				}
				for i := range a {
					if a[i].AnimalId != animals[i].Id {
						t.Errorf("Expected %v, got %v", a[0].AnimalId, animals[0].Id)
					}
					if a[i].AnimalIdHabitat.String() != a[i].HabitatId.String() {
						t.Errorf("Expected %v, got %v", a[i].AnimalIdHabitat.String(), a[i].HabitatId.String())
					}
				}
			},
		},
		{
			desc: "Select_Invalid_Scan",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Scan(a)
				if !errors.Is(err, goe.ErrInvalidScan) {
					t.Errorf("Expected goe.ErrInvalidScan, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Invalid_OrderBy",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).OrderByAsc(db.Animal.IdHabitat).Scan(&a)
				if !errors.Is(err, goe.ErrInvalidOrderBy) {
					t.Errorf("Expected goe.ErrInvalidOrderBy, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Invalid_Where",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.Animal).Where(db.Equals(db.Animal.Id, 1)).Scan(&a)
				if !errors.Is(err, goe.ErrInvalidWhere) {
					t.Errorf("Expected goe.ErrInvalidWhere, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Invalid_Arg",
			testCase: func(t *testing.T) {
				var a []Animal
				err = db.Select(db.DB).Join(&db.Animal.IdHabitat, &db.Habitat).Scan(&a)
				if !errors.Is(err, goe.ErrInvalidArg) {
					t.Errorf("Expected goe.ErrInvalidArg, got error: %v", err)
				}

				err = db.Select(nil).Join(db.Animal, db.Weather).Scan(&a)
				if !errors.Is(err, goe.ErrInvalidArg) {
					t.Errorf("Expected goe.ErrInvalidArg, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Context_Cancel",
			testCase: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				var a *int
				err = db.SelectContext(ctx, &db.Animal.Id).Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&a)
				if !errors.Is(err, context.Canceled) {
					t.Errorf("Expected a context.Canceled, got error: %v", err)
				}
			},
		},
		{
			desc: "Select_Context_Timeout",
			testCase: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond*1)
				defer cancel()
				var a *int
				err = db.SelectContext(ctx, &db.Animal.Id).Where(db.Equals(&db.Animal.Id, animals[0].Id)).Scan(&a)
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("Expected a context.DeadlineExceeded, got error: %v", err)
				}
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, tC.testCase)
	}
}
