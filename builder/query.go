package builder

import (
	"strings"
)

func NewQuery(str string, args []interface{}) *Query {
	q := &Query{
		args: args,
	}
	q.str.WriteString(str)
	return q
}

type Query struct {
	str  strings.Builder
	args []interface{}
}

func (q *Query) String() string {
	return q.str.String()
}

func (q *Query) Args() []interface{} {
	return q.args
}
