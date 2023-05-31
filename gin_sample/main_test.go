package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAlbumsList(t *testing.T) {

	router := gin.Default()
	router.GET("/albums", getAlbums)

	req := httptest.NewRequest(http.MethodGet, "/albums", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var albums []album

	json.Unmarshal(recorder.Body.Bytes(), &albums)
	// fmt.Printf("len of albums: %v\n", len(albums))
	// fmt.Printf("%+v", albums)

	if len(albums) != 3 {
		t.Fatal("album list failed")
	}
	// t.Log(recorder.Body.String())
}

func TestAlbumsGet(t *testing.T) {

	router := gin.Default()
	router.GET("/albums/:id", getAlbumByID)

	req := httptest.NewRequest(http.MethodGet, "/albums/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var a album

	json.Unmarshal(recorder.Body.Bytes(), &a)
	// fmt.Printf("len of albums: %v\n", len(albums))
	// fmt.Printf("%+v", albums)

	if a.ID != "1" {
		t.Fatal("album get failed")
	}
	// t.Log(recorder.Body.String())
}

func TestAlbumsPost(t *testing.T) {

	router := gin.Default()
	router.POST("/albums", postAlbums)

	var a album
	a.ID = "4"
	a.Title = "The Modern Sound of Betty Carter"
	a.Artist = "Betty Carter"
	a.Price = 49.99

	body, err := json.Marshal(a)
	if err == nil {
		// fmt.Printf("%v", body)
	} else {
		fmt.Printf("%s", err)
		t.Fatal("albums get failed: json.Marshal")
	}

	req := httptest.NewRequest(http.MethodPost, "/albums", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != 201 {
		fmt.Printf("http status code: %v\n", recorder.Result().StatusCode)
		t.Fatal("albums get failed")
	}
	// t.Log(recorder.Body.String())
}
