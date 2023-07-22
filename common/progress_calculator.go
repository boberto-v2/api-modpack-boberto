package common

type ProgressCalculator struct {
	maxValue     int
	currentValue float64
	Progress     float64
}

func New(max int) ProgressCalculator {
	ProgressCalculator := ProgressCalculator{
		maxValue:     max,
		currentValue: 0,
	}
	return ProgressCalculator
}

func (ProgressCalculator *ProgressCalculator) Increment() {
	ProgressCalculator.currentValue++
}

func (ProgressCalculator *ProgressCalculator) Calculate() float64 {
	percentage := (ProgressCalculator.currentValue / float64(ProgressCalculator.maxValue)) * 100
	ProgressCalculator.Progress = percentage
	return percentage
}
