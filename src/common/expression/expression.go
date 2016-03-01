package expression

import (
	"strings"

	"github.com/neutrinoapp/neutrino/src/common/models"

	"net/url"
	"strconv"
)

const (
	EXPRESSION_OP_EQUALS = "$eq"

	EXPRESSION_FILTER = "filter"
)

type ExpressionParam struct {
	Left, Op string
	Right    interface{}
}

type Expression struct {
	Params []ExpressionParam
}

//TODO: sort, group etc.
type ExpressionGroup struct {
	Filter Expression
}

func (e ExpressionGroup) ToMongo() (models.JSON, error) {
	query := models.JSON{}
	and := make([]models.JSON, 0)

	for _, p := range e.Filter.Params {
		if p.Op == EXPRESSION_OP_EQUALS {
			query := models.JSON{}
			query[p.Left] = p.Right
			and = append(and, query)
		}
	}

	if len(and) > 0 {
		query["$and"] = and
	}

	return query, nil
}

func ParseExpressionGroup(query url.Values) (ExpressionGroup, error) {
	g := ExpressionGroup{
		Filter: Expression{},
	}

	filterParams := query.Get(EXPRESSION_FILTER)
	if filterParams != "" {
		filterParams = strings.Replace(filterParams, `"`, ``, -1)
		splitPairs := strings.Split(filterParams, ",")
		for _, pair := range splitPairs {
			splitPair := strings.Split(pair, "=")
			leftParam := splitPair[0]
			rightParam := splitPair[1]

			var rightVal interface{}
			if asInt, err := strconv.ParseInt(rightParam, 10, 64); err == nil {
				rightVal = asInt
			} else if asFloat, err := strconv.ParseFloat(rightParam, 10); err == nil {
				rightVal = asFloat
			} else if asBool, err := strconv.ParseBool(rightParam); err == nil {
				rightVal = asBool
			} else {
				rightVal = rightParam
			}

			g.Filter.Params = append(g.Filter.Params, ExpressionParam{
				Left:  leftParam,
				Right: rightVal,
				Op:    EXPRESSION_OP_EQUALS,
			})
		}
	}

	return g, nil
}
