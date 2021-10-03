package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func rest(dao Dao, endpoint string, hostname string) (func(), error) {
	id := NewId(hostname)
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //looks important
	rapi := router.Group("/api")
	header := func(r *http.Request, name string) (string, error) {
		key := http.CanonicalHeaderKey(name)
		values := r.Header[key]
		if len(values) != 1 {
			return "", fmt.Errorf("Invalid header %s %v", name, values)
		}
		return values[0], nil
	}
	rapi.POST("/mail", func(c *gin.Context) {
		mid := id.Next()
		mfrom, err := header(c.Request, "Mail-From")
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		mto, err := header(c.Request, "Mail-To")
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		msubject, err := header(c.Request, "Mail-Subject")
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		mmime, err := header(c.Request, "Mail-Mime")
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		mbody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		err = mailSend(dao, mid, mfrom, mto, msubject, mmime, string(mbody))
		if err != nil {
			c.JSON(400, gin.H{"id": mid, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"id": mid,
		})
	})
	rapi.GET("/domain/:name", func(c *gin.Context) {
		name := c.Param("name")
		dro, err := dao.GetDomain(name)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"name": name,
			"pub":  dro.PublicKey,
			"key":  dro.PrivateKey,
		})
	})
	rapi.GET("/domain", func(c *gin.Context) {
		names, err := dao.GetDomains()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"names": names,
		})
	})
	rapi.POST("/domain/:name", func(c *gin.Context) {
		name := c.Param("name")
		pub, key, err := keygen()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		pubs := string(pub)
		keys := string(key)
		err = dao.AddDomain(name, pubs, keys)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"name": name,
			"pub":  pubs,
			"key":  key,
		})
	})
	rapi.DELETE("/domain/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := dao.DelDomain(name)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{})
	})
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	port := listen.Addr().(*net.TCPAddr).Port
	log.Println("port", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	exit := make(chan interface{})
	go func() {
		err = server.Serve(listen)
		if err != nil {
			log.Println(err)
		}
		close(exit)
	}()
	closer := func() {
		ctx := context.Background()
		server.Shutdown(ctx)
		<-exit
	}
	return closer, nil
}
