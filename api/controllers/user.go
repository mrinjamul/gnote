package controllers

// Create a authentication system using JWT
// To implement Multi-level Authentication

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/repository"
	"github.com/mrinjamul/gnote/utils"
)

var (
	jwtKey string
)

func init() {
	jwtKey = utils.GetEnv("JWT_SECRET")
	// if jwtKey == "" {
	// 	panic("JWT_SECRET not set")
	// }
}

// User is a controller for users
type User interface {
	// SignUp creates a new user
	SignUp(ctx *gin.Context)
	// LoginUser logs in a user
	SignIn(gin *gin.Context)
	// RefreshToken refreshes the token
	RefreshToken(ctx *gin.Context)
	// SignOut logs out a user
	SignOut(gin *gin.Context)
	// UserDetails returns the user details from token
	UserDetails(ctx *gin.Context)
	// ViewUser returns the public user
	ViewUser(ctx *gin.Context)
	// UpdateUser updates the user details
	UpdateUser(ctx *gin.Context)
	// DeleteUser deletes a user
	DeleteUser(ctx *gin.Context)
}

// user is a controller for users
type user struct {
	userRepo repository.UserRepo
}

// SignUp creates a new user
func (u *user) SignUp(ctx *gin.Context) {
	var user models.User
	// Get the JSON body and decode into user struct
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		ctx.Abort()
		return
	}

	// check if valid username
	if !utils.IsValidUserName(user.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid username",
		})
		ctx.Abort()
		return
	}
	user.Username = strings.ToLower(user.Username)
	user.Username = strings.TrimSpace(user.Username)

	// check if valid email
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	// Validate Password
	ok := utils.IsValidPassword(user.Password)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad request",
			"message": "bad password",
		})
		ctx.Abort()
		return
	}

	user.Role = "user"
	user.Level = 1

	// Hash the password before storing
	user.Password, err = utils.HashAndSalt(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	// Create the user
	err = u.userRepo.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	// First User will be admin
	if user.ID == 1 {
		user.Role = "admin"
		user.Level = 4
		err = u.userRepo.UpdateUser(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
	}

	userinfo := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"middle_name": user.MiddleName,
		"last_name":   user.LastName,
		"full_name":   user.FirstName + "" + user.MiddleName + "" + user.LastName,
		"dob":         user.DOB,
		"created_at":  user.CreatedAt,
		"deleted_at":  user.DeletedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user created successfully",
		"user":    userinfo,
	})
}

// SignIn logs in a user
func (u *user) SignIn(ctx *gin.Context) {
	var creds models.Credentials
	var user models.User
	// Get the JSON body and decode into creds struct
	err := ctx.BindJSON(&creds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Trim spaces
	user.Username = strings.ToLower(user.Username)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	if (creds.Email == "" && creds.Username == "") || creds.Password == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "email or password cannot be empty",
		})
		return
	}
	// Get the expected password from the database
	user, err = u.userRepo.GetUserByUsername(creds.Username)
	if err != nil || user.ID == 0 {
		user, err = u.userRepo.GetUserByEmail(creds.Email) // error
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			ctx.Abort()
			return
		}
	}

	// if user is deleted then return unauthorized
	if user.DeletedAt.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is deleted",
		})
		ctx.Abort()
		return
	}

	// Validate the password
	valid := utils.VerifyHash(creds.Password, user.Password)
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Password",
		})
		ctx.Abort()
		return
	}

	// expires in  5 minutes
	issuedAt := time.Now()
	expiresAt := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: creds.Username,
		Role:     user.Role,
		Level:    user.Level,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(issuedAt),
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get Hostname
	hostname := ctx.Request.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	// set cookie with name "token" and value "tokenString"
	ctx.SetCookie("token", tokenString, utils.ToMaxAge(expiresAt), "/", hostname, false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  tokenString,
	})
}

// RefreshToken refreshes the token
func (u *user) RefreshToken(ctx *gin.Context) {
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		ctx.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": "cookie not found",
	// 		})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "bad cookie",
	// 	})
	// }
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()

			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad token",
		})
		ctx.Abort()
		return
	}

	// check if token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "token expired",
		})
		ctx.Abort()
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 60*time.Second {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad request",
		})
		ctx.Abort()
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	issuedAt := time.Now()
	expiresAt := time.Now().Add(5 * time.Minute)
	claims.IssuedAt = jwt.NewNumericDate(issuedAt)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	hostname := ctx.Request.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}
	// Set cookie with name "token" and value "tokenString"
	ctx.SetCookie("token", tokenString, utils.ToMaxAge(expiresAt), "/", hostname, false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  tokenString,
	})
}

