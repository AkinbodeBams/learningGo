package main

import (
	"context"
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


type postKey string
const postCtx  postKey = "post"

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

	ctx := r.Context()
	user := getUserFromCtx(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserId:  int64(user.ID),
		Tags:    payload.Tags,
	}

	

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}


func (app *application)  getPostHandler(w http.ResponseWriter,r *http.Request){
	
	post:= getPostFromCtx(r)
	comments, err := app.store.Comments.GetByPostID(r.Context(),post.ID)
	if err != nil {
		app.internalServerError(w,r,err)
		return
	}
	post.Comments = comments
	if err := writeJSON(w,http.StatusOK,post);err!=nil{
		app.internalServerError(w,r,err)
		
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	if err = app.store.Posts.Delete(ctx, int(id)); err != nil {
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

type UpdatePostPayload struct {
	Title *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func(app *application) updatePostHandler(w http.ResponseWriter, r *http.Request){
	post:= getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w,r,&payload); err!= nil {
		app.badRequest(w,r,err)
		return 
	}

	if err := Validate.Struct(payload);err != nil {
		app.badRequest(w,r,err)
	}

	if payload.Content != nil  {
		post.Content = *payload.Content
	}
	if payload.Title != nil  {
		post.Title = *payload.Title
	}

	
if err := app.store.Posts.Update(r.Context(),post); err != nil {
	app.internalServerError(w,r,err)
}
	if err := writeJSON(w,http.StatusOK,post);err!= nil{
		app.internalServerError(w,r,err)
		 
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx= context.WithValue(ctx,postCtx,post)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post{
	post,_ := r.Context().Value("post").(*store.Post)
	return post
}