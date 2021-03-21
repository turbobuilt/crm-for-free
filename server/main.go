package main

import (
	"fmt"
	// "time"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Jesus loves you forever")

	GetDB()
	router := setupRouter()
	router.Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000","http://"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return origin == "https://github.com"
	// 	// },
	// 	MaxAge: 12 * time.Hour,
	// }))

	r.Static("/api/static", "./static")

	r.LoadHTMLGlob("templates/*")
	api := r.Group("/api/v1.0")
	{
		api.POST("/user", UserSignup)
		api.POST("/user/ConfirmAccount", ConfirmAccount)
		api.GET("/user/CreatePassword", CreatePassword)
		api.POST("/user/Login", UserLogin)
		api.POST("/user/SendResetPasswordEmail", SendResetPasswordEmail)
		api.GET("/user/ResetPassword", ResetPassword)
		api.POST("/user/ConfirmResetPassword", ConfirmResetPassword)

		authed := api.Group("/", AuthenticateUser)
		{
			authed.POST("/customer", CreateCustomer)
			authed.PUT("/customer/:id", UpdateCustomer)
			authed.GET("/customer", GetCustomersList)
			authed.GET("/customer/:id", GetCustomer)
		}
	}
	return r
}
