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
	t := make([]byte, len(T))
	for i := range t {
		t[i] = '_'
	}

	t[len(T)-1] = 'S'

	// Identificamos la posición como "S" o "L"
	for i := (len(T) - 1); i > 0; i-- {
		if T[i-1] == T[i] {
			t[i-1] = t[i]
		} else if T[i-1] < T[i] {
			t[i-1] = 'S'
		} else {
			t[i-1] = 'L'
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
		if t[i] == 'S' && t[i-1] == 'L' {
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
		if SA[i] > 0 && SA[i] < len(T) && SA[i]-1 >= 0 && t[SA[i]-1] == 'L' { //Posición válida
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
		if SA[i] > 0 && SA[i] < len(T) && SA[i]-1 >= 0 && t[SA[i]-1] == 'S' { //Posiciones válidas en SA
			symbol := T[SA[i]-1]
			count[symbol]++
			revoffset := count[symbol]
			pos := buckets[symbol][1] - revoffset
			if pos >= 0 {
				SA[pos] = SA[i] - 1
			}
		}
	}

	// Nombres
	namesp := make([]int, len(T))
	for i := range namesp {
		namesp[i] = -1
	}
	name := 0
	var prev *int
	for i := 0; i < len(T); i++ {
		if SA[i] <= 0 || SA[i] >= len(T) {
			continue
		}
		if _, ok := LMS[SA[i]]; !ok {
			continue
		}
		if prev != nil && *prev >= 0 && *prev < len(T) {
			prevEnd, currEnd := LMS[*prev], LMS[SA[i]]
			if prevEnd > len(T) || currEnd > len(T) {
				continue
			}
			previous := T[*prev:prevEnd]
			current := T[SA[i]:currEnd]
			equal := len(previous) == len(current)
			if equal {
				for j := range previous {
					if previous[j] != current[j] {
						equal = false
						break
					}
				}
			}
			if !equal {
				name++
			}
		}
		tmp := SA[i]
		prev = &tmp
		namesp[SA[i]] = name
	}

	names := make([]int, 0)
	SApIdx := make([]int, 0)
	for i := 0; i < len(T); i++ {
		if namesp[i] != -1 {
			names = append(names, namesp[i])
			SApIdx = append(SApIdx, i)
		}
	}

	if len(names) > 0 && name < len(names)-1 {
		names = sais(names)
	}

	//Revertimos los nombres
	for i, j := 0, len(names)-1; i < j; i, j = i+1, j-1 {
		names[i], names[j] = names[j], names[i]
	}

	//Sort final.
	SA = make([]int, len(T))
	for i := range SA {
		SA[i] = -1
	}
	count = make(map[int]int)

	for i := 0; i < len(names); i++ {
		if names[i] < 0 || names[i] >= len(SApIdx) {
			continue
		}
		pos := SApIdx[names[i]]
		if pos < 0 || pos >= len(T) {
			continue
		}
		count[T[pos]]++
		revoffset := count[T[pos]]
		idx := buckets[T[pos]][1] - revoffset
		if idx >= 0 && idx < len(T) {
			SA[idx] = pos
		}
	}

	//Inducción para "L"
	count = make(map[int]int)
	for i := 0; i < len(T); i++ {
		if SA[i] >= 0 && SA[i]-1 >= 0 && SA[i] < len(T) && SA[i]-1 >= 0 && t[SA[i]-1] == 'L' {
			symbol := T[SA[i]-1]
			offset := count[symbol]
			pos := buckets[symbol][0] + offset
			if pos < len(T) {
				SA[pos] = SA[i] - 1
			}
			count[symbol]++
		}
	}

	// Inducción para "S"
	count = make(map[int]int)
	for i := len(T) - 1; i > 0; i-- {
		if SA[i] > 0 && SA[i] < len(T) && SA[i]-1 >= 0 && t[SA[i]-1] == 'S' {
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
