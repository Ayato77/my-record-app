package controllers

import (
	"my-record-app/database"
	"my-record-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateRecord(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.Record
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Sugar().Errorf("Create Record: BindJSON failed!", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			logger.Sugar().Errorf("User ID %s does not exist", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		body.UserID = userID.(uint)

		recordId, err := database.CreateRecord(body)
		if err != nil {
			logger.Sugar().Errorf("error: Failed to create record", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}

		logger.Info("A new record created successfully")
		c.JSON(http.StatusOK, gin.H{"id": recordId})
	}
}

func GetRecords(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("GetRecords!")
		userID, exists := c.Get("userID")
		if !exists {
			logger.Sugar().Errorf("User ID %s does not exist", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		//TODO: Implement search with keyword in titles or contents

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		sort := c.DefaultQuery("sort", "created_at desc")
		tags := c.QueryArray("tag")

		if page < 1 {
			page = 1
		}
		//Limit max records to be returned
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		records, total, err := database.GetWithPaginationDB(userID.(int), page, limit, offset, tags, sort)
		if err != nil {
			logger.Sugar().Errorf("Getting records failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Getting records failed"})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": records,
			"meta": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})

	}
}

func DeleteRecord(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("DeleteRecord")

		//get user_id
		userID, exists := c.Get("userID")
		if !exists {
			logger.Sugar().Errorf("User ID %s does not exist", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		id, _ := strconv.Atoi(c.DefaultQuery("id", ""))
		rowsAffected, err := database.DeleteSingleRecord(userID.(int), id)

		if err != nil {
			logger.Sugar().Errorf("Deleting a record failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deleting a record failed"})
			return
		}

		//No rows are affected => requested id does not exist
		if rowsAffected == 0 {
			logger.Sugar().Error("Record to be deleted does not exist")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record's ID does not exist"})
			return
		}

		c.JSON(http.StatusNoContent, nil)

	}
}

func UpdateRecord(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("UpdateRecord")
	}
}
