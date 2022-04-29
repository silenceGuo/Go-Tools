package algorithm

import "testing"

func TestSelectSort(t *testing.T) {
	S := [6]int{3, 5, 12, 54, 11, 52}
	SelectSort(&S)
}
