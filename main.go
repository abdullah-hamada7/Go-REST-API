// Package declaration - 'main' makes this an executable program
package main

// Import dependencies
import (
	"net/http" // Standard library for HTTP status codes

	"github.com/gin-gonic/gin" // Web framework for building APIs
)

// =========================
// Models (Data Structures)
// =========================

// Book represents the main data structure stored in our system
type Book struct {
	ID       string  `json:"id"`       // Unique identifier, exported as "id" in JSON
	Title    string  `json:"title"`    // Book title, exported as "title" in JSON
	Author   *string `json:"author"`   // Pointer to string - allows null values in JSON
	Quantity int     `json:"quantity"` // Number of copies available
}

// BookCreateInput defines the expected structure for creating new books
// Includes validation rules using Gin's binding tags
type BookCreateInput struct {
	ID       string  `json:"id" binding:"required"`             // Must be provided
	Title    string  `json:"title" binding:"required,min=3"`    // Required & at least 3 chars
	Author   *string `json:"author" binding:"required"`         // Must be provided (can be null string)
	Quantity int     `json:"quantity" binding:"required,gte=1"` // Required & >= 1
}

// BookPutInput for FULL updates (PUT requests) - replaces entire book
type BookPutInput struct {
	Title    string  `json:"title" binding:"required,min=3"` // All fields required for full replacement
	Author   *string `json:"author" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gte=1"`
}

// BookPatchInput for PARTIAL updates (PATCH requests) - updates only provided fields
// All fields are pointers so we can detect which fields were actually provided
type BookPatchInput struct {
	Title    *string `json:"title"`    // Pointer - nil if not provided in request
	Author   *string `json:"author"`   // Pointer - nil if not provided in request
	Quantity *int    `json:"quantity"` // Pointer - nil if not provided in request
}

// =========================
// Initial Data
// =========================

// Helper function to create string pointers
// Needed because we can't directly take the address of string literals like &"hello"
func ptr(s string) *string {
	return &s
}

// Initial book data stored in memory (in real app, this would be a database)
var books = []Book{
	{ID: "1", Title: "Book One", Author: ptr("Author One"), Quantity: 1},
	{ID: "2", Title: "Book Two", Author: ptr("Author Two"), Quantity: 2},
	{ID: "3", Title: "Book Three", Author: ptr("Author Three"), Quantity: 3},
}

// =========================
// Handler Functions
// =========================

// getBooks returns all books in the system
func getBooks(c *gin.Context) {
	// c.JSON sends a JSON response with HTTP 200 status code
	// gin.Context contains request info and response methods
	c.JSON(http.StatusOK, books)
}

// getBook returns a specific book by ID
func getBook(c *gin.Context) {
	// Extract "id" parameter from URL path (e.g., /books/1 -> id = "1")
	id := c.Param("id")

	// Loop through all books to find matching ID
	// range returns (index, value) - we ignore index with _
	for _, b := range books {
		if b.ID == id {
			// Found book - return it with 200 OK
			c.JSON(http.StatusOK, b)
			return // Exit function early
		}
	}

	// If we get here, no book was found - return 404 Not Found
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// createBook adds a new book to the system
func createBook(c *gin.Context) {
	var input BookCreateInput // Declare variable to hold parsed JSON data

	// This is the Go error handling pattern explained earlier:
	// 1. c.ShouldBindJSON(&input) parses request JSON into input struct
	// 2. It returns an error if JSON is invalid or validation fails
	// 3. if err := ...; err != nil checks if error occurred
	// 4. If error, return 400 Bad Request with error message
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // Exit early on error
	}

	// Create new book from validated input data
	newBook := Book{
		ID:       input.ID,
		Title:    input.Title,
		Author:   input.Author,
		Quantity: input.Quantity,
	}

	// Add new book to our slice
	books = append(books, newBook)

	// Return 201 Created status with the new book data
	c.JSON(http.StatusCreated, newBook)
}

