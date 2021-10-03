package main

import "time"

//.tables
//.dump domain_dros | message_dros | attempt_dros
//.quit

type DomainDro struct {
	Name       string `gorm:"primaryKey"`
	PrivateKey string
	PublicKey  string
}

type MessageDro struct {
	ID      string
	From    string
	To      string
	Subject string
	Mime    string
	Body    string
	Created time.Time
}

type AttemptDro struct {
	ID      string
	Created time.Time
	Result  string
}
