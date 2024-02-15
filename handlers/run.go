package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
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
