package db

type TreeType string

const (
	TreeTypeAVL      TreeType = "avl"
	TreeTypeRedBlack TreeType = "redblack"
	TreeTypeBTree    TreeType = "btree"
)

type AVLNode struct {
	Key    string
	Value  string
	Height int
	Left   *AVLNode
	Right  *AVLNode
}

type AVLTree struct {
	Root *AVLNode
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

func height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func (t *AVLTree) rightRotate(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	return x
}

func (t *AVLTree) leftRotate(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	return y
}

func (t *AVLTree) insert(node *AVLNode, key string, value string) *AVLNode {
	if node == nil {
		return &AVLNode{Key: key, Value: value, Height: 1}
	}

	if key < node.Key {
		node.Left = t.insert(node.Left, key, value)
	} else if key > node.Key {
		node.Right = t.insert(node.Right, key, value)
	} else {
		node.Value = value
		return node
	}

	node.Height = max(height(node.Left), height(node.Right)) + 1

	balance := getBalance(node)

	if balance > 1 && key < node.Left.Key {
		return t.rightRotate(node)
	}

	if balance < -1 && key > node.Right.Key {
		return t.leftRotate(node)
	}

	if balance > 1 && key > node.Left.Key {
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}

	if balance < -1 && key < node.Right.Key {
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

func (t *AVLTree) search(node *AVLNode, key string) *AVLNode {
	if node == nil || node.Key == key {
		return node
	}

	if key < node.Key {
		return t.search(node.Left, key)
	}
	return t.search(node.Right, key)
}

func (t *AVLTree) searchRange(node *AVLNode, leftBound string, rightBound string, result *map[string]string) {
	if node == nil {
		return
	}

	if leftBound < node.Key {
		t.searchRange(node.Left, leftBound, rightBound, result)
	}

	if leftBound <= node.Key && node.Key <= rightBound {
		(*result)[node.Key] = node.Value
	}

	if rightBound > node.Key {
		t.searchRange(node.Right, leftBound, rightBound, result)
	}
}

func (t *AVLTree) minValueNode(node *AVLNode) *AVLNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (t *AVLTree) delete(node *AVLNode, key string) *AVLNode {
	if node == nil {
		return node
	}

	if key < node.Key {
		node.Left = t.delete(node.Left, key)
	} else if key > node.Key {
		node.Right = t.delete(node.Right, key)
	} else {
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		temp := t.minValueNode(node.Right)
		node.Key = temp.Key
		node.Value = temp.Value
		node.Right = t.delete(node.Right, temp.Key)
	}

	if node == nil {
		return node
	}

	node.Height = max(height(node.Left), height(node.Right)) + 1

	balance := getBalance(node)

	if balance > 1 && getBalance(node.Left) >= 0 {
		return t.rightRotate(node)
	}

	if balance > 1 && getBalance(node.Left) < 0 {
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}

	if balance < -1 && getBalance(node.Right) <= 0 {
		return t.leftRotate(node)
	}

	if balance < -1 && getBalance(node.Right) > 0 {
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

type Color bool

const (
	RED   Color = true
	BLACK Color = false
)

type RBNode struct {
	Key    string
	Value  string
	Color  Color
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
}

type RedBlackTree struct {
	Root *RBNode
	NIL  *RBNode
}

func NewRedBlackTree() *RedBlackTree {
	nil_node := &RBNode{Color: BLACK}
	return &RedBlackTree{
		NIL:  nil_node,
		Root: nil_node,
	}
}

func (t *RedBlackTree) leftRotate(x *RBNode) {
	y := x.Right
	x.Right = y.Left

	if y.Left != t.NIL {
		y.Left.Parent = x
	}

	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (t *RedBlackTree) rightRotate(y *RBNode) {
	x := y.Left
	y.Left = x.Right

	if x.Right != t.NIL {
		x.Right.Parent = y
	}

	x.Parent = y.Parent

	if y.Parent == t.NIL {
		t.Root = x
	} else if y == y.Parent.Right {
		y.Parent.Right = x
	} else {
		y.Parent.Left = x
	}

	x.Right = y
	y.Parent = x
}

func (t *RedBlackTree) insertFixup(z *RBNode) {
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y.Color == RED {
				// Case 1: Uncle is red
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.leftRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.rightRotate(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rightRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.leftRotate(z.Parent.Parent)
			}
		}
		if z == t.Root {
			break
		}
	}
	t.Root.Color = BLACK
}

func (t *RedBlackTree) insert(key string, value string) {
	z := &RBNode{
		Key:    key,
		Value:  value,
		Color:  RED,
		Left:   t.NIL,
		Right:  t.NIL,
		Parent: t.NIL,
	}

	y := t.NIL
	x := t.Root

	for x != t.NIL {
		y = x
		if z.Key < x.Key {
			x = x.Left
		} else if z.Key > x.Key {
			x = x.Right
		} else {
			x.Value = value
			return
		}
	}

	z.Parent = y
	if y == t.NIL {
		t.Root = z
	} else if z.Key < y.Key {
		y.Left = z
	} else {
		y.Right = z
	}

	t.insertFixup(z)
}

func (t *RedBlackTree) search(key string) *RBNode {
	x := t.Root
	for x != t.NIL && x.Key != key {
		if key < x.Key {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	if x == t.NIL {
		return nil
	}
	return x
}

func (t *RedBlackTree) searchRange(leftBound string, rightBound string, result *map[string]string) {
	var inorderTraversal func(*RBNode)
	inorderTraversal = func(node *RBNode) {
		if node == t.NIL {
			return
		}

		if leftBound < node.Key {
			inorderTraversal(node.Left)
		}

		if leftBound <= node.Key && node.Key <= rightBound {
			(*result)[node.Key] = node.Value
		}

		if rightBound > node.Key {
			inorderTraversal(node.Right)
		}
	}

	inorderTraversal(t.Root)
}

func (t *RedBlackTree) transplant(u *RBNode, v *RBNode) {
	if u.Parent == t.NIL {
		t.Root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	v.Parent = u.Parent
}

func (t *RedBlackTree) minimum(x *RBNode) *RBNode {
	for x.Left != t.NIL {
		x = x.Left
	}
	return x
}

func (t *RedBlackTree) deleteFixup(x *RBNode) {
	for x != t.Root && x.Color == BLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == BLACK && w.Right.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Right.Color == BLACK {
					w.Left.Color = BLACK
					w.Color = RED
					t.rightRotate(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Right.Color = BLACK
				t.leftRotate(x.Parent)
				x = t.Root
			}
		} else {
			w := x.Parent.Left
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Right.Color == BLACK && w.Left.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Left.Color == BLACK {
					w.Right.Color = BLACK
					w.Color = RED
					t.leftRotate(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Left.Color = BLACK
				t.rightRotate(x.Parent)
				x = t.Root
			}
		}
	}
	x.Color = BLACK
}

func (t *RedBlackTree) delete(key string) bool {
	z := t.search(key)
	if z == nil {
		return false
	}

	y := z
	y_original_color := y.Color
	var x *RBNode

	if z.Left == t.NIL {
		x = z.Right
		t.transplant(z, z.Right)
	} else if z.Right == t.NIL {
		x = z.Left
		t.transplant(z, z.Left)
	} else {
		y = t.minimum(z.Right)
		y_original_color = y.Color
		x = y.Right

		if y.Parent == z {
			x.Parent = y
		} else {
			t.transplant(y, y.Right)
			y.Right = z.Right
			y.Right.Parent = y
		}

		t.transplant(z, y)
		y.Left = z.Left
		y.Left.Parent = y
		y.Color = z.Color
	}

	if y_original_color == BLACK {
		t.deleteFixup(x)
	}

	return true
}

type BTreeNode struct {
	Keys     []string
	Values   []string
	Children []*BTreeNode
	Leaf     bool
	n        int
}

type BTree struct {
	Root   *BTreeNode
	MinDeg int
}

func NewBTree(degree int) *BTree {
	return &BTree{
		MinDeg: degree,
		Root: &BTreeNode{
			Leaf:     true,
			Keys:     make([]string, 2*degree-1),
			Values:   make([]string, 2*degree-1),
			Children: make([]*BTreeNode, 2*degree),
		},
	}
}

func NewBTreeNode(t int, leaf bool) *BTreeNode {
	return &BTreeNode{
		Leaf:     leaf,
		Keys:     make([]string, 2*t-1),
		Values:   make([]string, 2*t-1),
		Children: make([]*BTreeNode, 2*t),
		n:        0,
	}
}

func (node *BTreeNode) search(key string) (string, bool) {
	i := 0
	for i < node.n && key > node.Keys[i] {
		i++
	}

	if i < node.n && key == node.Keys[i] {
		return node.Values[i], true
	}

	if node.Leaf {
		return "", false
	}

	return node.Children[i].search(key)
}

func (t *BTree) Search(key string) (string, bool) {
	if t.Root == nil {
		return "", false
	}
	return t.Root.search(key)
}

func (node *BTreeNode) searchRange(leftBound, rightBound string, result *map[string]string) {
	i := 0

	for i < node.n && node.Keys[i] < leftBound {
		i++
	}

	for i < node.n && node.Keys[i] <= rightBound {
		if node.Keys[i] >= leftBound {
			(*result)[node.Keys[i]] = node.Values[i]
		}
		i++
	}

	if !node.Leaf {
		for j := 0; j < node.n+1 && j <= i; j++ {
			if j < node.n && node.Keys[j] > rightBound {
				break
			}
			node.Children[j].searchRange(leftBound, rightBound, result)
		}
	}
}

func (t *BTree) insertNonFull(node *BTreeNode, key string, value string) {
	i := node.n - 1

	if node.Leaf {
		for i >= 0 && key < node.Keys[i] {
			node.Keys[i+1] = node.Keys[i]
			node.Values[i+1] = node.Values[i]
			i--
		}
		node.Keys[i+1] = key
		node.Values[i+1] = value
		node.n++
	} else {
		for i >= 0 && key < node.Keys[i] {
			i--
		}
		i++

		if node.Children[i].n == 2*t.MinDeg-1 {
			t.splitChild(node, i)

			if key > node.Keys[i] {
				i++
			}
		}
		t.insertNonFull(node.Children[i], key, value)
	}
}

func (t *BTree) splitChild(parent *BTreeNode, i int) {
	minDeg := t.MinDeg
	child := parent.Children[i]
	newNode := NewBTreeNode(minDeg, child.Leaf)
	newNode.n = minDeg - 1

	for j := 0; j < minDeg-1; j++ {
		newNode.Keys[j] = child.Keys[j+minDeg]
		newNode.Values[j] = child.Values[j+minDeg]
	}

	if !child.Leaf {
		for j := 0; j < minDeg; j++ {
			newNode.Children[j] = child.Children[j+minDeg]
		}
	}

	child.n = minDeg - 1

	for j := parent.n; j >= i+1; j-- {
		parent.Children[j+1] = parent.Children[j]
	}

	parent.Children[i+1] = newNode

	for j := parent.n - 1; j >= i; j-- {
		parent.Keys[j+1] = parent.Keys[j]
		parent.Values[j+1] = parent.Values[j]
	}

	parent.Keys[i] = child.Keys[minDeg-1]
	parent.Values[i] = child.Values[minDeg-1]
	parent.n++
}

func (t *BTree) Insert(key string, value string) {
	root := t.Root

	if root.n == 2*t.MinDeg-1 {
		newRoot := NewBTreeNode(t.MinDeg, false)
		t.Root = newRoot
		newRoot.Children[0] = root
		t.splitChild(newRoot, 0)
		t.insertNonFull(newRoot, key, value)
	} else {
		t.insertNonFull(root, key, value)
	}
}

func (node *BTreeNode) findKey(key string) int {
	idx := 0
	for idx < node.n && node.Keys[idx] < key {
		idx++
	}
	return idx
}

func (t *BTree) Delete(key string) bool {
	if t.Root == nil {
		return false
	}

	t.delete(t.Root, key)

	if t.Root.n == 0 && !t.Root.Leaf {
		t.Root = t.Root.Children[0]
	}
	return true
}

func (t *BTree) delete(node *BTreeNode, key string) {
	idx := node.findKey(key)

	if idx < node.n && node.Keys[idx] == key {
		if node.Leaf {
			//  node is leaf(remove the key)
			for i := idx + 1; i < node.n; i++ {
				node.Keys[i-1] = node.Keys[i]
				node.Values[i-1] = node.Values[i]
			}
			node.n--
		} else {
			// node is internal
			if node.Children[idx].n >= t.MinDeg {
				// left child has at least t keys
				pred := t.getPred(node, idx)
				node.Keys[idx] = pred.Keys[pred.n-1]
				node.Values[idx] = pred.Values[pred.n-1]
				t.delete(node.Children[idx], node.Keys[idx])
			} else if node.Children[idx+1].n >= t.MinDeg {
				// right child has at least t keys
				succ := t.getSucc(node, idx)
				node.Keys[idx] = succ.Keys[0]
				node.Values[idx] = succ.Values[0]
				t.delete(node.Children[idx+1], node.Keys[idx])
			} else {
				// left and right children have t-1 keys
				t.merge(node, idx)
				t.delete(node.Children[idx], key)
			}
		}
	} else {
		if node.Leaf {
			return
		}

		// check iflast child needs to be traversed
		flag := idx == node.n

		// if child has less than t keys fill it
		if node.Children[idx].n < t.MinDeg {
			t.fill(node, idx)
		}

		// if last child has been merged
		if flag && idx > node.n {
			t.delete(node.Children[idx-1], key)
		} else {
			t.delete(node.Children[idx], key)
		}
	}
}

func (t *BTree) getPred(node *BTreeNode, idx int) *BTreeNode {
	curr := node.Children[idx]
	for !curr.Leaf {
		curr = curr.Children[curr.n]
	}
	return curr
}

func (t *BTree) getSucc(node *BTreeNode, idx int) *BTreeNode {
	curr := node.Children[idx+1]
	for !curr.Leaf {
		curr = curr.Children[0]
	}
	return curr
}

func (t *BTree) fill(node *BTreeNode, idx int) {
	if idx != 0 && node.Children[idx-1].n >= t.MinDeg {
		t.borrowFromPrev(node, idx)
	} else if idx != node.n && node.Children[idx+1].n >= t.MinDeg {
		t.borrowFromNext(node, idx)
	} else {
		if idx != node.n {
			t.merge(node, idx)
		} else {
			t.merge(node, idx-1)
		}
	}
}

func (t *BTree) borrowFromPrev(node *BTreeNode, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx-1]

	for i := child.n - 1; i >= 0; i-- {
		child.Keys[i+1] = child.Keys[i]
		child.Values[i+1] = child.Values[i]
	}

	if !child.Leaf {
		for i := child.n; i >= 0; i-- {
			child.Children[i+1] = child.Children[i]
		}
	}

	child.Keys[0] = node.Keys[idx-1]
	child.Values[0] = node.Values[idx-1]

	if !child.Leaf {
		child.Children[0] = sibling.Children[sibling.n]
	}

	node.Keys[idx-1] = sibling.Keys[sibling.n-1]
	node.Values[idx-1] = sibling.Values[sibling.n-1]

	child.n++
	sibling.n--
}

func (t *BTree) borrowFromNext(node *BTreeNode, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx+1]

	child.Keys[child.n] = node.Keys[idx]
	child.Values[child.n] = node.Values[idx]

	if !child.Leaf {
		child.Children[child.n+1] = sibling.Children[0]
	}

	node.Keys[idx] = sibling.Keys[0]
	node.Values[idx] = sibling.Values[0]

	for i := 1; i < sibling.n; i++ {
		sibling.Keys[i-1] = sibling.Keys[i]
		sibling.Values[i-1] = sibling.Values[i]
	}

	if !sibling.Leaf {
		for i := 1; i <= sibling.n; i++ {
			sibling.Children[i-1] = sibling.Children[i]
		}
	}

	child.n++
	sibling.n--
}

func (t *BTree) merge(node *BTreeNode, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx+1]

	child.Keys[t.MinDeg-1] = node.Keys[idx]
	child.Values[t.MinDeg-1] = node.Values[idx]

	for i := 0; i < sibling.n; i++ {
		child.Keys[i+t.MinDeg] = sibling.Keys[i]
		child.Values[i+t.MinDeg] = sibling.Values[i]
	}

	if !child.Leaf {
		for i := 0; i <= sibling.n; i++ {
			child.Children[i+t.MinDeg] = sibling.Children[i]
		}
	}

	for i := idx + 1; i < node.n; i++ {
		node.Keys[i-1] = node.Keys[i]
		node.Values[i-1] = node.Values[i]
	}

	for i := idx + 2; i <= node.n; i++ {
		node.Children[i-1] = node.Children[i]
	}

	child.n += sibling.n + 1
	node.n--
}

type TreeCollection struct {
	TreeType TreeType
	AVL      *AVLTree
	RB       *RedBlackTree
	BT       *BTree
}

func NewTreeCollection(treeType TreeType) *TreeCollection {
	tc := &TreeCollection{TreeType: treeType}
	switch treeType {
	case TreeTypeAVL:
		tc.AVL = NewAVLTree()
	case TreeTypeRedBlack:
		tc.RB = NewRedBlackTree()
	case TreeTypeBTree:
		tc.BT = NewBTree(3)
	}
	return tc
}

func (tc *TreeCollection) Set(key string, secondaryKey string, value string) string {
	switch tc.TreeType {
	case TreeTypeAVL:
		tc.AVL.Root = tc.AVL.insert(tc.AVL.Root, key, value)
		return "ok"
	case TreeTypeRedBlack:
		tc.RB.insert(key, value)
		return "ok"
	case TreeTypeBTree:
		tc.BT.Insert(key, value)
		return "ok"
	}
	return "error: unknown tree type"
}

func (tc *TreeCollection) Update(key string, value string) string {
	return tc.Set(key, key, value)
}

func (tc *TreeCollection) Get(key string) (string, string) {
	switch tc.TreeType {
	case TreeTypeAVL:
		node := tc.AVL.search(tc.AVL.Root, key)
		if node == nil {
			return "", "error: not found"
		}
		return node.Value, "ok"
	case TreeTypeRedBlack:
		node := tc.RB.search(key)
		if node == nil {
			return "", "error: not found"
		}
		return node.Value, "ok"
	case TreeTypeBTree:
		if value, found := tc.BT.Search(key); found {
			return value, "ok"
		}
		return "", "error: not found"
	}
	return "", "error: unknown tree type"
}

func (tc *TreeCollection) GetRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	switch tc.TreeType {
	case TreeTypeAVL:
		tc.AVL.searchRange(tc.AVL.Root, leftBound, rightBound, &result)
		return &result, "ok"
	case TreeTypeRedBlack:
		tc.RB.searchRange(leftBound, rightBound, &result)
		return &result, "ok"
	case TreeTypeBTree:
		tc.BT.Root.searchRange(leftBound, rightBound, &result)
		return &result, "ok"
	}
	return &result, "error: unknown tree type"
}

func (tc *TreeCollection) Delete(key string) string {
	switch tc.TreeType {
	case TreeTypeAVL:
		node := tc.AVL.search(tc.AVL.Root, key)
		if node == nil {
			return "error: not found"
		}
		tc.AVL.Root = tc.AVL.delete(tc.AVL.Root, key)
		return "ok"
	case TreeTypeRedBlack:
		if tc.RB.delete(key) {
			return "ok"
		}
		return "error: not found"
	case TreeTypeBTree:
		if tc.BT.Delete(key) {
			return "ok"
		}
		return "error: not found"
	}
	return "error: unknown tree type"
}
