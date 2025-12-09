package handlers

import (
	"net/http"

	"microservice_overview/models"
	"microservice_overview/storage"

	"github.com/gin-gonic/gin"
)

// EdgeHandler obsługuje żądania związane z relacjami
type EdgeHandler struct {
	storage storage.Storage
}

// NewEdgeHandler tworzy nowy EdgeHandler
func NewEdgeHandler(s storage.Storage) *EdgeHandler {
	return &EdgeHandler{storage: s}
}

// GetAllEdges zwraca listę wszystkich relacji
func (h *EdgeHandler) GetAllEdges(c *gin.Context) {
	edges, err := h.storage.GetAllEdges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, edges)
}

// GetEdgeByID zwraca relację po ID
func (h *EdgeHandler) GetEdgeByID(c *gin.Context) {
	id := c.Param("id")
	edge, err := h.storage.GetEdgeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "edge not found"})
		return
	}
	c.JSON(http.StatusOK, edge)
}

// CreateEdge tworzy nową relację
func (h *EdgeHandler) CreateEdge(c *gin.Context) {
	var edge models.Edge
	if err := c.ShouldBindJSON(&edge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if edge.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if edge.From == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from is required"})
		return
	}

	if edge.To == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to is required"})
		return
	}

	if err := h.storage.CreateEdge(&edge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, edge)
}

// UpdateEdge aktualizuje istniejącą relację
func (h *EdgeHandler) UpdateEdge(c *gin.Context) {
	id := c.Param("id")

	var edge models.Edge
	if err := c.ShouldBindJSON(&edge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	edge.ID = id

	if err := h.storage.UpdateEdge(&edge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, edge)
}

// DeleteEdge usuwa relację
func (h *EdgeHandler) DeleteEdge(c *gin.Context) {
	id := c.Param("id")

	if err := h.storage.DeleteEdge(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "edge deleted"})
}



