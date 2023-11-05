package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/util"
)

type UnwindOptions struct {
	path                       string
	includeArrayIndex          *string
	preserveNullAndEmptyArrays *bool
}

func NewUnwindOptions(path string) *UnwindOptions {
	return &UnwindOptions{path: path}
}

func (u *UnwindOptions) IncludeArrayIndex(includeArrayIndex string) *UnwindOptions {
	u.includeArrayIndex = util.ToPtr(includeArrayIndex)
	return u
}

func (u *UnwindOptions) PreserveNullAndEmptyArrays(preserveNullAndEmptyArrays bool) *UnwindOptions {
	u.preserveNullAndEmptyArrays = util.ToPtr(preserveNullAndEmptyArrays)
	return u
}

func (u *UnwindOptions) MarshalBSON() ([]byte, error) {
	m := bson.M(bson.E("path", u.path))
	if util.IsNonNil(u.includeArrayIndex) {
		m["includeArrayIndex"] = u.includeArrayIndex
	}
	if util.IsNonNil(u.preserveNullAndEmptyArrays) {
		m["preserveNullAndEmptyArrays"] = u.preserveNullAndEmptyArrays
	}
	return m.Marshal()
}
