package testutils

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB initializes a GORM DB with sqlmock
func SetupMockDB() (*gorm.DB, *sql.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return nil, nil, nil, err
	}

	return gormDB, sqlDB, mock, nil
}
