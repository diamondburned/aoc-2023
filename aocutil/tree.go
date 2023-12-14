package aocutil

import "slices"

// DFS is a depth-first search iterator. It takes a root node and a function
// that returns the children of a node. It returns an iterator that yields
// nodes in depth-first order.
func DFS[T any](root T, children func(T) []T) Iter[T] {
	return func(yield func(T) bool) {
		stack := []T{root}
		for len(stack) > 0 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if !yield(node) {
				break
			}

			for _, child := range children(node) {
				stack = append(stack, child)
			}
		}
	}
}

// BFS is a breadth-first search iterator. It takes a root node and a function
// that returns the children of a node. It returns an iterator that yields
// nodes in breadth-first order.
func BFS[T any](root T, children func(T) []T) Iter[T] {
	return func(yield func(T) bool) {
		var z T

		queue := []T{root}
		for len(queue) > 0 {
			node := queue[0]

			// Prevent memory leaks.
			queue[0] = z
			queue = slices.Delete(queue, 0, 1)

			if !yield(node) {
				break
			}

			for _, child := range children(node) {
				queue = append(queue, child)
			}
		}
	}
}

// TreeNode is a generic tree node.
type TreeNode[T any] struct {
	Value    T
	Children []TreeNode[T]
}

// NewTreeNode creates a new tree node with the given value and children.
func NewTreeNode[T any](value T, children ...TreeNode[T]) TreeNode[T] {
	return TreeNode[T]{
		Value:    value,
		Children: children,
	}
}

// Tree is a generic tree.
type Tree[T any] struct {
	Root TreeNode[T]
}

// NewTree creates a new tree with the given root and children function.
// The children function is recursively applied to the root node to create
// the tree.
func NewTree[T any](root T, children func(T) []T) Tree[T] {
	return Tree[T]{
		Root: TreeNode[T]{
			Value: root,
			Children: Map(children(root), func(child T) TreeNode[T] {
				return NewTree(child, children).Root
			}),
		},
	}
}

// DFS is a depth-first search iterator. It returns an iterator that yields
// nodes in depth-first order.
func (t Tree[T]) DFS() Iter[T] {
	dfs := DFS[TreeNode[T]](t.Root, func(node TreeNode[T]) []TreeNode[T] {
		return node.Children
	})
	return func(yield func(T) bool) {
		for node := range dfs {
			if !yield(node.Value) {
				break
			}
		}
	}
}

// BFS is a breadth-first search iterator. It returns an iterator that yields
// nodes in breadth-first order.
func (t Tree[T]) BFS() Iter[T] {
	bfs := BFS[TreeNode[T]](t.Root, func(node TreeNode[T]) []TreeNode[T] {
		return node.Children
	})
	return func(yield func(T) bool) {
		for node := range bfs {
			if !yield(node.Value) {
				break
			}
		}
	}
}
