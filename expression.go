package main

import (
	"math"
)

func parseExpression(tokens *exprList, t float64) (float64, bool) {
	var (
		left, right float64
	)

	ff, ok := term(tokens, t)
	if !ok {
		return 0, false
	}
	left = ff
	for temp := tokens.currentToken();
		temp.Match(PLUS) || temp.Match(MINUS);
	temp = tokens.currentToken() {
		tokens.matchToken(temp.TokenType)
		ff, ok = term(tokens, t)
		if !ok {
			return 0, false
		}
		right = ff

		left = clac(temp.TokenType, left, right)
	}

	return left, true
}

func term(tokens *exprList, t float64) (float64, bool) {
	var (
		left, right float64
	)

	ff, ok := factor(tokens, t)
	if !ok {
		return 0, false
	}
	left = ff
	for temp := tokens.currentToken();
		temp.Match(MUL) || temp.Match(DIV);
	temp = tokens.currentToken() {
		tokens.matchToken(temp.TokenType)
		ff, ok := factor(tokens, t)
		if !ok {
			return 0, false
		}
		right = ff

		left = clac(temp.TokenType, left, right)
	}

	return left, true
}

func factor(tokens *exprList, t float64) (result float64, ok bool) {
	var ff float64
	ok = true
	temp := tokens.currentToken()
	if temp.Match(PLUS) {
		tokens.matchToken(PLUS)
		ff, ok = factor(tokens, t)
		if !ok {
			return 0, false
		}
		result = ff
	} else if temp.Match(MINUS) {
		tokens.matchToken(MINUS)
		ff, ok = factor(tokens, t)
		if !ok {
			return 0, false
		}
		result = 0 - ff
	} else {
		ff, ok = component(tokens, t)
		if !ok {
			return 0, false
		}
		result = ff
	}

	return
}

func component(tokens *exprList, t float64) (float64, bool) {
	var (
		left, right float64
	)

	ff, ok := atom(tokens, t)
	if !ok {
		return 0, false
	}
	left = ff

	if temp := tokens.currentToken(); temp.Match(POWER) {
		tokens.matchToken(POWER)
		ff, ok = component(tokens, t)
		if !ok {
			return 0, false
		}
		right = ff

		left = clac(POWER, left, right)
	}
	return left, true
}

func atom(tokens *exprList, t float64) (result float64, ok bool) {
	temp := tokens.currentToken()
	var ff float64
	ok = true
	switch temp.TokenType {
	case CONST_ID :
		tokens.matchToken(CONST_ID)
		result = temp.ConstValue
	case T:
		tokens.matchToken(T)
		result = t
	case FUNC:
		tokens.matchToken(FUNC)
		tokens.matchToken(L_BRACKET)
		ff, ok = parseExpression(tokens, t)
		if !ok {
			return 0, false
		}
		result = temp.FuncHandler(ff)
		tokens.matchToken(R_BRACKET)
	case L_BRACKET:
		tokens.matchToken(L_BRACKET)
		ff, ok = parseExpression(tokens, t)
		if !ok {
			return 0, false
		}
		result = ff
		tokens.matchToken(R_BRACKET)
	default:
		return
	}

	return
}

func clac(optr int, x, y float64) (result float64) {
	switch optr {
	case PLUS:
		result = x + y
	case MINUS:
		result = x - y
	case MUL:
		result = x * y
	case DIV:
		result = x / y
	case POWER:
		result = math.Pow(x, y)
	}
	return
}