package controllers

import (
	"fmt"
	"net/http"
	"new_projects/initializers"
	"new_projects/models"
	"new_projects/utilities"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var userRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

func CreateUser(c *gin.Context) {
	// bind the request body
	err := c.Bind(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind the request"})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 14)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to hash the password"})
	}

	// create user
	user := models.User{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  string(hash),
		Role:      models.UserRole,
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create user"})
		return
	}

	// respond it
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUser(c *gin.Context) {
	var pgParam utilities.PaginationParam
	err := c.BindQuery(&pgParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind the pagination param"})
	}

	paginParam, err := utilities.ExtractPagination(pgParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to extruct pagination"})
	}
	// get the user data with search and filter
	var users []models.User
	db := initializers.DB

	if pgParam.Search != "" {
		initializers.DB.Where("first_name LIKE %%?%% OR lastName LIKE %%?%%", pgParam.Search, pgParam.Search)
	} else if paginParam.Filter != nil {
		for _, filter := range paginParam.Filter {
			db.Where(fmt.Sprintf("%s %s %v", filter.ComlumnName, filter.Operator, filter.Value))
		}
	}

	offset := (paginParam.Page - 1) * paginParam.PerPage

	result := db.Offset(offset).Limit(paginParam.PerPage).Order(paginParam.Sort.ComlumnName + " " + paginParam.Sort.Value).Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the user"})
		return
	}

	// respond it
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUserByID(c *gin.Context) {
	// get the url id
	id := c.Param("id")

	// get the user by the id
	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the user by the id"})
		return
	}

	// respond it
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	// get the url id
	id := c.Param("id")

	// bind the request body
	err := c.Bind(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind the request"})
		return
	}

	// get the data by the pk id
	var user models.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the data by the pk id"})
		return
	}

	// update the data
	result = initializers.DB.Model(&user).Updates(models.User{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  userRequest.Password,
		Role:      models.Role(userRequest.Role),
	})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to update the user data"})
		return
	}

	// respond it
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func DeleteUser(c *gin.Context) {
	// get the url id
	id := c.Param("id")

	// delete the user data by the pk
	initializers.DB.Delete(&models.User{}, id)

	// respond it
	c.JSON(http.StatusBadRequest, gin.H{})
}

func Login(c *gin.Context) {
	// bind the email and password
	var userRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.Bind(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind the email and password"})
		return
	}

	// lookup the request user
	var user models.User
	result := initializers.DB.First(&user, "email = ?", userRequest.Email)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to compare the email"})
		return
	}

	// compare the input passwor with the hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to compareb the password"})
		return
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// send and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to send the token"})
		return
	}

	// send it back
	c.JSON(http.StatusOK, gin.H{"Token": tokenString})
}
