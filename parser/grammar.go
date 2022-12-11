package parser

import (
	"fmt"
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
type foundSymbol struct{
	symbol string
	src *[]token.Token
	tree *SyntaxTree
}
type log struct{
	nesting int
	help int

}

var Log *log

func (l *log) printLog(){
	i :=0
	for i < l.nesting{
		fmt.Printf("_")
		i++
	}
}


type GrammarSymbol interface{
	Build(*[]token.Token, *SyntaxTree) bool
	GetValue() string
}

//We say that a grammar symbol is a non terminal if and only if it's the head of a production. A production being a set of options
type NonTerminal struct{
	head string
	options []Option
}

//A terminal is a token type and a grammar symbol
type Terminal token.Type

//Each option is a sequence of Grammar Symbols
type Option struct{
	grammarSymbols []GrammarSymbol
}

//Build verifies if a given token is syntactically valid in a context and, if it is, it's added as head of the syntax tree that is being built
//then we slices the slice of tokens to move forward
//In case of having more than one terminal in the same nesting level we only keep the last one as head of the (sub)tree
//the grammar is made in such a way that this only happens with "() [] {}",  and the if/else statement
//(which we use to distinguish it from the if statement)
/*
	Example
    code: {
           if 1 != 0{
			  let foo *byte
            }else{
			   let foo bool
			}
		  }
	tree:
	           EOF         (all programs starts with EOF as head)
	            |
		        }
			    |
	           else         (the else represents the if/else. The left's child represents the condition,
	       /    |   \        the middle's child, the if block; and the right's child, the else block.
	    !=      }      }       In case of not having a simple if the parent of the subtree is going to be "if"
		/ \		|      |           and it is going to have only two children)
	   1	0   let     let
	           / \     /  \
	        foo  *   foo bool
				 |
                bool
*/
func (t Terminal) Build(src *[]token.Token, tree *SyntaxTree) bool{
	//Log.printLog()
	//fmt.Printf("SOURCE: %s WAITING: %s EQUAL: %t line: %d\n", (*src)[0].Literal, token.Type(t), (*src)[0].Type == token.Type(t), (*src)[0].Line)

	if (*src)[0].Type == token.Type(t) {
		//New lines don't have a purpose in our tree, so we skip them
		if token.Type(t) != token.NEWLINE {
			tree.head.value = (*src)[0]
		}
		*src = (*src)[1:]
		return true
	}
	return false
}

func (t Terminal) GetValue() string{
	return string(t)
}

//Build verifies if a sequence of tokens and non-terminals are syntactically valid in a context.
//If they are valid, they are added to the syntax tree that is being built, building for that subtree
//Then we replace the slice of token being analysed for the auxiliary one to move forward
func (nonT NonTerminal) Build(src *[]token.Token, tree *SyntaxTree) bool{
	empty := token.NewToken("", "", 0)
	Log.nesting++
	Log.help++
	var found bool
	foundSymbols := make([]foundSymbol,0)

	for _, option := range nonT.options{
		srcAux := *src
		auxTree := NewSyntaxTree(NewNode(empty))

		for j, symbol := range option.grammarSymbols {
			symbolValue := symbol.GetValue()

			if len(foundSymbols) > j {
				if foundSymbols[j].symbol == symbol.GetValue(){
					found = true
					srcAux = *foundSymbols[j].src
					auxTree = foundSymbols[j].tree
					continue
				}

			}
			found = symbol.Build(&srcAux, auxTree)
			if found {
				foundSymbols = append(foundSymbols, foundSymbol{symbol: symbolValue, src: &srcAux, tree :auxTree})
			}else{
				break
			}
		}
		if found{
			//Representations of non-terminals don't have a purpose in our tree, so we skip them to avoid empty nodes
			if auxTree.head.value == empty{
				for _, child := range auxTree.head.children{
					tree.head.AddChild(child)

				}
			}else{
				tree.head.AddChild(auxTree.head)
			}
			*src = srcAux
			Log.nesting--
			return true
		}
	}
	Log.nesting--
	return false
}

func (nonT *NonTerminal) GetValue() string{
	return nonT.head
}
//GetGrammar creates the grammar of the language.
func GetGrammar() map[string]*NonTerminal {

	Log = new(log)

	const PROGRAM = "program"
	const BLOCK = "block"
	const STATEMENTS = "statements"
	const STATEMENT = "statement"
	const DECLARATION = "declaration"
	const PARAM_DECLARATION = "paramdeclaration"
	const VAR = "var"
	const LITERAL = "literal"
	const ADDRESS = "address"
	const IDENT = "ident"
	const CALL = "call"
	const PARAMS = "params"
	const DATATYPE= "datatype"
	const ARGS= "args"
	const FUNC_DATATYPE = "funcdatatype"
	
	const EXPRESSION = "expression"
	const EXPRESSION_P9 = "expression_p9"
	const EXPRESSION_P8 = "expression_p8"
	const EXPRESSION_P7 = "expression_p7"
	const EXPRESSION_P6 = "expression_p6"
	const EXPRESSION_P5 = "expression_p5"
	const EXPRESSION_P4 = "expression_p4"
	const EXPRESSION_P3 = "expression_p3"
	const EXPRESSION_P2 = "expression_p2"
	const EXPRESSION_P1 = "expression_p1"
	const EXPRESSION_P0 = "expression_p0"


	productions := make(map[string]*NonTerminal)
	productions[PROGRAM] = new(NonTerminal)
	productions[BLOCK] = new(NonTerminal)
	productions[STATEMENTS] = new(NonTerminal)
	productions[STATEMENT] = new(NonTerminal)
	productions[DECLARATION] = new(NonTerminal)
	productions[PARAM_DECLARATION] = new(NonTerminal)
	productions[VAR] = new(NonTerminal)
	productions[LITERAL] = new(NonTerminal)
	productions[ADDRESS] = new(NonTerminal)
	productions[IDENT] = new(NonTerminal)
	productions[CALL] = new(NonTerminal)
	productions[PARAMS] = new(NonTerminal)
	productions[DATATYPE] = new(NonTerminal)
	productions[ARGS] = new(NonTerminal)
	productions[FUNC_DATATYPE] = new(NonTerminal)

	productions[EXPRESSION] = new(NonTerminal)
	productions[EXPRESSION_P9] = new(NonTerminal)
	productions[EXPRESSION_P8] = new(NonTerminal)
	productions[EXPRESSION_P7] = new(NonTerminal)
	productions[EXPRESSION_P6] = new(NonTerminal)
	productions[EXPRESSION_P5] = new(NonTerminal)
	productions[EXPRESSION_P4] = new(NonTerminal)
	productions[EXPRESSION_P3] = new(NonTerminal)
	productions[EXPRESSION_P2] = new(NonTerminal)
	productions[EXPRESSION_P1] = new(NonTerminal)
	productions[EXPRESSION_P0] = new(NonTerminal)

	grammarSymbols := make([]GrammarSymbol, 0)

	//PROGRAM:
	options := make([]Option, 1)

	grammarSymbols = append(grammarSymbols, productions[BLOCK])
	grammarSymbols = append(grammarSymbols, Terminal(token.EOF))
	options[0].grammarSymbols = grammarSymbols

	productions[PROGRAM].options = options
	productions[PROGRAM].head = PROGRAM

	//BLOCK:

	options = make([]Option, 1)
	grammarSymbols = make([]GrammarSymbol, 0)

	grammarSymbols = append(grammarSymbols, Terminal(token.LBRACE))
	grammarSymbols = append(grammarSymbols, productions[STATEMENTS])
	grammarSymbols = append(grammarSymbols, Terminal(token.RBRACE))
	options[0].grammarSymbols = grammarSymbols

	productions[BLOCK].options = options
	productions[BLOCK].head = BLOCK


	//STATEMENTS:
	options = make([]Option, 2)
	grammarSymbols = make([]GrammarSymbol, 0)


	grammarSymbols = append(grammarSymbols, productions[STATEMENT])
	grammarSymbols = append(grammarSymbols, productions[STATEMENTS])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[STATEMENT])
	options[1].grammarSymbols = grammarSymbols

	productions[STATEMENTS].options = options
	productions[STATEMENTS].head = STATEMENTS

	//STATEMENT
	options = make([]Option, 5)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.NEWLINE))
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[DECLARATION])
	grammarSymbols = append(grammarSymbols, Terminal(token.NEWLINE))
	options[1].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.FUNCTION))
	grammarSymbols = append(grammarSymbols, productions[IDENT])
	grammarSymbols = append(grammarSymbols, productions[ARGS])
	grammarSymbols = append(grammarSymbols, productions[FUNC_DATATYPE])
	grammarSymbols = append(grammarSymbols, productions[BLOCK])
	grammarSymbols = append(grammarSymbols, Terminal(token.NEWLINE))
	options[2].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[CALL])
	grammarSymbols = append(grammarSymbols, Terminal(token.NEWLINE))
	options[3].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[VAR])
	grammarSymbols = append(grammarSymbols, Terminal(token.EQ))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION])
	grammarSymbols = append(grammarSymbols, Terminal(token.NEWLINE))
	options[4].grammarSymbols = grammarSymbols

	productions[STATEMENT].options = options
	productions[STATEMENT].head = STATEMENT

	//DECLARATION
	options = make([]Option, 1)
	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LET))
	grammarSymbols = append(grammarSymbols, productions[IDENT])
	grammarSymbols = append(grammarSymbols, productions[DATATYPE])
	options[0].grammarSymbols = grammarSymbols
	productions[DECLARATION].options = options
	productions[DECLARATION].head = DECLARATION


	// IDENT
	options = make([]Option, 1)
	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.IDENT))
	options[0].grammarSymbols = grammarSymbols

	productions[IDENT].options = options
	productions[IDENT].head = IDENT

	// CALL
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[VAR])
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[VAR])
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, productions[PARAMS])
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[1].grammarSymbols = grammarSymbols



	productions[CALL].options = options
	productions[CALL].head = CALL

	//PARAM DECLARATION

	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[DECLARATION])
	grammarSymbols = append(grammarSymbols, Terminal(token.COMMA))
	grammarSymbols = append(grammarSymbols, productions[PARAM_DECLARATION])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[DECLARATION])
	options[1].grammarSymbols = grammarSymbols

	productions[PARAM_DECLARATION].options = options
	productions[PARAM_DECLARATION].head = PARAM_DECLARATION
	//PARAMS
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION])
	grammarSymbols = append(grammarSymbols, Terminal(token.COMMA))
	grammarSymbols = append(grammarSymbols, productions[PARAMS])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION])
	options[1].grammarSymbols = grammarSymbols

	productions[PARAMS].options = options
	productions[PARAMS].head = PARAMS

	// DATATYPE
	options = make([]Option, 4)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.ASTERISK)) //Pointer
	grammarSymbols = append(grammarSymbols, productions[DATATYPE])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LBRACKET))
	grammarSymbols = append(grammarSymbols, productions[LITERAL])
	grammarSymbols = append(grammarSymbols, Terminal(token.RBRACKET))
	grammarSymbols = append(grammarSymbols, productions[DATATYPE])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.TYPEBOOL))
	options[2].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.TYPEBYTE))
	options[3].grammarSymbols = grammarSymbols


	productions[DATATYPE].options = options
	productions[DATATYPE].head = DATATYPE

	//FUNC DATA TYPE:
	options = make([]Option, 2)


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.VOID))
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[DATATYPE])
	options[1].grammarSymbols = grammarSymbols

	productions[FUNC_DATATYPE].options = options
	productions[FUNC_DATATYPE].head = FUNC_DATATYPE


	//VAR:
	options = make([]Option, 4)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.ASTERISK))
	grammarSymbols = append(grammarSymbols, productions[VAR])

	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LBRACKET))
	grammarSymbols = append(grammarSymbols, productions[LITERAL])
	grammarSymbols = append(grammarSymbols, Terminal(token.RBRACKET))
	grammarSymbols = append(grammarSymbols, productions[VAR])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, productions[VAR])
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[2].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.IDENT))
	options[3].grammarSymbols = grammarSymbols


	productions[VAR].options = options
	productions[VAR].head = VAR


	// ADDRESS:

	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.DOLLAR))
	grammarSymbols = append(grammarSymbols, productions[VAR])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.DOLLAR))
	grammarSymbols = append(grammarSymbols, productions[ADDRESS])
	options[1].grammarSymbols = grammarSymbols

	productions[ADDRESS].options = options
	productions[ADDRESS].head = ADDRESS

	// LITERAL:

	options = make([]Option, 3)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.BYTE))
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.BOOL))
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[ADDRESS])
	options[2].grammarSymbols = grammarSymbols


	productions[LITERAL].options = options
	productions[LITERAL].head = LITERAL


	//ARGS:
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, productions[PARAM_DECLARATION])
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[1].grammarSymbols = grammarSymbols


	productions[ARGS].options = options
	productions[ARGS].head = ARGS

	//EXPRESSION_P0:
	//options = make([]Option, 4)
	options = make([]Option, 4)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[LITERAL])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[CALL])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[VAR])
	options[2].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, Terminal(token.LPAREN))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION])
	grammarSymbols = append(grammarSymbols, Terminal(token.RPAREN))
	options[3].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P0].options = options
	productions[EXPRESSION_P0].head = EXPRESSION_P0

	//EXPRESSION_P1
	options = make([]Option, 4)


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P0])
	grammarSymbols = append(grammarSymbols, Terminal(token.ASTERISK))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P0])
	grammarSymbols = append(grammarSymbols, Terminal(token.SLASH))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P0])
	grammarSymbols = append(grammarSymbols, Terminal(token.PERCENT))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	options[2].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P0])
	options[3].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P1].options = options
	productions[EXPRESSION_P1].head = EXPRESSION_P1

	//EXPRESSION_P2:
	options = make([]Option, 3)


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	grammarSymbols = append(grammarSymbols, Terminal(token.PLUS))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P2])

	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	grammarSymbols = append(grammarSymbols, Terminal(token.MINUS))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P2])

	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P1])
	options[2].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P2].options = options
	productions[EXPRESSION_P2].head = EXPRESSION_P2

	//EXPRESSION_P3
	options = make([]Option, 3)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P2])
	grammarSymbols = append(grammarSymbols, Terminal(token.GTGT))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P3])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P2])
	grammarSymbols = append(grammarSymbols, Terminal(token.LTLT))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P3])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P2])
	options[2].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P3].options = options
	productions[EXPRESSION_P3].head = EXPRESSION_P3

	//EXPRESSION_P4
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P3])
	grammarSymbols = append(grammarSymbols, Terminal(token.AND))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P4])
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P3])
	options[1].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P4].options = options
	productions[EXPRESSION_P4].head = EXPRESSION_P4

	//EXPRESSION_P5
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P4])
	grammarSymbols = append(grammarSymbols, Terminal(token.CIRC))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P5])

	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P4])
	options[1].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P5].options = options
	productions[EXPRESSION_P5].head = EXPRESSION_P5

	//EXPRESSION_P6
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P5])
	grammarSymbols = append(grammarSymbols, Terminal(token.OR))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P5])
	options[1].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P6].options = options
	productions[EXPRESSION_P6].head = EXPRESSION_P6

	//EXPRESSION_P7
	options = make([]Option, 5)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	grammarSymbols = append(grammarSymbols, Terminal(token.GT))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	grammarSymbols = append(grammarSymbols, Terminal(token.LT))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	grammarSymbols = append(grammarSymbols, Terminal(token.GTEQ))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	options[2].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	grammarSymbols = append(grammarSymbols, Terminal(token.LTEQ))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	options[3].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P6])
	options[4].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P7].options = options
	productions[EXPRESSION_P7].head = EXPRESSION_P7


	//EXPRESSION_P8
	options = make([]Option, 3)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	grammarSymbols = append(grammarSymbols, Terminal(token.EQEQ))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P8])
	options[0].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	grammarSymbols = append(grammarSymbols, Terminal(token.NOTEQ))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P8])
	options[1].grammarSymbols = grammarSymbols

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P7])
	options[2].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P8].options = options
	productions[EXPRESSION_P8].head = EXPRESSION_P8

	//EXPRESSION_P9
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P8])
	grammarSymbols = append(grammarSymbols, Terminal(token.LAND))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P9])
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P8])
	options[1].grammarSymbols = grammarSymbols

	productions[EXPRESSION_P9].options = options
	productions[EXPRESSION_P9].head = EXPRESSION_P9

	//EXPRESSION
	options = make([]Option, 2)

	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P9])
	grammarSymbols = append(grammarSymbols, Terminal(token.LOR))
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION])
	options[0].grammarSymbols = grammarSymbols


	grammarSymbols = make([]GrammarSymbol, 0)
	grammarSymbols = append(grammarSymbols, productions[EXPRESSION_P9])
	options[1].grammarSymbols = grammarSymbols

	productions[EXPRESSION].options = options
	productions[EXPRESSION].head = EXPRESSION

	return productions
}


