package list

type ArrayList[T any] []T

const defaultInitialCapacity = 10

func NewArrayList[T any]() ArrayList[T] {
	return make(ArrayList[T], 0, defaultInitialCapacity)
}

func (a ArrayList[T]) Add(e T) ArrayList[T] {
	return append(a, e)
}

func (a ArrayList[T]) Contains(e T) bool {
	//idx := a.binarySearch(e, 0, a.Size())
	//return idx > -1
	return false
}

func (a ArrayList[T]) Foreach(handle func(e T)) {
	for _, e := range a {
		handle(e)
	}
}

func (a ArrayList[T]) Size() int {
	return len(a)
}

//func (a ArrayList[T]) binarySearch(target T, left, right int) int {
//	for left <= right {
//		middleIndex := left + (right-left)/2
//		if a[middleIndex] == target {
//			return middleIndex
//		} else if a[middleIndex] > target {
//			right = middleIndex - 1
//		} else {
//			left = middleIndex + 1
//		}
//	}
//}
