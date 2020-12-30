package main

import "fmt"

type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

func f(a *[3]int) {
	a[1] = 100
}

func main() {
	a := [3]int{1, 2, 3}
	f(&a)
	fmt.Println(a[1])

	primes := [6]int{2, 3, 5, 7, 11, 13}

	var s []int = primes[1:4]
	fmt.Println(s)

	board := [][][]string{
		[][]string{[]string{"A", "B"}, []string{"C", "D"}},
		[][]string{[]string{"E", "F"}, []string{"G", "H"}},
	}

	fmt.Println(board)

	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])

}
