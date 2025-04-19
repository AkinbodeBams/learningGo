package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akinbodeBams/go-bookstore/pkg/models"
	"github.com/akinbodeBams/go-bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(newBooks); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}



func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]

	ID, err := strconv.ParseInt(bookId, 10, 64) // base 10, 64-bit integer
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	bookDetails, _ := models.GetBooksById(int(ID))
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder( w).Encode(bookDetails); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook:= &models.Book{}
	utils.ParseBody(r, CreateBook)
	b := CreateBook.CreateBook()
	res,_ := json.Marshal(b)
	// if err != nil {
	// 	http.Error(w, "Could not create book", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header()
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

	
	// if err := json.NewEncoder(w).Encode(b); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// }
}
