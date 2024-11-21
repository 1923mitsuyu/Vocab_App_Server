package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Massad/gin-boilerplate/controllers"
	"github.com/Massad/gin-boilerplate/middleware"
)

// ルーターの設定を行う関数
func SetupRouter() *gin.Engine {

	// Ginインスタンスを初期化
	router := gin.Default()

	// ミドルウェアの追加
	router.Use(middleware.CORSMiddleware()) // CORS設定ミドルウェアを追加
	router.Use(middleware.RequestIDMiddleware()) // リクエストIDを設定するミドルウェアを追加

	// "/v1"で始まるルートグループを作成
	v1 := router.Group("/v1")
	{
		// コントローラーのインスタンスを作成
		user := new(controllers.UserController)
		deck := new(controllers.DeckController)
		word := new(controllers.WordController)

		// User-related endpoints
		v1.POST("/login", user.Login)
		v1.POST("/signup", user.Signup)
		// v1.GET("/user/logout", user.Logout)

		// Deck-related endpoints
		v1.GET("/fetchDecks/:userId", deck.FetchDecks)
		v1.POST("/saveDeck", deck.AddDeck)
		
		// Word-related endppints
		v1.GET("/fetchWords/:deckId", word.FetchWords)
		v1.POST("/saveWord", word.AddWord)
	}

	return router
}