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
