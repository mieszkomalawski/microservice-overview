package main

import (
	"log"

	"microservice_overview/handlers"
	"microservice_overview/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicjalizacja storage
	s, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Inicjalizacja routera
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Statyczne pliki (frontend)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	// Frontend - strona główna
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Inicjalizacja handlerów
	vertexHandler := handlers.NewVertexHandler(s)
	edgeHandler := handlers.NewEdgeHandler(s)
	graphHandler := handlers.NewGraphHandler(s)

	// API routes
	api := r.Group("/api")
	{
		// Wierzchołki
		api.GET("/vertices", vertexHandler.GetAllVertices)
		api.GET("/vertices/:id", vertexHandler.GetVertexByID)
		api.POST("/vertices", vertexHandler.CreateVertex)
		api.PUT("/vertices/:id", vertexHandler.UpdateVertex)
		api.DELETE("/vertices/:id", vertexHandler.DeleteVertex)

		// Relacje
		api.GET("/edges", edgeHandler.GetAllEdges)
		api.GET("/edges/:id", edgeHandler.GetEdgeByID)
		api.POST("/edges", edgeHandler.CreateEdge)
		api.PUT("/edges/:id", edgeHandler.UpdateEdge)
		api.DELETE("/edges/:id", edgeHandler.DeleteEdge)

		// Graf
		api.GET("/graph", graphHandler.GetGraph)
	}

	// Uruchomienie serwera
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}



