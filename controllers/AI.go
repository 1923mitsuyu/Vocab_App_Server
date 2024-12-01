package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIController ...
type AIController struct{}

// GenerateContent は Gin のハンドラーとして修正されたバージョン
func (ctrl AIController) GenerateContent(c *gin.Context) {
	var aiForm forms.GenerateContentForm

	// リクエストのJSONデータを AIPromptForm 構造体にバインド
	if err := c.ShouldBindJSON(&aiForm); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// バリデーション
	if aiForm.TargetWord == "" {
		log.Printf("Missing prompt: %+v", aiForm)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prompt is required"})
		return
	}

	// Gemini APIクライアントを作成する
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_APIKEY")))
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Gemini client: %v", err)})
		return
	}
	defer client.Close()

	// 質問を送信して回答を取得する
	model := client.GenerativeModel("gemini-pro")
	prompt := genai.Text(fmt.Sprintf("Please create an example sentence in English that includes the word '%s', based on the following definition: '%s'. Please make sure to include that target word inside the brackets {{}}. ", aiForm.TargetWord, aiForm.WordDefinition))
	resp, err := model.GenerateContent(context.Background(), prompt)
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call Gemini API: %v", err)})
		return
	}

	// 候補データを抽出
	var generatedContent []string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			generatedContent = append(generatedContent, fmt.Sprintf("%v", part))
		}
	}

	// デバッグ用ログ出力
	fmt.Printf("Generated Content: %v\n", generatedContent)

	// レスポンスを返却
	c.JSON(http.StatusOK, gin.H{
		"message": "Content generated successfully",
		"content": generatedContent,
	})
}

// GenerateContent は Gin のハンドラーとして修正されたバージョン
func (ctrl AIController) GenerateTranslation(c *gin.Context) {
	var aiForm forms.GenerateTranslationForm

	fmt.Println("Step1")
	// リクエストのJSONデータを AIPromptForm 構造体にバインド
	if err := c.ShouldBindJSON(&aiForm); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	fmt.Println("Step2")
	// バリデーション
	if aiForm.TargetSentence == "" {
		log.Printf("Missing prompt: %+v", aiForm)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prompt is required"})
		return
	}

	fmt.Println("Step3")
	// Gemini APIクライアントを作成する
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_APIKEY")))
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Gemini client: %v", err)})
		return
	}
	defer client.Close()

	// 質問を送信して回答を取得する
	model := client.GenerativeModel("gemini-pro")
	prompt := genai.Text(fmt.Sprintf("Please translate '%s' into natural Japanese. Do not have to use the brackets for the translation.", aiForm.TargetSentence))
	resp, err := model.GenerateContent(context.Background(), prompt)
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call Gemini API: %v", err)})
		return
	}

	// 候補データを抽出
	var generatedContent []string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			generatedContent = append(generatedContent, fmt.Sprintf("%v", part))
		}
	}

	// デバッグ用ログ出力
	fmt.Printf("Generated Content: %v\n", generatedContent)

	// レスポンスを返却
	c.JSON(http.StatusOK, gin.H{
		"message": "Content generated successfully",
		"content": generatedContent,
	})
}

