package main

import (
	//"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	// TODO: Finish struct
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Pages int64  `json:"pages"`
}

var bookId int64 = 1

var bookshelf = []Book{
	// TODO: Init bookshelf
	{1, "Blue Bird", 500},
}

func getBooks(c *gin.Context) {
	c.JSON(200, bookshelf) //OK
}
func getBook(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	for _, book := range bookshelf {
		if id == book.Id {
			c.JSON(200, book) //OK
			return
		}
	}
	c.JSON(404, gin.H{"message": "book not found"}) //Not found
}
func addBook(c *gin.Context) {
	var newBook Book
	err := c.BindJSON(&newBook)
	if err != nil {
		return
	}
	for _, book := range bookshelf {
		if newBook.Name == book.Name {
			c.JSON(409, gin.H{"message": "duplicate book name"}) //Conflict
			return
		}
	}
	bookId++
	newBook.Id = bookId
	bookshelf = append(bookshelf, newBook)
	c.JSON(201, newBook) //Created
}
func deleteBook(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	for i, book := range bookshelf {
		if id == book.Id {
			bookshelf = slices.Delete(bookshelf, i, i+1)
			c.JSON(204, nil) //No content
			return
		}
	}
	c.JSON(204, nil) //No content
}
func updateBook(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	for _, book := range bookshelf {
		if newBook.Name == book.Name {
			c.JSON(409, gin.H{"message": "duplicate book name"}) //Conflict
			return
		}
	}
	for i, book := range bookshelf {
		if id == book.Id {
			book.Name = newBook.Name
			book.Pages = newBook.Pages
			bookshelf[i] = book
			c.JSON(200, bookshelf[i]) //OK
			return
		}
	}
	c.JSON(404, gin.H{"message": "book not found"}) //Not found
}

func main() {
	r := gin.Default()
	r.RedirectFixedPath = true

	// TODO: Add routes
	r.GET("/bookshelf", getBooks)
	r.GET("/bookshelf/:id", getBook)
	r.POST("/bookshelf", addBook)
	r.DELETE("/bookshelf/:id", deleteBook)
	r.PUT("/bookshelf/:id", updateBook)

	err := r.Run(":8087")
	if err != nil {
		return
	}
}
