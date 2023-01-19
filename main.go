package main

import (
	"log"
	"mnc-portal/artikel"
	"mnc-portal/auth"
	"mnc-portal/handler"
	"mnc-portal/helper"
	"mnc-portal/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/buddyku?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	artikelRepository := artikel.NewRepository(db)

	userService := user.NewService(userRepository)
	artikelService := artikel.NewService(artikelRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	artikelHandler := handler.NewArtikelHandler(artikelService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	//user
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	//artikel
	api.GET("/artikels", artikelHandler.GetArtikels)
	api.GET("/artikel/:id", artikelHandler.GetArtikel)
	api.POST("/artikels", authMiddleware(authService, userService), artikelHandler.CreateArtikel)
	api.POST("/artikels/:id", authMiddleware(authService, userService), artikelHandler.UpdateArtikel)
	api.POST("/artikel-images", authMiddleware(authService, userService), artikelHandler.UploadImage)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHaeder := c.GetHeader("Authorization")

		if !strings.Contains(authHaeder, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHaeder, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorizen", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorizen", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorizen", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
