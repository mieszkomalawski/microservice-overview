package vertex_integration_test

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
	vertexHandler := handlers.NewVertexHandler(s)

	api := r.Group("/api")
	{
		api.GET("/vertices", vertexHandler.GetAllVertices)
		api.GET("/vertices/:id", vertexHandler.GetVertexByID)
		api.POST("/vertices", vertexHandler.CreateVertex)
		api.PUT("/vertices/:id", vertexHandler.UpdateVertex)
		api.DELETE("/vertices/:id", vertexHandler.DeleteVertex)
	}

	return r, s
}

func TestCreateVertex_Integration(t *testing.T) {
	r, _ := setupTestRouter()

	vertex := models.Vertex{
		ID:          "test-vertex",
		Name:        "Test Vertex",
		Description: "Test description",
	}

	jsonValue, _ := json.Marshal(vertex)
	req, _ := http.NewRequest("POST", "/api/vertices", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.Vertex
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != vertex.ID {
		t.Errorf("Expected ID %s, got %s", vertex.ID, response.ID)
	}
	if response.Name != vertex.Name {
		t.Errorf("Expected Name %s, got %s", vertex.Name, response.Name)
	}
}

func TestGetAllVertices_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz testowe wierzchołki
	vertex1 := &models.Vertex{ID: "v1", Name: "Vertex 1"}
	vertex2 := &models.Vertex{ID: "v2", Name: "Vertex 2"}
	s.CreateVertex(vertex1)
	s.CreateVertex(vertex2)

	req, _ := http.NewRequest("GET", "/api/vertices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var vertices []models.Vertex
	json.Unmarshal(w.Body.Bytes(), &vertices)

	if len(vertices) != 2 {
		t.Errorf("Expected 2 vertices, got %d", len(vertices))
	}
}

func TestGetVertexByID_Integration(t *testing.T) {
	r, s := setupTestRouter()

	vertex := &models.Vertex{ID: "test-id", Name: "Test Name"}
	s.CreateVertex(vertex)

	req, _ := http.NewRequest("GET", "/api/vertices/test-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Vertex
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != vertex.ID {
		t.Errorf("Expected ID %s, got %s", vertex.ID, response.ID)
	}
}

func TestGetVertexByID_NotFound_Integration(t *testing.T) {
	r, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/vertices/non-existent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateVertex_Integration(t *testing.T) {
	r, s := setupTestRouter()

	// Utwórz wierzchołek
	vertex := &models.Vertex{ID: "update-test", Name: "Original Name"}
	s.CreateVertex(vertex)

	// Aktualizuj
	updated := models.Vertex{
		Name:        "Updated Name",
		Description: "Updated description",
	}

	jsonValue, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", "/api/vertices/update-test", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Vertex
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != updated.Name {
		t.Errorf("Expected Name %s, got %s", updated.Name, response.Name)
	}
}

func TestDeleteVertex_Integration(t *testing.T) {
	r, s := setupTestRouter()

	vertex := &models.Vertex{ID: "delete-test", Name: "To Delete"}
	s.CreateVertex(vertex)

	req, _ := http.NewRequest("DELETE", "/api/vertices/delete-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Sprawdź czy został usunięty
	_, err := s.GetVertexByID("delete-test")
	if err == nil {
		t.Error("Vertex should be deleted")
	}
}

func TestCreateVertex_Validation_Integration(t *testing.T) {
	r, _ := setupTestRouter()

	tests := []struct {
		name           string
		vertex         models.Vertex
		expectedStatus int
	}{
		{
			name:           "missing ID",
			vertex:         models.Vertex{Name: "Test"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing Name",
			vertex:         models.Vertex{ID: "test-id"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.vertex)
			req, _ := http.NewRequest("POST", "/api/vertices", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

