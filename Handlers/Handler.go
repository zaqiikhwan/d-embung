package Handlers

import (
	"backend-d-embung/Database"
	"backend-d-embung/Entities"
	"bytes"
	"html"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gosimple/slug"
	storage_go "github.com/supabase-community/storage-go"
	"golang.org/x/crypto/bcrypt"
)

func PostNewTestimoni(c *gin.Context) {
	var input Entities.Testimoni

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "input is invalid",
			"success": false,
			"statusCode": http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}

	newTestimoni := Entities.Testimoni {
		Identitas: input.Identitas,
		Testimoni: input.Testimoni,
	}

	if err := Database.Open().Create(&newTestimoni); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can't create new operational",
			"success": false,
			"statusCode": http.StatusInternalServerError,
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"success": true,
		"error": nil,
		"statusCode": http.StatusCreated,
		"data": newTestimoni.Testimoni,
	})
}

func GetAllTestimoni(c *gin.Context) {
	var allTestimoni []Entities.Testimoni

	if err := Database.Open().Order("id desc").Limit(3).Find(&allTestimoni); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can't find testimonies",
			"success": false,
			"statusCode": http.StatusNotFound,
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"error": nil,
		"statusCode": http.StatusOK,
		"data": allTestimoni,
	})
}

func RegisterAdmin(c *gin.Context) {
	var input Entities.Admin

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "input should bind json",
			"statusCode": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := Entities.Admin {
		Nickname: input.Nickname,
		Password: string(hashedPassword),
	}

	if err := Database.Open().Create(&user); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "failed when creating a new data user",
			"success": false,
			"statusCode": http.StatusInternalServerError,
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"message": "registered successfully",
		"success": true,
		"statusCode": http.StatusCreated,
		"error": nil,
	})
}

func LoginAdmin(c *gin.Context) {
	var input Entities.Admin

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "input must bind with json",
			"statusCode": http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}

	var user Entities.Admin

	if err := Database.Open().Where("nickname = ?", input.Nickname).Take(&user); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "nickname Anda tidak sesuai.",
			"statusCode": http.StatusBadRequest,
			"error":   err.Error.Error(),
		})
		return
	}

	hashedInput, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	error := bcrypt.CompareHashAndPassword([]byte(user.Password), hashedInput)
	
	if error != nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(time.Hour * 30 * 24).Unix(),
		})
		strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"statusCode": http.StatusInternalServerError,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"success": true,
			"message": "Welcome, here's your token. don't lose it ;)",
			"data": gin.H{
				"data": user.Nickname,
				"token": strToken,
			},
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"error": error,
			"success": false,
			"statusCode": http.StatusForbidden,
			"message": "password Anda salah.",
		})
		return
	}
}

func Authorization(c *gin.Context) {
	id, err := c.Get("id")
	
	var Auth Entities.Admin

	if !err {
		c.JSON(http.StatusUnauthorized, gin.H {
			"statusCode": http.StatusUnauthorized,
			"success": false,
			"message": "id is not exist",
			"error": err,
		})
		return
	}

	if err := Database.Open().Where("id = ?", id).Take(&Auth); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"statusCode": http.StatusInternalServerError,
			"success": false,
			"message": "error when querying user from database",
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"statusCode": http.StatusOK,
		"success": true,
		"message": "berhasil login",
		"error": nil,
	})
}

func PostNewOperational(c *gin.Context) {
	var input Entities.Operasional

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "input is invalid",
			"statusCode": http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}
	newOperasional := Entities.Operasional {
		HariOperasional: input.HariOperasional,
		JamOperasional: input.JamOperasional,
		Harga: input.Harga,
	}

	if err := Database.Open().Create(&newOperasional); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "can't create new operational",
			"statusCode": http.StatusInternalServerError,
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "a new operational has successfully created",
		"statusCode": http.StatusCreated,
		"data":    newOperasional.CreatedAt,
	})
}

func GetLatestOperational(c *gin.Context) {

	var operasional Entities.Operasional

	if err := Database.Open().Order("id desc").Take(&operasional); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "operational isn't available",
			"statusCode": http.StatusNotFound,
			"error": err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"message": "success querying latest operational",
		"statusCode": http.StatusOK,
		"data":    operasional,
	})
}

