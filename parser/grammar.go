package parser

import (
	"github.com/NoetherianRing/c8-compiler/token"
)
/*
Definition of Grammars according to The Dragon Book:

1. A set of terminal symbols, sometimes referred to as "tokens." The
terminals are the elementary symbols of the language denned by the grammar.

2. A set of nonterminals, sometimes called "syntactic variables." Each
nonterminal represents a set of strings of terminals.

3. A set of productions, where each production consists of a nonterminal,
called the head or left side of the production, an arrow, and a sequence of terminals and/or
non terminals, called the body or right side of the production.
The intuitive intent of a production is to specify one of the written
forms of a construct; if the head nonterminal represents a construct, then
the body represents a written form of the construct.

4. A designation of one of the nonterminals as the start symbol.

Example:
(1)	list -> list + digit | list - digit | digit

(2)	digit  -> 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

In here "list" and "digit" are nonterminals;  + - and the numbers from 0 to 9 are terminals; and (1) and (2) are productions.
The different options of bodies in each production are separated by the symbol |.
*/
/*
type Grammar struct{
	head *NonTerminal
	productions map[string]*NonTerminal
}
*/
/*
type log struct{
	nesting int

}

var Log *log

func (l *log) printLog(){
	i :=0
	for i < l.nesting{
		fmt.Printf("_")
		i++
	}
}
*/

type GrammarSymbol interface{
	//Check([]token.Token, *int) bool
	Build(*[]token.Token, *SyntaxTree) bool
}

//All nonterminals are heads of a production with one or more options of bodies. A nonterminal is a grammar symbol
type NonTerminal struct{
	options []Option
}

//A terminal is a token and a node
type Terminal token.TokenType

//Each option is a sequence of Grammar Symbols
type Option struct{
	grammarSymbols []GrammarSymbol
}
/*
//Build verifies if a given token is syntactically valid in a context.
func (t Terminal)  Check(src []token.Token, index *int) bool{

//	fmt.Printf("%d)", *index)
	if *index >= len(src){
		return false
	}
	Log.printLog()

	fmt.Printf("index: %d SOURCE: %s WAITING: %s EQUAL: %t\n", *index ,src[*index].Literal, token.TokenType(t), src[*index].Type == token.TokenType(t))
	if src[*index].Type == token.TokenType(t){

		*index++
		return true
	}
//	return src[*index].Type == token.TokenType(t)
	return false
}

//Build verifies if a sequence of tokens and nonterminals are syntactically valid in a context.
//If its valid, it moves the index of the slice of tokens that is being analysed
func (nonT NonTerminal) Check(src []token.Token, index *int) bool{
	Log.nesting++
	var auxIndex int
	var found bool
	for _, option := range nonT.options{
		auxIndex = *index
		found = true
		for _, node := range option.grammarSymbols{
			if !node.Build(src, &auxIndex) {
				found = false
				break
			}
			//auxIndex++
			fmt.Println(auxIndex)

		}

		if found {
			fmt.Println("found")

			*index = auxIndex
			Log.nesting--
			return true
		}
	}
	Log.nesting--
	return false
}
*/
//Check verifies if a given token is syntactically valid in a context and, if it is, it's added to the syntax tree that is being built
func (t Terminal) Build(src *[]token.Token, tree *SyntaxTree) bool{

	if len(*src)==0{
		return false
	}
	//Log.printLog()

	//fmt.Printf("SOURCE: %s WAITING: %s EQUAL: %t\n", (*src)[0].Literal, token.TokenType(t), (*src)[0].Type == token.TokenType(t))
	if (*src)[0].Type == token.TokenType(t){
		tree.head.value = (*src)[0]
		*src = (*src)[1:]
		return true
	}
	return false
}

//Check verifies if a sequence of tokens and nonterminals are syntactically valid in a context.
//If they are valid, they are added to the syntax tree that is being built
func (nonT NonTerminal) Build(src *[]token.Token, tree *SyntaxTree) bool{
//	Log.nesting++
	var found bool
	for _, option := range nonT.options{
		srcAux := *src
		dumb := token.NewToken("", "", 0)
		auxTree := NewSyntaxTree(NewNode(dumb))
		for _, node := range option.grammarSymbols {
			found = node.Build(&srcAux, auxTree)
			if !found {
				break
			}
		}
		if found {
		//	fmt.Println("found")
			if auxTree.head.value != dumb{
				tree.head.AddChild(auxTree.head)
			}else{
				for _, child := range auxTree.head.children{
					tree.head.AddChild(child)
				}
			}
			*src = srcAux
		//	Log.nesting--
			return true
		}
	}
//	Log.nesting--
	return false
}

//GetGrammar creates the grammar of the language.
func GetGrammar() map[string]*NonTerminal {

	//Log = new(log)
	productions := make(map[string]*NonTerminal)
	productions["statements"] = new(NonTerminal)
	productions["expression"] = new(NonTerminal)
	productions["term"] = new(NonTerminal)
	productions["factor"] = new(NonTerminal)

	options := make([]Option, 0)
	grammarSymbols := make([]GrammarSymbol, 0)


	//STATEMENTS
	options = make([]Option, 2)
	grammarSymbols = make([]GrammarSymbol, 0)


	grammarSymbols = append(grammarSymbols, productions["expression"]) //TODO: Change to statement
	grammarSymbols = append(grammarSymbols, productions["statements"])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.EOF))
	options[1].grammarSymbols = grammarSymbols
	productions["program"].options = options



	//EXPRESSION:
	options = make([]Option, 3)


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["term"])
	grammarSymbols = append(grammarSymbols, Terminal(token.PLUS))
	grammarSymbols = append(grammarSymbols, productions["expression"])

	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["term"])

	grammarSymbols = append(grammarSymbols, Terminal(token.MINUS))
	grammarSymbols = append(grammarSymbols, productions["expression"])

	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["term"])
	options[2].grammarSymbols = grammarSymbols

	productions["expression"].options = options

	//TERM
	options = make([]Option, 4)


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["factor"])

	grammarSymbols = append(grammarSymbols, Terminal(token.ASTERISK))
	grammarSymbols = append(grammarSymbols, productions["term"])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["factor"])
	grammarSymbols = append(grammarSymbols, Terminal(token.SLASH))
	grammarSymbols = append(grammarSymbols, productions["term"])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["factor"])
	grammarSymbols = append(grammarSymbols, Terminal(token.PERCENT))
	grammarSymbols = append(grammarSymbols, productions["term"])
	options[2].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions["factor"])
	options[3].grammarSymbols = grammarSymbols

	productions["term"].options = options

	//FACTOR:
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.INT))
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, productions["expression"])
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[1].grammarSymbols = grammarSymbols



	productions["factor"].options = options

	//	productions["statement"] = new(NonTerminal)
	//	productions["statement"].options = make([]Option,0)
/*
	return &Grammar{
	head:
		productions["stmts"],
	productions:
		productions,
	}
*/
	return productions
}


