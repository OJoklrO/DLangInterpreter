package main

import (
	"github.com/OJoklrO/Interpreter/dfa"
	"math"
	"strconv"
	"strings"

)

const (
	ID = iota
	CONST_ID

	ORIGIN
	SCALE
	ROT
	IS
	TO
	STEP
	DRAW
	FOR
	FROM

	T

	DELIMITER
	SEMICO
	L_BRACKET
	R_BRACKET
	COMMA

	OPERATOR
	PLUS
	MINUS
	MUL
	DIV
	POWER

	FUNC

	RESET

	NONTOKEN
)

var TokenMap = map[string]tokenValue{
	"origin" : {
		tokenType: ORIGIN,
	},
	"scale" : {
		tokenType: SCALE,
	},
	"rot" : {
		tokenType: ROT,
	},
	"is" : {
		tokenType: IS,
	},
	"for" : {
		tokenType: FOR,
	},
	"from" : {
		tokenType: FROM,
	},
	"to" : {
		tokenType: TO,
	},
	"step" : {
		tokenType: STEP,
	},
	"draw" : {
		tokenType: DRAW,
	},
	"reset" : {
		tokenType: RESET,
	},
	"sin" : {
		tokenType: FUNC,
		funcHandler: math.Sin,
	},
	"cos" : {
		tokenType: FUNC,
		funcHandler: math.Cos,
	},
	"tan" : {
		tokenType: FUNC,
		funcHandler: math.Tan,
	},
	"ln" : {
		tokenType: FUNC,
		funcHandler: math.Log,
	},
	"exp" : {
		tokenType: FUNC,
		funcHandler: math.Exp,
	},
	"sqrt" : {
		tokenType: FUNC,
		funcHandler: math.Sqrt,
	},

	"pi" : {
		tokenType: CONST_ID,
		constValue: math.Pi,
	},
	"e" : {
		tokenType: CONST_ID,
		constValue: math.E,
	},

	"t" : {
		tokenType: T,
	},

	"," : {
		tokenType: COMMA,
	},
	";" : {
		tokenType: SEMICO,
	},
	"(" : {
		tokenType: L_BRACKET,
	},
	")" : {
		tokenType: R_BRACKET,
	},

	"+" : {
		tokenType: PLUS,
	},
	"-" : {
		tokenType: MINUS,
	},
	"*" : {
		tokenType: MUL,
	},
	"/" : {
		tokenType: DIV,
	},
	"**" : {
		tokenType: POWER,
	},
}

type tokenValue struct {
	tokenType int
	constValue float64
	funcHandler func(float64)float64
}

type Token struct {
	TokenType   int
	ConstValue  float64
	FuncHandler func(float64)float64
	Value       string
}

func NewToken(tokenType int, value string) interface{} {
	switch tokenType {
	case DELIMITER, OPERATOR:
		return Token{
			TokenType: TokenMap[strings.ToLower(value)].tokenType,
			Value:     value,
		}
	case ID:
		info, ok := TokenMap[strings.ToLower(value)]
		if !ok {
			return Token{
				TokenType: NONTOKEN,
				Value:     value,
			}
		} else if info.tokenType == FUNC {
			return Token{
				TokenType:   FUNC,
				FuncHandler: info.funcHandler,
				Value:       value,
			}
		} else {
			return Token{
				TokenType: info.tokenType,
				Value:     value,
			}
		}
	case CONST_ID:
		info, ok := TokenMap[strings.ToLower(value)]
		if ok {
			return Token{
				TokenType:  CONST_ID,
				ConstValue: info.constValue,
				Value:      value,
			}
		} else {
			f, _ := strconv.ParseFloat(value, 64)
			return Token{
				TokenType:  CONST_ID,
				ConstValue: f,
				Value:      value,
			}
		}
	case dfa.ErrorToken:
		return Token{
			TokenType: NONTOKEN,
			Value:     value,
		}
	}
	return nil
}

func (t *Token) Match(targetType int) bool {
	return t.TokenType == targetType
}

type Lexer struct {
	d *dfa.DFA
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Input(s string) (result []Token) {
	verified := true
	for _, r := range s {
		if r != ' ' {
			if l.d.Input(r) {
				l.d.Input(r)                // back
			}
			verified = false
		} else {
			if !verified {
				l.d.Verify()
				verified = true
			}
		}
	}
	l.d.Verify()

	for _, t := range l.d.GetResult() {
		result = append(result, t.(Token))
	}

	return
}

// create dfa states and transitions
func (l *Lexer) Init() {
	l.d = dfa.NewDFA(0, NewToken)
	l.d.AddState(1, ID)
	l.d.AddState(2, CONST_ID)
	l.d.AddState(3, ID)
	l.d.AddState(4, -1)
	l.d.AddState(5, CONST_ID)
	l.d.AddState(6, OPERATOR)
	l.d.AddState(7, OPERATOR)
	l.d.AddState(8, DELIMITER)
	l.d.AddState(9, CONST_ID)
	l.d.AddState(10, -1)
	l.d.AddState(11, CONST_ID)

	// state 0 transition
	l.d.AddTransition(0, 1, func(r rune) bool {
		if isLetter(r) {
			if r == 'e' || r == 'p' {
				return false
			}
			return true
		}
		return false
	})
	l.d.AddTransition(0, 2, isDigit)
	l.d.AddTransition(0, 6, func(r rune) bool {
		return r == '+' || r == '-' || r == '/'
	})
	l.d.AddTransition(0, 7, func(r rune) bool {
		return r == '*'
	})
	l.d.AddTransition(0, 8, func(r rune) bool {
		return r == ',' || r == ';' || r == '(' || r == ')'
	})
	l.d.AddTransition(0, 9, func(r rune) bool {
		return r == 'e'
	})
	l.d.AddTransition(0, 10, func(r rune) bool {
		return r == 'p'
	})

	// state 1 transition
	l.d.AddTransition(1, 3, func(r rune) bool {
		return isLetter(r) || isDigit(r)
	})

	// state 2 transition
	l.d.AddTransition(2, 2, func(r rune) bool {
		return isDigit(r)
	})
	l.d.AddTransition(2, 4, func(r rune) bool {
		return r == '.'
	})

	// state 3 transition
	l.d.AddTransition(3, 3, func(r rune) bool {
		return isLetter(r) || isDigit(r)
	})

	// state 4 transition
	l.d.AddTransition(4, 5, func(r rune) bool {
		return isDigit(r)
	})

	// state 5 transition
	l.d.AddTransition(5, 5, func(r rune) bool {
		return isDigit(r)
	})

	// state 7 transition
	l.d.AddTransition(7, 6, func(r rune) bool {
		return r == '*'
	})

	// state 9 transition
	l.d.AddTransition(9, 1, func(r rune) bool {
		return isDigit(r) || isLetter(r)
	})

	// state 10 transition
	l.d.AddTransition(10, 1, func(r rune) bool {
		if isLetter(r) {
			if r == 'i'{
				return false
			}
			return true
		} else if isDigit(r) {
			return true
		}
		return false
	})
	l.d.AddTransition(10, 11, func(r rune) bool {
		return r == 'i'
	})

	l.d.Reset()
}

// input check helper
func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}