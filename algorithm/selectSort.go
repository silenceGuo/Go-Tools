package algorithm

import "fmt"

func SelectSort(arr *[6]int) {
	for j := 0; j < len(arr)-1; j++ {
		max := arr[j]
		maxIndex := j
		for i := j + 1; i < len(arr); i++ {
			if max > arr[i] {
				max = arr[i]
				maxIndex = i
			}
		}
		if maxIndex != j {
			arr[j], arr[maxIndex] = arr[maxIndex], arr[j]
		}

	}
	fmt.Println(arr)
	a := 1
	b := 2
	a, b = 2, 1
	fmt.Println(a, b)
}

func InsertSort(arr *[6]int) {
	for i := 1; i < len(arr); i++ {

	}
}
