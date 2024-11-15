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
		// ユーザー関連のエンドポイント
		v1.POST("/login", user.Login)
		v1.POST("/signup", user.Signup)
		// v1.GET("/user/logout", user.Logout)
	}

	return router
}
