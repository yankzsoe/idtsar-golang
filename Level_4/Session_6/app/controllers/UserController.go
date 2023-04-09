package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserController struct{}

// GetUser godoc
// @Summary      Get User Data
// @Description  Get Dummy User Account
// @Tags         User
// @Accept       json
// @Produce      json
// @Param id path int false "ID"
// @Router       /user/{id} [get]
func (ctrl *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, gin.H{"user_id": id})
}

// PostUser godoc
// @Summary      Post User Data
// @Description  Add new fake User Account
// @Tags         User
// @Accept       json
// @Produce      json
// @Param request body CreateUserRequest true "User"
// @Router       /user [post]
func (ctrl *UserController) PostUser(ctx *gin.Context) {
	req := CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully"})
}
