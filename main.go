package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ch-go-poc/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	r := gin.Default()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	dsn := "clickhouse://clickhouse:clickhouse@localhost:19000/clickhouse?dial_timeout=10s&read_timeout=20s"
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&otel.OtelLogRecord{})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
