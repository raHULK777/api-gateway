package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-gateway/internal/repository"
	"api-gateway/internal/service"
)

func Register(keyStore repository.KeyStore, tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
			Name     string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user exists
		existingUser, err := keyStore.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		// Hash password
		hashedPassword, err := tokenService.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		// Create user
		user := &repository.User{
			Email:    req.Email,
			Password: hashedPassword,
			Name:     req.Name,
		}

		if err := keyStore.CreateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		// Generate JWT token
		token, err := tokenService.GenerateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"name":  user.Name,
			},
			"token": token,
		})
	}
}

func Login(keyStore repository.KeyStore, tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := keyStore.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if err := tokenService.CheckPassword(user.Password, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := tokenService.GenerateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"name":  user.Name,
			},
			"token": token,
		})
	}
}

func GetMe(keyStore repository.KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		user, err := keyStore.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		})
	}
}

func GetUser(keyStore repository.KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		user, err := keyStore.GetUserByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		})
	}
}

func ListItems(keyStore repository.KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		items, err := keyStore.GetUserItems(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch items"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}

func CreateItem(keyStore repository.KeyStore, tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item := &repository.Item{
			Name:        req.Name,
			Description: req.Description,
			UserID:      userID,
		}

		if err := keyStore.CreateItem(item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func GetItem(keyStore repository.KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		item, err := keyStore.GetItem(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func UpdateItem(keyStore repository.KeyStore, tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		item, err := keyStore.GetItem(uint(id))
		if err != nil || item.UserID != userID {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found or unauthorized"})
			return
		}

		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Name != "" {
			item.Name = req.Name
		}
		if req.Description != "" {
			item.Description = req.Description
		}

		if err := keyStore.UpdateItem(item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func DeleteItem(keyStore repository.KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		item, err := keyStore.GetItem(uint(id))
		if err != nil || item.UserID != userID {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found or unauthorized"})
			return
		}

		if err := keyStore.DeleteItem(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
	}
}