// replaceBook completely replaces an existing book (PUT)
func replaceBook(c *gin.Context) {
	id := c.Param("id")    // Get book ID from URL
	var input BookPutInput // Struct for full replacement data

	// Parse and validate input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find book by ID and update all fields
	for i, b := range books {
		if b.ID == id {
			// Update the book in the slice (using index i)
			books[i].Title = input.Title
			books[i].Author = input.Author
			books[i].Quantity = input.Quantity

			// Return updated book
			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	// Book not found
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// updateBook partially updates a book (PATCH)
func updateBook(c *gin.Context) {
	id := c.Param("id")
	var input BookPatchInput // All fields are pointers

	// Parse input - only provided fields will be non-nil
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find and update only the provided fields
	for i, b := range books {
		if b.ID == id {
			// Only update Title if provided (pointer not nil)
			if input.Title != nil {
				books[i].Title = *input.Title // Dereference pointer to get actual value
			}
			// Only update Author if provided
			if input.Author != nil {
				books[i].Author = input.Author
			}
			// Only update Quantity if provided
			if input.Quantity != nil {
				books[i].Quantity = *input.Quantity // Dereference pointer
			}

			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// deleteBook removes a book from the system
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	// Find book by ID
	for i, b := range books {
		if b.ID == id {
			// Remove book from slice using slice manipulation:
			// books[:i] = elements from start to index i-1
			// books[i+1:] = elements from index i+1 to end
			// append(...) combines them, effectively removing element at index i
			// ... unpacks the second slice into individual elements
			books = append(books[:i], books[i+1:]...)

			// Return 204 No Content (successful deletion, no response body)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// checkoutBook handles checking out a book (decreasing quantity by 1)
// Now uses the route /books/checkout/:id instead of /books/:id/checkout
func checkoutBook(c *gin.Context) {
	id := c.Param("id") // Extract "id" parameter from URL path (e.g., /books/checkout/1 -> id = "1")

	// Loop through all books to find the one with matching ID
	// We need both index (i) and value (b) because:
	// - b (value) is used for reading/comparison (checking ID and current Quantity)
	// - i (index) is used for modification (updating Quantity in the original slice)
	for i, b := range books {
		if b.ID == id {
			if b.Quantity > 0 {
				// Book found AND has available copies
				// Decrement the quantity by 1 (check out one copy)
				// Must use books[i].Quantity-- NOT b.Quantity-- because:
				// - b is a COPY of the book from the range loop
				// - books[i] accesses the ORIGINAL book in the slice
				books[i].Quantity--

				// Return 200 OK status (successful checkout)
				c.JSON(http.StatusOK, gin.H{
					"message": "Book checked out successfully",
					"book":    books[i],
				})
				return // Exit function early
			}
			// Book found but no copies available (quantity is 0)
			// Use 409 Conflict to indicate the request conflicts with current state
			c.JSON(http.StatusConflict, gin.H{
				"message": "Book is out of stock",
				"book":    b,
			})
			return
		}
	}

	// If we get here, no book was found with the given ID
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// returnBook handles returning a book (increasing quantity by 1)
// Now uses the route /books/return/:id instead of /books/:id/return
func returnBook(c *gin.Context) {
	id := c.Param("id") // Extract "id" parameter from URL path (e.g., /books/return/1 -> id = "1")

	// Loop through all books to find the one with matching ID
	for i, b := range books {
		if b.ID == id {
			// Book found - increment the quantity by 1 (return one copy)
			// Using books[i].Quantity++ to modify the original book in the slice
			books[i].Quantity++

			// Return 200 OK status with success message and updated book
			c.JSON(http.StatusOK, gin.H{
				"message": "Book returned successfully",
				"book":    books[i],
			})
			return // Exit function early
		}
	}

	// If we get here, no book was found with the given ID
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// =========================
// Main Function - Application Entry Point
// =========================

func main() {
	// Create Gin router with default middleware (logging, panic recovery)
	router := gin.Default()

	// Register routes - map HTTP methods and paths to handler functions
	router.GET("/books", getBooks)                   // Get all books
	router.GET("/books/:id", getBook)                // Get single book by ID
	router.POST("/books", createBook)                // Create new book
	router.PUT("/books/:id", replaceBook)            // Fully replace book
	router.PATCH("/books/:id", updateBook)           // Partially update book
	router.DELETE("/books/:id", deleteBook)          // Delete book
	router.POST("/books/checkout/:id", checkoutBook) // Check out a book (decrease quantity) - ROUTE CHANGED
	router.POST("/books/return/:id", returnBook)     // Return a book (increase quantity) - ROUTE CHANGED

	// Start HTTP server on port 8080
	// This blocks and keeps the server running until terminated
	router.Run(":8080")
}
