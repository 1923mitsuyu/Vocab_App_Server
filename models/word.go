package models

import (
	"fmt"
	//"encoding/json"
	"github.com/Massad/gin-boilerplate/db"
	"database/sql" 
	//"os"
	//"time"
)

// Word モデル
type Word struct {
    ID       uint     `json:"id"`
	DeckID   int 	  `json:"deckId"`
    WordName string   `json:"word"`
	WordDefinition string `json:"definition"`
	WordExample string 	  `json:"example"`
	WordTranslation string `json:"translation"`
	CorrectTimes int `json:"correctTimes"`
    WordOrder int     `json:"word_order"`
}

type WordModel struct{}

// DeckIdに当てはまる全てのWordを取得
func (d *WordModel) GetWords(deckID *int) ([]Word, error) {
	var words []Word
	var rows *sql.Rows
	var err error

    fmt.Println("Step3")
	// If deckID is provided, filter by deckID, else get all decks
	if deckID != nil {
		rows, err = db.GetDB().Query("SELECT id, deck_id, word, definition, example, translation, word_order, correctTimes FROM words WHERE deck_id = ?", *deckID)
	} else {
		rows, err = db.GetDB().Query("SELECT id, deck_id, word, definition, example, translation, word_order, correctTimes FROM words")
	}

    fmt.Println("Step4")
	if err != nil {
		fmt.Printf("Error querying words: %v", err)
		return nil, err
	}
	defer rows.Close()

    fmt.Println("Step5")
	// Iterate over rows and append to words slice
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.DeckID, &word.WordName, &word.WordDefinition, &word.WordExample, &word.WordTranslation, &word.WordOrder, &word.CorrectTimes); err != nil {
			fmt.Printf("Error scanning row: %v", err)
			return nil, err
		}
		words = append(words, word)
	}

    fmt.Println("Step6")
	// Check for any error that occurred during the row iteration
	if err := rows.Err(); err != nil {
		fmt.Printf("Error after iterating rows: %v", err)
		return nil, err
	}

	return words, nil
}


// 新しいWordをDBに保存
func CreateWord(word Word) error {
    fmt.Println("Step6")
	_, err := db.GetDB().Exec(
		"INSERT INTO Vocab_App.words (deck_id, word, definition, example, translation, word_order, correctTimes) VALUES (?, ?, ?, ?, ?, ?, ?)", 
		word.DeckID, 
		word.WordName, 
		word.WordDefinition, 
		word.WordExample, 
		word.WordTranslation, 
		word.WordOrder,
		word.CorrectTimes,
		) 

    fmt.Println("Step7")
	if err != nil {
		return fmt.Errorf("could not create word: %v", err)
	}
	return nil
}
