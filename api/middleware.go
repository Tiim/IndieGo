package api

// request_logger.go
import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func trailingSlash(e *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path != "/" && strings.HasSuffix(path, "/") {
			c.Request.URL.Path = strings.TrimSuffix(path, "/")
			e.HandleContext(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			logger.Printf("Error: %v", c.Errors)
			status := c.Writer.Status()
			if status == 0 || status < 400 {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"status": status, "error": c.Errors})
		}
	}
}

var (
	durations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Histogram request duration in seconds",
		Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 5, 10},
	},
		[]string{},
	)
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_hits_total",
		Help: "Total number of requests",
	},
		[]string{"path"},
	)
)

func metricsMiddleware(reg *prometheus.Registry) gin.HandlerFunc {

	reg.MustRegister(durations)
	reg.MustRegister(hits)

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)
		durations.WithLabelValues().Observe(elapsed)
		hits.WithLabelValues(c.Request.URL.Path).Inc()

	}
}
