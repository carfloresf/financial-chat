package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/financial-chat/config"
	"github.com/carfloresf/financial-chat/internal/constants"
	"github.com/carfloresf/financial-chat/internal/controller/middleware"
	"github.com/carfloresf/financial-chat/internal/queue"
	"github.com/carfloresf/financial-chat/internal/websocket"
)

type User struct {
	service userService
}

type userService interface {
	Authenticate(username, password string) error
	Register(username, password string) error
}

func NewRouter(
	conf *config.Config,
	user userService,
	queueClient queue.Queuer,
	websocketServer *melody.Melody) (*gin.Engine, error) {
	router := gin.Default()
	userHandler := User{service: user}

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Use(sessions.Sessions("session", cookie.NewStore([]byte(conf.Auth.Cookie))))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	// Public routes
	public := router.Group("/")
	public.GET("/login", LoginHandler())
	public.POST("/login", userHandler.LoginPostHandler())
	public.POST("/v1/user", userHandler.RegisterHandler())

	// Private routes
	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	private.GET("/", IndexHandler())
	private.GET("/logout", LogoutHandler())
	private.GET("/channel/:name", ChannelHandler())
	// Websocket
	private.GET("/channel/:name/ws", websocket.Handler(websocketServer, queueClient))

	return router, nil
}

func LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get(constants.Userkey)
		if user != nil {
			ctx.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Please logout first",
					"user":    user,
				})

			return
		}

		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get(constants.Userkey)

		if user == nil {
			log.Println("Invalid session token")

			return
		}

		session.Delete(constants.Userkey)

		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)

			return
		}

		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func IndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(constants.Userkey)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"user": user,
		})
	}
}

func ChannelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(constants.Userkey)

		c.HTML(http.StatusOK, "channel.html", gin.H{
			"user": user,
		})
	}
}
