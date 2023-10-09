package handlers

import (
	"context"
	"fmt"
	"golang-api/models"
	"net/http"
	"os"

	"golang-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Client *mongo.Client
}

var dbName string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env")
	}
	dbName = os.Getenv("DB_NAME")
}

// ch for client handler
func (ch *Handler) RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	collection := ch.Client.Database(dbName).Collection("users")
	existingUser := models.User{}

	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Email already used",
		})
		return
	}
	// else if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "failed",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"user":    user,
	})
}

func (ch *Handler) LoginUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	collection := ch.Client.Database(dbName).Collection("users")
	existingUser := models.User{}

	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		fmt.Println("Stored Password:", existingUser.Password)
		fmt.Println("Provided Password:", user.Password)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Invalid password",
		})
		return
	}

	token, err := utils.CreateToken(existingUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to create token",
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"token":   token,
		"message": "Login Successfully",
		"user":    existingUser,
	})
	return

}

func (ch *Handler) LogoutUser(c *gin.Context) {
	tokenStirng := string(c.GetHeader("Authorization")[7:])
	utils.RevokeToken(tokenStirng)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status":  "failed",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user logout successfully",
	})
	return
}

func (ch *Handler) GetAllUsers(c *gin.Context) {

	// claims, ok := c.Get("claims")

	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status":  "failed",
	// 		"message": "no claims found",
	// 	})
	// 	return
	// }

	// jwtClaims, ok := claims.(jwt.MapClaims)

	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status":  "faield",
	// 		"message": "invalid claim format",
	// 	})
	// 	return
	// }

	// username, exist := jwtClaims["username"].(string)

	// if !exist {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status":  "failed",
	// 		"message": "username not found in claims",
	// 	})
	// 	return
	// }

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "username not found claims",
		})
		return
	}

	collection := ch.Client.Database(dbName).Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "No users were found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Users found successfully",
		"login_as": username,
		"users":    users,
	})
}
