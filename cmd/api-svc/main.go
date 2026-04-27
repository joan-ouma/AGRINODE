package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// DailyAverage represents the exact data structure from our Week 6 PostgreSQL View
type DailyAverage struct {
	NodeName           string  `json:"node_name"`
	Zone               string  `json:"zone"`
	ReadingDate        string  `json:"reading_date"`
	AvgTemp            float32 `json:"avg_temp"`
	AvgHumidity        float32 `json:"avg_humidity"`
	AvgMoisture        float32 `json:"avg_moisture"`
	TotalDailyReadings int     `json:"total_daily_readings"`
}

var db *sql.DB

func main() {
	var err error

	// 1. Connect to PostgreSQL
	connStr := "postgres://agrinode_admin:supersecretpassword@agrinode-postgres:5432/agrinode?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open db connection: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping db: ", err)
	}
	log.Println("api service connected to postgresql...")

	// 2. Initialize the Gin Router
	// Gin automatically handles request logging and crash recovery
	router := gin.Default()

	// 3. Define our API Routes
	// We group them under /api/v1 so we can easily upgrade to v2 later without breaking older dashboards
	v1 := router.Group("/api/v1")
	{
		v1.GET("/analytics/daily", getDailyAverages)
	}

	// 4. Start the Server on Port 8080
	log.Println("starting api gateway on :8080...")
	router.Run(":8080")
}

// getDailyAverages handles the incoming web requests and queries the database View
func getDailyAverages(c *gin.Context) {
	// Query the heavily optimized View we built in Week 6
	query := `
		SELECT node_name, zone, CAST(reading_date AS TEXT), avg_temp, avg_humidity, avg_moisture, total_daily_readings 
		FROM daily_node_averages 
		ORDER BY reading_date DESC 
		LIMIT 14
	`
	rows, err := db.Query(query)
	if err != nil {
		// If the database fails, return a secure 500 error instead of crashing
		log.Printf("db query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch analytics data"})
		return
	}
	defer rows.Close()

	var results []DailyAverage

	// Loop through the database rows and map them to our Go struct
	for rows.Next() {
		var d DailyAverage
		if err := rows.Scan(&d.NodeName, &d.Zone, &d.ReadingDate, &d.AvgTemp, &d.AvgHumidity, &d.AvgMoisture, &d.TotalDailyReadings); err != nil {
			log.Printf("row scan error: %v", err)
			continue
		}
		results = append(results, d)
	}

	// Send a 200 OK HTTP response with the fully formatted JSON payload
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(results),
		"data":   results,
	})
}
