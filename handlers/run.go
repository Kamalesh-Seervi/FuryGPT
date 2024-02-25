package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/kamalesh-seervi/simpleGPT/service"
	"google.golang.org/api/option"
)

func Run(c *gin.Context) {

	prompt := struct {
		Input  string `json:"input"`
		APIKey string `json:"apiKey"`
	}{}

	err := c.BindJSON(&prompt)
	if err != nil {
		log.Printf("Error processing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// adding apiKey in Gemini
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(prompt.APIKey))
	if err != nil {
		log.Printf("Error processing request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	//  gemini promodel
	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt.Input))
	if err != nil {
		log.Printf("Error generating content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	formattedContent := formatResponse(resp)
	// c.JSON(http.StatusOK, gin.H{"input": prompt.Input, "response": formattedContent})

	// I am storing both the input prompt and the generated response in Redis
	service.StoreHistory(prompt.APIKey, prompt.Input, formattedContent)
	history := service.GetHistory(prompt.APIKey)
	c.JSON(http.StatusOK, gin.H{
		"input":    prompt.Input,
		"response": formattedContent,
		"history":  history,
	})
}
func HandleFetchHistory(c *gin.Context) {
	prompt := struct {
		Input  string `json:"input"`
		APIKey string `json:"apiKey"`
	}{}
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Fetch history from Redis
	history := service.GetHistory(prompt.APIKey)

	// Convert history strings to structured JSON objects
	var historyObjects []map[string]string
	for _, historyStr := range history {
		var historyObject map[string]string
		err := json.Unmarshal([]byte(historyStr), &historyObject)
		if err != nil {
			continue
		}
		historyObjects = append(historyObjects, historyObject)
	}

	// Send the structured JSON objects to the frontend
	c.JSON(http.StatusOK, gin.H{
		"history": historyObjects,
	})
}

// format resposne
func formatResponse(resp *genai.GenerateContentResponse) string {
	var formattedContent strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formattedContent.String()
}
