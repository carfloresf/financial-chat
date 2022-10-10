package controller

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/financial-chat/internal/constants"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) RegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input RegisterInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		err := user.service.Register(input.Username, input.Password)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "registration success"})
	}
}

func (user *User) LoginPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userSession := session.Get(constants.Userkey)
		if userSession != nil {
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first"})

			return
		}

		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Parameters can't be empty"})

			return
		}

		err := user.service.Authenticate(username, password)
		if err != nil {
			log.Errorf("error authenticating user: %v", err)
			ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password"})

			return
		}

		session.Set(constants.Userkey, username)

		if err := session.Save(); err != nil {
			log.Errorf("failed to save session: %v", err)
			ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})

			return
		}

		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
