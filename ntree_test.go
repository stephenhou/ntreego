package ntree_test

import (
	"fmt"
	"testing"

	ntree "github.com/blazingorb/ntreego"
)

type MockData struct {
	Id    string
	Value int
}

var benchroot = GenerateBenchmarkTree(ntree.New(&MockData{"Root", 1}), 6, 12)

func GenerateTree() map[string]*ntree.Node {
	a_1 := &MockData{"a_1", 0}
	a_2 := &MockData{"a_2", 0}

	a_1_1 := &MockData{"a_1_1", 0}
	a_1_2 := &MockData{"a_1_2", 0}
	a_1_3 := &MockData{"a_1_3", 0}

	a_2_1 := &MockData{"a_2_1", 0}
	a_2_2 := &MockData{"a_2_2", 0}

	nodes := make(map[string]*ntree.Node)
	nodes["root"] = ntree.New(&MockData{"armor", 0})
	nodes["a_1"] = ntree.AppendChild(nodes["root"], ntree.New(a_1))
	nodes["a_2"] = ntree.AppendChild(nodes["root"], ntree.New(a_2))

	nodes["a_1_1"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_1))
	nodes["a_1_2"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_2))
	nodes["a_1_3"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_3))

	nodes["a_2_1"] = ntree.AppendChild(nodes["a_2"], ntree.New(a_2_1))
	nodes["a_2_2"] = ntree.AppendChild(nodes["a_2"], ntree.New(a_2_2))

	return nodes
}

func TestCreateTree(t *testing.T) {
	var initialValue int = 1
	tree := ntree.New(initialValue)
	value, ok := tree.Value.(int)
	if !ok || value != initialValue {
		t.Errorf("Tree Creation Failed")
	}
}

func TestTreeSearch(t *testing.T) {
	nodes := GenerateTree()
	toSearch := nodes["a_2_2"]
	parentOfToSearch := nodes["a_2"]

	var resultNode *ntree.Node
	searchFunc := func(n *ntree.Node, value interface{}) bool {
		v := value.(string)
		if n.Value.(*MockData).Id == v {
			resultNode = n
			return true
		}
		return false
	}

	ntree.Traverse(nodes["root"], ntree.TraverseInOrder, ntree.TraverseAll, -1, searchFunc, toSearch.Value.(*MockData).Id)
	if toSearch != resultNode {
		t.Error("Search Failed", toSearch.Value, resultNode)
	}
	if toSearch.Parent != parentOfToSearch {
		t.Error("Unexpected Parent", toSearch, parentOfToSearch)
	}
}

func TestTreeModify(t *testing.T) {
	nodes := GenerateTree()
	addFunc := func(n *ntree.Node, value interface{}) bool {
		v := value.(int)
		n.Value.(*MockData).Value += v
		return false
	}

	ntree.Traverse(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, -1, addFunc, 5)

	for _, node := range nodes {
		if node.Value.(*MockData).Value != 5 {
			t.Error("Value of node is wrong!")
		}
	}

	str := nodes["root"].String()
	fmt.Println(str)
}

func TestNodeCount(t *testing.T) {
	nodes := GenerateTree()
	count := ntree.NodeCount(nodes["root"], ntree.TraverseAll)
	if count != len(nodes) {
		t.Error("Mismatched node count!")
	}

	if ntree.NodeCount(nil, ntree.TraverseAll) != 0 {
		t.Error("NodeCount should be zero when the tree assigned is nil!")
	}

	if ntree.NodeCount(nodes["root"], ntree.TraverseMask+1) != 0 {
		t.Error("NodeCount should be zero when flags is larger than TraverseMask")
	}
}

