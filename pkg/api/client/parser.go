package client

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"slices"
)

type CompositeExpression struct {
	Operator CompositeOperator
	Body     map[string]any
	Operands []any
}

type CompositeOperator string

const (
	And CompositeOperator = "$and"
	Or  CompositeOperator = "$or"
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
	In                 FieldOperator = "in"
	NotIn              FieldOperator = "not-in"
	ArrayContains      FieldOperator = "array-contains"
	ArrayContainsAny   FieldOperator = "array-contains-any"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

func toFirestoreDirection(direction Direction) firestore.Direction {
	if direction == Descending {
		return firestore.Desc
	}
	return firestore.Asc
}

func CreateCompositeExpression(operator CompositeOperator, body map[string]any) *CompositeExpression {
	return &CompositeExpression{
		Operator: operator,
		Body:     body,
		Operands: make([]any, 0),
	}
}

func (c *CompositeExpression) toEntityFilter() firestore.EntityFilter {
	var filter firestore.CompositeFilter
	if c.Operator == And {
		filter = &firestore.AndFilter{}
	} else if c.Operator == Or {
		filter = &firestore.OrFilter{}
	}

	for _, operand := range c.Operands {
		switch operand.(type) {
		case *CompositeExpression:
			if c.Operator == And {
				f := filter.(*firestore.AndFilter)
				f.Filters = append(f.Filters, operand.(*CompositeExpression).toEntityFilter())
			} else if c.Operator == Or {
				f := filter.(*firestore.OrFilter)
				f.Filters = append(f.Filters, operand.(*CompositeExpression).toEntityFilter())
			}
		case *FieldExpression:
			field := operand.(*FieldExpression)
			if c.Operator == And {
				f := filter.(*firestore.AndFilter)
				f.Filters = append(f.Filters, firestore.PropertyFilter{
					Path:     field.Field,
					Operator: string(field.Operator),
					Value:    field.Value,
				})
			} else if c.Operator == Or {
				f := filter.(*firestore.OrFilter)
				f.Filters = append(f.Filters, firestore.PropertyFilter{
					Path:     field.Field,
					Operator: string(field.Operator),
					Value:    field.Value,
				})
			}
		}
	}

	return filter
}

type compositeType string

const (
	compositeTypeAnd         compositeType = "and"
	compositeTypeImplicitAnd compositeType = "implicit_and"
	compositeTypeOr          compositeType = "or"
	compositeTypeUnknown     compositeType = "unknown"
)

func determineCompositeType(body map[string]any) compositeType {
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

func createRootExpression(body map[string]any) (*CompositeExpression, error) {
	switch determineCompositeType(body) {
	case compositeTypeAnd:
		return CreateCompositeExpression(And, body[string(And)].(map[string]any)), nil
	case compositeTypeImplicitAnd:
		return CreateCompositeExpression(And, body), nil
	case compositeTypeOr:
		return CreateCompositeExpression(Or, body[string(Or)].(map[string]any)), nil
	default:
		return nil, fmt.Errorf("invalid query format, see help for more information on query syntax")
	}
}

func parse(parent *CompositeExpression) {
	for k, v := range parent.Body {
		if k == string(And) {
			child := CreateCompositeExpression(And, v.(map[string]any))
			parent.Operands = append(parent.Operands, child)
			parse(child)
		} else if k == string(Or) {
			child := CreateCompositeExpression(Or, v.(map[string]any))
			parent.Operands = append(parent.Operands, child)
			parse(child)
		} else {
			parseFieldOperator(parent, k, v)
		}
	}
}

func parseFieldOperator(parent *CompositeExpression, k string, v any) {
	operator := Equal
	var value any

	switch v.(type) {
	case map[string]any:
		p := v.(map[string]any)
		if len(p) != 1 {
			fmt.Println("invalid query format, field operator requires exactly one operator; see help for more information")
			return
		}
		for pk, pv := range p {
			operator = FieldOperator(pk)
			value = pv
			break
		}
	default:
		value = v
	}

	if !slices.Contains([]FieldOperator{Equal, NotEqual, LessThan, LessThanOrEqual, GreaterThan, GreaterThanOrEqual, In, NotIn, ArrayContains, ArrayContainsAny}, operator) {
		fmt.Printf("invalid query format, unknown field operator %s; see help for more information on query syntax", operator)
		return
	}

	parent.Operands = append(parent.Operands, &FieldExpression{
		Field:    k,
		Operator: operator,
		Value:    value,
	})
}
