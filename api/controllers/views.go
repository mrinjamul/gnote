package controllers

import (
	"io/fs"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/utils"
)

type Views interface {
	App(ctx *gin.Context, fsRoot fs.FS)
	Login(ctx *gin.Context, fsRoot fs.FS)
	Register(ctx *gin.Context, fsRoot fs.FS)
	MyAccount(ctx *gin.Context, fsRoot fs.FS)
	NotFound(ctx *gin.Context, fsRoot fs.FS)
	Delete(ctx *gin.Context, fsRoot fs.FS)
	DeleteNote(ctx *gin.Context, fsRoot fs.FS)
}

type views struct {
}

// App returns a app page
func (v *views) App(ctx *gin.Context, fsRoot fs.FS) {
	// check if token is present
	// Get cookie "token"
	_, err := ctx.Cookie("token")
	if err != nil {
		_, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			// Get welcome.html from fsRoot
			welcome, err := fsRoot.Open("welcome.html")
			if err != nil {
				panic(err)
			}
			defer welcome.Close()
			// Read the file
			b, err := ioutil.ReadAll(welcome)
			if err != nil {
				panic(err)
			}
			// Write the content to the response
			// ctx.Writer.Write(b)
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
			return
		}
	}

	// Get index.html from fsRoot
	index, err := fsRoot.Open("index.html")
	if err != nil {
		panic(err)
	}
	defer index.Close()
	// Read the file
	b, err := ioutil.ReadAll(index)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	// ctx.Writer.Write(b)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)

}

// Login returns a login page
func (v *views) Login(ctx *gin.Context, fsRoot fs.FS) {
	// Get login.html from fsRoot
	login, err := fsRoot.Open("login.html")
	if err != nil {
		panic(err)
	}
	defer login.Close()
	// Read the file
	b, err := ioutil.ReadAll(login)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	// ctx.Writer.Write(b)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// Register returns a register page
func (v *views) Register(ctx *gin.Context, fsRoot fs.FS) {
	// Get register.html from fsRoot
	register, err := fsRoot.Open("register.html")
	if err != nil {
		panic(err)
	}
	defer register.Close()
	// Read the file
	b, err := ioutil.ReadAll(register)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// MyAccount returns a my account page
func (v *views) MyAccount(ctx *gin.Context, fsRoot fs.FS) {
	// check if token is present
	// Get cookie "token"
	_, err := ctx.Cookie("token")
	if err != nil {
		_, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			// redirect to login page if not logged in
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
	}
	// Get account.html from fsRoot
	account, err := fsRoot.Open("account.html")
	if err != nil {
		panic(err)
	}
	defer account.Close()
	// Read the file
	b, err := ioutil.ReadAll(account)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	// ctx.Writer.Write(b)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// NotFound returns a 404 page
func (v *views) NotFound(ctx *gin.Context, fsRoot fs.FS) {
	// Get 404.html from fsRoot
	notFound, err := fsRoot.Open("404.html")
	if err != nil {
		panic(err)
	}
	defer notFound.Close()
	// Read the file
	b, err := ioutil.ReadAll(notFound)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// Delete returns a delete page
func (v *views) Delete(ctx *gin.Context, fsRoot fs.FS) {
	// check if token is present
	// Get cookie "token"
	_, err := ctx.Cookie("token")
	if err != nil {
		_, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			// redirect to login page if not logged in
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
	}
	// Get delete.html from fsRoot
	delete, err := fsRoot.Open("delete.html")
	if err != nil {
		panic(err)
	}
	defer delete.Close()
	// Read the file
	b, err := ioutil.ReadAll(delete)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	// ctx.Writer.Write(b)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// Delete returns a delete page
func (v *views) DeleteNote(ctx *gin.Context, fsRoot fs.FS) {
	// check if token is present
	// Get cookie "token"
	_, err := ctx.Cookie("token")
	if err != nil {
		_, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			// redirect to login page if not logged in
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
	}
	// Get delete_note.html from fsRoot
	delete, err := fsRoot.Open("delete_note.html")
	if err != nil {
		panic(err)
	}
	defer delete.Close()
	// Read the file
	b, err := ioutil.ReadAll(delete)
	if err != nil {
		panic(err)
	}
	// Write the content to the response
	// ctx.Writer.Write(b)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

// NewViews returns a new Views
func NewViews() Views {
	return &views{}
}
