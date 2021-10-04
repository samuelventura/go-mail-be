package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type daoDso struct {
	db *gorm.DB
}

type Dao interface {
	Close() error
	GetDomains() ([]string, error)
	AddDomain(name string, privkey string, pubkey string) error
	GetDomain(name string) (*DomainDro, error)
	DelDomain(name string) error
	AddMessage(dro *MessageDro) error
	AddAttempt(dro *AttemptDro) error
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
	config := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
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

func (dso *daoDso) GetDomains() ([]string, error) {
	var dros []DomainDro
	result := dso.db.Find(&dros)
	if result.Error != nil {
		return nil, result.Error
	}
	names := make([]string, 0, result.RowsAffected)
	for _, dro := range dros {
		names = append(names, dro.Name)
	}
	return names, nil
}

func (dso *daoDso) AddDomain(name string, pubkey string, privkey string) error {
	dro := &DomainDro{Name: name, PrivateKey: privkey, PublicKey: pubkey}
	result := dso.db.Create(dro)
	return result.Error
}

func (dso *daoDso) GetDomain(name string) (*DomainDro, error) {
	dro := &DomainDro{}
	result := dso.db.Where("name = ?", name).First(dro)
	return dro, result.Error
}

func (dso *daoDso) DelDomain(name string) error {
	dro := &DomainDro{}
	result := dso.db.Where("name = ?", name).Delete(dro)
	if result.Error == nil && result.RowsAffected != 1 {
		return fmt.Errorf("row not found")
	}
	return result.Error
}

func (dso *daoDso) AddMessage(dro *MessageDro) error {
	result := dso.db.Create(dro)
	return result.Error
}

func (dso *daoDso) AddAttempt(dro *AttemptDro) error {
	result := dso.db.Create(dro)
	return result.Error
}
