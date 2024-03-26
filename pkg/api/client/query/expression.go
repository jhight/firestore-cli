package query

import (
	"fmt"
)

type Expression struct {
	Operator LogicOperator
	Body     map[string]any
	Operands []any
}

type LogicOperator string

const (
	And LogicOperator = "$and"
	Or  LogicOperator = "$or"
)

type FieldExpression struct {
	Field    string
	Operator FieldOperator
	Value    any
}

type FieldOperator string

const (
	Equal              FieldOperator = "=="
	NotEqual           FieldOperator = "!="
	GreaterThan        FieldOperator = ">"
	GreaterThanOrEqual FieldOperator = ">="
	LessThan           FieldOperator = "<"
	LessThanOrEqual    FieldOperator = "<="
	In                 FieldOperator = "$in"
	NotIn              FieldOperator = "$not-in"
	ArrayContains      FieldOperator = "$array-contains"
	ArrayContainsAny   FieldOperator = "$array-contains-any"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

const (
	FunctionTimestamp string = "$timestamp"
	FunctionNow       string = "$now"
)

func CreateExpression(body map[string]any) (*Expression, error) {
	var e *Expression
	switch determineType(body) {
	case compositeTypeAnd:
		e = create(And, body[string(And)].(map[string]any))
	case compositeTypeImplicitAnd:
		e = create(And, body)
	case compositeTypeOr:
		e = create(Or, body[string(Or)].(map[string]any))
	default:
		return nil, fmt.Errorf("invalid query format, see help for more information on query syntax")
	}

	parse(e)
	return e, nil
}

func determineType(body map[string]any) Type {
	for k, _ := range body {
		switch {
		case k == string(And):
			return compositeTypeAnd
		case k == string(Or):
			return compositeTypeOr
		default:
			return compositeTypeImplicitAnd
		}
	}

	return compositeTypeUnknown
}

func create(operator LogicOperator, body map[string]any) *Expression {
	return &Expression{
		Operator: operator,
		Body:     body,
		Operands: make([]any, 0),
	}
}

type Type string

const (
	compositeTypeAnd         Type = "and"
	compositeTypeImplicitAnd Type = "implicit_and"
	compositeTypeOr          Type = "or"
	compositeTypeUnknown     Type = "unknown"
)

const (
	SelectionDocumentID   string = "$id"
	SelectionDocumentPath string = "$path"
)
