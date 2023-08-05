package text

import (
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Langeuage string

const (
	Danish  Langeuage = "da"
	Dutch   Langeuage = "nl"
	English Langeuage = "en"
)

type Search struct {
	search             string
	language           *Langeuage
	caseSensitive      *bool
	diacriticSensitive *bool
}

func New(search string) *Search {
	return &Search{search: search}
}

func (s *Search) Language(lang Langeuage) *Search {
	s.language = &lang
	return s
}

func (s *Search) CaseSensitive() *Search {
	sens := true
	s.caseSensitive = &sens
	return s
}

func (s *Search) DiacriticSensitive() *Search {
	sens := true
	s.diacriticSensitive = &sens
	return s
}

func (s *Search) Marshal() bson.D {
	res := bson.D{
		{Key: operator.Search, Value: s.search},
	}
	if s.language != nil {
		res = append(res, bson.E{
			Key: operator.Language, Value: *s.language,
		})
	}
	if s.caseSensitive != nil {
		res = append(res, bson.E{
			Key: operator.CaseSensitive, Value: *s.caseSensitive,
		})
	}
	if s.diacriticSensitive != nil {
		res = append(res, bson.E{
			Key: operator.DiacriticSensitive, Value: *s.diacriticSensitive,
		})
	}
	return res
}
