package query

import (
	"cloud.google.com/go/firestore"
	"strings"
)

func (c *Expression) FirestoreFilter() firestore.EntityFilter {
	var filter firestore.CompositeFilter
	if c.Operator == And {
		filter = &firestore.AndFilter{}
	} else if c.Operator == Or {
		filter = &firestore.OrFilter{}
	}

	for _, operand := range c.Operands {
		switch operand.(type) {
		case *Expression:
			if c.Operator == And {
				f := filter.(*firestore.AndFilter)
				f.Filters = append(f.Filters, operand.(*Expression).FirestoreFilter())
			} else if c.Operator == Or {
				f := filter.(*firestore.OrFilter)
				f.Filters = append(f.Filters, operand.(*Expression).FirestoreFilter())
			}
		case *FieldExpression:
			field := operand.(*FieldExpression)
			if c.Operator == And {
				f := filter.(*firestore.AndFilter)
				f.Filters = append(f.Filters, firestore.PropertyFilter{
					Path:     field.Field,
					Operator: strings.TrimPrefix(string(field.Operator), "$"),
					Value:    field.Value,
				})
			} else if c.Operator == Or {
				f := filter.(*firestore.OrFilter)
				f.Filters = append(f.Filters, firestore.PropertyFilter{
					Path:     field.Field,
					Operator: strings.TrimPrefix(string(field.Operator), "$"),
					Value:    field.Value,
				})
			}
		}
	}

	return filter
}

func (d Direction) FirestoreDirection() firestore.Direction {
	if d == Descending {
		return firestore.Desc
	}
	return firestore.Asc
}
