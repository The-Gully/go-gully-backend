package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

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

	c.JSON(resp.StatusCode, result)
}
