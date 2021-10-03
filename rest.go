package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func rest(dao Dao, endpoint string) (func(), error) {
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //looks important
	rapi := router.Group("/api")
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
