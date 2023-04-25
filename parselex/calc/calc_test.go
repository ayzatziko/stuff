package calc

import (
	"fmt"
	"strconv"
	"testing"
)

func TestTest(t *testing.T) {
	tt := []struct {
		i string
		v int
	}{
		{i: "5", v: 5},
		{i: "5+2", v: 7},
		{i: "5 + 2*3", v: 11},
		{i: "5*2 + 3", v: 13},
		{i: "5*(2+3)", v: 25},
		{i: "(5*2)+3", v: 13},
		{i: "((5*2)+3)", v: 13},
		{i: "7+((5*2)+3)", v: 20},
		{i: "((5*2)+3)+7", v: 20},
	}

	for _, tc := range tt {
		tokens, err := parse(tc.i)
		if err != nil {
			t.Fatal(err)
		}
		v, err := eval(tokens)
		if err != nil {
			t.Fatal(err)
		}

		if *v != tc.v {
			t.Fatalf("want %v, got %v", tc.v, *v)
		}
	}
}

func parse(s string) ([]token, error) {
	return parseRecursive(s, []token{})
}

func parseRecursive(s string, tt []token) ([]token, error) {
	switch {
	case len(s) == 0:
		return tt, nil
	case s[0] == '(':
		tt = append(tt, token{"(", openTok})
		return parseRecursive(s[1:], tt)
	case s[0] == ')':
		tt = append(tt, token{")", closeTok})
		return parseRecursive(s[1:], tt)
	case isArithChar(s[0]):
		return parseRecursive(s[1:], append(tt, arithChars[s[0]]))
	case s[0] == ' ':
		return parseRecursive(s[1:], tt)
	case isNextNum(s):
		n, numtok := parseNum(s)
		return parseRecursive(s[n:], append(tt, numtok))
	}

	return nil, fmt.Errorf("cannot parse %s", s)
}

type token struct {
	v string

	typ int
}

const (
	openTok  int = 0
	closeTok int = 1
	arithTok int = 2
	numToken int = 3
)

func isArithChar(b byte) bool {
	_, ok := arithChars[b]
	return ok
}

var arithChars map[byte]token = map[byte]token{
	'+': {"+", arithTok},
	'-': {"-", arithTok},
	'/': {"/", arithTok},
	'*': {"*", arithTok},
}

func isNextNum(s string) bool {
	for _, v := range []byte(s) {
		if v == ' ' || v == '(' || v == ')' || isArithChar(v) {
			return true
		}
		if _, ok := numset[v]; !ok {
			return false
		}
	}

	return true
}

var numset = map[byte]struct{}{
	'0': {},
	'1': {},
	'2': {},
	'3': {},
	'4': {},
	'5': {},
	'6': {},
	'7': {},
	'8': {},
	'9': {},
}

func parseNum(s string) (int, token) {
	n := 0
	for _, v := range []byte(s) {
		if v == ' ' || v == '(' || v == ')' || isArithChar(v) {
			return n, token{s[:n], numToken}
		}
		if _, ok := numset[v]; !ok {
			panic("must not happen")
		}
		n++
	}

	return len(s), token{s, numToken}
}

// “
// 5         = 5
// 5+2       = 7
// 5 + 2*3   = 11
// 5*2 + 3   = 13
// 5*(2+3)   = 25
// (5*2)+3   = 13
// ((5*2)+3) = 13
func eval(tokens []token) (_ *int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("tokens %v, error %v", tokens, err)
		}
	}()

	switch {
	case len(tokens) == 0:
		return nil, nil
	case tokens[0].typ == openTok:
		i, err := findClosingRecursive(tokens[1:], 0, 0)
		if err != nil {
			return nil, err
		}
		i += 1
		v0, err := eval(tokens[1:i])
		if err != nil {
			return nil, err
		}
		switch {
		case len(tokens) == i+1:
			return v0, nil
		default:
			return eval(append([]token{{typ: numToken, v: strconv.Itoa(*v0)}}, tokens[i+1:]...))
		}
	case len(tokens) == 1 && tokens[0].typ == numToken:
		v, err := strconv.Atoi(tokens[0].v)
		if err != nil {
			return nil, fmt.Errorf("token %v, error %v", tokens[0], err)
		}

		return &v, nil
	case len(tokens) > 2 && tokens[0].typ == numToken && (tokens[1].v == "+" || tokens[1].v == "-"):
		v0, err := strconv.Atoi(tokens[0].v)
		if err != nil {
			return nil, fmt.Errorf("token %v, error %v", tokens[0], err)
		}
		v1, err := eval(tokens[2:])
		if err != nil {
			return nil, err
		}
		v := ops[tokens[1].v](v0, *v1)
		return &v, nil
	case len(tokens) > 2 && tokens[0].typ == numToken && (tokens[1].v == "*" || tokens[1].v == "/"):
		v0, err := strconv.Atoi(tokens[0].v)
		if err != nil {
			return nil, fmt.Errorf("token %v, error %v", tokens[0], err)
		}
		switch {
		case tokens[2].typ == numToken:
			v1, err := strconv.Atoi(tokens[2].v)
			if err != nil {
				return nil, fmt.Errorf("token %v, error %v", tokens[2], err)
			}

			v := ops[tokens[1].v](v0, v1)
			if len(tokens) == 3 {
				return &v, nil
			}

			return eval(append([]token{{v: strconv.Itoa(v), typ: numToken}}, tokens[3:]...))
		case tokens[2].typ == openTok:
			i, err := findClosingRecursive(tokens[3:], 0, 0)
			if err != nil {
				return nil, err
			}
			i += 3

			v1, err := eval(tokens[3:i])
			if err != nil {
				return nil, err
			}

			v := ops[tokens[1].v](v0, *v1)
			if len(tokens) == i+1 {
				return &v, nil
			}

			return eval(append([]token{{v: strconv.Itoa(v), typ: numToken}}, tokens[i+1:]...))
		default:
			return nil, fmt.Errorf("unhandled inner case inside of num (*,/) [tokens...] case")
		}
	default:
		return nil, fmt.Errorf("unhandled case")
	}
}

func findClosingRecursive(tokens []token, v, stack int) (int, error) {
	switch {
	case len(tokens) == 0:
		return -1, fmt.Errorf("no closing")
	case tokens[0].typ == numToken || tokens[0].typ == arithTok:
		return findClosingRecursive(tokens[1:], v+1, stack)
	case tokens[0].typ == closeTok && stack == 0:
		return v, nil
	case tokens[0].typ == closeTok:
		return findClosingRecursive(tokens[1:], v+1, stack-1)
	case tokens[0].typ == openTok:
		return findClosingRecursive(tokens[1:], v+1, stack+1)
	default:
		return 0, fmt.Errorf("tokens %v, unhandled case", tokens)
	}
}

var ops = map[string]func(int, int) int{
	"+": func(i1, i2 int) int { return i1 + i2 },
	"-": func(i1, i2 int) int { return i1 - i2 },
	"/": func(i1, i2 int) int { return i1 / i2 },
	"*": func(i1, i2 int) int { return i1 * i2 },
}
