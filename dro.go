package main

import "time"

type DomainDro struct {
	Name       string `gorm:"primaryKey"`
	PrivateKey string
	PublicKey  string
}

type MessageDro struct {
	Mid     string `gorm:"primaryKey"`
	From    string
	To      string
	Subject string
	Mime    string
	Body    string
	Created time.Time
}

type AttemptDro struct {
	Mid     string
	Created time.Time
	Addr    string
	Dial    bool
	Error   string
}
