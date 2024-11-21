package forms

// import (
//     "encoding/json"
//     "github.com/go-playground/validator/v10"
// )

// WordForm handles the forms for deck-related operations.
type WordForm struct{}

// GetDeckForm is used to fetch a deck by its name.
type GetWordForm struct {
	DeckID   int    `json:"deckId" binding:"required"`
    WordName string `json:"word" binding:"required"`
	WordDefinition string `json:"definition" binding:"required"`
	WordExample string `json:"example" binding:"required"`		
	WordTranslation string `json:"translation" binding:"required"`
    WordOrder int   `json:"word_order" binding:"required"`
	CorrectTimes int `json:"correctTimes" binding:"required"`
}

// AddDeckForm is used to add a new deck.
type AddWordForm struct {
	DeckID   int    `json:"deckId" binding:"required"`
    WordName string `json:"word" binding:"required"`
	WordDefinition string `json:"definition" binding:"required"`
	WordExample string `json:"example" binding:"required"`		
	WordTranslation string `json:"translation" binding:"required"`
	WordOrder int   `json:"word_order" binding:"required"`
	CorrectTimes int `json:"correctTimes"`
}

// UpdateDeckForm is used to update an existing deck.
type UpdateWordForm struct {
    WordID    int      `json:"deckId" binding:"required"`
    WordName  string   `json:"deckName" binding:"required"`
    WordOrder int      `json:"deckOrder"`
    DeckID    int      `json:"userId" binding:"required"`
}
