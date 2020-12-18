package main

import (
	"fmt"
	"github.com/OJoklrO/Interpreter/vector2"
	"math"
)

const (
	FORSTMT = iota
	ORIGINSTMT
	SCALESTMT
	ROTSTMT
	RESETSTMT
)

type Parser struct {
	tokens   *exprList
	StmtType int
	F        ForStmtParam
	O        OriginStmtParam
	S        ScaleStmtParam
	R        RotStmtParam
}

func NewVector2(x, y float64) vector2.Vector2 {
	return vector2.Vector2{
		X: x,
		Y: y,
	}
}

type exprList struct {
	tokens []Token
	tokenIndex int
}

func (e *exprList) currentToken() *Token {
	if e.tokenIndex >= len(e.tokens) {
		return nil
	}
	return &e.tokens[e.tokenIndex]
}

func (e *exprList) matchToken(targetType int) (res bool) {
	res =  e.currentToken().Match(targetType)
	e.tokenIndex++
	return
}

func newExprList(tokens []Token) *exprList {
	return &exprList{
		tokens:     tokens,
		tokenIndex: 0,
	}
}

type ForStmtParam struct {
	start  float64
	end    float64
	step   float64
	Points []vector2.Vector2
	tvalue []float64
}

func (f *ForStmtParam) clacT() {
	for i := math.Min(f.start, f.end); i <= math.Max(f.start, f.end); i += math.Abs(f.step) {
		f.tvalue = append(f.tvalue, i)
	}
}

type OriginStmtParam struct {
	Pos vector2.Vector2
}

type ScaleStmtParam struct {
	Scale vector2.Vector2
}

type RotStmtParam struct {
	Rot float64
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   newExprList(tokens),
		StmtType: -1,
		F:        ForStmtParam{},
		O:        OriginStmtParam{},
		S:        ScaleStmtParam{},
		R:        RotStmtParam{},
	}
}

func (p *Parser) Parse() bool {
	switch p.tokens.currentToken().TokenType {
	case FOR:
		if !p.forParse() {
			return false
		}
		p.StmtType = FORSTMT
	case ORIGIN:
		if !p.originParse() {
			return false
		}
		p.StmtType = ORIGINSTMT
	case SCALE:
		if !p.scaleParse() {
			return false
		}
		p.StmtType = SCALESTMT
	case ROT:
		if !p.rotParse() {
			return false
		}
		p.StmtType = ROTSTMT
	case RESET:
		p.StmtType = RESETSTMT
	default:
		return false
	}
	return true
}

func (p *Parser) forParse() bool {
	if !p.checkQueue([]int{FOR, T, FROM}) {
		return false
	}

	ff, ok := parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.F.start = ff

	if !p.checkQueue([]int{TO}) {
		return false
	}

	ff, ok = parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.F.end = ff

	if !p.checkQueue([]int{STEP}) {
		return false
	}

	ff, ok = parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.F.step = ff

	if !p.checkQueue([]int{DRAW, L_BRACKET}) {
		return false
	}

	p.F.clacT()
	pos := p.tokens.tokenIndex
	for _, v := range p.F.tvalue {
		p.tokens.tokenIndex = pos
		x, ok := parseExpression(p.tokens, v)
		if !ok {
			return false
		}

		if !p.checkQueue([]int{COMMA}) {
			return false
		}

		y, ok := parseExpression(p.tokens, v)
		if !ok {
			return false
		}
		p.F.Points = append(p.F.Points, NewVector2(x, y))
	}

	if !p.checkQueue([]int{R_BRACKET, SEMICO}) {
		return false
	}

	return true
}

func (p *Parser) originParse() bool {
	if !p.checkQueue([]int{ORIGIN, IS, L_BRACKET}) {
		return false
	}

	ff, ok := parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.O.Pos.X = ff

	if !p.checkQueue([]int{COMMA}) {
		return false
	}

	ff, ok = parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.O.Pos.Y = ff

	if !p.checkQueue([]int{R_BRACKET, SEMICO}) {
		return false
	}
	return true
}

func (p *Parser) scaleParse() bool {
	if !p.checkQueue([]int{SCALE, IS, L_BRACKET}) {
		return false
	}

	ff, ok := parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.S.Scale.X = ff

	if !p.checkQueue([]int{COMMA}) {
		return false
	}

	ff, ok = parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.S.Scale.Y = ff

	if !p.checkQueue([]int{R_BRACKET, SEMICO}) {
		return false
	}

	return true
}

func (p *Parser) rotParse() bool {
	if !p.checkQueue([]int{ROT, IS}) {
		return false
	}

	ff, ok := parseExpression(p.tokens, 0)
	if !ok {
		return false
	}
	p.R.Rot = ff

	if !p.checkQueue([]int{SEMICO}) {
		return false
	}

	return true
}

func (p *Parser) checkQueue(queue []int) bool {
	for _, q := range queue {
		if !p.tokens.matchToken(q) {
			return false
		}
	}
	return true
}

func (p *Parser) LogError() {
	i := 0
	for ; i < p.tokens.tokenIndex; i++ {
		fmt.Printf("%s ", p.tokens.tokens[i].Value)
	}
	fmt.Printf("\033[1;31;40m%s\033[0m","!")
	for ; i < len(p.tokens.tokens); i++ {
		fmt.Printf("%s ", p.tokens.tokens[i].Value)
	}
}