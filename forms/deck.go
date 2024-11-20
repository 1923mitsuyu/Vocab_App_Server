package forms

// import (
//     "encoding/json"
//     "github.com/go-playground/validator/v10"
// )

// DeckForm handles the forms for deck-related operations.
type DeckForm struct{}

// GetDeckForm is used to fetch a deck by its name.
type GetDeckForm struct {
    DeckName string `json:"name" binding:"required"`
    DeckOrder int   `json:"deckOrder" binding:"required"`
    UserID   int    `json:"userId" binding:"required"`
}

// AddDeckForm is used to add a new deck.
type AddDeckForm struct {
    DeckName  string   `json:"name" binding:"required"`
    DeckOrder int      `json:"deckOrder" binding:"required"`
    UserID    int      `json:"userId" binding:"required"`
}

// UpdateDeckForm is used to update an existing deck.
type UpdateDeckForm struct {
    DeckID    int      `json:"deckId" binding:"required"`
    DeckName  string   `json:"deckName" binding:"required"`
    DeckOrder int      `json:"deckOrder"`
    UserID    int      `json:"userId" binding:"required"`
}
