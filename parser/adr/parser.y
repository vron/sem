%{
package adr

import (
	"fmt"
)

// Global that we can store the outermost node in
var fullNode *Node
var err error
%}


%union{ 
	node *Node
	/* Below this for the lexer */
	val int
	reg string
}


/* Tokens */

%token HASH NUMBER REG DOLLAR DOT ADRMARK
%token PLUS MINUS
%token COMMA SEMI
%token pBSTART pBEND

/* Define precedence order */
%left COMMA SEMI
%left PLUS MINUS

%start fulladr

%% /* Grammar follows */

fulladr : adr
			{fullNode = $<node>1}
		;

adr 	: adr COMMA adr
				{$<node>$ = &Node{Type: COMMA, Left: $<node>1, Right: $<node>3 }}
		| COMMA adr
				{
					$<node>$ = &Node{
									Type: COMMA,
									Left: &Node{Type:NUMBER, Val: 0}, 
									Right: $<node>2,
								}
				}
		| adr COMMA
				{
					$<node>$ = &Node{
									Type: COMMA,
									Left: $<node>1,
									Right: &Node{Type:DOLLAR},
								}
				}
		| COMMA
				{
					$<node>$ = &Node{
									Type: COMMA,
									Left: &Node{Type:NUMBER, Val: 0}, 
									Right: &Node{Type:DOLLAR},
								}
				}
				
		| adr SEMI adr
				{$<node>$ = &Node{Type: SEMI, Left: $<node>1, Right: $<node>3 }}
		| SEMI adr
				{
					$<node>$ = &Node{
									Type: SEMI,
									Left: &Node{Type:NUMBER, Val: 0}, 
									Right: $<node>2,
								}
				}
		| adr SEMI
				{
					$<node>$ = &Node{
									Type: SEMI,
									Left: $<node>1,
									Right: &Node{Type:DOLLAR}, 
								}
				}
		| SEMI
				{
					$<node>$ = &Node{
									Type: SEMI,
									Left: &Node{Type:NUMBER, Val: 0}, 
									Right: &Node{Type:DOLLAR}, 
								}
				}
				
		| adr PLUS adr
				{$<node>$ = &Node{Type: PLUS, Left: $<node>1, Right: $<node>3 }}
		| PLUS adr
				{
					$<node>$ = &Node{
									Type: PLUS,
									Left: &Node{Type:DOT}, 
									Right: $<node>2,
								}
				}
		| adr PLUS
				{
					$<node>$ = &Node{
									Type: PLUS,
									Left: $<node>1,
									Right: &Node{Type:NUMBER, Val: 1}, 
								}
				}
		| PLUS
				{
					$<node>$ = &Node{
									Type: PLUS,
									Left: &Node{Type:DOT}, 
									Right: &Node{Type:NUMBER, Val: 1}, 
								}
				}
				
		| adr MINUS adr
				{$<node>$ = &Node{Type: MINUS, Left: $<node>1, Right: $<node>3 }}
		| MINUS adr
				{
					$<node>$ = &Node{
									Type: MINUS,
									Left: &Node{Type:DOT}, 
									Right: $<node>2,
								}
				}
		| adr MINUS
				{
					$<node>$ = &Node{
									Type: MINUS,
									Left: $<node>1,
									Right: &Node{Type:NUMBER, Val: 1}, 
								}
				}
		| MINUS
				{
					$<node>$ = &Node{
									Type: MINUS,
									Left: &Node{Type:DOT}, 
									Right: &Node{Type:NUMBER, Val: 1}, 
								}
				}
		| simpadr
				{$<node>$ = $<node>1}
		| pBSTART adr pBEND
				{$<node>$ = $<node>1}
		; 
	
simpadr	: NUMBER
				{$<node>$ = &Node{Type:NUMBER, Val: $<val>1}}
		| HASH NUMBER
				{$<node>$ = &Node{Type:HASH, Val: $<val>2}}
		| REG
				{$<node>$ = &Node{Type:REG, Reg: $<reg>1}}
		| DOLLAR
				{$<node>$ = &Node{Type:DOLLAR}}
		| DOT
				{$<node>$ = &Node{Type:DOT}}
		| ADRMARK
				{$<node>$ = &Node{Type:ADRMARK}}
		;



%% /* Other stuff */