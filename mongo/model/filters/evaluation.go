package filters

import (
	"github.com/qianwj/typed/mongo/model/filters/regex"
	"github.com/qianwj/typed/mongo/model/filters/text"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func Expr(expression any) *Filter {
	return New().Expr(expression)
}

func Mod(key string, divisor, remainder float64) *Filter {
	return New().Mod(key, divisor, remainder)
}

func Like(key string, matcher *regex.Matcher) *Filter {
	return New().Like(key, matcher)
}

func Text(search *text.Search) *Filter {
	return New().Text(search)
}

func Where(key, expression string) *Filter {
	return New().Where(key, expression)
}

func (f *Filter) Expr(expression any) *Filter {
	f.put(operator.Expr, expression)
	return f
}

func (f *Filter) Mod(key string, divisor, remainder float64) *Filter {
	f.put(key, bson.M{operator.Mod: []float64{divisor, remainder}})
	return f
}

func (f *Filter) Like(key string, matcher *regex.Matcher) *Filter {
	f.put(key, bson.M{operator.Regex: matcher.Compile()})
	return f
}

func (f *Filter) Text(search *text.Search) *Filter {
	f.put(operator.Text, search.Marshal())
	return f
}

func (f *Filter) Where(key, expression string) *Filter {
	f.put(key, bson.M{operator.Where: expression})
	return f
}