func PatchOperational(c *gin.Context) {
	var input Entities.Operasional

	if query := Database.Open().Order("id desc").Take(&input); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Error when querying the database.",
			"statusCode": http.StatusNotFound,
			"error":   query.Error.Error(),
		})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "input is invalid",
			"statusCode": http.StatusBadRequest,
			"error": err.Error(),
		})
	}

	patchOperasional := Entities.Operasional {
		HariOperasional: input.HariOperasional,
		JamOperasional: input.JamOperasional,
		Harga: input.Harga,
	}

	result := Database.Open().Where("id = ?", input.ID).Model(&patchOperasional).Updates(patchOperasional)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when updating new operational.",
			"statusCode": http.StatusInternalServerError,
			"error":   result.Error.Error(),
		})
		return
	}

	if result = Database.Open().Order("id desc").Take(&patchOperasional); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Error when querying the database.",
			"statusCode": http.StatusNotFound,
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "operational not found.",
			"statusCode": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Update successful.",
		"statusCode": http.StatusOK,
		"data":    patchOperasional,
	})
}

func DeleteOperational(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var operasional Entities.Operasional

	if err := Database.Open().Where("id = ?", id).Delete(&operasional); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when deleting from the database.",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"statusCode": http.StatusOK,
		"message": "Delete successful.",
	})
}

func PostNewNews(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "get form err: " + err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		return
	}
	imageIo, _ := image.Open()
	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)

	// p := bluemonday.NewPolicy()
	client.UploadFile("images", image.Filename, imageIo)

	enText := slug.MakeLang(c.PostForm("title"), "en")

	excerpt := HTML(c.PostForm("body"))
	if (len(excerpt) > 120) {
		excerpt = excerpt[:120]
	} 

	newArticle := Entities.Artikel {
		Title: c.PostForm("title"),
		Slug: enText,
		Image: os.Getenv("BASE_URL") + image.Filename,
		Excerpt: excerpt,
		Body: c.PostForm("body"),
	}

	if err := Database.Open().Create(&newArticle); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when inserting a new article",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":     true,
		"message":     "a new article has successfully created",
		"statusCode": http.StatusCreated,
		"error":       nil,
		"judul_article": newArticle.Title,
	})
}

func PatchNews(c *gin.Context) {
	search, _ := c.Params.Get("slug")

	var article Entities.Artikel

	if res := Database.Open().Where("slug = ?", search).Take(&article); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Something went wrong",
			"statusCode": http.StatusInternalServerError,
			"error":   res.Error.Error(),
		})
		return
	}
	// p := bluemonday.NewPolicy()

	excerpt := HTML(c.PostForm("body"))
	if (len(excerpt) > 120) {
		excerpt = excerpt[:120]
	} 

	var newArticle Entities.Artikel

	image, _ := c.FormFile("image")

	if image == nil {
		newArticle = Entities.Artikel{
			Title: c.PostForm("title"),
			Excerpt: excerpt,
			Body: c.PostForm("body"),
			Image: article.Image,
		}
	} else {
		imageIo, _ := image.Open()
		client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)
		client.DeleteBucket(article.Image)
		client.UploadFile("images", image.Filename, imageIo)

		excerpt := HTML(c.PostForm("body"))
		if (len(excerpt) > 120) {
			excerpt = excerpt[:120]
		} 

		newArticle = Entities.Artikel{
			Title: c.PostForm("title"),
			Excerpt: excerpt,
			Body: c.PostForm("body"),
			Image: os.Getenv("BASE_URL") + image.Filename,
		}
	}

	if err := Database.Open().Where("slug = ?", search).Model(&newArticle).Updates(newArticle); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when inserting a new agenda",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	if res := Database.Open().Where("slug = ?", search).Take(&article); res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Something went wrong",
			"statusCode": http.StatusNotFound,
			"error":   res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "a new article has successfully updated",
		"error":   nil,
		"statusCode": http.StatusOK,
		"data":    article,
	})
}

func GetAllNews(c *gin.Context) {
	var allArticle []Entities.Artikel

	limit, _ := c.GetQuery("limit")

	if limit == "true" {
		if res := Database.Open().Order("id desc").Limit(3).Find(&allArticle); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"success": false,
				"message": "failed when query all article",
				"statusCode": http.StatusInternalServerError,
				"error": res.Error.Error(),
			})
			return
		}
	} else {
		if res := Database.Open().Order("id desc").Find(&allArticle); res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"success": false,
				"message": "failed when query all article",
				"statusCode": http.StatusInternalServerError,
				"error": res.Error.Error(),
			})
			return
		}
	}

	
	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"statusCode": http.StatusOK,
		"error": nil,
		"data": allArticle,
	})
}

