package main

import "log"

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
	dao.Close()
	return nil
}
