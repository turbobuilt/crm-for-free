package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TenantId string
	Email    string

	Name         string
	PasswordHash string
}
type DBUser struct {
	TenantId string
	Email    string
	JSON     string
}

func UserLogin(c *gin.Context) {
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, Error{Message: "Invalid submission"})
		return
	}
	if body["Email"] == "" || body["Password"] == "" {
		c.JSON(400, Error{Message: "Email and password are required"})
		return
	}
	email := body["Email"]
	password := body["Password"]

	userKey, err := db.Get(ctx, "email:"+email).Result()
	if err == redis.Nil {
		c.JSON(400, Error{Message: "Email is not yet registered"})
		return
	}
	user, _ := db.HGetAll(ctx, userKey).Result()
	userData := gin.H{}
	json.Unmarshal([]byte(user["JSON"]), &userData)

	if err := bcrypt.CompareHashAndPassword([]byte(userData["PasswordHash"].(string)), []byte(password)); err != nil {
		c.JSON(400, Error{Message: "Invalid Password"})
		return
	}
	expiry := GetLoginExpirationMillisecondsString()
	hashStr := JWT_SECRET + "_" + userKey + "_" + expiry
	hash := sha256.Sum256([]byte(hashStr))
	c.JSON(200, gin.H{
		"UserKey":         userKey,
		"LoginExpiration": expiry,
		"Verification":    base64.RawStdEncoding.EncodeToString(hash[:]), //   hex.EncodeToString(hash[:]),
	})
}

func UserSignup(c *gin.Context) {
	var body map[string]string

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, Error{Message: "Invalid submission"})
		return
	}
	if body["Email"] == "" {
		c.JSON(400, Error{Message: "Email is required"})
		return
	}
	email := strings.ToLower(body["Email"])
	emailKey := "email:" + email
	newUserKey := "new_user:" + email
	if _, err := db.Get(ctx, emailKey).Result(); err == nil {
		c.JSON(400, Error{Message: "Email is already in use"})
		return
	}
	token := RandStringBytes(32)
	encoded, _ := json.Marshal(gin.H{
		"token":  token,
		"tenant": "",
	})
	db.Set(ctx, newUserKey, string(encoded), 0)
	SendEmail(FROM_EMAIL, []string{email}, "verify-account.html", "Verify Email", gin.H{
		"email": email,
		"token": token,
	})
}

func CreatePassword(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")
	value, _ := db.Get(ctx, "new_user:"+email).Result()

	var data map[string]interface{}
	_ = json.Unmarshal([]byte(value), &data)

	fmt.Println("token is ", token, "data is ", data["token"])
	if token != data["token"] {
		c.JSON(400, gin.H{"Message": "Invalid token for that email."})
		return
	}

	c.HTML(http.StatusOK, "create-password.html", gin.H{
		"email": email,
		"token": token,
	})
}

