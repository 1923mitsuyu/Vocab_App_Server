package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	// "github.com/Massad/gin-boilerplate/controllers"
	"github.com/Massad/gin-boilerplate/db"
	// uuid "github.com/google/uuid"
	"github.com/joho/godotenv"
 	"github.com/Massad/gin-boilerplate/routes"
	"github.com/gin-gonic/gin" // Ginフレームワークのインポート
)


// var auth = new(controllers.AuthController)

// アプリケーションのエントリーポイント
func main() {

	// ルーターの設定 (ルートとミドルウェアを設定)
	router := routes.SetupRouter() 

	// ポート番号を環境変数から取得。設定されていなければデフォルトで8080を使用
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" 
	}

	// .envファイルを読み込む
	err := godotenv.Load(".env")
	if err != nil {
		// .envファイルが読み込めなかった場合、エラーログを出力して終了
		log.Fatal("error: failed to load the env file")
	}

	// 環境変数が"PRODUCTION"なら、Ginフレームワークを本番モードに設定
	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// DBの初期化処理
	db.Init()

	// HTMLテンプレートファイルをロード
	router.LoadHTMLGlob("./public/html/*")
	// 静的ファイル（例えばCSSや画像）を"/public"パスで提供
	router.Static("/public", "./public")

	// ルートURL（"/"）のGETリクエストを処理
	router.GET("/", func(c *gin.Context) {
		// レスポンスとしてHTMLを返す
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03", // Ginのバージョンを表示
			"goVersion":             runtime.Version(), // // Goのバージョンを表示
		})
	})

	// どのルートにもマッチしない場合、404ページを表示
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	// SSL/TLS設定: 環境変数"SSL"が"TRUE"の場合、HTTPSでサーバーを起動
	if os.Getenv("SSL") == "TRUE" {

		// SSL証明書と秘密鍵のパス
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer", // SSL証明書
			KEY:  "./cert/myCA.key", // SSL秘密鍵
		}

		// HTTPSサーバーをポート番号と証明書ファイルを指定して起動
		router.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		// SSL設定がない場合、通常のHTTPサーバーを起動
		router.Run(":" + port)
	}

}
