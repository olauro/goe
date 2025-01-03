package goe

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/olauro/goe/utils"
)

var ErrStructWithoutPrimaryKey = errors.New("goe")

func Open(db any, driver Driver, config Config) error {
	valueOf := reflect.ValueOf(db)
	if valueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("%v: the target value needs to be pass as a pointer", pkg)
	}
	dbTarget := new(DB)
	valueOf = valueOf.Elem()

	dbTarget.AddrMap = make(map[uintptr]Field)

	// set value for Fields
	for i := 0; i < valueOf.NumField(); i++ {
		if valueOf.Field(i).IsNil() {
			valueOf.Field(i).Set(reflect.ValueOf(reflect.New(valueOf.Field(i).Type().Elem()).Interface()))
		}
	}

	var err error
	// init Fields
	for i := 0; i < valueOf.NumField(); i++ {
		if valueOf.Field(i).Elem().Type().Name() != "DB" {
			err = initField(valueOf, valueOf.Field(i).Elem(), dbTarget, driver)
			if err != nil {
				return err
			}
		}
	}

	dbTarget.Driver = driver
	dbTarget.Driver.Init(dbTarget)
	dbTarget.Config = &config
	valueOf.FieldByName("DB").Set(reflect.ValueOf(dbTarget))
	return nil
}

func initField(tables reflect.Value, valueOf reflect.Value, db *DB, driver Driver) error {
	pks, FieldNames, err := getPk(valueOf.Type(), driver)
	if err != nil {
		return err
	}

	for i := range pks {
		db.AddrMap[uintptr(valueOf.FieldByName(FieldNames[i]).Addr().UnsafePointer())] = pks[i]
	}
	var Field reflect.StructField

	for i := 0; i < valueOf.NumField(); i++ {
		Field = valueOf.Type().Field(i)
		//skip primary key
		if slices.Contains(FieldNames, Field.Name) {
			//TODO: Check this
			table, prefix := checkTablePattern(tables, Field)
			if table == "" && prefix == "" {
				continue
			}
		}
		switch valueOf.Field(i).Kind() {
		case reflect.Slice:
			err := handlerSlice(tables, valueOf.Field(i).Type().Elem(), valueOf, i, pks, db, driver)
			if err != nil {
				return err
			}
		case reflect.Struct:
			handlerStruct(valueOf.Field(i).Type(), valueOf, i, pks[0], db, driver)
		case reflect.Ptr:
			table, prefix := checkTablePattern(tables, valueOf.Type().Field(i))
			if table != "" {
				if mto := isManyToOne(tables, valueOf.Type(), driver, table, prefix); mto != nil {
					switch v := mto.(type) {
					case *manyToOne:
						if v == nil {
							newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
							break
						}
						key := driver.KeywordHandler(utils.TableNamePattern(table))
						db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
						for _, pk := range pks {
							if pk.structAttributeName == prefix || pk.structAttributeName == prefix+table {
								pk.fks[key] = mto
								v.pk = pk
							}
						}
					case *oneToOne:
						if v == nil {
							newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
							break
						}
						key := driver.KeywordHandler(utils.TableNamePattern(table))
						db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
						for _, pk := range pks {
							if pk.structAttributeName == prefix || pk.structAttributeName == prefix+table {
								pk.fks[key] = mto
								v.pk = pk
							}
						}
					}
					continue
				}
			}
			newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
		default:
			table, prefix := checkTablePattern(tables, valueOf.Type().Field(i))
			if table != "" {
				if mto := isManyToOne(tables, valueOf.Type(), driver, table, prefix); mto != nil {
					switch v := mto.(type) {
					case *manyToOne:
						if v == nil {
							newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
							break
						}
						key := driver.KeywordHandler(utils.TableNamePattern(table))
						db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
						for _, pk := range pks {
							if pk.structAttributeName == prefix {
								pk.fks[key] = mto
								v.pk = pk
							} else if pk.structAttributeName == v.structAttributeName {
								pk.fks[key] = mto
								pk.autoIncrement = false
								v.pk = pk
							}
						}
					case *oneToOne:
						if v == nil {
							newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
							break
						}
						key := driver.KeywordHandler(utils.TableNamePattern(table))
						db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
						for _, pk := range pks {
							if pk.structAttributeName == prefix || pk.structAttributeName == prefix+table {
								pk.fks[key] = mto
								v.pk = pk
							}
						}
					}
					continue
				}
			}
			newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
		}
	}
	return nil
}

