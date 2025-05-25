package main

import (
	"log"
	"net/http"

	"github.com/akinbodeBams/social/internal/store"
)
 
func (app *application) getUserFeedHandler(w http.ResponseWriter , r *http.Request){
	fq := store.PaginatedFeedQuery{
Limit: 20,
Offset: 0,
Sort: "desc",
	}
	fq, err :=  fq.Parse(r)
	
	if err != nil {
		app.badRequest(w,r,err)
		return 
	}
	
	if err := Validate.Struct(fq); err != nil {
		app.badRequest(w,r, err)
	}
	ctx := r.Context()
	user:= getUserFromCtx(r)
	

	feed,err := app.store.Posts.GetUserFeed(ctx,int(user.ID),fq)
	if err != nil {
		log.Fatal(err)
		app.internalServerError(w,r,err)
		return 
	}

	if err:= app.jsonResponse(w,http.StatusOK,feed);err!= nil{
		app.internalServerError(w,r,err)
	}
}