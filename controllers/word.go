package controllers

import (
	"fmt"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	//"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"  
)

//UserController ...
type WordController struct{}

var wordModel = new(models.WordModel)
var wordForm = new(forms.WordForm)

// 1. Fetch all words in the db
func (ctrl WordController) FetchWords(c *gin.Context) {
    deckIdParam := c.Param("deckId") 
    var deckId *int

    fmt.Println("Step1")
    if deckIdParam != "" {
        id, err := strconv.Atoi(deckIdParam)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid deckId parameter"})
            return
        }
        deckId = &id
    }

    fmt.Println("Step2")
    // デッキを取得
    words, err := (&models.WordModel{}).GetWords(deckId) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch words", "error": err.Error()})
        return
    }

    fmt.Println("Step7")
    // 成功した場合、ワードリストを返す
    c.JSON(http.StatusOK, gin.H{"message": "Words fetched successfully", "words": words}) 
}

// ２. Save a new word to the db
func (ctrl WordController) AddWord(c *gin.Context) {

	var wordForm forms.AddWordForm

	fmt.Println("Adding process start....")

	fmt.Println("Step1")
	// リクエストのJSONデータをSignupForm構造体にバインド
	if err := c.ShouldBindJSON(&wordForm); err != nil {
		fmt.Println("Error binding JSON: %v", err)  
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	fmt.Println("Step2")
	// バリデーション
	if wordForm.WordName == "" {
		fmt.Println("Missing word name: %v", wordForm)  // エラー内容をログに出力
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck name is required"})
		return
	}

	// Wordの保存
	word := models.Word{
		DeckID: wordForm.DeckID,
		WordName: wordForm.WordName, 
		WordDefinition: wordForm.WordDefinition,
		WordExample: wordForm.WordExample,
		WordTranslation: wordForm.WordTranslation,
		CorrectTimes: wordForm.CorrectTimes,
		WordOrder: wordForm.WordOrder,  
	}

	fmt.Println("Step5")
	// 新しいDeckをDBに保存
	err := models.CreateWord(word)
	if err != nil {
		fmt.Println("Error creating word in DB: %v", err)  // ユーザー作成エラーをログに出力
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add word"})
		return
	}

	// 成功した場合、JWTトークンを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "Adding the word to the db successfully",
		"word": word,
	})
}

