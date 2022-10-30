package Controller

import (
	"backend-d-embung/Auth"
	"backend-d-embung/Handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OperasionalController(r *gin.Engine) {
	// post a new operational
	r.POST("/operasional", Auth.Authorization(), Handlers.PostNewOperational)

	// get operational
	r.GET("/operasional", Handlers.GetLatestOperational)

	// patch operational time
	r.PATCH("/operasional", Auth.Authorization(), Handlers.PatchOperational)

	// delete operational time by id
	r.DELETE("/operasional/:id", Auth.Authorization(), Handlers.DeleteOperational)
}

func NewsController(r *gin.Engine) {
	// post new article
	r.POST("/article", Auth.Authorization(), Handlers.PostNewNews)

	// get all article
	r.GET("/articles", Handlers.GetAllNews)

	// search article by query title
	r.GET("/article", Handlers.SearchNews)

	// -> 3 api below using slug params
	// get detail article 
	r.GET("/article/:slug", Handlers.GetDetailNewsBySlug)

	// patch article 
	r.PATCH("/article/:slug", Auth.Authorization(), Handlers.PatchNews)

	// delete article
	r.DELETE("/article/:slug", Auth.Authorization(), Handlers.DeleteNewsBySlug)
}

func TestimoniController(r *gin.Engine) {
	r.POST("/testimoni", Auth.Authorization(), Handlers.PostNewTestimoni)

	r.GET("/testimoni", Handlers.GetAllTestimoni)
}

func Authorization(db *gorm.DB, r *gin.Engine) {
	r.POST("/register", Handlers.RegisterAdmin)

	r.POST("/login", Handlers.LoginAdmin)

	r.POST("/authToken", Auth.Authorization(), Handlers.Authorization)
}

func Post(r *gin.Engine) {
	r.POST("/picture", Auth.Authorization(), Handlers.PostPicture)

	r.GET("/picture/:id", Handlers.GetPictureByID)

	r.GET("/pictures", Handlers.GetAllPicture)

	r.PATCH("/picture/:id", Auth.Authorization(), Handlers.PatchPicture)

	r.DELETE("/picture/:id", Auth.Authorization(), Handlers.DeletePicture)
}

func AdditionalInfo(r *gin.Engine) {
	r.POST("/info", Auth.Authorization(), Handlers.PostInformation)

	r.GET("/info", Handlers.GetAllInformation)

	r.PATCH("/info", Auth.Authorization(), Handlers.PatchInformation)
}

