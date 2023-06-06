package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type Album struct {
	ID     string  `json:"id" binding:"min=1,max=6"`
	Title  string  `json:"title" binding:"min=1,max=9999"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price" binding:"required,gt=0,lte=99"`
}

type AlbumQuery struct {
	ID string `form:"id" binding:"required,min=1,max=6"`
}

// albums slice to seed record album data.
var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func globalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("middleware globalMiddleware begin, path:%s --------\n", c.Request.URL.Path)
		c.Next()
		fmt.Println("middleware globalMiddleware end --------")
	}
}

func userMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware user begin ----")
		c.Next()
		fmt.Println("middleware user end ----")
	}
}

func getUserInfo(c *gin.Context) {
	fmt.Printf("recv /user/get_info get request\n")

	c.JSON(http.StatusOK, gin.H{"name": "wolf", "age": 18})
}

func getUserStatus(c *gin.Context) {
	fmt.Printf("recv /user/get_status get request\n")

	c.JSON(http.StatusOK, gin.H{"status": "running"})
}

func main() {
	router := gin.Default()
	router.Use(globalMiddleware())
	router.GET("/album_list", getAlbums)
	router.GET("/album_query", getAlbumByID)
	router.POST("/album_add", postAlbums)

	userGroup := router.Group("/user")
	userGroup.GET("/get_status", getUserStatus)
	userGroup.Use(userMiddleware())
	userGroup.GET("/get_info", getUserInfo)
	router.Run("localhost:8081")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	fmt.Printf("recv /albums get request\n")

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"message": "params not valid!"})
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	var r AlbumQuery
	err := c.ShouldBindQuery(&r)
	if err != nil {
		fmt.Printf("ShouldBindQuery error:%s", err)
		c.JSON(http.StatusOK, gin.H{"message": "id not valid!"})
		return
	}
	id := r.ID

	fmt.Printf("recv id %v", r.ID)

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
