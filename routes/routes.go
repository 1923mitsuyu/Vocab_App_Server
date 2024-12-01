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
		ai := new(controllers.AIController)

		// User-related endpoints
		v1.POST("/login", user.Login)
		v1.POST("/signup", user.Signup)
		// v1.GET("/user/logout", user.Logout)

		// Deck-related endpoints
		v1.GET("/fetchDecks/:userId", deck.FetchDecks)
		v1.POST("/saveDeck", deck.AddDeck)
		v1.PUT("/modifyDeck", deck.EditDeck)
		v1.PUT("/modifyDeckOrders", deck.EditDeckOrder)
		v1.DELETE("/removeDeck", deck.DeleteDeck)
		
		// Word-related endppints
		v1.GET("/fetchWords/:deckId", word.FetchWords)
		v1.POST("/saveWord", word.AddWord)
		v1.PUT("/modifyWord", word.EditWord)
		v1.PUT("/modifyCorrectCount", word.EditCorrectCount)
		v1.PUT("/modifyWordOrders", word.EditWordOrder)
		v1.DELETE("/removeWord", word.DeleteWord)

		// AI related endpoints
		v1.POST("/generate", ai.GenerateContent)
	}

	return router
}
