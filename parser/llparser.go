package parser

import (
	"fmt"
	"github.com/OJoklrO/Interpreter/scanner"
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

type Vector2 struct {
	X, Y float64
}

func NewVector2(x, y float64) Vector2 {
	return Vector2{
		X: x,
		Y: y,
	}
}

type exprList struct {
	tokens []scanner.Token
	tokenIndex int
}

func (e *exprList) currentToken() *scanner.Token {
	if e.tokenIndex >= len(e.tokens) {
		return nil
	}
	return &e.tokens[e.tokenIndex]
}

func (e *exprList) matchToken(targetType int) {
	res :=  e.currentToken().Match(targetType)
	e.tokenIndex++

	if !res {
		logError("exprList.matchToken")
	}
}

func newExprList(tokens []scanner.Token) *exprList {
	return &exprList{
		tokens:     tokens,
		tokenIndex: 0,
	}
}

type ForStmtParam struct {
	start  float64
	end    float64
	step   float64
	Points []Vector2
	tvalue []float64
}

func (f *ForStmtParam) clacT() {
	for i := math.Min(f.start, f.end); i <= math.Max(f.start, f.end); i += math.Abs(f.step) {
		f.tvalue = append(f.tvalue, i)
	}
}

type OriginStmtParam struct {
	Pos Vector2
}

type ScaleStmtParam struct {
	Scale Vector2
}

type RotStmtParam struct {
	Rot float64
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:   newExprList(tokens),
		StmtType: -1,
		F:        ForStmtParam{},
		O:        OriginStmtParam{},
		S:        ScaleStmtParam{},
		R:        RotStmtParam{},
	}
}

func (p *Parser) Test() {
	switch p.StmtType {
	case FORSTMT:
		fmt.Println("for steatment: ", p.F)
	case ORIGINSTMT:
		fmt.Println("origin steatment: ", p.O)
	case SCALESTMT:
		fmt.Println("Scale steatment: ", p.S)
	case ROTSTMT:
		fmt.Println("Rot steatment: ", p.R)
	default:
		logError("parser.Test")
	}
}

func (p *Parser) Execute() {

}

func (p *Parser) Parse() {
	switch p.tokens.currentToken().TokenType {
	case scanner.FOR:
		p.forParse()
		p.StmtType = FORSTMT
	case scanner.ORIGIN:
		p.originParse()
		p.StmtType = ORIGINSTMT
	case scanner.SCALE:
		p.scaleParse()
		p.StmtType = SCALESTMT
	case scanner.ROT:
		p.rotParse()
		p.StmtType = ROTSTMT
	case scanner.RESET:
		p.StmtType = RESETSTMT
	}
}

func (p *Parser) forParse() {
	p.tokens.matchToken(scanner.FOR)
	p.tokens.matchToken(scanner.T)
	p.tokens.matchToken(scanner.FROM)

	p.F.start = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.TO)

	p.F.end = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.STEP)

	p.F.step = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.DRAW)
	p.tokens.matchToken(scanner.L_BRACKET)

	p.F.clacT()
	pos := p.tokens.tokenIndex
	for _, v := range p.F.tvalue {
		p.tokens.tokenIndex = pos
		x := parseExpression(p.tokens, v)

		p.tokens.matchToken(scanner.COMMA)

		y := parseExpression(p.tokens, v)
		p.F.Points = append(p.F.Points, NewVector2(x, y))
	}

	p.tokens.matchToken(scanner.R_BRACKET)
	p.tokens.matchToken(scanner.SEMICO)
}

func (p *Parser) originParse() {
	p.tokens.matchToken(scanner.ORIGIN)
	p.tokens.matchToken(scanner.IS)
	p.tokens.matchToken(scanner.L_BRACKET)

	p.O.Pos.X = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.COMMA)

	p.O.Pos.Y = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.R_BRACKET)
	p.tokens.matchToken(scanner.SEMICO)
}

