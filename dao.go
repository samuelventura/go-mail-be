package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type daoDso struct {
	db *gorm.DB
}

type Dao interface {
	Close() error
	AddDomain(name string, privkey string, pubkey string) error
	GetDomainPrivateKey(name string) (*DomainDro, error)
}

func Dialector(driver string, source string) (gorm.Dialector, error) {
	switch driver {
	case "sqlite":
		return sqlite.Open(source), nil
	case "postgres":
		return postgres.Open(source), nil
	}
	return nil, fmt.Errorf("unknown driver %s", driver)
}

func NewDao(driver string, source string) (Dao, error) {
	config := &gorm.Config{}
	dialector, err := Dialector(driver, source)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&DomainDro{}, &MessageDro{}, &AttemptDro{})
	if err != nil {
		return nil, err
	}
	return &daoDso{db}, nil
}

func (dso *daoDso) Close() error {
	sqlDB, err := dso.db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (dso *daoDso) AddDomain(name string, pubkey string, privkey string) error {
	row := &DomainDro{Name: name, PrivateKey: privkey, PublicKey: pubkey}
	result := dso.db.Create(row)
	return result.Error
}

func (dso *daoDso) GetDomainPrivateKey(name string) (*DomainDro, error) {
	row := &DomainDro{}
	result := dso.db.Where("name = ?", name).First(row)
	return row, result.Error
}
