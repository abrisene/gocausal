package normalize

func ScaleNormalArray(collection []float64, scale float64) []float64 {
	scaled := make([]float64, len(collection))
	var sum float64

	for _, v := range collection {
		sum += v
	}

	for i, v := range collection {
		scaled[i] = (v / sum) * scale
	}
	return scaled
}

func NormalizeArray(collection []float64) []float64 {
	return ScaleNormalArray(collection, 1)
}

func ScaleNormalMap(collection map[string]float64, scale float64) map[string]float64 {
	scaled := make(map[string]float64)
	var sum float64

	for _, v := range collection {
		sum += v
	}

	for k, v := range collection {
		scaled[k] = ((v / sum) * scale)
	}

	return scaled
}

func NormalizeMap(collection map[string]float64) map[string]float64 {
	return ScaleNormalMap(collection, 1)
}