func SearchNews(c *gin.Context) {
	query, _ := c.GetQuery("q")

	var allArticle []Entities.Artikel

	if res := Database.Open().Where("title ILIKE ?", "%"+query+"%").Find(&allArticle); res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Hasil Pencarian Tidak Ditemukan",
			"statusCode": http.StatusNotFound,
			"error":   res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Search successful",
		"statusCode": http.StatusOK,
		"query": query,
		"data":    allArticle,
	})
}

func GetDetailNewsBySlug(c *gin.Context) {
	slug, _ := c.Params.Get("slug")

	var article Entities.Artikel

	if result := Database.Open().Where("slug = ?", slug).Take(&article); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when querying the database.",
			"statusCode": http.StatusInternalServerError,
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "query article successful.",
		"statusCode": http.StatusOK,
		"error":   nil,
		"data":    article,
	})
	// get all data except body article
}

func DeleteNewsBySlug(c *gin.Context) {
	slug, _ := c.Params.Get("slug")

	var article Entities.Artikel

	if res := Database.Open().Where("slug = ?", slug).Delete(&article); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when deleting from the database.",
			"statusCode": http.StatusInternalServerError,
			"error":   res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Delete successful.",
		"statusCode": http.StatusOK,
	})
}

func HTML(s string) (output string) {

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		// First remove line breaks etc as these have no meaning outside html tags (except pre)
		// this means pre sections will lose formatting... but will result in less unintentional paras.
		s = strings.Replace(s, "\n", "", -1)

		// Then replace line breaks with newlines, to preserve that formatting
		s = strings.Replace(s, "</p>", " ", -1)
		s = strings.Replace(s, "<br>", " ", -1)
		s = strings.Replace(s, "</br>", " ", -1)
		s = strings.Replace(s, "<br/>", " ", -1)
		s = strings.Replace(s, "<br />", " ", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}

	// Remove a few common harmless entities, to arrive at something more like plain text
	output = strings.Replace(output, "&#8216;", "'", -1)
	output = strings.Replace(output, "&#8217;", "'", -1)
	output = strings.Replace(output, "&#8220;", "\"", -1)
	output = strings.Replace(output, "&#8221;", "\"", -1)
	output = strings.Replace(output, "&nbsp;", " ", -1)
	output = strings.Replace(output, "&quot;", "\"", -1)
	output = strings.Replace(output, "&apos;", "'", -1)

	// Translate some entities into their plain text equivalent (for example accents, if encoded as entities)
	output = html.UnescapeString(output)

	// In case we have missed any tags above, escape the text - removes <, >, &, ' and ".
	output = template.HTMLEscapeString(output)

	// After processing, remove some harmless entities &, ' and " which are encoded by HTMLEscapeString
	output = strings.Replace(output, "&#34;", "\"", -1)
	output = strings.Replace(output, "&#39;", "'", -1)
	output = strings.Replace(output, "&amp; ", "& ", -1)     // NB space after
	output = strings.Replace(output, "&amp;amp; ", "& ", -1) // NB space after

	return output
}

func PostPicture(c *gin.Context) {
	image, err := c.FormFile("image")
	deskripsi := c.PostForm("deskripsi")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error":      "get form err: " + err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	imageIo, err := image.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error":      "get error when open the image: " + err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		return
	}
	client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)
	client.UploadFile("images", image.Filename, imageIo)

	newPhoto := Entities.Photo{
		LinkFoto:  os.Getenv("BASE_URL") + image.Filename,
		Deskripsi: deskripsi,
	}

	if err := Database.Open().Create(&newPhoto); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"message":    "error when inserting a new photo",
			"statusCode": http.StatusInternalServerError,
			"error":      err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"message":    "successfully upload file and description",
		"linkImage":  os.Getenv("BASE_URL") + image.Filename,
	})
}

func GetPictureByID(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var getPicture Entities.Photo

	if err := Database.Open().Where("id = ?", id).Take(&getPicture); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when querying the database.",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "query article successful.",
		"statusCode": http.StatusOK,
		"error":   nil,
		"data":    getPicture,
	})
}

