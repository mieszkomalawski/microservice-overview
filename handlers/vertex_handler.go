package handlers

import (
	"net/http"

	"microservice_overview/models"
	"microservice_overview/storage"

	"github.com/gin-gonic/gin"
)

// VertexHandler obsługuje żądania związane z wierzchołkami
type VertexHandler struct {
	storage storage.Storage
}

// NewVertexHandler tworzy nowy VertexHandler
func NewVertexHandler(s storage.Storage) *VertexHandler {
	return &VertexHandler{storage: s}
}

// GetAllVertices zwraca listę wszystkich wierzchołków
func (h *VertexHandler) GetAllVertices(c *gin.Context) {
	vertices, err := h.storage.GetAllVertices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vertices)
}

// GetVertexByID zwraca wierzchołek po ID
func (h *VertexHandler) GetVertexByID(c *gin.Context) {
	id := c.Param("id")
	vertex, err := h.storage.GetVertexByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vertex not found"})
		return
	}
	c.JSON(http.StatusOK, vertex)
}

// CreateVertex tworzy nowy wierzchołek
func (h *VertexHandler) CreateVertex(c *gin.Context) {
	var vertex models.Vertex
	if err := c.ShouldBindJSON(&vertex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vertex.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if vertex.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if err := h.storage.CreateVertex(&vertex); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vertex)
}

// UpdateVertex aktualizuje istniejący wierzchołek
func (h *VertexHandler) UpdateVertex(c *gin.Context) {
	id := c.Param("id")

	var vertex models.Vertex
	if err := c.ShouldBindJSON(&vertex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vertex.ID = id

	if err := h.storage.UpdateVertex(&vertex); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vertex)
}

// DeleteVertex usuwa wierzchołek
func (h *VertexHandler) DeleteVertex(c *gin.Context) {
	id := c.Param("id")

	if err := h.storage.DeleteVertex(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vertex deleted"})
}



