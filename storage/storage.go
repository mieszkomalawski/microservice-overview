package storage

import (
	"fmt"
	"os"

	"microservice_overview/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Storage interfejs do operacji na danych
type Storage interface {
	// Wierzchołki
	GetAllVertices() ([]models.Vertex, error)
	GetVertexByID(id string) (*models.Vertex, error)
	CreateVertex(vertex *models.Vertex) error
	UpdateVertex(vertex *models.Vertex) error
	DeleteVertex(id string) error
	HasChildren(vertexID string) (bool, error)  // Sprawdza czy wierzchołek ma dzieci
	IsLeafVertex(vertexID string) (bool, error) // Sprawdza czy wierzchołek jest na najniższym poziomie

	// Relacje
	GetAllEdges() ([]models.Edge, error)
	GetEdgeByID(id string) (*models.Edge, error)
	CreateEdge(edge *models.Edge) error
	UpdateEdge(edge *models.Edge) error
	DeleteEdge(id string) error

	// Graf
	GetGraph() (*models.Graph, error)
}

// DBStorage implementacja Storage używająca GORM
type DBStorage struct {
	db *gorm.DB
}

// NewStorage tworzy nową instancję Storage
func NewStorage() (Storage, error) {
	devMode := os.Getenv("DEV_MODE") == "true"

	var db *gorm.DB
	var err error

	if devMode {
		// Tryb developerski - SQLite w pamięci
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to in-memory database: %w", err)
		}
	} else {
		// Tryb produkcyjny - PostgreSQL
		dsn := buildPostgresDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
	}

	// Automatyczna migracja schematu
	err = db.AutoMigrate(&models.Vertex{}, &models.Edge{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DBStorage{db: db}, nil
}

// buildPostgresDSN buduje connection string dla PostgreSQL
func buildPostgresDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "microservice_overview")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// getEnv pobiera zmienną środowiskową lub zwraca wartość domyślną
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Wierzchołki

func (s *DBStorage) GetAllVertices() ([]models.Vertex, error) {
	var vertices []models.Vertex
	err := s.db.Find(&vertices).Error
	return vertices, err
}

func (s *DBStorage) GetVertexByID(id string) (*models.Vertex, error) {
	var vertex models.Vertex
	err := s.db.First(&vertex, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &vertex, nil
}

func (s *DBStorage) CreateVertex(vertex *models.Vertex) error {
	// Walidacja: jeśli ParentID jest ustawione, sprawdź czy rodzic istnieje
	if vertex.ParentID != nil && *vertex.ParentID != "" {
		var parent models.Vertex
		if err := s.db.First(&parent, "id = ?", *vertex.ParentID).Error; err != nil {
			return fmt.Errorf("parent vertex not found: %w", err)
		}
		// Walidacja: sprawdź czy nie tworzymy cyklu (rodzic nie może być potomkiem tego wierzchołka)
		if err := s.validateNoCycle(*vertex.ParentID, vertex.ID); err != nil {
			return err
		}
	}
	return s.db.Create(vertex).Error
}

// validateNoCycle sprawdza czy dodanie parentID nie tworzy cyklu
func (s *DBStorage) validateNoCycle(parentID, currentID string) error {
	// Sprawdź czy parentID nie jest potomkiem currentID
	visited := make(map[string]bool)
	return s.checkCycle(parentID, currentID, visited)
}

// checkCycle rekurencyjnie sprawdza czy istnieje ścieżka od start do target
func (s *DBStorage) checkCycle(start, target string, visited map[string]bool) error {
	if start == target {
		return fmt.Errorf("cannot create cycle: vertex %s would be ancestor of itself", target)
	}
	if visited[start] {
		return nil // Już sprawdziliśmy tę gałąź
	}
	visited[start] = true

	var vertex models.Vertex
	if err := s.db.First(&vertex, "id = ?", start).Error; err != nil {
		return nil // Wierzchołek nie istnieje, nie ma cyklu
	}

	if vertex.ParentID != nil && *vertex.ParentID != "" {
		return s.checkCycle(*vertex.ParentID, target, visited)
	}

	return nil
}

func (s *DBStorage) UpdateVertex(vertex *models.Vertex) error {
	// Walidacja: jeśli ParentID jest ustawione, sprawdź czy rodzic istnieje
	if vertex.ParentID != nil && *vertex.ParentID != "" {
		var parent models.Vertex
		if err := s.db.First(&parent, "id = ?", *vertex.ParentID).Error; err != nil {
			return fmt.Errorf("parent vertex not found: %w", err)
		}
		// Walidacja: sprawdź czy nie tworzymy cyklu
		if err := s.validateNoCycle(*vertex.ParentID, vertex.ID); err != nil {
			return err
		}
	}
	return s.db.Save(vertex).Error
}

func (s *DBStorage) DeleteVertex(id string) error {
	return s.db.Delete(&models.Vertex{}, "id = ?", id).Error
}

// HasChildren sprawdza czy wierzchołek ma dzieci
func (s *DBStorage) HasChildren(vertexID string) (bool, error) {
	var count int64
	err := s.db.Model(&models.Vertex{}).Where("parent_id = ?", vertexID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsLeafVertex sprawdza czy wierzchołek jest na najniższym poziomie (nie ma dzieci)
func (s *DBStorage) IsLeafVertex(vertexID string) (bool, error) {
	hasChildren, err := s.HasChildren(vertexID)
	if err != nil {
		return false, err
	}
	return !hasChildren, nil
}

// Relacje

func (s *DBStorage) GetAllEdges() ([]models.Edge, error) {
	var edges []models.Edge
	err := s.db.Find(&edges).Error
	return edges, err
}

func (s *DBStorage) GetEdgeByID(id string) (*models.Edge, error) {
	var edge models.Edge
	err := s.db.First(&edge, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &edge, nil
}

func (s *DBStorage) CreateEdge(edge *models.Edge) error {
	// Sprawdź czy wierzchołki istnieją
	var fromVertex, toVertex models.Vertex
	if err := s.db.First(&fromVertex, "id = ?", edge.From).Error; err != nil {
		return fmt.Errorf("source vertex not found: %w", err)
	}
	if err := s.db.First(&toVertex, "id = ?", edge.To).Error; err != nil {
		return fmt.Errorf("target vertex not found: %w", err)
	}

	// Sprawdź czy oba wierzchołki są na najniższym poziomie (nie mają dzieci)
	fromIsLeaf, err := s.IsLeafVertex(edge.From)
	if err != nil {
		return fmt.Errorf("failed to check if source vertex is leaf: %w", err)
	}
	if !fromIsLeaf {
		return fmt.Errorf("source vertex %s has children - edges can only be created between leaf vertices", edge.From)
	}

	toIsLeaf, err := s.IsLeafVertex(edge.To)
	if err != nil {
		return fmt.Errorf("failed to check if target vertex is leaf: %w", err)
	}
	if !toIsLeaf {
		return fmt.Errorf("target vertex %s has children - edges can only be created between leaf vertices", edge.To)
	}

	return s.db.Create(edge).Error
}

func (s *DBStorage) UpdateEdge(edge *models.Edge) error {
	// Sprawdź czy wierzchołki istnieją
	var fromVertex, toVertex models.Vertex
	if err := s.db.First(&fromVertex, "id = ?", edge.From).Error; err != nil {
		return fmt.Errorf("source vertex not found: %w", err)
	}
	if err := s.db.First(&toVertex, "id = ?", edge.To).Error; err != nil {
		return fmt.Errorf("target vertex not found: %w", err)
	}

	// Sprawdź czy oba wierzchołki są na najniższym poziomie (nie mają dzieci)
	fromIsLeaf, err := s.IsLeafVertex(edge.From)
	if err != nil {
		return fmt.Errorf("failed to check if source vertex is leaf: %w", err)
	}
	if !fromIsLeaf {
		return fmt.Errorf("source vertex %s has children - edges can only be created between leaf vertices", edge.From)
	}

	toIsLeaf, err := s.IsLeafVertex(edge.To)
	if err != nil {
		return fmt.Errorf("failed to check if target vertex is leaf: %w", err)
	}
	if !toIsLeaf {
		return fmt.Errorf("target vertex %s has children - edges can only be created between leaf vertices", edge.To)
	}

	return s.db.Save(edge).Error
}

func (s *DBStorage) DeleteEdge(id string) error {
	return s.db.Delete(&models.Edge{}, "id = ?", id).Error
}

// Graf

func (s *DBStorage) GetGraph() (*models.Graph, error) {
	vertices, err := s.GetAllVertices()
	if err != nil {
		return nil, err
	}

	edges, err := s.GetAllEdges()
	if err != nil {
		return nil, err
	}

	return &models.Graph{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}
