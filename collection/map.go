package collection

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Keys() []K {
	res := make([]K, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func (m Map[K, V]) ContainsKey(key K) bool {
	_, ok := m[key]
	return ok
}

func (m Map[K, V]) Values() []V {
	res := make([]V, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func (m Map[K, V]) Foreach(handle func(K, V)) {
	for k, v := range m {
		handle(k, v)
	}
}
