package controllers

import (
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
type DeckController struct{}

var deckModel = new(models.DeckModel)
var deckForm = new(forms.DeckForm)

// 1. Fetch all decks in the db
func (ctrl DeckController) FetchDecks(c *gin.Context) {
	userId := c.Param("userId")
	var userID *int

	fmt.Println("Step1")
	if userId != "" {
		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid userId parameter"})
			return
		}
		userID = &id
	}
	fmt.Println("Step2")
	// デッキを取得
	decks, err := (&models.DeckModel{}).GetDecks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch decks", "error": err.Error()})
		return
	}
	fmt.Println("Step7")
	// 成功した場合、デッキリストを返す
	c.JSON(http.StatusOK, gin.H{"message": "Decks fetched successfully", "decks": decks})
}

// 2. Add a new deck to the db
func (ctrl DeckController) AddDeck(c *gin.Context) {

	var deckForm forms.AddDeckForm

	fmt.Println("Adding process start....")

	fmt.Println("Step1")
	// リクエストのJSONデータをSignupForm構造体にバインド
	if err := c.ShouldBindJSON(&deckForm); err != nil {
		fmt.Println("Error binding JSON: %v", err)  
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	fmt.Println("Step2")
	// バリデーション
	if deckForm.DeckName == "" {
		fmt.Println("Missing deck name: %v", deckForm)  // エラー内容をログに出力
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck name is required"})
		return
	}

	fmt.Println("Step3")
	// Deck名が既に存在するか確認
	existingDeck, err := models.GetDeckByName(deckForm.DeckName, deckForm.UserID)
	if err != nil && err.Error() != "deck not found" {
        // 他のエラー（DBエラー等）の場合
        fmt.Printf("Error checking if deck exists: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
        return
    }

	fmt.Println("Step4")
	if existingDeck != nil && existingDeck.DeckID != 0 {
		fmt.Println("Deck %s is already taken", deckForm.DeckName)  // 既に存在するユーザー名をログに出力
		c.JSON(http.StatusConflict, gin.H{"message": "Deck is already taken"})
		return
	}

	// Deckの保存
	deck := models.Deck{
		DeckName: deckForm.DeckName,  
		DeckOrder: deckForm.DeckOrder,  
		UserID: deckForm.UserID,  
	}

	fmt.Println("Step5")
	// 新しいDeckをDBに保存
	err = models.CreateDeck(deck)
	if err != nil {
		fmt.Println("Error creating deck in DB: %v", err)  // ユーザー作成エラーをログに出力
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create deck"})
		return
	}

	// 成功した場合、JWTトークンを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "Adding the deck to the db successfully",
		"deck": deck,
	})
}

// 2. Edit an existing deck 
func (ctrl DeckController) EditDeck(c *gin.Context) {
    var deckForm forms.EditDeckForm

    fmt.Println("Adding process start....")

    // JSON データをバインド
	fmt.Println("Step1")
    if err := c.ShouldBindJSON(&deckForm); err != nil {
        fmt.Printf("Error binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
        return
    }

    // バリデーション
	fmt.Println("Step2")
    if deckForm.DeckName == "" {
        fmt.Printf("Missing deck name: %+v\n", deckForm)
        c.JSON(http.StatusBadRequest, gin.H{"message": "Deck name is required"})
        return
    }

    // データベースクエリの実行
    query := "UPDATE decks SET name = ? WHERE id = ?"

	fmt.Println("Step3")
    stmt, err := db.GetDB().Prepare(query)
    if err != nil {
        fmt.Printf("Error preparing query: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare update query", "error": err.Error()})
        return
    }
    defer stmt.Close()

	fmt.Println("Step4")
    result, err := stmt.Exec(deckForm.DeckName, deckForm.DeckID)
    if err != nil {
        fmt.Printf("Error updating deck: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update deck", "error": err.Error()})
        return
    }

	fmt.Println("Step5")
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Printf("Error retrieving rows affected: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve update result", "error": err.Error()})
        return
    }

	fmt.Println("Step6")
    if rowsAffected == 0 {
        fmt.Printf("No deck found with id: %v\n", deckForm.DeckID)
        c.JSON(http.StatusNotFound, gin.H{"message": "Deck not found"})
        return
    }

	fmt.Println("Step7")
    fmt.Println("Deck updated successfully")
    c.JSON(http.StatusOK, gin.H{"message": "Deck updated successfully"})
}


