package handlers

import (
	"fmt"
	"net/http"
	"usermanagement/internal/models"
	"usermanagement/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	router      *gin.Engine
	userService services.UserServiceInterface
}

func NewServer(userService services.UserServiceInterface) *Server {
	return &Server{
		router:      gin.Default(),
		userService: userService,
	}
}

// login database
func (s *Server) LoginDB() {
	s.userService.LoginDB()
}

// SetupRoute sets up routes on the server
func (s *Server) SetupRoute() {
	s.router.GET("/users", s.handleGetAllUsers)
	s.router.POST("/register", s.handleRegister)
	s.router.POST("/login", s.handleLogin)
	s.router.GET("/search", s.handleSearchUser)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Run: run server
func (s *Server) RunServer() {
	s.router.Run()
}

// ----- APIs start -----

// handleRegister handles the user registration process for the POST /register API endpoint.
// It expects a JSON payload containing a username and password.
func (s *Server) handleRegister(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// hash password
	password := data.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	data.Password = string(hashedPassword)

	user := models.NewUser(data.Username, data.Password)
	if err := s.userService.CreateUser(*user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// handleLogin handles the user authentication process for the POST /login API endpoint.
// It expects a JSON payload containing a username and password.
func (s *Server) handleLogin(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	foundUser, err := s.userService.SearchUserByUsername(userInput.Username)
	if err != nil {
		// not found user in the database
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user doesn't exists",
		})
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userInput.Password))
	if err != nil {
		// wrong password
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
	})
}

// handleGetAllUsers handles the GET /users API endpoint.
// It retrieves all users from the userService and responds with a 200 OK status
// and a JSON array of all user details.
func (s *Server) handleGetAllUsers(c *gin.Context) {
	users, err := s.userService.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

// handleSearchUser handles the GET /search API endpoint with query parameters for username or id.
// The function searches for a user based on the provided username or id.
// - If a username is provided, it attempts to find the user by username.
// - If an id is provided, it attempts to find the user by id.
// If neither a username nor an id is provided, it responds with a 200 OK status and an empty JSON object.
func (s *Server) handleSearchUser(c *gin.Context) {
	username := c.Query("username")
	id := c.Query("id")

	if username != "" {
		// search by username
		foundUser, err := s.userService.SearchUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, foundUser)

	} else if id != "" {
		// search by id
		foundUser, err := s.userService.SearchUserByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// ----- APIs end -----
