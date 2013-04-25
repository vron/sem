//line parser.y:2
package adr

import (
	"fmt"
)

// Global that we can store the outermost node in
var fullNode *Node
var err error

//line parser.y:14
type yySymType struct {
	yys  int
	node *Node
	/* Below this for the lexer */
	val int
	reg string
}

const (
	HASH    = 57346
	NUMBER  = 57347
	REG     = 57348
	DOLLAR  = 57349
	DOT     = 57350
	ADRMARK = 57351
	PLUS    = 57352
	MINUS   = 57353
	COMMA   = 57354
	SEMI    = 57355
	pBSTART = 57356
	pBEND   = 57357
)

var yyToknames = []string{
	"HASH",
	"NUMBER",
	"REG",
	"DOLLAR",
	"DOT",
	"ADRMARK",
	"PLUS",
	"MINUS",
	"COMMA",
	"SEMI",
	"pBSTART",
	"pBEND",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:170

/* Other stuff */
//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 26
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 57

var yyAct = []int{

	2, 17, 18, 24, 19, 20, 21, 22, 7, 23,
	17, 18, 15, 16, 1, 29, 25, 26, 27, 28,
	10, 9, 11, 12, 13, 14, 5, 6, 3, 4,
	8, 10, 9, 11, 12, 13, 14, 5, 6, 0,
	0, 8, 10, 9, 11, 12, 13, 14, 0, 0,
	0, 0, 8, 17, 18, 15, 16,
}
var yyPact = []int{

	16, -1000, 43, 27, 27, 38, 38, -1000, 16, -1000,
	-2, -1000, -1000, -1000, -1000, 27, 27, 38, 38, -9,
	-9, -1000, -1000, 0, -1000, -9, -9, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 14, 0, 8,
}
var yyR1 = []int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 3, 2, 2, 1, 3, 2, 2, 1,
	3, 2, 2, 1, 3, 2, 2, 1, 1, 3,
	1, 2, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -1, -2, 12, 13, 10, 11, -3, 14, 5,
	4, 6, 7, 8, 9, 12, 13, 10, 11, -2,
	-2, -2, -2, -2, 5, -2, -2, -2, -2, 15,
}
var yyDef = []int{

	0, -2, 1, 5, 9, 13, 17, 18, 0, 20,
	0, 22, 23, 24, 25, 4, 8, 12, 16, 3,
	7, 11, 15, 0, 21, 2, 6, 10, 14, 19,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line parser.y:38
		{
			fullNode = yyS[yypt-0].node
		}
	case 2:
		//line parser.y:42
		{
			yyVAL.node = &Node{Type: COMMA, Left: yyS[yypt-2].node, Right: yyS[yypt-0].node}
		}
	case 3:
		//line parser.y:44
		{
			yyVAL.node = &Node{
				Type:  COMMA,
				Left:  &Node{Type: NUMBER, Val: 0},
				Right: yyS[yypt-0].node,
			}
		}
	case 4:
		//line parser.y:52
		{
			yyVAL.node = &Node{
				Type:  COMMA,
				Left:  yyS[yypt-1].node,
				Right: &Node{Type: DOLLAR},
			}
		}
	case 5:
		//line parser.y:60
		{
			yyVAL.node = &Node{
				Type:  COMMA,
				Left:  &Node{Type: NUMBER, Val: 0},
				Right: &Node{Type: DOLLAR},
			}
		}
	case 6:
		//line parser.y:69
		{
			yyVAL.node = &Node{Type: SEMI, Left: yyS[yypt-2].node, Right: yyS[yypt-0].node}
		}
	case 7:
		//line parser.y:71
		{
			yyVAL.node = &Node{
				Type:  SEMI,
				Left:  &Node{Type: NUMBER, Val: 0},
				Right: yyS[yypt-0].node,
			}
		}
	case 8:
		//line parser.y:79
		{
			yyVAL.node = &Node{
				Type:  SEMI,
				Left:  yyS[yypt-1].node,
				Right: &Node{Type: DOLLAR},
			}
		}
	case 9:
		//line parser.y:87
		{
			yyVAL.node = &Node{
				Type:  SEMI,
				Left:  &Node{Type: NUMBER, Val: 0},
				Right: &Node{Type: DOLLAR},
			}
		}
	case 10:
		//line parser.y:96
		{
			yyVAL.node = &Node{Type: PLUS, Left: yyS[yypt-2].node, Right: yyS[yypt-0].node}
		}
	case 11:
		//line parser.y:98
		{
			yyVAL.node = &Node{
				Type:  PLUS,
				Left:  &Node{Type: DOT},
				Right: yyS[yypt-0].node,
			}
		}
	case 12:
		//line parser.y:106
		{
			yyVAL.node = &Node{
				Type:  PLUS,
				Left:  yyS[yypt-1].node,
				Right: &Node{Type: NUMBER, Val: 1},
			}
		}
	case 13:
		//line parser.y:114
		{
			yyVAL.node = &Node{
				Type:  PLUS,
				Left:  &Node{Type: DOT},
				Right: &Node{Type: NUMBER, Val: 1},
			}
		}
	case 14:
		//line parser.y:123
		{
			yyVAL.node = &Node{Type: MINUS, Left: yyS[yypt-2].node, Right: yyS[yypt-0].node}
		}
	case 15:
		//line parser.y:125
		{
			yyVAL.node = &Node{
				Type:  MINUS,
				Left:  &Node{Type: DOT},
				Right: yyS[yypt-0].node,
			}
		}
	case 16:
		//line parser.y:133
		{
			yyVAL.node = &Node{
				Type:  MINUS,
				Left:  yyS[yypt-1].node,
				Right: &Node{Type: NUMBER, Val: 1},
			}
		}
	case 17:
		//line parser.y:141
		{
			yyVAL.node = &Node{
				Type:  MINUS,
				Left:  &Node{Type: DOT},
				Right: &Node{Type: NUMBER, Val: 1},
			}
		}
	case 18:
		//line parser.y:149
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 19:
		//line parser.y:151
		{
			yyVAL.node = yyS[yypt-2].node
		}
	case 20:
		//line parser.y:155
		{
			yyVAL.node = &Node{Type: NUMBER, Val: yyS[yypt-0].val}
		}
	case 21:
		//line parser.y:157
		{
			yyVAL.node = &Node{Type: HASH, Val: yyS[yypt-0].val}
		}
	case 22:
		//line parser.y:159
		{
			yyVAL.node = &Node{Type: REG, Reg: yyS[yypt-0].reg}
		}
	case 23:
		//line parser.y:161
		{
			yyVAL.node = &Node{Type: DOLLAR}
		}
	case 24:
		//line parser.y:163
		{
			yyVAL.node = &Node{Type: DOT}
		}
	case 25:
		//line parser.y:165
		{
			yyVAL.node = &Node{Type: ADRMARK}
		}
	}
	goto yystack /* stack new state and value */
}
