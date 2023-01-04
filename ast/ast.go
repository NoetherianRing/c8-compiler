package ast

import (
	"github.com/NoetherianRing/c8-compiler/token"
)

type SyntaxTree struct{
	Head    *Node
	Current *Node
}

type Node struct{
	Children []*Node
	Parent *Node
	Value    token.Token
}

func NewSyntaxTree(head *Node) *SyntaxTree {
	tree := new(SyntaxTree)
	tree.Head = head
	tree.Current = head
	return tree
}

func NewNode(value token.Token) *Node {
	node := new(Node)
	node.Value = value
	node.Parent = nil
	node.Children = make([]*Node,0)
	return node
}

func (node *Node) AddChild(child *Node){
	child.Parent = node
	node.Children = append(node.Children, child)
}
