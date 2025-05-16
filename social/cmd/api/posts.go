package main

import (
	"errors"

	"net/http"
	"strconv"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/go-chi/chi/v5"
)


type CreatePostPayload struct {
	Title string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=100"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	userId := 1

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserId:  int64(userId),
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}


func (app *application)  getPostHandler(w http.ResponseWriter,r *http.Request){
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID,10,64)
	if err!= nil{
		app.internalServerError(w,r,err)	
	}
	ctx:= r.Context()
	post,err:= app.store.Posts.GetById(ctx, int(id));
	
	if err!= nil{
		switch {
		case errors.Is(err,store.ErrNotFound):
			app.notFound(w,r,err)
		default:
			app.internalServerError(w,r,err)

		}
	return
	}
	if err := writeJSON(w,http.StatusOK,post);err!=nil{
		app.internalServerError(w,r,err)
		
		return
	}
}