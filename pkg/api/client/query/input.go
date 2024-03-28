package query

type Input struct {
	Path    string
	Fields  []string
	Filter  map[string]any
	OrderBy []OrderBy
	Limit   int
	Offset  int
	Count   bool
}

type OrderBy struct {
	Field     string
	Direction Direction
}
