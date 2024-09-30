package test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ra1nz0r/effective_mobile-1/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestTableExists(t *testing.T) {
	type args struct {
		tableName   string
		buildEXPECT func(mock sqlmock.Sqlmock)
		boolASSERT  func(t assert.TestingT, exists bool)
		errorASSERT func(t assert.TestingT, err error)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success.",
			args: args{
				tableName: "test_table",
				buildEXPECT: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`
						SELECT EXISTS \(
							SELECT FROM pg_tables
							WHERE schemaname = 'public' OR schemaname = 'private'
							AND tablename = \$1
						\);`).WithArgs("test_table").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
				},
				boolASSERT: func(t assert.TestingT, exists bool) {
					assert.True(t, exists)
				},
				errorASSERT: func(t assert.TestingT, err error) {
					assert.NoError(t, err)
				},
			},
		},
		{
			name: "Table does not exist.",
			args: args{
				tableName: "nonexistent_table",
				buildEXPECT: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`
						SELECT EXISTS \(
							SELECT FROM pg_tables
							WHERE schemaname = 'public' OR schemaname = 'private'
							AND tablename = \$1
						\);`).WithArgs("nonexistent_table").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
				},
				boolASSERT: func(t assert.TestingT, exists bool) {
					assert.False(t, exists)
				},
				errorASSERT: func(t assert.TestingT, err error) {
					assert.NoError(t, err)
				},
			},
		},
		{
			name: "Error Query.",
			args: args{
				tableName: "test_table",
				buildEXPECT: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`
						SELECT EXISTS \(
							SELECT FROM pg_tables
							WHERE schemaname = 'public' OR schemaname = 'private'
							AND tablename = \$1
						\);`).WithArgs("test_table").WillReturnError(errors.New("query error"))
				},
				boolASSERT: func(t assert.TestingT, exists bool) {
					assert.False(t, exists)
				},
				errorASSERT: func(t assert.TestingT, err error) {
					assert.Error(t, err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание mock объекта для базы данных
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Ожидание запроса и его успешный результат
			tt.args.buildEXPECT(mock)

			// Вызов тестируемой функции
			exists, err := services.TableExists(db, tt.args.tableName)

			// Проверка результата
			tt.args.errorASSERT(t, err)
			tt.args.boolASSERT(t, exists)

			// Проверка, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
