package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID       string `json:"id"`
	Title    string	`json:"title"`
	Author   string	`json:"author"`
	Quantity int	`json:"quantity"`
}

var books = []Book{
	{ID : "1", Title : "Harry Potter 1", Author : "Marcel Proust", Quantity : 2},
	{ID : "2", Title : "Harry Potter 2", Author : "JK Rowling ", Quantity : 2},
	{ID : "3", Title : "Harry Potter 3", Author : "Someone Famous", Quantity : 2},
}

func getBooks(c * gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c * gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		log.Fatal(err)
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*Book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not Found")
}

func bookById(c * gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "Book not Found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkOutBook(c * gin.Context){
	// id, ok := c.GetQuery("id")
	id := c.Param("id")

	// if !ok {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message" : "Missing ID in Query"})
	// 	return
	// }

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "Book not Found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message" : "Book not Available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("books/:id", bookById)

	router.POST("/books", createBook)

	router.PATCH("/checkout/:id", checkOutBook)
	router.Run("localhost:3000")
}