func (p *Parser) scaleParse() {
	p.tokens.matchToken(scanner.SCALE)
	p.tokens.matchToken(scanner.IS)
	p.tokens.matchToken(scanner.L_BRACKET)

	p.S.Scale.X = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.COMMA)

	p.S.Scale.Y = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.R_BRACKET)
	p.tokens.matchToken(scanner.SEMICO)
}

func (p *Parser) rotParse() {
	p.tokens.matchToken(scanner.ROT)
	p.tokens.matchToken(scanner.IS)

	p.R.Rot = parseExpression(p.tokens, 0)

	p.tokens.matchToken(scanner.SEMICO)
}

func logError(who string) {
	fmt.Println(who, " log error")
	//panic("asd")
}

func parseExpression(tokens *exprList, t float64) float64 {
	//fmt.Println("expr: ", tokens.tokenIndex)

	var (
		left, right float64
	)

	left = term(tokens, t)
	for temp := tokens.currentToken();
		temp.Match(scanner.PLUS) || temp.Match(scanner.MINUS);
		temp = tokens.currentToken() {
		tokens.matchToken(temp.TokenType)
		right = term(tokens, t)

		left = clac(temp.TokenType, left, right)
	}

	return left
}

func term(tokens *exprList, t float64) float64 {
	//fmt.Println("term: ", tokens.tokenIndex)

	var (
		left, right float64
	)

	left = factor(tokens, t)
	for temp := tokens.currentToken();
		temp.Match(scanner.MUL) || temp.Match(scanner.DIV);
		temp = tokens.currentToken() {
		tokens.matchToken(temp.TokenType)
		right = factor(tokens, t)

		left = clac(temp.TokenType, left, right)
	}

	return left
}

func factor(tokens *exprList, t float64) (result float64) {
	//fmt.Println("factor: ", tokens.tokenIndex)

	temp := tokens.currentToken()
	if temp.Match(scanner.PLUS) {
		tokens.matchToken(scanner.PLUS)
		result = factor(tokens, t)
	} else if temp.Match(scanner.MINUS) {
		tokens.matchToken(scanner.MINUS)
		result = 0 - factor(tokens, t)
	} else {
		result = component(tokens, t)
	}

	return
}

func component(tokens *exprList, t float64) float64 {
	//fmt.Println("comp: ", tokens.tokenIndex)

	var (
		left, right float64
	)

	left = atom(tokens, t)

	if temp := tokens.currentToken(); temp.Match(scanner.POWER) {
		tokens.matchToken(scanner.POWER)
		right = component(tokens, t)

		left = clac(scanner.POWER, left, right)
	}
	return left
}

func atom(tokens *exprList, t float64) (result float64) {
	//fmt.Println("atom: ", tokens.tokenIndex)

	temp := tokens.currentToken()

	switch temp.TokenType {
	case scanner.CONST_ID :
		tokens.matchToken(scanner.CONST_ID)
		result = temp.ConstValue
	case scanner.T:
		tokens.matchToken(scanner.T)
		result = t
	case scanner.FUNC:
		tokens.matchToken(scanner.FUNC)
		tokens.matchToken(scanner.L_BRACKET)
		result = temp.FuncHandler(parseExpression(tokens, t))
		tokens.matchToken(scanner.R_BRACKET)
	case scanner.L_BRACKET:
		tokens.matchToken(scanner.L_BRACKET)
		result = parseExpression(tokens, t)
		tokens.matchToken(scanner.R_BRACKET)
	default:
		logError("atom")
	}

	return
}

func clac(optr int, x, y float64) (result float64) {
	switch optr {
	case scanner.PLUS:
		result = x + y
	case scanner.MINUS:
		result = x - y
	case scanner.MUL:
		result = x * y
	case scanner.DIV:
		result = x / y
	case scanner.POWER:
		result = math.Pow(x, y)
	}
	return
}