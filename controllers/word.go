package controllers

import (
	"strings"
	"fmt"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/Massad/gin-boilerplate/db"
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
	// ワードリストが空の場合、特別なメッセージを返す
    if len(words) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No words found for this deck", "words": words})
        return
    }
	
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

// 3. Edit an existing word 
func (ctrl WordController) EditWord(c *gin.Context) {
    var wordForm forms.EditWordForm

	fmt.Printf("Step1")
    // JSON データをバインド
    if err := c.ShouldBindJSON(&wordForm); err != nil {
        fmt.Printf("Error binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
        return
    }

	fmt.Printf("Step2")
    // 動的にクエリを構築
    query := "UPDATE words SET "
    updates := []string{}
    params := []interface{}{}

    if wordForm.WordName != "" {
        updates = append(updates, "word = ?")
        params = append(params, wordForm.WordName)
    }
    if wordForm.Definition != "" {
        updates = append(updates, "definition = ?")
        params = append(params, wordForm.Definition)
    }
    if wordForm.Example != "" {
        updates = append(updates, "example = ?")
        params = append(params, wordForm.Example)
    }
    if wordForm.Translation != "" {
        updates = append(updates, "translation = ?")
        params = append(params, wordForm.Translation)
    }

	fmt.Printf("Step3")
    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "No fields to update"})
        return
    }

	fmt.Printf("Step4")
    query += strings.Join(updates, ", ") + " WHERE id = ?"
    params = append(params, wordForm.WordID)

	fmt.Printf("Step5")
    // クエリの実行
    stmt, err := db.GetDB().Prepare(query)
    if err != nil {
        fmt.Printf("Error preparing query: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare update query"})
        return
    }
    defer stmt.Close()

	fmt.Printf("Step6")
    result, err := stmt.Exec(params...)
    if err != nil {
        fmt.Printf("Error updating word: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update word"})
        return
    }

	fmt.Printf("Step7")
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Printf("Error retrieving rows affected: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve update result"})
        return
    }

	fmt.Printf("Step8")
    if rowsAffected == 0 {
        fmt.Printf("No word found with id: %v\n", wordForm.WordID)
        c.JSON(http.StatusNotFound, gin.H{"message": "Word not found"})
        return
    }

	fmt.Printf("Step9!")
    fmt.Println("Word updated successfully")
    c.JSON(http.StatusOK, gin.H{"message": "Word updated successfully"})
}

