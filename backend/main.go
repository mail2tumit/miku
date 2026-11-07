package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// News represents a news article.
type News struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

// Connect to PostgreSQL.
func connectDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is not configured")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL!")

	return db
}

// Retrieve all news articles.
func getNews(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, title, content, image_url, created_at
			FROM news
			ORDER BY created_at DESC
		`)

		if err != nil {
			log.Println("Query error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve news",
			})
			return
		}

		defer rows.Close()

		newsList := make([]News, 0)

		for rows.Next() {
			var news News

			err := rows.Scan(
				&news.ID,
				&news.Title,
				&news.Content,
				&news.ImageURL,
				&news.CreatedAt,
			)

			if err != nil {
				log.Println("Scan error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to process news",
				})
				return
			}

			newsList = append(newsList, news)
		}

		if err := rows.Err(); err != nil {
			log.Println("Rows error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read news",
			})
			return
		}

		c.JSON(http.StatusOK, newsList)
	}
}

func main() {
	// Connect to database.
	db := connectDB()
	defer db.Close()

	// Create Gin router.
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://miku-lime-eta.vercel.app/",
		},
		AllowMethods: []string{
			"GET",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
		},
	}))

	// API routes.
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "News API is running!",
		})
	})

	router.GET("/api/news", getNews(db))

	// Start server.
	log.Println("Server running on port 8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
