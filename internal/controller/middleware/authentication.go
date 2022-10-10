package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AuthRequired(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if user := session.Get("user"); user == nil {
		log.Errorf("unauthorized access to %s", ctx.Request.URL.Path)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		ctx.Abort()

		return
	}

	ctx.Next()
}
