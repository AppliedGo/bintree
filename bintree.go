/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "A Binary Search Tree"
description = "A simple binary search tree in Go."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-08-04"
publishdate = "2016-08-04"
domains = ["Algorithms And Data Structures"]
tags = ["Tree", "Search Tree", "Binary Tree"]
categories = ["Tutorial"]
+++

This is about a simple binary search tree.

<!--more-->

## The Problem

Linear data containers (lists, arrays) have O(n) search complexity.

## The Solution

Search trees have O(log(n)) to O(n).

## Binary Search Tree Basics

* Node: Value, left child, right child
* Tree: One Root node, 0 or more child nodes
* All left children have smaller values, all right children have larger values than the root node's value
* Same applied to each subnode -> each subtree is also a tree

*/

// ## Imports and globals
package main

import (
	"errors"
	"fmt"
	"log"
)

// Aliases for the results of comparing two node values.
const (
	less = iota
	equal
	greater
)

/*
## The External Data Container

For flexibility we don't store the value itself but rather an index to an external data container. (I got this "trick" from [here](https://github.com/natefinch/tree).)

The external data container can hold any value of any type, as long as the type is comparable, and as long as it implements a comparison function that the tree data type can use.

The comparison function has its own type, `Comparer`. The external data container just needs to implement a function with the same function signature as `Comparer`. We then can pass the comparison function to every tree operation that requires comparison.

Let's define the Comparer function now, as well as a simple data container type, including its comparison function.

*/

// `Comparer` is a function that returns these constants:
// * `less` if the data at index i is less than that at index j,
// * `equal` if both are the same, and
// * `greater` otherwise.
type Comparer func(i, j int) int

// `Container` is the external data container. The underlying type can by any kind of container (array, slice, map,...) as long as its elements are comparable. Here, we simply use a slice of strings.
type Container []string

// `Compare` implements type `Comparer`.
func (c Container) Compare(i, j int) int {
	if c[i] < c[j] {
		return less
	}
	if c[i] == c[j] {
		return equal
	}
	return greater
}

/*
## A Tree Node

In general, a tree node consists of
* a value,
* a left subtree, and
* a right subtree.

This is a recursive data structure: Each subtree of a node is also a node containing subtrees.

In our setup, the "value" is in fact just the index to the value in the external data container, as explained above.
*/

// `Node` contains the data index, a left child node, and a right child node.
type Node struct {
	index int
	Left  *Node
	Right *Node
}

/* ## Node Operations

### Insert

Insert inserts a value...

*/

// `Insert` takes an index and a Comparer.
func (n *Node) Insert(i int, compare Comparer) error {

	if n == nil {
		return errors.New("Cannot insert a value into a nil tree")
	}

	// The Comparer must exist.
	if compare == nil {
		return errors.New("compare must not be nil")
	}

	// Compare the data at index `i` with the data at index `n.index`.
	switch compare(i, n.index) {
	// If the data is already in the tree, return.
	case equal:
		return nil
	// If the data value is less than the current node's value, and if the left child node is `nil`, insert a new left child node. Else call `Insert` on the left subtree.
	case less:
		if n.Left == nil {
			n.Left = &Node{index: i}
			return nil
		}
		return n.Left.Insert(i, compare)
	// If the data value is greater than the current node's value, do the same but for the right subtree.
	case greater:
		if n.Right == nil {
			n.Right = &Node{index: i}
			return nil
		}
		return n.Right.Insert(i, compare)
	}
	return nil
}

/*
### Find

Finding a value...

*/

// `Find` takes an index and a Comparer. It returns
// * The node containing the value, or
// * `nil` if no such node exists.
// The algorithm is pretty much the same as that of `Insert`.
func (n *Node) Find(i int, compare Comparer) (*Node, error) {

	if n == nil {
		return nil, nil
	}

	// As with `Insert`, the Comparer must exist.
	if compare == nil {
		return nil, errors.New("compare must not be nil")
	}

	// Compare the data at index `i` with the data at index `n.index`.
	switch compare(i, n.index) {
	// If the current node contains the value, return the node.
	case equal:
		return n, nil
	// If the data value is less than the current node's value, call `Find` for the left child node,
	case less:
		return n.Left.Find(i, compare)
	// else call `Find` for the right child node.
	case greater:
		return n.Right.Find(i, compare)
	}
	return nil, errors.New("compare() returned an unexpected value")
}

/*
### Delete

Deleting works like this...

...
*/

