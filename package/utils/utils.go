package utils

func Gcd(a, b int) int {
	// Нахождение наибольшего общего делителя (алгоритм Евклида)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
