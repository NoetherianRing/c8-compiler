package parser

import (
	"github.com/NoetherianRing/c8-compiler/token"
)

type SyntaxTree struct{
	head *Node
	current *Node
}

type Node struct{
	children []*Node
	value token.Token
}

func NewSyntaxTree(head *Node) *SyntaxTree{
	tree := new(SyntaxTree)
	tree.head = head
	tree.current = head
	return tree
}

func NewNode(value token.Token) *Node{
	node := new(Node)
	node.value = value
	node.children = make([]*Node,0)
	return node
}

func (node *Node) AddChild(child *Node){
	node.children = append(node.children, child)
}
