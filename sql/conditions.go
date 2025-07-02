package sql

//functions and structures for SQL conditions

import (
	"fmt"
	"strings"
)

const (
	SGN_SQL_E       SQLCondition = "="
	SGN_SQL_L       SQLCondition = "<"
	SGN_SQL_G       SQLCondition = ">"
	SGN_SQL_LE      SQLCondition = "<="
	SGN_SQL_GE      SQLCondition = ">="
	SGN_SQL_LK      SQLCondition = "LIKE"
	SGN_SQL_NE      SQLCondition = "<>"
	SGN_SQL_I       SQLCondition = "IS"
	SGN_SQL_IN      SQLCondition = "IS NOT"
	SGN_SQL_INCL    SQLCondition = "IN"
	SGN_SQL_ANY     SQLCondition = "=ANY"
	SGN_SQL_OVERLAP SQLCondition = "&&"
)

type SQLCondition string

// FilterCond is an SQL condition for filter
// FieldID - field id
// Value - field value
// If Expression is set, FieldID/Value can be empty. It is a validated && sanatized sql expression.
type FilterCond struct {
	FieldID    string
	Value      interface{}
	Sign       SQLCondition
	InsCase    bool
	Expression string //validated,sanatized expression
}

type FilterCondCollection []*FilterCond

// adds where expression to sql string
//
//	field - sql field id
//	sgn - signe of sql.SQLCondition
//	ic - bool, true if insensetive case
//	cond - SQL bind condition AND/OR
//	sql - string to add expression to
func AddCondExpr(field string, sgn SQLCondition, ic bool, condInd int, parInd int, cond string, sql *strings.Builder) {
	if condInd > 0 {
		sql.WriteString(" " + cond + " ")
	}
	if ic {
		sql.WriteString(fmt.Sprintf("(lower(%s::text) %s lower($%d::text))", field, sgn, parInd+1))

	} else {
		sql.WriteString(fmt.Sprintf("(%s %s $%d)", field, sgn, parInd+1))
	}
}
