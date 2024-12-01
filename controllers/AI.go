package controllers

import (
	"bufio"
	"fmt"
	//"context"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

//UserController ...
type AIController struct{}

// GenerateContent は Gin のハンドラーとして修正されたバージョン
func (ctrl AIController) GenerateContent(c *gin.Context) {
	// コンソールから質問を入力する
	fmt.Print("質問を入力してください: ")
	reader := bufio.NewReader(os.Stdin)
	question, err := reader.ReadString('\n')
	if err != nil {
		// エラーが発生した場合、レスポンスとしてエラーを返す
		c.JSON(400, gin.H{"error": fmt.Sprintf("質問の読み取りに失敗しました: %v", err)})
		return
	}
	
	// Gemini APIクライアントを作成する
	client, err := genai.NewClient(c, option.WithAPIKey(os.Getenv("GEMINI_APIKEY")))
	if err != nil {
		// クライアントの作成に失敗した場合、レスポンスとしてエラーを返す
		c.JSON(500, gin.H{"error": fmt.Sprintf("Geminiクライアントの作成に失敗しました: %v", err)})
		return
	}
	defer client.Close()

	// 質問を送信して回答を取得する
	model := client.GenerativeModel("gemini-pro")
	prompt := genai.Text(question)
	resp, err := model.GenerateContent(c, prompt)
	if err != nil {
		// API呼び出しに失敗した場合、レスポンスとしてエラーを返す
		c.JSON(500, gin.H{"error": fmt.Sprintf("Gemini APIの呼び出しに失敗しました: %v", err)})
		return
	}

	// 回答を表示する
	printCandidates(resp.Candidates)

	// レスポンスとして成功を返す（必要に応じて）
	c.JSON(200, gin.H{"message": "Content generated successfully"})
}


func printCandidates(cs []*genai.Candidate) {
	for _, c := range cs {
		for _, p := range c.Content.Parts {
			fmt.Println(p)
		}
	}
}