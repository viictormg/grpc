package main

func main() {
	// TODO
	n := 8
	b := 8

	for i := 0; i < n; i++ {
		mark := ""
		for j := 0; j < n; j++ {
			if j >= (b-1)-i {
				mark += "#"
			}
			b--
		}

		println(mark)
	}
}
