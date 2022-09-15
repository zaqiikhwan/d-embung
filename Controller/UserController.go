package Controller

import (
	"backend-d-embung/Entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestimoniController(db *gorm.DB, r *gin.Engine) {
	r.POST("/testimoni", func(c *gin.Context) {
		var input Entities.Testimoni

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "input is invalid",
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		newTestimoni := Entities.Testimoni {
			Identitas: input.Identitas,
			Testimoni: input.Testimoni,
		}

		if err := db.Create(&newTestimoni); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "can't create new operational",
				"success": false,
				"error": err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"error": nil,
			"data": newTestimoni.Testimoni,
		})
	})

	r.GET("/testimoni", func(c *gin.Context) {
		var allTestimoni []Entities.Testimoni

		if err := db.Find(&allTestimoni); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "can't create new operational",
				"success": false,
				"error": err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"error": nil,
			"data": allTestimoni,
		})
	})
}