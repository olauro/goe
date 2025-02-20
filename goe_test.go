package goe_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/olauro/goe"
)

type MockDriver struct {
}

func (md *MockDriver) Name() string {
	return "Mock"
}

func (md *MockDriver) MigrateContext(context.Context, *goe.Migrator, goe.Connection) (string, error) {
	return "", nil
}

func (md *MockDriver) DropTable(string, goe.Connection) (string, error) {
	return "", nil
}

func (md *MockDriver) DropColumn(string, string, goe.Connection) (string, error) {
	return "", nil
}

func (md *MockDriver) RenameColumn(string, string, string, goe.Connection) (string, error) {
	return "", nil
}

func (md *MockDriver) Init() {
}

func (md *MockDriver) KeywordHandler(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}

func (md *MockDriver) Returning([]byte) []byte {
	return nil
}

func (md *MockDriver) Select() []byte {
	return nil
}

func (md *MockDriver) From() []byte {
	return nil
}

func (md *MockDriver) Where() []byte {
	return nil
}

func (md *MockDriver) Insert() []byte {
	return nil
}

func (md *MockDriver) Values() []byte {
	return nil
}

func (md *MockDriver) Update() []byte {
	return nil
}

func (md *MockDriver) Set() []byte {
	return nil
}

func (md *MockDriver) Delete() []byte {
	return nil
}

func (md *MockDriver) NewConnection() goe.Connection {
	return nil
}

func (md *MockDriver) NewTransaction(ctx context.Context, opts *sql.TxOptions) (goe.Transaction, error) {
	return nil, nil
}

func (md *MockDriver) Stats() sql.DBStats {
	return sql.DBStats{}
}

func TestMapDatabase(t *testing.T) {
	type User struct {
		Id   uint
		Name string
	}

	type UserLog struct {
		Id     uint
		Action string
		IdUser uint
	}

	type Database struct {
		User    *User
		UserLog *UserLog
	}

	_, err := goe.Open[Database](&MockDriver{}, goe.Config{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMapDatabaseErrorPrimaryKey(t *testing.T) {
	type User struct {
		IdUser string
		Name   string
	}

	type Database struct {
		User *User
	}

	_, err := goe.Open[Database](&MockDriver{}, goe.Config{})
	if !errors.Is(err, goe.ErrStructWithoutPrimaryKey) {
		t.Fatal("Was expected a goe.ErrStructWithoutPrimaryKey but get:", err)
	}
}
