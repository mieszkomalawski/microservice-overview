package models

// Graph reprezentuje pełny graf (wszystkie wierzchołki i relacje)
type Graph struct {
	Vertices []Vertex `json:"vertices"`
	Edges    []Edge   `json:"edges"`
}


