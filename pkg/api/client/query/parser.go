package query

import (
	"fmt"
	"slices"
)

func parse(parent *Expression) {
	for k, v := range parent.Body {
		if k == string(And) {
			child := create(And, v.(map[string]any))
			parent.Operands = append(parent.Operands, child)
			parse(child)
		} else if k == string(Or) {
			child := create(Or, v.(map[string]any))
			parent.Operands = append(parent.Operands, child)
			parse(child)
		} else {
			parseField(parent, k, v)
		}
	}
}

func parseField(parent *Expression, k string, v any) {
	operator := Equal
	var value any

	switch v.(type) {
	case map[string]any:
		p := v.(map[string]any)
		if len(p) != 1 {
			fmt.Println("field operator requires exactly one operator; see help for more information")
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
		fmt.Printf("unknown field operator %s; see help for more information on query syntax", operator)
		return
	}

	parent.Operands = append(parent.Operands, &FieldExpression{
		Field:    k,
		Operator: operator,
		Value:    value,
	})
}
