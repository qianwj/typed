package valrange

// Of Obtains a fixed value range.
// This factory obtains a range where the minimum and maximum values are fixed.
// For example, the ISO month-of-year always runs from 1 to 12.
func Of(min, max int64) *ValueRange {
	return nil
}

func MaxOf(min, maxSmallest, maxLargest int64) *ValueRange {
	return nil
}

func Full(minSmallest, minLargest, maxSmallest, maxLargest int64) *ValueRange {
	return nil
}

type ValueRange struct {
}
