package main

import (
	"backend-d-embung/Controller"
	"backend-d-embung/Database"
	_ "crypto/sha512"
	_ "encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// if using file .env use code below
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal(err.Error())
	// }

	//Database
	db := Database.Open()
	if db != nil {
		println("Database Terhubung..")
	}

	// Gin Framework
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.SetTrustedProxies(
		[]string{
			os.Getenv("PROXY_1"),
			os.Getenv("PROXY_2"),
			os.Getenv("PROXY_3"),
		},
	)

	//CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})

	//Routers
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Running...",
			"success": true,
		})
	})
	Controller.OperasionalController(r)
	Controller.TestimoniController(r)
	Controller.NewsController(r)
	Controller.Authorization(r)
	Controller.Post(r)
	Controller.AdditionalInfo(r)
	if err := r.Run(); err != nil {
		log.Fatal(err.Error())
		return
	}
}
