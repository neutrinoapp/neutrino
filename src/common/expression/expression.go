package expression

import (
	"strings"

	"net/url"
)

const (
	EXPRESSION_OP_EQUALS = "equals"

	EXPRESSION_PLACEHOLDER = "{exp}"
	EXPRESSION_FILTER      = "filter"

	EXPRESSION_REGEXP = `((\?|\&)` + EXPRESSION_PLACEHOLDER + `)=\"(?P<` + EXPRESSION_PLACEHOLDER + `>.*?)\"`
)

type ExpressionParam struct {
	Left, Right, Op string
}

type Expression struct {
	Params []ExpressionParam
}

//TODO: sort, group etc.
type ExpressionGroup struct {
	Filter Expression
}

func ParseExpressionGroup(query url.Values) (ExpressionGroup, error) {
	g := ExpressionGroup{
		Filter: Expression{},
	}

	filterParams := query.Get(EXPRESSION_FILTER)
	if filterParams != "" {
		filterParams = strings.Replace(filterParams, `"`, ``, -1)
		splitPairs := strings.Split(filterParams, ",")
		for i := range splitPairs {
			pair := splitPairs[i]
			splitPair := strings.Split(pair, "=")
			leftParam := splitPair[0]
			rightParam := splitPair[1]
			g.Filter.Params = append(g.Filter.Params, ExpressionParam{
				Left:  leftParam,
				Right: rightParam,
				Op:    EXPRESSION_OP_EQUALS,
			})
		}
	}

	return g, nil
}
