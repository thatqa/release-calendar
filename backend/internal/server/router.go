package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"release-calendar/backend/internal/models"
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
		releases := api.Group("/releases")
		{
			releases.GET("", func(c *gin.Context) { listReleases(c, db) })
			releases.POST("", func(c *gin.Context) { createRelease(c, db) })
			releases.GET(":id", func(c *gin.Context) { getRelease(c, db) })
			releases.PUT(":id", func(c *gin.Context) { updateRelease(c, db) })
			releases.DELETE(":id", func(c *gin.Context) { deleteRelease(c, db) })

			// comments
			releases.GET(":id/comments", func(c *gin.Context) { listComments(c, db) })
			releases.POST(":id/comments", func(c *gin.Context) { addComment(c, db) })
			releases.PUT(":id/comments/:commentId", func(c *gin.Context) { updateComment(c, db) })
			releases.DELETE(":id/comments/:commentId", func(c *gin.Context) { deleteComment(c, db) })
		}
	}

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })
	return r
}

type ReleaseDTO struct {
	Title     string               `json:"title" binding:"required"`
	Date      time.Time            `json:"date" binding:"required"`
	Status    models.ReleaseStatus `json:"status" binding:"required"`
	Notes     string               `json:"notes"`
	DutyUsers []string             `json:"dutyUsers"`
	Links     []LinkDTO            `json:"links"`
}

type LinkDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

func listReleases(c *gin.Context, db *gorm.DB) {
	var releases []models.Release
	q := db.Preload("Links").Order("date ASC")

	if day := c.Query("date"); day != "" {
		if t, err := time.Parse("2006-01-02", day); err == nil {
			start := t
			end := t.Add(24 * time.Hour)
			q = q.Where("date >= ? AND date < ?", start, end)
		}
	}

	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Find(&releases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if duty := c.Query("duty"); duty != "" {
		filtered := make([]models.Release, 0, len(releases))
		for _, r := range releases {
			var du []string
			_ = json.Unmarshal(r.DutyUsers, &du)
			for _, u := range du {
				if u == duty {
					filtered = append(filtered, r)
					break
				}
			}
		}
		releases = filtered
	}

	c.JSON(http.StatusOK, releases)
}

func createRelease(c *gin.Context, db *gorm.DB) {
	var req ReleaseDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	du, _ := json.Marshal(req.DutyUsers)
	rel := models.Release{
		Title:     req.Title,
		Date:      req.Date,
		Status:    req.Status,
		Notes:     req.Notes,
		DutyUsers: du,
	}

	if err := db.Create(&rel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, l := range req.Links {
		lnk := models.Link{ReleaseID: rel.ID, Name: l.Name, URL: l.URL}
		db.Create(&lnk)
	}

	getReleaseByID(c, db, rel.ID)
}

func getRelease(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	getReleaseByID(c, db, uint(id))
}

func getReleaseByID(c *gin.Context, db *gorm.DB, id uint) {
	var rel models.Release
	if err := db.Preload("Links").Preload("Comments").First(&rel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var du []string
	_ = json.Unmarshal(rel.DutyUsers, &du)
	c.JSON(http.StatusOK, gin.H{
		"id":        rel.ID,
		"title":     rel.Title,
		"date":      rel.Date,
		"status":    rel.Status,
		"notes":     rel.Notes,
		"dutyUsers": du,
		"links":     rel.Links,
		"comments":  rel.Comments,
		"createdAt": rel.CreatedAt,
		"updatedAt": rel.UpdatedAt,
	})
}

func updateRelease(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	var rel models.Release
	if err := db.Preload("Links").First(&rel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req ReleaseDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	du, _ := json.Marshal(req.DutyUsers)
	rel.Title = req.Title
	rel.Date = req.Date
	rel.Status = req.Status
	rel.Notes = req.Notes
	rel.DutyUsers = du

	if err := db.Save(&rel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	existing := map[uint]models.Link{}
	for _, l := range rel.Links {
		existing[l.ID] = l
	}
	seen := map[uint]bool{}

	for _, pl := range req.Links {
		if pl.ID == 0 {
			lnk := models.Link{ReleaseID: rel.ID, Name: pl.Name, URL: pl.URL}
			db.Create(&lnk)
		} else {
			if l, ok := existing[pl.ID]; ok {
				l.Name = pl.Name
				l.URL = pl.URL
				db.Save(&l)
				seen[pl.ID] = true
			}
		}
	}
	for id, l := range existing {
		if !seen[id] {
			db.Delete(&l)
		}
	}

	getReleaseByID(c, db, rel.ID)
}

func deleteRelease(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Release{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Comments ---
type CommentDTO struct {
	Author  string `json:"author" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func listComments(c *gin.Context, db *gorm.DB) {
	releaseID, _ := strconv.Atoi(c.Param("id"))
	var items []models.Comment
	if err := db.Where("release_id = ?", releaseID).Order("created_at ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func addComment(c *gin.Context, db *gorm.DB) {
	releaseID, _ := strconv.Atoi(c.Param("id"))
	var req CommentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cm := models.Comment{ReleaseID: uint(releaseID), Author: req.Author, Message: req.Message}
	if err := db.Create(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cm)
}

func updateComment(c *gin.Context, db *gorm.DB) {
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	var req CommentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var cm models.Comment
	if err := db.First(&cm, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	cm.Author = req.Author
	cm.Message = req.Message
	if err := db.Save(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}

func deleteComment(c *gin.Context, db *gorm.DB) {
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	if err := db.Delete(&models.Comment{}, commentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
