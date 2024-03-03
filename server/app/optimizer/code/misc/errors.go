package misc

import "fmt"

type ValidationError struct {
}

func (ve ValidationError) Error() string {
	return "validation error"
}

type NoSpaceError struct {
}

func (nse NoSpaceError) Error() string {
	return "no space"
}

type NoBucketError struct {
	Price float64
}

func (nbe NoBucketError) Error() string {
	return fmt.Sprintf("no bucket %f", nbe.Price)
}

type PriceRangeError struct {
	Min float64
	Max float64
}

func (pre PriceRangeError) Error() string {
	return fmt.Sprintf("invalid price range [%f, %f]", pre.Min, pre.Max)
}

type InvalidLevelError struct {
	Level int
}

func (ile InvalidLevelError) Error() string {
	return fmt.Sprintf("invalid level error %d", ile.Level)
}

type BucketSizeError struct {
	BucketSize int
}

func (bse BucketSizeError) Error() string {
	return fmt.Sprintf("invalid bucket size %d", bse.BucketSize)
}

type BufferSizeError struct {
	BufferSize int
}

func (buse BufferSizeError) Error() string {
	return fmt.Sprintf("buffer size error %d", buse.BufferSize)
}

type DiscountFactorError struct {
	Discount float64
}

func (dfe DiscountFactorError) Error() string {
	return fmt.Sprintf("discount functor error %f", dfe.Discount)
}

type UnfeasiblePriceError struct {
	Price float64
	Min   float64
	Max   float64
}

func (upe UnfeasiblePriceError) Error() string {
	return fmt.Sprintf("unfeasible price %f [%f, %f]", upe.Price, upe.Min, upe.Max)
}

type FailedToSampleExplorationBucket struct {
}

func (feb FailedToSampleExplorationBucket) Error() string {
	return "Failed to sample bucket for exploration"
}
