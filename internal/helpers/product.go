package helpers

func SumFloatSlice(numbers []float64) float64 {
	var total float64
	for _, num := range numbers {
		total += num
	}
	return total
}
