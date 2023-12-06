package test

// segment-tree

type SegNode struct {
	left, right, sum int
	tag              bool
}

type SegTree struct {
	arr []*SegNode
}

func NewSegTree(nums []int) *SegTree {
	n := len(nums)
	tree := &SegTree{
		arr: make([]*SegNode, n*4+1),
	}
	build(tree.arr, 1, 0, n-1, nums)
	return tree
}

func build(arr []*SegNode, idx, left, right int, nums []int) {
	node := &SegNode{
		left:  left,
		right: right,
		tag:   false,
	}
	arr[idx] = node
	if left == right {
		node.sum = nums[left]
		return
	}

	mid := left + ((right - left) >> 1)
	build(arr, idx*2, left, mid, nums)
	build(arr, idx*2+1, mid+1, right, nums)

	node.sum = leftOf(arr, idx).sum + rightOf(arr, idx).sum
}

func (tree *SegTree) IndexOf(index int) *SegNode {
	return tree.arr[index]
}

func (tree *SegTree) LeftOf(index int) *SegNode {
	return leftOf(tree.arr, index)
}

func leftOf(arr []*SegNode, index int) *SegNode {
	return arr[index*2]
}

func (tree *SegTree) RightOf(index int) *SegNode {
	return rightOf(tree.arr, index)
}

func rightOf(arr []*SegNode, index int) *SegNode {
	return arr[index*2+1]
}

func pushdown(arr []*SegNode, idx int) {
	node := arr[idx]
	if node.tag {
		lNode, rNode := leftOf(arr, idx), rightOf(arr, idx)
		lNode.tag = node.tag
		lNode.sum = lNode.right - lNode.left + 1 - lNode.sum

		rNode.tag = node.tag
		rNode.sum = rNode.right - rNode.left + 1 - rNode.sum

		node.tag = false
	}
}

// 区间反转
func reverse(arr []*SegNode, idx, left, right int) {
	node := arr[idx]
	if node.left >= left && node.right <= right {
		node.sum = node.right - node.left + 1 - node.sum
		node.tag = !node.tag
		return
	}
	if node.left > right || node.right < left {
		return
	}

	// Pushdown
	if node.tag {
		pushdown(arr, idx)
	}

	mid := node.left + ((node.right - node.left) >> 1)
	// left.right >= left, range in left
	if left <= mid {
		reverse(arr, idx*2, left, right)
	}
	// right.left <= right, range in right
	if right > mid {
		reverse(arr, idx*2+1, left, right)
	}
	node.sum = leftOf(arr, idx).sum + rightOf(arr, idx).sum
}

func (tree *SegTree) Reverse(left, right int) {
	reverse(tree.arr, 1, left, right)
}

// 区间求和
func sum(arr []*SegNode, idx, left, right int) int {
	node := arr[idx]
	if node.left >= left && node.right <= right {
		return node.sum
	}
	if node.left > right || node.right < left {
		return 0
	}

	// Pushdown
	if node.tag {
		pushdown(arr, idx)
	}

	res := 0
	// left.right >= left, range in left
	if leftOf(arr, idx).right >= left {
		res += sum(arr, idx*2, left, right)
	}
	// right.left <= right, range in right
	if rightOf(arr, idx).left <= right {
		res += sum(arr, idx*2+1, left, right)
	}
	return res
}

func (tree *SegTree) Sum(left, right int) int {
	return sum(tree.arr, 1, left, right)
}

// algorithm impl
func handleQuery(nums1 []int, nums2 []int, queries [][]int) []int64 {

	sum := int64(0)
	for _, v := range nums2 {
		sum += int64(v)
	}

	tree := NewSegTree(nums1)
	res := []int64{}
	for i := 0; i < len(queries); i++ {
		if queries[i][0] == 1 {
			tree.Reverse(queries[i][1], queries[i][2])
		} else if queries[i][0] == 2 {
			nums1Sum := tree.Sum(0, len(nums1)-1)
			sum += int64(nums1Sum) * int64(queries[i][1])
		} else if queries[i][0] == 3 {
			res = append(res, sum)
		}
	}
	return res
}