// `findMax` finds the maximum element in a (sub-)tree. Its value replaces the value of the
// to-be-deleted node.
// Return values: the node itself and its parent node.
func (n *Node) findMax(parent *Node) (*Node, *Node) {
	if n.Right == nil {
		return n, parent
	}
	return n.Right.findMax(n)
}

// `replaceNode` replaces the `parent`'s child pointer to `n` with a pointer to the `replacement` node.
func (n *Node) replaceNode(parent, replacement *Node) error {
	if n == nil {
		return errors.New("replaceNode() not allowed on a nil node")
	}
	if parent == nil {
		return nil
	}
	if n == parent.Left {
		parent.Left = replacement
	}
	parent.Right = replacement
	return nil
}

// `Delete` removes an element from the tree.
// It is an error to try deleting an element that does not exist.
func (n *Node) Delete(i int, parent *Node, compare Comparer) error {

	if n == nil {
		return errors.New("Value to be deleted does not exist in the tree")
	}

	// Again, the Comparer must exist.
	if compare == nil {
		return errors.New("compare must not be nil")
	}

	// Search the node to be deleted.
	switch compare(i, n.index) {
	case less:
		return n.Left.Delete(i, n, compare)
	case greater:
		return n.Right.Delete(i, n, compare)
	case equal:
		// We found the node to be deleted.
		// If the node has no children, simply remove it from its parent.
		if n.Left == nil && n.Right == nil {
			n.replaceNode(parent, nil)
			return nil
		}

		// If the node has one child: Replace the node with its child.
		if n.Left == nil {
			n.replaceNode(parent, n.Right)
			return nil
		}
		if n.Right == nil {
			n.replaceNode(parent, n.Left)
			return nil
		}

		// If the node has two children:
		// Find the maximum element in the left subtree...
		replacement, replParent := n.Left.findMax(n)

		//...and replace the node's data with the replacement's data.
		n.index = replacement.index

		// Then remove the replacement node by calling Delete on it.
		return replacement.Delete(i, replParent, compare)
	}
	return errors.New("compare() returned an unexpected value")
}

/*
## The Tree

One of a binary tree's nodes is the root node - the "entry point" of the tree.

The Tree data type wraps the root node into some bits of special treatment. Especially, it handles the case where the tree is completely empty.

The Tree data type also provides an additional function for traversing the whole tree.

*/

// A `Tree` basically consists of a root node.
type Tree struct {
	Root *Node
}

// `Insert` calls `Node.Insert` unless the root node is `nil`
func (t *Tree) Insert(i int, compare Comparer) error {
	// If the tree is empty, create a new node,...
	if t.Root == nil {
		t.Root = &Node{index: i}
		return nil
	}
	// ...else call `Node.Insert`.
	return t.Root.Insert(i, compare)
}

// `Find` calls `Node.Find` unless the root node is `nil`
func (t *Tree) Find(i int, compare Comparer) (*Node, error) {
	if t.Root == nil {
		return nil, nil
	}
	return t.Root.Find(i, compare)
}

// `Delete` calls `Node.Delete` unless the root node is `nil`
func (t *Tree) Delete(i int, compare Comparer) error {
	if t.Root == nil {
		return errors.New("Cannot delete from an empty tree")
	}
	return t.Root.Delete(i, nil, compare)
}

// `Traverse` is a simple method that traverses the tree in left-to-right order
// (which, *by pure incidence* ;-), is the same as traversing from smallest to
// largest value) and calls a custom function on each node.
func (t *Tree) Traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}

/* ## A Couple Of Tree Operations
 */

// `main`
func main() {

	// Set up a small container.
	container := Container{"delta", "bravo", "charlie", "echo", "alpha"}
	fmt.Println("Container: ", container)

	// Create a tree and fill it from the container.
	tree := &Tree{}
	for i := range container {
		err := tree.Insert(i, container.Compare)
		if err != nil {
			log.Fatal("Error inserting value at index ", i, ": ", err)
		}
	}

	// Print the sorted container.
	fmt.Print("Sorted container: [")
	tree.Traverse(tree.Root, func(n *Node) { fmt.Print(container[n.index], " ") })
	fmt.Println("]")

	// Find values by index.
	fmt.Print("Find by index: ")
	for i := range container {
		node, err := tree.Find(i, container.Compare)
		if err != nil {
			log.Fatal("Error during Find() at index ", i, ": ", err)
		}
		fmt.Print(container[node.index], " ")
	}
	fmt.Println()
}
