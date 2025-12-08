package handlers

import (
	"net/http"

	"microservice_overview/storage"

	"github.com/gin-gonic/gin"
)

// GraphHandler obsługuje żądania związane z grafem
type GraphHandler struct {
	storage storage.Storage
}

// NewGraphHandler tworzy nowy GraphHandler
func NewGraphHandler(s storage.Storage) *GraphHandler {
	return &GraphHandler{storage: s}
}

// GetGraph zwraca pełny graf (wszystkie wierzchołki i relacje)
func (h *GraphHandler) GetGraph(c *gin.Context) {
	graph, err := h.storage.GetGraph()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, graph)
}


