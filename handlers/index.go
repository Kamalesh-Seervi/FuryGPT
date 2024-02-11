package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func index(c *gin.Context) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(c.Writer, nil)
}

func run(c *gin.Context) {
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
	response, err := getPromptFromCache(prompt.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response != "" {
		// If prompt is found in the cache, return the cached response
		c.JSON(http.StatusOK, gin.H{"input": prompt.Input, "response": response})
		return
	}

	// Create the LLM
	llm, err := openai.NewChat(openai.WithModel(viper.GetString("OPENAI_MODEL")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatmsg := []schema.ChatMessage{
		schema.SystemChatMessage{Content: "Hello, I am a friendly AI assistant."},
		schema.HumanChatMessage{Content: prompt.Input},
	}
	aimsg, err := llm.Call(context.Background(), chatmsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the prompt and response in the cache
	err = setPromptInCache(prompt.Input, aimsg.GetContent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"input": prompt.Input, "response": aimsg.GetContent()})
}
