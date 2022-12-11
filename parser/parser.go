package parser
/*
import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/token"
)

type Parser struct{
	tokens      []token.Token
	lookUpTable map[token.TokenType] *ast.Node
	index       int
}

func NewParser(tokens []token.Token) *Parser {
	p := new(Parser)
	p.lookUpTable = make(map[token.TokenType] *ast.Node, 0)
	let := ast.NewNode(token.LET)
	p.lookUpTable[token.LET] = let
	p.tokens = tokens
	return p
}

func (p *Parser) nextStatement()  ([]*ast.Node, error){
	statement := make([]*ast.Node,0)
	current := p.tokens[p.index]
	next := p.tokens[p.index+1]
	node := p.lookUpTable[current.GetValue]

	//TODO: En el if, el while y las funciones puedo tener lookUpTable
	for node.Value != token.SEMICOLON{
		if !node.IsConnectedTo(next.GetValue){
			//TODO: Explicar mejor el error
			return statement, errors.New("error in line: " + string(rune(current.Line)) + ". '" + current.Literal + " " + next.Literal+"'")
		}
		statement = append(statement, node)
		p.index++
		current = p.tokens[p.index]
		next = p.tokens[p.index+1]
		node = p.lookUpTable[current.GetValue]

	}
	p.index++

	return statement, nil
}


func (p *Parser) GetStatements() ([][]*ast.Node, error){
	statements := make([][]*ast.Node,0)
	for p.tokens[p.index].GetValue != token.EOF{
		nextStatement, err := p.nextStatement()
		if err != nil{
			return statements, err
		}
		statements = append(statements, nextStatement)
	}
	return statements, nil
}
*/