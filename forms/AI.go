package forms

type GenerateContentForm struct {
    TargetWord    string     `json:"word" binding:"required"`
	WordDefinition string  `json:"definition" binding:"required"`
}

type GenerateTranslationForm struct {
    TargetSentence    string     `json:"example" binding:"required"`
}

