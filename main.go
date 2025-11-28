package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Author   *string `json:"author"`
	Quantity int     `json:"quantity"`
}

type BookCreateInput struct {
	ID       string  `json:"id" binding:"required"`
	Title    string  `json:"title" binding:"required,min=3"`
	Author   *string `json:"author" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gte=1"`
}

type BookPutInput struct {
	Title    string  `json:"title" binding:"required,min=3"`
	Author   *string `json:"author" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gte=1"`
}

type BookPatchInput struct {
	Title    *string `json:"title"`
	Author   *string `json:"author"`
	Quantity *int    `json:"quantity"`
}

func ptr(s string) *string { return &s }

var books = []Book{
	{ID: "1", Title: "Book One", Author: ptr("Author One"), Quantity: 1},
	{ID: "2", Title: "Book Two", Author: ptr("Author Two"), Quantity: 2},
	{ID: "3", Title: "Book Three", Author: ptr("Author Three"), Quantity: 3},
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")

	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

func createBook(c *gin.Context) {
	var input BookCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook := Book{
		ID:       input.ID,
		Title:    input.Title,
		Author:   input.Author,
		Quantity: input.Quantity,
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func replaceBook(c *gin.Context) {
	id := c.Param("id")
	var input BookPutInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, b := range books {
		if b.ID == id {

			books[i].Title = input.Title
			books[i].Author = input.Author
			books[i].Quantity = input.Quantity

			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var input BookPatchInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, b := range books {
		if b.ID == id {

			if input.Title != nil {
				books[i].Title = *input.Title
			}
			if input.Author != nil {
				books[i].Author = input.Author
			}
			if input.Quantity != nil {
				books[i].Quantity = *input.Quantity
			}

			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", createBook)
	router.PUT("/books/:id", replaceBook)
	router.PATCH("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)

	router.Run(":8080")
}