func ConfirmAccount(c *gin.Context) {
	email := c.Request.FormValue("email")
	token := c.Request.FormValue("token")
	password := c.Request.FormValue("password")

	if email == "" || token == "" || password == "" {
		c.JSON(400, "email, token and password are required")
		return
	}

	fmt.Println("the email is ", email, "The token is ", token, "password is ", password)
	// return

	value, _ := db.Get(ctx, "new_user:"+email).Result()

	var data map[string]interface{}
	_ = json.Unmarshal([]byte(value), &data)

	if token != data["token"] {
		c.JSON(400, gin.H{"Message": "Invalid token for that email."})
		return
	}

	tenantId := data["tenant"].(string)
	if tenantId == "" {
		tenantIdInt, _ := db.Incr(ctx, "TenantId").Result()
		tenantId = strconv.FormatInt(tenantIdInt, 10)

		sdb := redisearch.NewClient("localhost:6379", "cust:"+tenantId+":Index")
		// Create redisearch index
		schema := redisearch.NewSchema(redisearch.DefaultOptions).
			AddField(redisearch.NewTextField("NotesSearch")).
			AddField(redisearch.NewTextFieldOptions("Name", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))

		indexDefinition := redisearch.NewIndexDefinition().AddPrefix("cust:" + tenantId)
		sdb.CreateIndexWithIndexDefinition(schema, indexDefinition)
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	id, _ := db.Incr(ctx, "UserId:"+tenantId).Result()
	user := User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	key := "user:" + tenantId + ":" + strconv.FormatInt(id, 10)
	emailKey := "email:" + strings.ToLower(user.Email)

	// check if user is already registered
	if _, err := db.Get(ctx, emailKey).Result(); err == nil {
		c.JSON(400, Error{Message: "This user has already registered"})
		return
	}

	// register user
	db.HSet(ctx, key, user.SerializeUserForDatabase())
	db.Set(ctx, emailKey, key, 0)

	c.Redirect(302, "/#/login")
	// c.JSON(200, gin.H{"id": id, "tenantId": tenantId})
}

func SendResetPasswordEmail(c *gin.Context) {
	body := make(map[string]string)
	c.ShouldBindJSON(&body)
	email := body["Email"]
	if email == "" {
		c.JSON(400, gin.H{"Message": "Email is required"})
		return
	}

	if _, err := db.Get(ctx, "email:"+email).Result(); err == redis.Nil {
		c.JSON(400, gin.H{"Message": "That email is not registered."})
		return
	}

	resetToken := RandStringBytes(32)
	db.Set(ctx, "reset_password_token:"+email, resetToken, time.Minute*30)

	SendEmail(FROM_EMAIL, []string{email}, "reset-password-email.html", "Reset Password", gin.H{
		"email": email,
		"token": resetToken,
	})
}

func ResetPassword(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")
	correctToken, _ := db.Get(ctx, "reset_password_token:"+email).Result()

	if token != correctToken {
		c.JSON(400, gin.H{"Message": "Invalid token for that email."})
		return
	}

	c.HTML(http.StatusOK, "reset-password.html", gin.H{
		"email": email,
		"token": token,
	})
}
func ConfirmResetPassword(c *gin.Context) {
	email := c.Request.FormValue("email")
	token := c.Request.FormValue("token")
	password := c.Request.FormValue("password")

	if email == "" || token == "" || password == "" || len(password) < 8 {
		c.JSON(400, gin.H{"Message": "email, token and password are required, and password must be at least 8 characters"})
		return
	}

	correctToken, _ := db.Get(ctx, "reset_password_token:"+email).Result()
	if token != correctToken {
		c.JSON(400, gin.H{"Message": "Invalid token for that email."})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := GetUserByEmail(email)
	user["PasswordHash"] = string(hash)
	SaveUser(user)

	c.JSON(200, "great work")
}
func GetUserById(userKey string) map[string]interface{} {
	if exists, _ := db.Exists(ctx, userKey).Result(); exists == 0 {
		return nil
	}
	userRaw, err := db.HGetAll(ctx, userKey).Result()
	if err == redis.Nil {
		return nil
	}
	jsonData := gin.H{}
	json.Unmarshal([]byte(userRaw["JSON"]), &jsonData)
	jsonData["Email"] = userRaw["Email"]
	return jsonData
}
func GetUserByEmail(email string) map[string]interface{} {
	userKey, _ := db.Get(ctx, "email:"+email).Result()
	if userKey == "" {
		return nil
	}
	return GetUserById(userKey)
}
func SaveUser(user map[string]interface{}) {
	j, _ := json.Marshal(user)
	userRaw := map[string]interface{}{
		"Email": user["Email"],
		"JSON":  string(j),
	}
	userKey, _ := db.Get(ctx, "email:"+user["Email"].(string)).Result()
	db.HSet(ctx, userKey, userRaw).Result()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (user *User) SerializeUserForDatabase() map[string]interface{} {
	j, _ := json.Marshal(map[string]string{
		"Name":         user.Name,
		"PasswordHash": user.PasswordHash,
	})
	data := map[string]interface{}{
		"Email": user.Email,
		"JSON":  string(j),
	}
	return data
}

func AuthenticateUser(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.AbortWithStatusJSON(401, gin.H{"Message": "You are not authenticated.  You need a JWT."})
		return
	}
	parts := strings.Split(authToken, "_")
	hashString := JWT_SECRET + "_" + parts[0] + "_" + parts[1]
	hash := sha256.Sum256([]byte(hashString))
	hashB64 := base64.RawStdEncoding.EncodeToString(hash[:])
	if hashB64 != parts[2] {
		c.AbortWithStatusJSON(401, gin.H{"Message": "You are not authenticated. You have an invalid JWT"})
		return
	}
	// expiration
	expirationTimestamp, _ := strconv.ParseInt(parts[1], 10, 64)
	if time.Now().UnixNano()/1000/1000 > expirationTimestamp {
		c.AbortWithStatusJSON(401, gin.H{"Message": "Your login has expired"})
		return
	}
	user := GetUserById(parts[0])
	if user == nil {
		c.AbortWithStatusJSON(401, gin.H{"Message": "A user was not found with your userId.  Please contact support."})
		return
	}
	userParts := strings.Split(parts[0], ":")
	c.Set("tenantId", userParts[1])
	c.Set("user", user)
}
