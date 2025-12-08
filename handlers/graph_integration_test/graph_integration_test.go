package graph_integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"microservice_overview/handlers"
	"microservice_overview/models"
	"microservice_overview/storage"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, storage.Storage) {
	gin.SetMode(gin.TestMode)

	// Ustaw tryb developerski dla testów
	os.Setenv("DEV_MODE", "true")

	// Utwórz storage z bazą w pamięci
	s, err := storage.NewStorage()
	if err != nil {
		os.Unsetenv("DEV_MODE")
		panic("failed to create storage: " + err.Error())
	}

	// Utwórz router
	r := gin.New()
	graphHandler := handlers.NewGraphHandler(s)

	api := r.Group("/api")
	{
		api.GET("/graph", graphHandler.GetGraph)
	}

	return r, s
}

func TestGetGraph_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołki
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	vertex3 := &models.Vertex{ID: "v3", Name: "Vertex 3"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)
	s.CreateVertex(vertex3)

	// Utwórz relacje
	edge1 := &models.Edge{ID: "e1", From: "v1", To: "v2", Type: "calls"}
	edge2 := &models.Edge{ID: "e2", From: "v2", To: "v3", Type: "requires"}
	s.CreateEdge(edge1)
	s.CreateEdge(edge2)

	req, _ := http.NewRequest("GET", "/api/graph", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var graph models.Graph
	json.Unmarshal(w.Body.Bytes(), &graph)

	if len(graph.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(graph.Vertices))
	}

	if len(graph.Edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(graph.Edges))
	}
}

func TestGetGraph_Empty_Integration(t *testing.T) {
	r, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/graph", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var graph models.Graph
	json.Unmarshal(w.Body.Bytes(), &graph)

	if len(graph.Vertices) != 0 {
		t.Errorf("Expected 0 vertices, got %d", len(graph.Vertices))
	}

	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(graph.Edges))
	}
}

func TestGetGraph_WithHierarchy_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołek rodzica
	parent := &models.Vertex{ID: "parent", Name: "Parent"}
	s.CreateVertex(parent)

	// Utwórz wierzchołki potomne
	child1 := &models.Vertex{ID: "child1", Name: "Child 1", ParentID: stringPtr("parent")}
	child2 := &models.Vertex{ID: "child2", Name: "Child 2", ParentID: stringPtr("parent")}
	s.CreateVertex(child1)
	s.CreateVertex(child2)

	// Utwórz relację tylko między wierzchołkami bez dzieci
	edge := &models.Edge{ID: "e1", From: "child1", To: "child2", Type: "calls"}
	s.CreateEdge(edge)

	req, _ := http.NewRequest("GET", "/api/graph", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var graph models.Graph
	json.Unmarshal(w.Body.Bytes(), &graph)

	if len(graph.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(graph.Vertices))
	}

	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}
}

func stringPtr(s string) *string {
	return &s
}

