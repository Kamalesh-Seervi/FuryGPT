package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/kamalesh-seervi/simpleGPT/service"
	"github.com/kamalesh-seervi/simpleGPT/utils"
	"google.golang.org/api/option"
)

var apiKey string

func init() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	apiKey = config.APIKey
}

func Run(c *gin.Context) {

	// user, exists := c.Get("user")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	// 	return
	// }
	// username := user.(string)
	// apiKey := c.GetHeader("X-API-Key") // Use a custom header for API key
	// if apiKey == "" {
	// 	// If API key is not in headers, try to get it from the request body
	// 	var apiKeyFromBody struct {
	// 		APIKey string `json:"apiKey"`
	// 	}

	// 	if err := c.ShouldBindJSON(&apiKeyFromBody); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
	// 		return
	// 	}

	// 	apiKey = apiKeyFromBody.APIKey
	// }

	// // Use the provided API key or the user's API key
	// if apiKey == "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is missing"})
	// 	return
	// }

	prompt := struct {
		Input string `json:"input"`
	}{}

	// Decode JSON from the client
	err := c.BindJSON(&prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the prompt is already in the cache
	response, err := service.GetPromptFromCache(prompt.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response != "" {
		// If prompt is found in the cache, return the cached response
		c.JSON(http.StatusOK, gin.H{"input": prompt.Input, "response": response})
		return
	}

	// Create the genai client using Viper to get the API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
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

	// Save the prompt and response in the cache
	err = service.SetPromptInCache(prompt.Input, formattedContent)
	if err != nil {
		log.Printf("Error generating cache: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"input": prompt.Input, "response": formattedContent})
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
