package edge_integration_test

import (
	"bytes"
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
	edgeHandler := handlers.NewEdgeHandler(s)

	api := r.Group("/api")
	{
		api.GET("/edges", edgeHandler.GetAllEdges)
		api.GET("/edges/:id", edgeHandler.GetEdgeByID)
		api.POST("/edges", edgeHandler.CreateEdge)
		api.PUT("/edges/:id", edgeHandler.UpdateEdge)
		api.DELETE("/edges/:id", edgeHandler.DeleteEdge)
	}

	return r, s
}

func TestCreateEdge_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołki (wymagane dla edge)
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)

	edge := models.Edge{
		ID:   "edge-1",
		From: "v1",
		To:   "v2",
		Type: "calls",
	}

	jsonValue, _ := json.Marshal(edge)
	req, _ := http.NewRequest("POST", "/api/edges", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response models.Edge
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != edge.ID {
		t.Errorf("Expected ID %s, got %s", edge.ID, response.ID)
	}
}

func TestCreateEdge_WithParentVertex_ShouldFail_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołek rodzica
	parent := &models.Vertex{ID: "parent", Name: "Parent"}
	s.CreateVertex(parent)

	// Utwórz wierzchołek potomny
	child := &models.Vertex{ID: "child", Name: "Child", ParentID: stringPtr("parent")}
	s.CreateVertex(child)

	// Próba utworzenia edge z wierzchołkiem mającym dzieci powinna się nie powieść
	edge := models.Edge{
		ID:   "edge-1",
		From: "parent",
		To:   "child",
		Type: "calls",
	}

	jsonValue, _ := json.Marshal(edge)
	req, _ := http.NewRequest("POST", "/api/edges", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestGetAllEdges_Integration(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/api/edges", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var edges []models.Edge
	json.Unmarshal(w.Body.Bytes(), &edges)

	if len(edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(edges))
	}
}

func TestGetEdgeByID_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołki
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)

	edge := &models.Edge{ID: "test-edge", From: "v1", To: "v2", Type: "calls"}
	s.CreateEdge(edge)

	req, _ := http.NewRequest("GET", "/api/edges/test-edge", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Edge
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != edge.ID {
		t.Errorf("Expected ID %s, got %s", edge.ID, response.ID)
	}
}

func TestUpdateEdge_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołki
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	vertex3 := &models.Vertex{ID: "v3", Name: "Vertex 3"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)
	s.CreateVertex(vertex3)

	// Utwórz edge
	edge := &models.Edge{ID: "update-edge", From: "v1", To: "v2", Type: "calls"}
	s.CreateEdge(edge)

	// Aktualizuj
	updated := models.Edge{
		From: "v1",
		To:   "v3",
		Type: "requires",
	}

	jsonValue, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", "/api/edges/update-edge", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response models.Edge
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.To != updated.To {
		t.Errorf("Expected To %s, got %s", updated.To, response.To)
	}
}

func TestDeleteEdge_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołki
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)

	edge := &models.Edge{ID: "delete-edge", From: "v1", To: "v2", Type: "calls"}
	s.CreateEdge(edge)

	req, _ := http.NewRequest("DELETE", "/api/edges/delete-edge", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Sprawdź czy został usunięty
	_, err := s.GetEdgeByID("delete-edge")
	if err == nil {
		t.Error("Edge should be deleted")
	}
}

func TestCreateEdge_Validation_Integration(t *testing.T) {
	r, _ := setupTestRouter()

	tests := []struct {
		name           string
		edge           models.Edge
		expectedStatus int
	}{
		{
			name:           "missing ID",
			edge:           models.Edge{From: "v1", To: "v2"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing From",
			edge:           models.Edge{ID: "e1", To: "v2"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing To",
			edge:           models.Edge{ID: "e1", From: "v1"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.edge)
			req, _ := http.NewRequest("POST", "/api/edges", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

