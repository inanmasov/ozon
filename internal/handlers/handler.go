package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"ozon/internal/cache"
	"ozon/internal/database"

	"github.com/gin-gonic/gin"
)

type Link struct {
	OriginalURL string `json:"original_url"`
}

func ShortenURLHandler(c *gin.Context) {
	var link Link

	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if link.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
		return
	}

	var shortURL string
	for {
		shortURL = generateShortURL()
		if shortURL != "" {
			break
		}
	}

	// Получение значения переменной среды
	storage := os.Getenv("STORAGE")

	if storage == "postgres" {
		db, err := database.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := db.SaveURL(link.OriginalURL, shortURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
			return
		}
	} else if storage == "cache" {
		cache := cache.InitCache()
		if err := cache.SaveURL(link.OriginalURL, shortURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"short_url": shortURL})
}

func RedirectHandler(c *gin.Context) {
	shortURL := c.Param("shortURL")
	var originalURL string

	// Получение значения переменной среды
	storage := os.Getenv("STORAGE")

	if storage == "postgres" {
		db, err := database.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		originalURL, err = db.GetOriginalURL(shortURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
	} else if storage == "cache" {
		cache := cache.InitCache()

		var err error
		originalURL, err = cache.GetOriginalURL(shortURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
	}

	c.Redirect(http.StatusFound, originalURL)
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	shortURLLength := 10
	b := make([]byte, shortURLLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	var shortURL string
	// Получение значения переменной среды
	storage := os.Getenv("STORAGE")

	if storage == "postgres" {
		db, err := database.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.CheckShortUrl(string(b))
		if err != nil {
			return ""
		} else {
			shortURL = string(b)
		}
	} else if storage == "cache" {
		cache := cache.InitCache()
		err := cache.CheckShortUrl(string(b))
		if err != nil {
			return ""
		} else {
			shortURL = string(b)
		}
	}
	return shortURL
}
