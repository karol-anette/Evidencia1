package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Adaptamos la función a Golang
func getBuckets(T []int) map[int][2]int {
	count := make(map[int]int)
	buckets := make(map[int][2]int)
	keys := make([]int, 0, len(count)) //Necesitamos keys
	//Cuenta las ocurrencias de cada símboolo
	for _, c := range T {
		count[c]++
	}
	//Sustituímos el sorted(count.keys()) por un sort manual
	for k := range count {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	//Hacemos las Buckets
	start := 0
	for _, c := range keys {
		buckets[c] = [2]int{start, start + count[c]}
		start += count[c]
	}
	return buckets
}

// Función sais para crear el suffix array
func sais(T []int) []int {
	isS := make([]bool, len(T)) //Arreglo bool
	isS[len(T)-1] = true        //El último carácter siempre es "S"

	// Identificamos la posición como "S" o "L"
	for i := (len(T) - 1); i > 0; i-- {
		if T[i-1] < T[i] {
			isS[i-1] = true
		} else if T[i-1] == T[i] && isS[i] {
			isS[i-1] = true //true es "S" y false es "L"
		} else {
			isS[i-1] = false
		}
	}

	buckets := getBuckets(T)

	count := make(map[int]int)
	SA := make([]int, len(T))
	for i := range SA {
		SA[i] = -1
	}
	LMS := make(map[int]int)
	var end *int

	// Colocamos los substrings LMS
	for i := (len(T) - 1); i > 0; i-- {
		if isS[i] && !isS[i-1] {
			count[T[i]]++
			revoffset := count[T[i]]
			pos := buckets[T[i]][1] - revoffset
			if pos >= 0 && pos < len(T) { //Checamos que sea una posición válida
				SA[pos] = i
			}
			if end != nil {
				LMS[i] = *end
			}
			tmp := i
			end = &tmp
		}
	}
	LMS[len(T)-1] = len(T) - 1

	// Hacemos el sort para "L"
	count = make(map[int]int)
	for i := 0; i < len(T); i++ {
		if SA[i] > 0 && SA[i]-1 >= 0 && !isS[SA[i]-1] { //Posición válida
			symbol := T[SA[i]-1]
			offset := count[symbol]
			pos := buckets[symbol][0] + offset
			if pos < len(T) {
				SA[pos] = SA[i] - 1
			}
			count[symbol]++
		}
	}

	// Hacemos el sort para "S"
	count = make(map[int]int)
	for i := len(T) - 1; i > 0; i-- {
		if SA[i] > 0 && SA[i]-1 >= 0 && isS[SA[i]-1] { //Posiciones válidas en SA
			symbol := T[SA[i]-1]
			count[symbol]++
			revoffset := count[symbol]
			pos := buckets[symbol][1] - revoffset
			if pos >= 0 && pos < len(T) {
				SA[pos] = SA[i] - 1
			}
		}
	}

	return SA
}

// Función main para poder leer el string
func main() {
	file, err := os.Open("hamlet.txt") //Cambiamos string por el archivo.
	if err != nil {
		fmt.Println("No se puede abrir el archivo", err)
		return
	}
	defer file.Close() //Avisa si no se puede abrir el archivo y lo cierra.

	var T []int
	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune() // ReadRune() lee caracter por caracter
		if err != nil {
			break
		}
		T = append(T, int(r))
	}

	SA := sais(T)
	fmt.Println(SA)
}
