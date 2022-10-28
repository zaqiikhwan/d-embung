package Handlers

import (
	"backend-d-embung/Database"
	"backend-d-embung/Entities"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	storage_go "github.com/supabase-community/storage-go"
)

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