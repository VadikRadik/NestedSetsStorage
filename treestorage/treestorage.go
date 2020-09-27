package treestorage

// NestedSetsNode is a tree node
type NestedSetsNode struct {
	Name  string
	Left  int
	Right int
}

// NestedSetsStorage is an interface for data base table
type NestedSetsStorage struct {
	DbConnectionString string
	DbDriver           string
}

// GetParents returns parents for the node name
func (s *NestedSetsStorage) GetParents(name string) []NestedSetsNode {
	return []NestedSetsNode{}
}

// GetChildren returns children for the node name
func (s *NestedSetsStorage) GetChildren(name string) []NestedSetsNode {
	return []NestedSetsNode{}
}

// AddNode adds new child node with name name for parent node with name parent
func (s *NestedSetsStorage) AddNode(name string, parent string) {

}

// RemoveNode removes node with name name with all its children
func (s *NestedSetsStorage) RemoveNode(name string) {

}

// MoveNode moves node with name name with all its children
func (s *NestedSetsStorage) MoveNode(name string, newParent string) {

}

// RenameNode renames node with name name
func (s *NestedSetsStorage) RenameNode(name string, newName string) {

}
