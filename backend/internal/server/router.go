package server

import (
	"release-calendar/backend/internal/server/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type"},
	}))

	api := r.Group("/api")
	{
		api.GET("/release-days", func(c *gin.Context) { handlers.NewGetReleaseDays(db).Handle(c) })
		releases := api.Group("/releases")
		{
			releases.GET("", func(c *gin.Context) { handlers.NewListRelease(db).Handle(c) })
			releases.POST("", func(c *gin.Context) { handlers.NewAddRelease(db).Handle(c) })
			releases.GET(":id", func(c *gin.Context) { handlers.NewGetRelease(db).Handle(c) })
			releases.PUT(":id", func(c *gin.Context) { handlers.NewUpdateRelease(db).Handle(c) })
			releases.DELETE(":id", func(c *gin.Context) { handlers.NewDeleteRelease(db).Handle(c) })

			releases.GET(":id/comments", func(c *gin.Context) { handlers.NewGetComments(db).Handle(c) })
			releases.POST(":id/comments", func(c *gin.Context) { handlers.NewAddComment(db).Handle(c) })
			releases.PUT(":id/comments/:commentId", func(c *gin.Context) { handlers.NewUpdateComment(db).Handle(c) })
			releases.DELETE(":id/comments/:commentId", func(c *gin.Context) { handlers.NewDeleteComment(db).Handle(c) })

			releases.GET(":id/summary", func(c *gin.Context) { handlers.NewGetSummary(db).Handle(c) })
		}
	}

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })
	return r
}