func TestUnlink(t *testing.T) {
	nodes := GenerateTree()

	ntree.Unlink(nil)

	if ntree.NodeCount(nodes["root"], ntree.TraverseAll) != len(nodes) {
		t.Error("NodeCount should remained the same when calling ntree.Unlink(nil)")
	}

	wrapperFunc := func(n *ntree.Node) {
		ntree.Unlink(n)
		if n.Parent != nil || n.Previous != nil || n.Next != nil {
			t.Error("Incomplete unlink!")
		}

		if ntree.FindNode(nodes["root"], ntree.TraverseInOrder, ntree.TraverseAll, n.Value) != nil {
			t.Error("Incomplete unlink!")
		}

		for _, node := range nodes {
			if node.Children == n || node.Previous == n || node.Next == n {
				t.Error("Incomplete unlink!")
			}
		}
	}

	wrapperFunc(nodes["a_1"])
	wrapperFunc(nodes["a_2_2"])

}

func TestGetRoot(t *testing.T) {
	nodes := GenerateTree()
	root, depth := ntree.GetRoot(nil)
	if root != nil || depth != 0 {
		t.Error("root should be nil and depth should be 0 when ntree.GetRoot(nil) is called!")
	}

	wrapperFunc := func(node *ntree.Node) {
		root, depth = ntree.GetRoot(node)
		if root == nil {
			t.Error("root == nil")
		} else {
			if root.Value != nodes["root"].Value {
				t.Error("Mismatched root node!")
			}
		}

		if depth != ntree.Depth(node) {
			t.Errorf("Depth of %s expected to be %d but return %d", node.Value.(*MockData).Id, ntree.Depth(node), depth)
		}
	}

	wrapperFunc(nodes["root"])
	wrapperFunc(nodes["a_1"])
	wrapperFunc(nodes["a_1_1"])
}

func TestDepth(t *testing.T) {
	nodes := GenerateTree()
	depth := 0

	wrapperFunc := func(node *ntree.Node, expectedLength int) {
		depth = ntree.Depth(node)
		if depth != expectedLength {
			t.Errorf("Depth of %s expected to be %d but return %d", node.Value.(*MockData).Id, expectedLength, depth)
		}
	}

	wrapperFunc(nodes["root"], 1)
	wrapperFunc(nodes["a_1"], 2)
	wrapperFunc(nodes["a_1_1"], 3)
}

func TestInsert(t *testing.T) {
	nodes := GenerateTree()
	a_1 := nodes["a_1"]
	a_2 := nodes["a_2"]
	a_2_1 := nodes["a_2_1"]

	newNode := ntree.New(1)

	if ntree.Insert(nil, newNode) != nil || ntree.Insert(nodes["root"], nil) != nil {
		t.Error("Result should be nil when one of the arguments is nil")
	}

	if ntree.Insert(a_1, a_2) != nil {
		t.Error("Result should be nil when a non-root node is inserted as a child of other node")
	}

	ntree.Insert(nodes["root"], newNode)
	if a_2.Next != newNode || newNode.Previous != a_2 {
		t.Error("newNode should be a sibling of a_2")
	}

	newNode2 := ntree.New(2)
	ntree.Insert(a_2_1, newNode2)
	if a_2_1.Children != newNode2 || newNode2.Parent != a_2_1 {
		t.Error("newNode2 should be a child of a_2_1")
	}
}

func TestAppendChild(t *testing.T) {
	nodes := GenerateTree()
	a_1 := nodes["a_1"]
	a_2 := nodes["a_2"]
	a_2_1 := nodes["a_2_1"]

	newNode := ntree.New(1)

	if ntree.AppendChild(nil, newNode) != nil || ntree.AppendChild(nodes["root"], nil) != nil {
		t.Error("Result should be nil when one of the arguments is nil")
	}

	if ntree.AppendChild(a_1, a_2) != nil {
		t.Error("Result should be nil when a non-root node is appended as a child of other node")
	}

	ntree.AppendChild(nodes["root"], newNode)
	if a_2.Next != newNode || newNode.Previous != a_2 {
		t.Error("newNode should be a sibling of a_2")
	}

	newNode2 := ntree.New(2)
	ntree.AppendChild(a_2_1, newNode2)
	if a_2_1.Children != newNode2 || newNode2.Parent != a_2_1 {
		t.Error("newNode2 should be a child of a_2_1")
	}
}

func TestIsRoot(t *testing.T) {
	nodes := GenerateTree()
	if !ntree.IsRoot(nodes["root"]) {
		t.Error("Result should be true when root node is passed in")
	}

	if ntree.IsRoot(nodes["a_1"]) || ntree.IsRoot(nodes["a_2"]) {
		t.Error("Result should be false when non-root nodes are passed in")
	}
}

