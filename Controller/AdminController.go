package Controller

import (
	_ "backend-d-embung/Auth"
	"backend-d-embung/Entities"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func OperasionalController(db *gorm.DB, r *gin.Engine) {
	// post a new operational
	r.POST("/operasional", func(c *gin.Context) {
		var input Entities.Operasional

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "input is invalid",
				"error":   err.Error(),
			})
			return
		}
		newOperasional := Entities.Operasional {
			HariOperasional: input.HariOperasional,
			JamOperasional: input.JamOperasional,
			Harga: input.Harga,
		}

		if err := db.Create(&newOperasional); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "can't create new operational",
				"error": err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "a new operational has successfully created",
			"data":    newOperasional.CreatedAt,
		})
	})

	// get operational
	r.GET("/operasional", func(c *gin.Context) {

		var operasional Entities.Operasional

		if err := db.Order("id desc").Take(&operasional); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "operational isn't available",
				"error": err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"message": "success querying latest operational",
			"data":    operasional,
		})
	})

	// patch operational time
	r.PATCH("/operasional/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")

		var input Entities.Operasional

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H {
				"success": false,
				"message": "input is invalid",
				"error": err.Error(),
			})
		}

		patchOperasional := Entities.Operasional {
			HariOperasional: input.HariOperasional,
			JamOperasional: input.JamOperasional,
			Harga: input.Harga,
		}

		result := db.Where("id = ?", id).Model(&patchOperasional).Updates(patchOperasional)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating new operational.",
				"error":   result.Error.Error(),
			})
			return
		}

		if result = db.Where("id = ?", id).Take(&patchOperasional); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "operational not found.",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    patchOperasional,
		})
	})

	// delete operational time by id
	r.DELETE("/operasional/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")

		var operasional Entities.Operasional

		if err := db.Where("id = ?", id).Delete(&operasional); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})
}

func ArticleController(db *gorm.DB, r *gin.Engine) {
	r.Static("/article/image", "./Images")

	// post new article
	r.POST("/article", func(c *gin.Context) {
		image, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"success": false,
				"error":   "get form err: " + err.Error(),
			})
			return
		}

		rand.Seed(time.Now().Unix())

		str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		shuff := []rune(str)

		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})
		image.Filename = string(shuff)
		
		godotenv.Load("../.env")

		if err := c.SaveUploadedFile(image, "./Images/"+image.Filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"error":   "upload file err: " + err.Error(),
			})
			return
		}
		enText := slug.MakeLang(c.PostForm("title"), "en")

		newArticle := Entities.Artikel {
			Title: c.PostForm("title"),
			Slug: enText,
			Image: os.Getenv("BASE_URL") + "/article/image/" + image.Filename,
			Excerpt: c.PostForm("excerpt"),
			Body: c.PostForm("body"),
		}

		if err := db.Create(&newArticle); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error when inserting a new article",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":     true,
			"message":     "a new article has successfully created",
			"error":       nil,
			"judul_article": newArticle.Title,
		})
	})

	// get all article
	r.GET("/article", func(c *gin.Context) {
		var allArticle []Entities.Artikel

		if res := db.Find(&allArticle); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"success": false,
				"message": "failed when query all article",
				"error": res.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"error": nil,
			"data": allArticle,
		})
	})

	// -> 3 api below using slug params
	// get detail article 
	r.GET("/article/:slug", func(c *gin.Context) {
		slug, _ := c.Params.Get("slug")

		var article Entities.Artikel

		if result := db.Where("slug = ?", slug).Take(&article); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query article successful.",
			"error":   nil,
			"data":    article,
		})
		// get all data except body article
	})

	// patch article 
	r.PATCH("/article/:slug", func(c *gin.Context) {
		slug, _ := c.FormFile("slug")

		var article Entities.Artikel

		if res := db.Where("id = ?", slug).Take(&article); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   res.Error.Error(),
			})
			return
		}

		newArticle := Entities.Artikel {
			Title: c.PostForm("title"),
			Slug: strings.ToLower(strings.ReplaceAll(c.PostForm("title"), " ", "-")),
			Excerpt: c.PostForm("excerpt"),
			Body: c.PostForm("body"),
		}

		image, _ := c.FormFile("image")

		if image == nil {
			newArticle = Entities.Artikel{
				Image: article.Image,
			}
		} else {
			rand.Seed(time.Now().Unix())

			str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

			shuff := []rune(str)

			rand.Shuffle(len(shuff), func(i, j int) {
				shuff[i], shuff[j] = shuff[j], shuff[i]
			})
			image.Filename = string(shuff)

			godotenv.Load("../.env")

			if err := c.SaveUploadedFile(image, "./Images/"+image.Filename); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"Success": false,
					"error":   "upload file err: " + err.Error(),
				})
				return
			}

			newArticle = Entities.Artikel{
				Image: os.Getenv("BASE_URL") + "/article/image/" + image.Filename,
			}
		}

		if err := db.Where("slug = ?", slug).Model(&article).Updates(newArticle); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error when inserting a new agenda",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "a new article has successfully updated",
			"error":   nil,
			"data":    newArticle,
		})
	})

	// delete article
	r.DELETE("/article/:slug", func(c *gin.Context) {
		slug, _ := c.Params.Get("slug")

		var article Entities.Artikel

		if res := db.Where("slug = ?", slug).Delete(&article); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   res.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})
}