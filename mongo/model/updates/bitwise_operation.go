package updates

import (
	"go.mongodb.org/mongo-driver/bson"
)

func (u *Update) BitAnd(field string, val int) *Update {
	u.bit[field] = bson.M{"and": val}
	return u
}

func (u *Update) BitOr(field string, val int) *Update {
	u.bit[field] = bson.M{"or": val}
	return u
}

func (u *Update) BitXor(field string, val int) *Update {
	u.bit[field] = bson.M{"xor": val}
	return u
}
