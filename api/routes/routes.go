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

	// Get Views
	fsRoot, err := fs.Sub(ViewsFs, "views")
	if err != nil {
		log.Println(err)
	}
	// Get static files
	fsStatic, err := fs.Sub(ViewsFs, "views/static")
	if err != nil {
		log.Println(err)
	}
	// Serve static files
	routes.StaticFS("/static", http.FS(fsStatic))

	// Serve the frontend
	// Home Page
	routes.GET("/", func(ctx *gin.Context) {
		svc.ViewService().App(ctx, fsRoot)
	})
	// Login Page
	routes.GET("/login", func(ctx *gin.Context) {
		svc.ViewService().Login(ctx, fsRoot)
	})
	// Register Page
	routes.GET("/register", func(ctx *gin.Context) {
		svc.ViewService().Register(ctx, fsRoot)
	})
	// My Account Page
	routes.GET("/account", func(ctx *gin.Context) {
		svc.ViewService().MyAccount(ctx, fsRoot)
	})

	// Delete Page
	routes.GET("/delete/user", func(ctx *gin.Context) {
		svc.ViewService().Delete(ctx, fsRoot)
	})
	routes.GET("/delete/notes", func(ctx *gin.Context) {
		svc.ViewService().DeleteNote(ctx, fsRoot)
	})

	// Add 404 page
	routes.NoRoute(func(ctx *gin.Context) {
		svc.ViewService().NotFound(ctx, fsRoot)
	})

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
		auth.POST("/logout", func(c *gin.Context) {
			svc.UserService().SignOut(c)
		})

	}

	userRoute := routes.Group("/user")
	{
		userRoute.GET("/:username", func(ctx *gin.Context) {
			svc.UserService().ViewUser(ctx)
		})
		userRoute.GET("/search", func(ctx *gin.Context) {
			// TODO: implement search
			// response not implemented
			ctx.JSON(200, gin.H{
				"message": "Search not implemented",
			})
		})
		userRoute.GET("/me", middleware.JWTAuth(), func(ctx *gin.Context) {
			svc.UserService().UserDetails(ctx)
		})
		userRoute.PATCH("/me", middleware.JWTAuth(), func(ctx *gin.Context) {
			svc.UserService().UpdateUser(ctx)
		})
		userRoute.DELETE("/me", middleware.JWTAuth(), func(ctx *gin.Context) {
			svc.UserService().DeleteUser(ctx)
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
		api.DELETE("/notes", func(c *gin.Context) {
			svc.NoteService().DeleteByUsername(c)
		})
	}
}
