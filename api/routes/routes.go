package routes

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/api/services"
	"github.com/mrinjamul/gnote/middleware"
)

// ViewsFs for static files
var ViewsFs embed.FS

var (
	StartTime time.Time
	BootTime  time.Duration
)

func InitRoutes(routes *gin.Engine) {
	// Initialize services
	svc := services.NewServices()

	// Serve the frontend
	// This will ensure that the web pages are served correctly

	// routes.NoRoute(func(c *gin.Context) {
	// 	dir, file := path.Split(c.Request.RequestURI)
	// 	ext := filepath.Ext(file)
	// 	if file == "" || ext == "" {
	// 		c.File("./views/index.html")
	// 	} else {
	// 		c.File("./views/" + path.Join(dir, file))
	// 	}
	// })

	fsRoot, err := fs.Sub(ViewsFs, "views")
	if err != nil {
		log.Println(err)
	}
	routes.NoRoute(gin.WrapH(http.FileServer(http.FS(fsRoot))))

	// Backend API

	// health check
	routes.GET("/api/health", func(c *gin.Context) {
		svc.HealthCheckService().HealthCheck(c, StartTime, BootTime)
	})

	auth := routes.Group("/auth")
	{
		auth.POST("/signup", func(c *gin.Context) {
			svc.UserService().SignUp(c)
		})
		auth.POST("/login", func(c *gin.Context) {
			svc.UserService().SignIn(c)
		})
		auth.POST("/refresh", func(c *gin.Context) {
			svc.UserService().RefreshToken(c)
		})
		auth.GET("/logout", func(c *gin.Context) {
			svc.UserService().SignOut(c)
		})

	}

	api := routes.Group("/api")
	api.Use(middleware.CORSMiddleware())
	api.Use(middleware.JWTAuth())
	{
		api.GET("/notes", func(c *gin.Context) {
			svc.NoteService().ReadAll(c)
		})
		api.GET("/notes/:id", func(c *gin.Context) {
			svc.NoteService().Read(c)
		})
		api.POST("/notes", func(c *gin.Context) {
			svc.NoteService().Create(c)
		})
		api.PUT("/notes/:id", func(c *gin.Context) {
			svc.NoteService().Update(c)
		})
		api.DELETE("/notes/:id", func(c *gin.Context) {
			svc.NoteService().Delete(c)
		})
	}
}
