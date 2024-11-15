package controllers

import (
	"fmt"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/gin-gonic/gin"
)

//UserController ...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)

//getUserID ...
func getUserID(c *gin.Context) (userID int64) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	return c.MustGet("userID").(int64)
}

// 1. Authenticate users with usernames and passwords
func (ctrl UserController) Login(c *gin.Context) {

	var loginForm forms.LoginForm

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	user, err := models.GetUserByUsername(loginForm.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
	if err != nil {
		// パスワードが一致しない場合
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// JWTトークンの生成
	token, err := models.GenerateJWT(user.ID)
	if err != nil {
		// JWTの生成に失敗した場合
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// 成功した場合、トークンを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": user,
		"token":   token,
	})
}

// 2. Register a new user with username and password 
func (ctrl UserController) Signup(c *gin.Context) {
	var signupForm forms.SignupForm

	// リクエストのJSONデータをSignupForm構造体にバインド
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		fmt.Println("Error binding JSON: %v", err)  
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// バリデーション
	if signupForm.Username == "" || signupForm.Password == "" {
		fmt.Println("Missing username or password: %v", signupForm)  // エラー内容をログに出力
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username, password and email are required"})
		return
	}

	// ユーザー名が既に存在するか確認
	existingUser, err := models.GetUserByUsername(signupForm.Username)
	if err != nil && err.Error() != "user not found" {
        // 他のエラー（DBエラー等）の場合
        fmt.Printf("Error checking if user exists: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
        return
    }

	if existingUser != nil && existingUser.ID != 0 {
		fmt.Println("Username %s is already taken", signupForm.Username)  // 既に存在するユーザー名をログに出力
		c.JSON(http.StatusConflict, gin.H{"message": "Username is already taken"})
		return
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password: %v", err)  // パスワードハッシュ化エラーをログに出力
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password"})
		return
	}

	// ユーザーの保存
	user := models.User{
		Username: signupForm.Username,
		Password: string(hashedPassword), 
	}

	// 新しいユーザーをDBに保存
	err = models.CreateUser(user)
	if err != nil {
		fmt.Println("Error creating user in DB: %v", err)  // ユーザー作成エラーをログに出力
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	// サインアップが成功した場合、JWTトークンを生成
	token, err := models.GenerateJWT(user.ID)
	if err != nil {
		fmt.Println("Error generating JWT: %v", err)  // JWT生成エラーをログに出力
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// 成功した場合、JWTトークンを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "Signup successful",
		"user": user,
		"token":   token,
	})
}


//Logout ...
// func (ctrl UserController) Logout(c *gin.Context) {

// 	au, err := authModel.ExtractTokenMetadata(c.Request)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
// 		return
// 	}

// 	deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
// 	if delErr != nil || deleted == 0 { //if any goes wrong
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
// }
