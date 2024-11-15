package models

import (
	"fmt"
	"github.com/Massad/gin-boilerplate/db"
	"database/sql" 
	"github.com/dgrijalva/jwt-go" 
	"os"
	"time"
)

// User モデル
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct{}

// ユーザー名でユーザーを検索する関数
func GetUserByUsername(username string) (*User, error) {
    var user User
    // QueryRow is used to select a single row from the database
    err := db.GetDB().QueryRow("SELECT id, username, password FROM Vocab_App.users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    return &user, nil
}

// 新しいユーザーをDBに保存
func CreateUser(user User) error {
	// ユーザー情報をDBに挿入
	_, err := db.GetDB().Exec("INSERT INTO Vocab_App.users (username, password) VALUES (?, ?)",
		user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}

// GenerateJWT トークンを生成する関数
func GenerateJWT(userID uint) (string, error) {
	// JWTの署名に使う秘密鍵
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// JWTトークンのクレーム設定
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // トークンの有効期限を1日後に設定
	}

	// トークンの作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return signedToken, nil
}

// //One ...
// func (m UserModel) One(userID int64) (user User, err error) {
// 	err = db.GetDB().SelectOne(&user, "SELECT id, email, name FROM public.user WHERE id=$1 LIMIT 1", userID)
// 	return user, err
// }
