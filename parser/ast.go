package parser

import (
	"fmt"
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

func (tree *SyntaxTree) debug(){
	debug(tree, 0)
}
func debug(tree *SyntaxTree, nesting int){

	for _, child := range tree.head.children{

		fmt.Printf("[%d]", nesting)
		fmt.Printf(" %s\n", child.value.Literal)


		auxTree := NewSyntaxTree(child)
		debug(auxTree, nesting + 1)

	}
}