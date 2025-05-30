package main

import (
	"net/http"
	"restApi/db"
	"restApi/models"

	"github.com/gin-gonic/gin"
)

func main(){
	db.InitDb()
	server:= gin.Default()
	server.GET("/events",getEvents)
	server.POST("/events", createEvent)
	server.Run(":8080")
}

func getEvents(context *gin.Context){
	events := models.GetAllEvents()
context.JSON(http.StatusOK,events)
}

func createEvent(context *gin.Context){
var event  models.Event
 err := context.ShouldBindJSON(&event)

 if err != nil {
	context.JSON(http.StatusBadRequest,gin.H{"message":"Could not parse "})
return 
}

 event.ID = 1
 event.UserId = 1
event.Save()
 context.JSON(http.StatusCreated,gin.H{"Message":"Event created","event":event})

}