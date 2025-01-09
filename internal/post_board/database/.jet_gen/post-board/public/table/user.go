//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var User = newUserTable("public", "user", "")

type userTable struct {
	postgres.Table

	// Columns
	ID            postgres.ColumnString
	Email         postgres.ColumnString
	Name          postgres.ColumnString
	EncryptedPass postgres.ColumnString
	CreatedAt     postgres.ColumnTimestamp
	UpdatedAt     postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UserTable struct {
	userTable

	EXCLUDED userTable
}

// AS creates new UserTable with assigned alias
func (a UserTable) AS(alias string) *UserTable {
	return newUserTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserTable with assigned schema name
func (a UserTable) FromSchema(schemaName string) *UserTable {
	return newUserTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserTable with assigned table prefix
func (a UserTable) WithPrefix(prefix string) *UserTable {
	return newUserTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserTable with assigned table suffix
func (a UserTable) WithSuffix(suffix string) *UserTable {
	return newUserTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserTable(schemaName, tableName, alias string) *UserTable {
	return &UserTable{
		userTable: newUserTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newUserTableImpl("", "excluded", ""),
	}
}

func newUserTableImpl(schemaName, tableName, alias string) userTable {
	var (
		IDColumn            = postgres.StringColumn("id")
		EmailColumn         = postgres.StringColumn("email")
		NameColumn          = postgres.StringColumn("name")
		EncryptedPassColumn = postgres.StringColumn("encrypted_pass")
		CreatedAtColumn     = postgres.TimestampColumn("created_at")
		UpdatedAtColumn     = postgres.TimestampColumn("updated_at")
		allColumns          = postgres.ColumnList{IDColumn, EmailColumn, NameColumn, EncryptedPassColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns      = postgres.ColumnList{EmailColumn, NameColumn, EncryptedPassColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return userTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:            IDColumn,
		Email:         EmailColumn,
		Name:          NameColumn,
		EncryptedPass: EncryptedPassColumn,
		CreatedAt:     CreatedAtColumn,
		UpdatedAt:     UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
