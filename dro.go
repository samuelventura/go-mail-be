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
	ID       string
	Mime     string
	Body     string
	From     string
	To       string
	Created  time.Time
	Updated  time.Time
	Dropped  *time.Time
	Sent     *time.Time
	Attempts int
	Result   string
}

type AttemptDro struct {
	ID      string
	Attempt int
	Created time.Time
	Result  string
}