// SignOut logs out a user
func (u *user) SignOut(ctx *gin.Context) {

	hostname := ctx.Request.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}
	// remove the token from the cookies
	ctx.SetCookie("token", "", -1, "/", hostname, false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out",
	})
	// Redirect to the login page
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

// UserDetails returns the user details
func (u *user) UserDetails(ctx *gin.Context) {
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	user, err := u.userRepo.GetUserByUsername(claims.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	userinfo := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"middle_name": user.MiddleName,
		"last_name":   user.LastName,
		"full_name":   strings.TrimSpace(user.FirstName + " " + user.MiddleName + " " + user.LastName),
		"dob":         user.DOB,
		"created_at":  user.CreatedAt,
		"deleted_at":  user.DeletedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Welcome, " + user.FirstName + "!",
		"user":    userinfo,
	})
}

// ViewUser returns the public user details
func (u *user) ViewUser(ctx *gin.Context) {
	// get the username param from context
	username := ctx.Param("username")
	// get the user from the database
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	// check if user is found or not
	if user.ID > 0 && !user.DeletedAt.Valid {
		// if user is found, return the user info
		userinfo := map[string]string{
			"username":   user.Username,
			"full_name":  user.FirstName + " " + user.MiddleName + " " + user.LastName,
			"email":      user.Email,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   userinfo,
		})
	} else {
		// if user is not found, return error
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
	}
}

// UpdateUser updates the user details
func (u *user) UpdateUser(ctx *gin.Context) {
	var userinfo map[string]interface{}
	// Get the JSON body and decode into userinfo struct
	err := ctx.BindJSON(&userinfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	user, err := u.userRepo.GetUserByUsername(claims.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	if userinfo["first_name"] != nil {
		user.FirstName = userinfo["first_name"].(string)
	}
	if userinfo["middle_name"] != nil {
		user.FirstName = userinfo["middle_name"].(string)
	}
	if userinfo["last_name"] != nil {
		user.LastName = userinfo["last_name"].(string)
	}
	if userinfo["email"] != nil {
		user.Email = userinfo["email"].(string)
	}
	if userinfo["username"] != nil {
		user.Username = userinfo["username"].(string)
	}
	if userinfo["dob"] != nil {
		user.DOB = userinfo["dob"].(time.Time)
	}

	// Update the user
	err = u.userRepo.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}
	// Return the updated user
	userinfo = map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"middle_name": user.MiddleName,
		"last_name":   user.LastName,
		"full_name":   strings.TrimSpace(user.FirstName + " " + user.MiddleName + " " + user.LastName),
		"dob":         user.DOB,
		"created_at":  user.CreatedAt,
		"deleted_at":  user.DeletedAt,
	}

	// Generate new JWT Token
	issuedAt := time.Now()
	expiresAt := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims = &models.Claims{
		Username: user.Username,
		Role:     user.Role,
		Level:    user.Level,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(issuedAt),
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err = token.SignedString([]byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	hostname := ctx.Request.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}
	// Set cookie with name "token" and value "tokenString"
	ctx.SetCookie("token", tokenString, utils.ToMaxAge(expiresAt), "/", hostname, false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
		"token":   tokenString,
		"user":    userinfo,
	})
}

// DeleteUser deletes a user
func (u *user) DeleteUser(ctx *gin.Context) {
	var creds models.Credentials
	var user models.User
	// Get the JSON body and decode into creds struct
	err := ctx.BindJSON(&creds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	user, err = u.userRepo.GetUserByUsername(claims.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}

	// Verify the password before deleting the user
	if creds.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "password is required",
		})
		ctx.Abort()
		return
	}
	valid := utils.VerifyHash(creds.Password, user.Password)
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Password",
		})
		ctx.Abort()
		return
	}

	err = u.userRepo.DeleteUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
	// Redirect to the home page
	ctx.Redirect(http.StatusMovedPermanently, "/")
}

// NewUser initializes a new user controller
func NewUser(userRepo repository.UserRepo) User {
	return &user{
		userRepo: userRepo,
	}
}
