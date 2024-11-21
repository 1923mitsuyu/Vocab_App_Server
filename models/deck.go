package models

import (
	"fmt"
	//"encoding/json"
	"github.com/Massad/gin-boilerplate/db"
	"database/sql" 
	//"os"
	//"time"
)

// Deck モデル
type Deck struct {
    ID       uint     `json:"id"`
    DeckName string   `json:"name"`
    DeckOrder int     `json:"deckOrder"`
    UserID    int     `json:"userId"`
}

type DeckModel struct{}

// Deck名でDeckを検索
func GetDeckByName(deckname string, user_id int) (*Deck, error) {
    var deck Deck
    // QueryRow is used to select a single row from the database
	err := db.GetDB().QueryRow(
		"SELECT id, name FROM Vocab_App.decks WHERE name = ? AND userId = ?",
		deckname, user_id,
	).Scan(&deck.ID, &deck.DeckName)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("deck not found")
        }
        return nil, err
    }
    return &deck, nil
}

// UserIdに当てはまる全てのDeckを取得
func (d *DeckModel) GetDecks(userID *int) ([]Deck, error) {
	var decks []Deck
	var rows *sql.Rows
	var err error

    fmt.Println("Step3")
	// If userID is provided, filter by userID, else get all decks
	if userID != nil {
		rows, err = db.GetDB().Query("SELECT id, name, deckOrder, userId FROM decks WHERE userid = ?", *userID)
	} else {
		rows, err = db.GetDB().Query("SELECT id, name, deckOrder, userId FROM decks")
	}

    fmt.Println("Step4")
	if err != nil {
		fmt.Printf("Error querying decks: %v", err)
		return nil, err
	}
	defer rows.Close()

    fmt.Println("Step5")
	// Iterate over rows and append to decks slice
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(&deck.ID, &deck.DeckName, &deck.DeckOrder, &deck.UserID); err != nil {
			fmt.Printf("Error scanning row: %v", err)
			return nil, err
		}
		decks = append(decks, deck)
	}

    fmt.Println("Step6")
	// Check for any error that occurred during the row iteration
	if err := rows.Err(); err != nil {
		fmt.Printf("Error after iterating rows: %v", err)
		return nil, err
	}

	return decks, nil
}


// 新しいDeckをDBに保存
func CreateDeck(deck Deck) error {
    fmt.Println("Step6")
	_, err := db.GetDB().Exec(
		"INSERT INTO Vocab_App.decks (name, deckOrder, userId) VALUES (?, ?, ?)", 
		deck.DeckName, deck.DeckOrder, deck.UserID)

    fmt.Println("Step7")
	if err != nil {
		return fmt.Errorf("could not create deck: %v", err)
	}
	return nil
}
