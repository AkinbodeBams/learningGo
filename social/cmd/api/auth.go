package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type RegisterUserPayload struct {
	Email string `json:"email" validate:"required,max=100"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type UserWithToken struct{
*store.User
TOken string `json:"token"`
}



// registerUserHandler godoc
//
//	@Summary		Register a user
//	@Description	get a user by ID
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequest(w, r, err)
		return
	}



	// userId := 1

	user := &store.User{
		Email:   payload.Email,
		Username: payload.Username,
		
		
	}

	if err:=  user.Password.Set(payload.Password); err!= nil {
		app.internalServerError(w,r,err)
		return
	}

	ctx := r.Context()

	plainTOken := uuid.New().String()
	hash:= sha256.Sum256([]byte(plainTOken))
	hashToken := hex.EncodeToString(hash[:])
	
	
	err:= app.store.Users.CreateAndInvite(ctx,user,hashToken,app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequest(w,r,err)
		case store.ErrDuplicateUsername:
			app.badRequest(w,r,err)
		default:
			app.internalServerError(w,r,err)
		}
		return
	}

userWithTOken := UserWithToken{
	User:user,
	TOken: plainTOken,
}

	if err := app.jsonResponse(w, http.StatusCreated, userWithTOken); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateTokenPayload struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTokenPayload	true	"User credentials"
//	@Success		200		{string}	string				"Token"
//	@Failure		400		{string}	error
//	@Failure		401		{string}	error
//	@Failure		500		{string}	error
//	@Router			/authentication/token [post]
func (app *application)createTokenHandler(w http.ResponseWriter , r *http.Request){

	var payload CreateTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
if err := Validate.Struct(&payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.store.Users.GetByEmail(r.Context(),payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedError(w,r,err)
		default:
			app.internalServerError(w,r,err)

		} 
		return
	}
claims := jwt.MapClaims{
	"sub": user.ID,
	"exp": time.Now().Add(app.config.auth.token.exp).Unix(),

	"iat":time.Now().Unix(),
	"nbf":time.Now().Unix(),
	"iss": app.config.auth.token.iss,
	"aud": app.config.auth.token.iss,
}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w,r,err)
	}
	 
	if err := app.jsonResponse(w, http.StatusCreated,token);err!= nil{
		app.internalServerError(w,r,err)
	}

}

