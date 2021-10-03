package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	log.Println("starting...")
	defer log.Println("exit")
	closer, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	exit := make(chan interface{})
	go func() {
		defer close(exit)
		ioutil.ReadAll(os.Stdin)
	}()
	select {
	case <-ctrlc:
	case <-exit:
	}
}

func run() (func(), error) {
	srcdef, err := withext("db3")
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	//FIXME gorm setup logging
	//FIXME gin setup logging
	driver := getenv("MAIL_DB_DRIVER", "sqlite")
	source := getenv("MAIL_DB_SOURCE", srcdef)
	endpoint := getenv("MAIL_ENDPOINT", "127.0.0.1:31650")
	dao, err := NewDao(driver, source)
	if err != nil {
		return nil, err
	}
	closer, err := rest(dao, endpoint, hostname)
	if err != nil {
		return nil, err
	}
	return func() {
		closer()
		err := dao.Close()
		if err != nil {
			log.Println(err)
		}
	}, nil
}
