package query

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type Sort struct {
	Field string
	Desc  bool
}

type QueryOptions struct {
	Filters []Filter
	Sorts   []Sort
	Fields  []string
}