func TestFindNode(t *testing.T) {
	nodes := GenerateTree()

	if ntree.FindNode(nil, ntree.TraverseInOrder, ntree.TraverseAll, nodes["a_1_1"].Value) != nil {
		t.Error("Result should be nil when nil is passed as root argument!")
	}

	nodeFound := ntree.FindNode(nodes["root"], ntree.TraverseInOrder, ntree.TraverseAll, nodes["a_1_1"].Value)

	if nodeFound != nodes["a_1_1"] {
		t.Error("Wrong node has be found!")
	}

	nodeFound = ntree.FindNode(nodes["a_1"], ntree.TraverseInOrder, ntree.TraverseAll, nodes["a_1_1"].Value)

	if nodeFound != nodes["a_1_1"] {
		t.Error("Wrong node has be found!")
	}
}

func TestTraverseFail(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	depth_a_1_1 := ntree.Depth(nodes["a_1_1"])

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFailFunc := func(root *ntree.Node, order ntree.TraverseType, flag ntree.TraverseFlags, depth int, traverseFunc func(n *ntree.Node, data interface{}) bool) {
		visitCount = 0
		ntree.Traverse(root, order, flag, depth, traverseFunc, 0)

		if visitCount != 0 {
			t.Error("Traverse Error! Visit count should be zero when argument meets the fail condition!")
		}
	}

	wrappedFailFunc(nil, ntree.TraversePreOrder, ntree.TraverseAll, depth_a_1_1, traverseFunc)
	wrappedFailFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, depth_a_1_1, nil)
	wrappedFailFunc(nodes["root"], ntree.TraverseLevelOrder+1, ntree.TraverseAll, depth_a_1_1, traverseFunc)
	wrappedFailFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseMask+1, depth_a_1_1, traverseFunc)
	wrappedFailFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, -2, traverseFunc)
	wrappedFailFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, 0, traverseFunc)
}
func TestTraverseAll(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseAll, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2_2"], len(nodes))
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2_2"], len(nodes))
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["root"], len(nodes))
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseAllwithConditions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseAll, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, depth_a_1_1)
}

func TestTraverseLeaves(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2_2"], 5)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2_2"], 5)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_2_2"], 5)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nil, 0)
	wrappedFunc(ntree.TraverseInOrder, 2, nil, 0)
	wrappedFunc(ntree.TraversePostOrder, 2, nil, 0)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseLeavesWithConditions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nil, 0)
	wrappedFunc(ntree.TraverseInOrder, 2, nil, 0)
	wrappedFunc(ntree.TraversePostOrder, 2, nil, 0)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseNonLeaves(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseNonLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseNonLeavesWithCondistions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseNonLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Errorf("Traverse Error! The expected visit count should be %d but return %d", expectedLength, visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return %s", expectedNode.Value.(*MockData).Id, lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be nil but return %s", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Errorf("Traverse Error! The expected node that traverseFunc last visited should be %s but return nil", expectedNode.Value.(*MockData).Id)
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func GenerateBenchmarkTree(parentNode *ntree.Node, depthLimit int, childCount int) *ntree.Node {
	if ntree.Depth(parentNode) < depthLimit {
		for i := 0; i < childCount; i++ {
			childNode := ntree.New(&MockData{"Mock", 1})
			ntree.AppendChild(parentNode, childNode)
			GenerateBenchmarkTree(childNode, depthLimit, childCount)
		}
	}

	return parentNode
}

func BenchmarkTreeTraverse(b *testing.B) {
	counter := 0
	traverseFunc := func(n *ntree.Node, value interface{}) bool {
		counter++
		return false
	}

	b.Logf("Node Count: %d \n", ntree.NodeCount(benchroot, ntree.TraverseAll))

	for n := 0; n < b.N; n++ {
		counter = 0
		ntree.Traverse(benchroot, ntree.TraverseInOrder, ntree.TraverseAll, -1, traverseFunc, "Mock")
	}
}
