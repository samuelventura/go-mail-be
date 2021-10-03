package main

import (
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	srcdef, err := withext("db3")
	if err != nil {
		return err
	}
	driver := getenv("MAIL_DB_DRIVER", "sqlite")
	source := getenv("MAIL_DB_SOURCE", srcdef)
	dao, err := NewDao(driver, source)
	if err != nil {
		return err
	}
	defer dao.Close()
	// pub, key, err := keygen()
	// if err != nil {
	// 	return err
	// }
	//dao.AddDomain("laurelview.io", string(pub), string(key))
	err = sendText(dao, "samuel@laurelview.io", "samuel.ventura@yeico.com", "go dkim test", "this is a test!")
	if err != nil {
		return err
	}
	return nil
}
