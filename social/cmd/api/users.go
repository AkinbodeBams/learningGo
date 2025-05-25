package main

import (
	"context"
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string
const userCtx  userKey = "user"

// ActivateUser godoc
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserhandler(w http.ResponseWriter, r *http.Request){
token := chi.URLParam(r, "token")

err:= app.store.Users.Activate(r.Context(), token)

if err != nil {
	switch err {
	case store.ErrNotFound:
		app.notFound(w,r,err)
	default:
		app.internalServerError(w,r,err)
	}
	return 
}
}


// Get User godoc
//	@Summary		Fetches a user profile
//	@Description	get a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application)  getUserHandler(w http.ResponseWriter,r *http.Request){
	user:= getUserFromCtx(r)
	user, err := app.store.Users.GetById(r.Context(),user.ID)
	if err != nil {
		app.internalServerError(w,r,err)
		return
	}
	
	if err := app.jsonResponse(w,http.StatusOK,user);err!=nil{
		app.internalServerError(w,r,err)
		
		return
	}
}

// DeleteUser godoc
//	@Summary		Delete a user
//	@Description	Deletes a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"No Content"
//	@Failure		400		{object}	Error
//	@Failure		404		{object}	Error
//	@Failure		500		{object}	Error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [delete]
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	if err = app.store.Users.Delete(ctx, int(id)); err != nil {
		switch{
			case errors.Is(err,store.ErrNotFound):
				
				app.notFound(w,r,err)
			default:
				app.internalServerError(w, r, err)

		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateUserPayload struct {
	Email *string `json:"email" validate:"omitempty,max=100"`
	Username *string `json:"username" validate:"omitempty,max=1000"`
}

// UpdateUser godoc
//	@Summary		Update a user's profile
//	@Description	Updates email or username for the authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateUserPayload	true	"User update payload"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}  [patch]
func(app *application) updateUserHandler(w http.ResponseWriter, r *http.Request){
	user:= getUserFromCtx(r)
	
if user == nil {
	app.notFound(w, r, errors.New("user not found in context"))
	return
}

	var payload UpdateUserPayload
	if err := readJSON(w,r,&payload); err!= nil {
		app.badRequest(w,r,err)
		return 
	}

	if err := Validate.Struct(payload);err != nil {
		app.badRequest(w,r,err)
		return 
	}

	if payload.Username != nil  {
		user.Username = *payload.Username
	}
	if payload.Email != nil  {
		user.Email = *payload.Email
	}

	
if err := app.store.Users.Update(r.Context(),user); err != nil {
	
	app.internalServerError(w,r,err)
	return
}
	if err := app.jsonResponse(w,http.StatusOK,user);err!= nil{
		app.internalServerError(w,r,err)
		 return 
	}

}
type followUser struct {
	UserId int `json:"user_id"`
}

// FollowUser godoc
//	@Summary		Follow a user
//	@Description	Follow another user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		followUser	true	"User to follow"
//	@Success		204		{string}	string		"No Content"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request){
	
	followerUser := getUserFromCtx(r)
	
	var payload followUser

		if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	
if err := app.store.Followers.Follow(r.Context(), followerUser.ID,payload.UserId); err!= nil {
	switch err {
	case store.ErrConflict:
		app.Conflict(w,r,err)
		
	default:
		app.internalServerError(w,r,err)
	}
	return
}
w.WriteHeader(http.StatusNoContent)
	
}

// UnfollowUser godoc
//	@Summary		Unfollow a user
//	@Description	Unfollow another user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		followUser	true	"User to unfollow"
//	@Success		204		{string}	string		"No Content"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/unfollow [post]
func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request){
unfollowedUser := getUserFromCtx(r)
	
	var payload followUser

		if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
if err := app.store.Followers.UnFollow(r.Context(), int(unfollowedUser.ID),payload.UserId); err!= nil {
	app.internalServerError(w,r,err)
	return

}
w.WriteHeader(http.StatusNoContent)
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		
		if _, err := strconv.Atoi(userID); err != nil {
    app.badRequest(w, r, fmt.Errorf("invalid or missing userID in URL path"))
    return
}
		fmt.Println(userID)
		id, err := strconv.ParseInt(userID,10,64)
		if err!= nil{
			app.internalServerError(w,r,err)	
			return
		}
		ctx:= r.Context()
		user,err:= app.store.Users.GetById(ctx, id);
		
		if err!= nil{
			switch {
			case errors.Is(err,store.ErrNotFound):
				app.notFound(w,r,err)
			default:
				app.internalServerError(w,r,err)
	
			}
		return
		
		}

		ctx= context.WithValue(ctx,userCtx,user)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User{
	user,_ := r.Context().Value(userCtx).(*store.User)
	fmt.Println(user)
	return user
}