func PatchPicture(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var getPicture Entities.Photo

	if err := Database.Open().Where("id = ?", id).Take(&getPicture); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when querying the database.",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	image, _ := c.FormFile("image")
	descripsi := c.PostForm("deskripsi")

	var patchPicture Entities.Photo

	if (image == nil) {
		patchPicture = Entities.Photo {
			LinkFoto: getPicture.LinkFoto,
			Deskripsi: descripsi,
		}
	} else {
		imageIo, _ := image.Open()
		client := storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SERVICE_TOKEN"), nil)
		client.DeleteBucket(getPicture.LinkFoto)
		client.UploadFile("images", image.Filename, imageIo)

		patchPicture = Entities.Photo{
			LinkFoto: os.Getenv("BASE_URL") + image.Filename,
			Deskripsi: descripsi,
		}
	}

	if err := Database.Open().Where("id = ?", id).Model(&patchPicture).Updates(patchPicture); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when inserting a new agenda",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	if res := Database.Open().Where("id = ?", id).Take(&getPicture); res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Something went wrong",
			"statusCode": http.StatusNotFound,
			"error":   res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "a picture has successfully updated",
		"error":   nil,
		"statusCode": http.StatusOK,
		"data":    getPicture,
	})
}

func DeletePicture(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var getPicture Entities.Photo

	if err := Database.Open().Where("id = ?", id).Delete(&getPicture); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error when deleting picture",
			"statusCode": http.StatusInternalServerError,
			"error":   err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"statusCode": http.StatusOK,
		"message": "delete successful",
	})
}

func GetAllPicture(c *gin.Context) {
	var allPicture []Entities.Photo

	limit, _ := c.GetQuery("limit")

		if limit == "true" {
			if res := Database.Open().Order("id desc").Limit(3).Find(&allPicture); res.Error != nil {
				c.JSON(http.StatusNotFound, gin.H {
					"success": false,
					"message": "failed when query all picture",
					"statusCode": http.StatusNotFound,
					"error": res.Error.Error(),
				})
				return
			}
		} else {
			if res := Database.Open().Order("id desc").Find(&allPicture); res.Error != nil {
				c.JSON(http.StatusNotFound, gin.H {
					"success": false,
					"message": "failed when query all picture",
					"statusCode": http.StatusNotFound,
					"error": res.Error.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H {
			"success": true,
			"statusCode": http.StatusOK,
			"error": nil,
			"data": allPicture,
		})
}

func PostInformation(c *gin.Context) {
	type linkArray struct {
		Description string `json:"description"`
		LinkImage []string `json:"linkImage"`
	}

	var linkInputDesc linkArray

	if err := c.BindJSON(&linkInputDesc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "input should bind json",
			"statusCode": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	combined := strings.Join(linkInputDesc.LinkImage, ";")

	newInfo := Entities.AdditionalInfo {
		Description: linkInputDesc.Description,
		LinkImages: combined,
	}

	if err := Database.Open().Create(&newInfo); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"message":    "error when inserting a new photo",
			"statusCode": http.StatusInternalServerError,
			"error":      err.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"success": true,
		"message": "a new additional info has successfully created",
		"statusCode": http.StatusCreated,
		"data": newInfo.ID,
	})
}

func GetAllInformation(c *gin.Context) {
	var getInfo Entities.AdditionalInfo

	if err := Database.Open().Order("id desc").Take(&getInfo); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "input should bind json",
			"statusCode": http.StatusBadRequest,
			"error": err.Error.Error(),
		})
		return
	}

	splitedLink := strings.Split(getInfo.LinkImages, ";")

	c.JSON(http.StatusOK, gin.H {
		"success": true,
		"statusCode": http.StatusOK,
		"description": getInfo.Description,
		"linkImage": splitedLink, 
	})
}

func PatchInformation(c *gin.Context) {
	var info Entities.AdditionalInfo

	if query := Database.Open().Order("id desc").Take(&info); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Error when querying the database.",
			"statusCode": http.StatusNotFound,
			"error":   query.Error.Error(),
		})
		return
	}

	type linkArray struct {
		Description string `json:"description"`
		LinkImage []string `json:"linkImage"`
	}

	var linkInputDesc linkArray

	if err := c.BindJSON(&linkInputDesc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"success": false,
			"message": "input should bind json",
			"statusCode": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	combined := strings.Join(linkInputDesc.LinkImage, ";")

	patchInfo := Entities.AdditionalInfo {
		Description: linkInputDesc.Description,
		LinkImages: combined,
	}

	result := Database.Open().Where("id = ?", info.ID).Model(&patchInfo).Updates(patchInfo)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when updating new operational.",
			"statusCode": http.StatusInternalServerError,
			"error":   result.Error.Error(),
		})
		return
	}

	if result = Database.Open().Order("id desc").Take(&patchInfo); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Error when querying the database.",
			"statusCode": http.StatusNotFound,
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "info not found.",
			"statusCode": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Update successful.",
		"statusCode": http.StatusOK,
		"data":    patchInfo,
	})
}