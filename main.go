package main

import "fmt"

func main() {
	// TODO
	fmt.Println(createString(5, 4))

	// for i := 0; i < n; i++ {

	// }
}

func createString(n, s int) string {
	cadena := ""

	for i := 0; i < n; i++ {
		if i >= s {
			cadena += "#"
		}
		cadena += "/"
	}

	return cadena
}
