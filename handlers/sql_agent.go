package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Astrasv/go-gully-backend/database"
	"github.com/Astrasv/go-gully-backend/models"
	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Response string `json:"response"`
}

func QueryAgent(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	baseURL := os.Getenv("SQL_AGENT_URL")
	if baseURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SQL_AGENT_URL not configured"})
		return
	}

	agentURL := baseURL + "/agents/query?query=" + url.QueryEscape(req.Query)
	resp, err := http.Get(agentURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to SQL agent"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "failed to parse response",
			"raw_response": string(body),
		})
		return
	}

	query := models.Query{
		UserID:   userID.(uint),
		Query:    req.Query,
		Response: result.Response,
	}
	database.GetDB().Create(&query)

	c.JSON(resp.StatusCode, result)
}

func GetQueryHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var queries []models.Query
	database.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&queries)

	c.JSON(http.StatusOK, gin.H{"queries": queries})
}

func GetEntities(c *gin.Context) {
	
	// Comment this block to test without auth
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	baseURL := os.Getenv("SQL_AGENT_URL")
	if baseURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SQL_AGENT_URL not configured"})
		return
	}

	agentURL := baseURL + "/entities"
	resp, err := http.Get(agentURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to SQL agent"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}