func handlerStruct(targetTypeOf reflect.Type, valueOf reflect.Value, i int, p *pk, db *DB, driver Driver) {
	switch targetTypeOf.Name() {
	case "Time":
		newAttr(valueOf, i, p, uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
	}
}

func handlerSlice(tables reflect.Value, targetTypeOf reflect.Type, valueOf reflect.Value, i int, pks []*pk, db *DB, driver Driver) error {
	switch targetTypeOf.Kind() {
	case reflect.Uint8:
		table, prefix := checkTablePattern(tables, valueOf.Type().Field(i))
		if table != "" {
			if mto := isManyToOne(tables, valueOf.Type(), driver, table, prefix); mto != nil {
				switch v := mto.(type) {
				case *manyToOne:
					if v == nil {
						break
					}
					key := driver.KeywordHandler(utils.TableNamePattern(table))
					db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
					for _, pk := range pks {
						if pk.structAttributeName == prefix || pk.structAttributeName == prefix+table {
							pk.fks[key] = mto
							v.pk = pk
						}
					}
				case *oneToOne:
					if v == nil {
						break
					}
					key := driver.KeywordHandler(utils.TableNamePattern(table))
					db.AddrMap[uintptr(valueOf.Field(i).Addr().UnsafePointer())] = v
					for _, pk := range pks {
						if pk.structAttributeName == prefix || pk.structAttributeName == prefix+table {
							pk.fks[key] = mto
							v.pk = pk
						}
					}
				}
				break
			}
		}
		//TODO: Check this
		valueOf.Field(i).SetBytes([]byte{})
		newAttr(valueOf, i, pks[0], uintptr(valueOf.Field(i).Addr().UnsafePointer()), db, driver)
	}
	return nil
}

func newAttr(valueOf reflect.Value, i int, p *pk, addr uintptr, db *DB, d Driver) {
	at := createAtt(
		valueOf.Type().Field(i).Name,
		p,
		d,
	)
	db.AddrMap[addr] = at
}

func getPk(typeOf reflect.Type, driver Driver) ([]*pk, []string, error) {
	var pks []*pk
	var FieldsNames []string

	id, valid := typeOf.FieldByName("Id")
	if valid {
		pks := make([]*pk, 1)
		FieldsNames = make([]string, 1)
		pks[0] = createPk([]byte(typeOf.Name()), id.Name, isAutoIncrement(id), driver)
		FieldsNames[0] = id.Name
		return pks, FieldsNames, nil
	}

	Fields := fieldsByTags("pk", typeOf)
	if len(Fields) == 0 {
		return nil, nil, fmt.Errorf("%w: struct %q don't have a primary key setted", ErrStructWithoutPrimaryKey, typeOf.Name())
	}

	pks = make([]*pk, len(Fields))
	FieldsNames = make([]string, len(Fields))
	for i := range Fields {
		pks[i] = createPk([]byte(typeOf.Name()), Fields[i].Name, isAutoIncrement(Fields[i]), driver)
		FieldsNames[i] = Fields[i].Name
	}

	return pks, FieldsNames, nil
}

func isAutoIncrement(id reflect.StructField) bool {
	return strings.Contains(id.Type.Kind().String(), "int")
}

func isManyToOne(tables reflect.Value, typeOf reflect.Type, driver Driver, table, prefix string) Field {
	for c := 0; c < tables.NumField(); c++ {
		if tables.Field(c).Elem().Type().Name() == table {
			for i := 0; i < tables.Field(c).Elem().NumField(); i++ {
				// check if there is a slice to typeOf
				if tables.Field(c).Elem().Field(i).Kind() == reflect.Slice {
					if tables.Field(c).Elem().Field(i).Type().Elem().Name() == typeOf.Name() {
						return createManyToOne(tables.Field(c).Elem().Type(), typeOf, driver, prefix)
					}
				}
			}
			if tableMtm := strings.ReplaceAll(typeOf.Name(), table, ""); tableMtm != typeOf.Name() {
				typeOfMtm := tables.FieldByName(tableMtm)
				if typeOfMtm.IsValid() && !typeOfMtm.IsZero() {
					typeOfMtm = typeOfMtm.Elem()
					for i := 0; i < typeOfMtm.NumField(); i++ {
						if typeOfMtm.Field(i).Kind() == reflect.Slice && typeOfMtm.Field(i).Type().Elem().Name() == table {
							return createManyToOne(typeOfMtm.Field(i).Type().Elem(), typeOf, driver, prefix)
						}
					}
				}
			}
			return createOneToOne(tables.Field(c).Elem().Type(), typeOf, driver, prefix)
		}
	}
	return nil
}

func primaryKeys(str reflect.Type) (pks []reflect.StructField) {
	Field, exists := str.FieldByName("Id")
	if exists {
		pks := make([]reflect.StructField, 1)
		pks[0] = Field
		return pks
	} else {
		//TODO: Return anonymous pk para len(pks) == 0
		return fieldsByTags("pk", str)
	}
}

func fieldsByTags(tag string, str reflect.Type) (f []reflect.StructField) {
	f = make([]reflect.StructField, 0)

	for i := 0; i < str.NumField(); i++ {
		if strings.Contains(str.Field(i).Tag.Get("goe"), tag) {
			f = append(f, str.Field(i))
		}
	}
	return f
}

func getTagValue(FieldTag string, subTag string) string {
	values := strings.Split(FieldTag, ";")
	for _, v := range values {
		if after, found := strings.CutPrefix(v, subTag); found {
			return after
		}
	}
	return ""
}

func checkTablePattern(tables reflect.Value, Field reflect.StructField) (table, prefix string) {
	table = getTagValue(Field.Tag.Get("goe"), "table:")
	if table != "" {
		prefix = strings.ReplaceAll(Field.Name, table, "")
		return table, prefix
	}
	if table == "" {
		for r := len(Field.Name) - 1; r > 1; r-- {
			if Field.Name[r] < 'a' {
				table = Field.Name[r:]
				prefix = Field.Name[:r]
				if tables.FieldByName(table).IsValid() {
					return table, prefix
				}
			}
		}
		if !tables.FieldByName(table).IsValid() {
			table = ""
		}
	}
	return table, prefix
}
