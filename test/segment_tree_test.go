package test

import (
	"fmt"
	"testing"
)

func TestHandleQuery(t *testing.T) {
	nums1 := []int{1, 0, 1}
	nums2 := []int{0, 0, 0}
	queries := [][]int{
		{1, 1, 1},
		{2, 1, 9},
		{3, 0, 0},
	}
	res := handleQuery(nums1, nums2, queries)
	fmt.Println(res)
	t.Log(res)
}
