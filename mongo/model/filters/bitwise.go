package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BitsAllClear(key string, mask int) *Filter {
	return New().BitsAllClear(key, mask)
}

func BitsAllClearByBinary(key string, bin primitive.Binary) *Filter {
	return New().BitsAllClearByBinary(key, bin)
}

func BitsAllClearByPositions(key string, pos []int) *Filter {
	return New().BitsAllClearByPositions(key, pos)
}

func BitsAllSet(key string, mask int) *Filter {
	return New().BitsAllSet(key, mask)
}

func BitsAllSetByBinary(key string, bin primitive.Binary) *Filter {
	return New().BitsAllSetByBinary(key, bin)
}

func BitsAllSetByPositions(key string, pos []int) *Filter {
	return New().BitsAllSetByPositions(key, pos)
}

func BitsAnyClear(key string, mask int) *Filter {
	return New().BitsAnyClear(key, mask)
}

func BitsAnyClearByBinary(key string, bin primitive.Binary) *Filter {
	return New().BitsAnyClearByBinary(key, bin)
}

func BitsAnyClearByPositions(key string, pos []int) *Filter {
	return New().BitsAnyClearByPositions(key, pos)
}

func BitsAnySet(key string, mask int) *Filter {
	return New().BitsAnySet(key, mask)
}

func BitsAnySetByBinary(key string, bin primitive.Binary) *Filter {
	return New().BitsAnySetByBinary(key, bin)
}

func BitsAnySetByPositions(key string, pos []int) *Filter {
	return New().BitsAnySetByPositions(key, pos)
}

func (f *Filter) BitsAllClear(key string, mask int) *Filter {
	f.put(key, bson.M{operator.BitsAllClear: mask})
	return f
}

func (f *Filter) BitsAllClearByBinary(key string, bin primitive.Binary) *Filter {
	f.put(key, bson.M{operator.BitsAllClear: bin})
	return f
}

func (f *Filter) BitsAllClearByPositions(key string, pos []int) *Filter {
	f.put(key, bson.M{operator.BitsAllClear: pos})
	return f
}

func (f *Filter) BitsAllSet(key string, mask int) *Filter {
	f.put(key, bson.M{operator.BitsAllSet: mask})
	return f
}

func (f *Filter) BitsAllSetByBinary(key string, bin primitive.Binary) *Filter {
	f.put(key, bson.M{operator.BitsAllSet: bin})
	return f
}

func (f *Filter) BitsAllSetByPositions(key string, pos []int) *Filter {
	f.put(key, bson.M{operator.BitsAllSet: pos})
	return f
}

func (f *Filter) BitsAnyClear(key string, mask int) *Filter {
	f.put(key, bson.M{operator.BitsAnyClear: mask})
	return f
}

func (f *Filter) BitsAnyClearByBinary(key string, bin primitive.Binary) *Filter {
	f.put(key, bson.M{operator.BitsAnyClear: bin})
	return f
}

func (f *Filter) BitsAnyClearByPositions(key string, pos []int) *Filter {
	f.put(key, bson.M{operator.BitsAnySet: pos})
	return f
}

func (f *Filter) BitsAnySet(key string, mask int) *Filter {
	f.put(key, bson.M{operator.BitsAnySet: mask})
	return f
}

func (f *Filter) BitsAnySetByBinary(key string, bin primitive.Binary) *Filter {
	f.put(key, bson.M{operator.BitsAnySet: bin})
	return f
}

func (f *Filter) BitsAnySetByPositions(key string, pos []int) *Filter {
	f.put(key, bson.M{operator.BitsAnySet: pos})
	return f
}
