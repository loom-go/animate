package animate

func clamp(x, low, hight float64) float64 {
	return max(low, min(x, hight))
}
