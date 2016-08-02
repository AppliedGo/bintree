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

// Aliases for the results of comparing two values.
const (
	less    = -1
	equal   = 0
	greater = 1
)

/*
## A Tree Node

In general, a tree node consists of
* a value,
* a left subtree, and
* a right subtree.

This is a recursive data structure: Each subtree of a node is also a node containing subtrees.

In this minimal setup, the tree contains simple string data.
*/

// `Node` contains the data, a left child node, and a right child node.
type Node struct {
	Value string
	Left  *Node
	Right *Node
}

/* ## Node Operations

### Insert

Insert inserts a value...

*/

// `Insert` inserts a new string value into the tree.
// Return values:
// * `true` if the string was successfully inserted,
// * `false` if the string value already exists in the tree.
func (n *Node) Insert(s string) error {

	if n == nil {
		return errors.New("Cannot insert a value into a nil tree")
	}

	// Compare the data.
	switch {
	// If the data is already in the tree, return.
	case s == n.Value:
		return nil
	// If the data value is less than the current node's value, and if the left child node is `nil`, insert a new left child node. Else call `Insert` on the left subtree.
	case s < n.Value:
		if n.Left == nil {
			n.Left = &Node{Value: s}
			return nil
		}
		return n.Left.Insert(s)
	// If the data value is greater than the current node's value, do the same but for the right subtree.
	case s > n.Value:
		if n.Right == nil {
			n.Right = &Node{Value: s}
			return nil
		}
		return n.Right.Insert(s)
	}
	return nil
}

/*
### Find

Finding a value...

*/

// `Find` searches for a string. It returns:
// * The node containing the value, or
// * `nil` if no such node exists.
func (n *Node) Find(s string) (*Node, error) {

	if n == nil {
		return nil, nil
	}

	switch {
	// If the current node contains the value, return the node.
	case s == n.Value:
		return n, nil
	// If the data value is less than the current node's value, call `Find` for the left child node,
	case s < n.Value:
		return n.Left.Find(s)
		// else call `Find` for the right child node.
	default:
		return n.Right.Find(s)
	}
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
// In order to remove an element properly, `Delete` needs to know the node's parent node.
func (n *Node) Delete(s string, parent *Node) error {

	if n == nil {
		return errors.New("Value to be deleted does not exist in the tree")
	}

	// Search the node to be deleted.
	switch {
	case s < n.Value:
		return n.Left.Delete(s, n)
	case s > n.Value:
		return n.Right.Delete(s, n)
	default:
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
		n.Value = replacement.Value

		// Then remove the replacement node.
		return replacement.Delete(s, replParent)
	}
	return errors.New("compare() returned an unexpected value")
}

/*
## The Tree

One of a binary tree's nodes is the root node - the "entry point" of the tree.

The Tree data type wraps the root node and applies some special treatment. Especially, it handles the cases where the tree is completely empty or consists of a single node.

The Tree data type also provides an additional function for traversing the whole tree.

*/

// A `Tree` basically consists of a root node.
type Tree struct {
	Root *Node
}

// `Insert` calls `Node.Insert` unless the root node is `nil`
func (t *Tree) Insert(s string) error {
	// If the tree is empty, create a new node,...
	if t.Root == nil {
		t.Root = &Node{Value: s}
		return nil
	}
	// ...else call `Node.Insert`.
	return t.Root.Insert(s)
}

// `Find` calls `Node.Find` unless the root node is `nil`
func (t *Tree) Find(s string) (*Node, error) {
	if t.Root == nil {
		return nil, nil
	}
	return t.Root.Find(s)
}

// `Delete` calls `Node.Delete` unless the root node is `nil`
func (t *Tree) Delete(s string) error {

	// Special case 1: empty tree.
	if t.Root == nil {
		return errors.New("Cannot delete from an empty tree")
	}

	// Special case 2: tree consists of root node only, and root node contains
	// the value to be removed.
	if s == t.Root.Value && t.Root.Left == nil && t.Root.Right == nil {
		t.Root = nil
	}

	// In all other cases, call `Node.Delete`.
	return t.Root.Delete(s, nil)
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

	// Set up a slice of strings.
	values := []string{"delta", "bravo", "charlie", "echo", "alpha"}
	fmt.Println("Values: ", values)

	// Create a tree and fill it from the values.
	tree := &Tree{}
	for _, s := range values {
		err := tree.Insert(s)
		if err != nil {
			log.Fatal("Error inserting value '", s, "': ", err)
		}
	}

	// Print the sorted values.
	fmt.Print("Sorted values: [")
	tree.Traverse(tree.Root, func(n *Node) { fmt.Print(n.Value, " ") })
	fmt.Println("]")

	// Find values.
	s := "delta"
	fmt.Print("Find node of '", s, "': ")
	node, err := tree.Find(s)
	if err != nil {
		log.Fatal("Error during Find(): ", err)
	}
	fmt.Printf("%v\n", node)
}
