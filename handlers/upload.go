package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhanavadh/sorkorsor-backend/storage"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage *storage.R2Storage
}

func NewUploadHandler(storage *storage.R2Storage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)

	url, err := h.storage.UploadFile(c.Request.Context(), file, fileName, contentType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "upload successfully",
		"url":     url,
	})
}
