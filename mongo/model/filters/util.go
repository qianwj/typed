package filters

import "go.mongodb.org/mongo-driver/bson"

func _map[K, V any](collection []K, convert func(k K) V) bson.A {
	res := make(bson.A, len(collection))
	for i, k := range collection {
		res[i] = convert(k)
	}
	return res
}
