package linestogo

// Lerp, linear interpolation.
// https://magcius.github.io/xplain/article/rast1.html#lerp
func lerp(a, b, t float32) float32 {
	return (a * (1.0 - t)) + (b * t)
}
