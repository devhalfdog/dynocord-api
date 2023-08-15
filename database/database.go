package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(host, port, user, password, dbname string) *Model {
	return &Model{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}
}

func (m *Model) Connect() error {
	var err error

	m.Connector, err = gorm.Open(mysql.Open(m.getDsn()), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("database connection ok")

	return nil
}

func (m *Model) getDsn() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)

	return dsn
}
