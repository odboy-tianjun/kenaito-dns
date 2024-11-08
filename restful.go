package main

import (
	"github.com/gin-gonic/gin"
)

func initRestFunc(r *gin.Engine) {
	//records := selectResolveRecords()
	//jsonData, err := json.Marshal(records)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(jsonData))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
