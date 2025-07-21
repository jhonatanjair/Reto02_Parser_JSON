package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// JSONParser representa el estado del parser
type JSONParser struct {
	input string
	pos   int
	line  int
	col   int
}

// JSONError representa un error de parsing con contexto
type JSONError struct {
	Message string
	Line    int
	Column  int
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("JSON parsing error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewJSONParser(input string) *JSONParser {
	return &JSONParser{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (p *JSONParser) peek() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(p.input[p.pos:])
	return r
}

func (p *JSONParser) advance() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	r, size := utf8.DecodeRuneInString(p.input[p.pos:])
	p.pos += size
	if r == '\n' {
		p.line++
		p.col = 1
	} else {
		p.col++
	}
	return r
}

func (p *JSONParser) skipWhitespace() {
	for {
		r := p.peek()
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			p.advance()
		} else {
			break
		}
	}
}

func ParseJSON(input string) (interface{}, error) {
	parser := NewJSONParser(input)
	result, err := parser.parse()
	if err != nil {
		return nil, err
	}
	parser.skipWhitespace()
	if parser.pos < len(parser.input) {
		return nil, parser.newError("caracter inesperado después del valor JSON")
	}
	return result, nil
}

func (p *JSONParser) parse() (interface{}, error) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return nil, p.newError("entrada vacía o solo espacios en blanco")
	}
	switch r := p.peek(); r {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't', 'f':
		return p.parseBoolean()
	case 'n':
		return p.parseNull()
	default:
		if r == '-' || unicode.IsDigit(r) {
			return p.parseNumber()
		}
		return nil, p.newError(fmt.Sprintf("carácter inesperado '%c'", r))
	}
}

func (p *JSONParser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	p.advance()
	p.skipWhitespace()
	if p.peek() == '}' {
		p.advance()
		return obj, nil
	}
	for {
		p.skipWhitespace()
		if p.peek() != '"' {
			return nil, p.newError("se esperaba una clave string")
		}
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if p.peek() != ':' {
			return nil, p.newError("se esperaba ':' después de la clave")
		}
		p.advance()
		p.skipWhitespace()
		value, err := p.parse()
		if err != nil {
			return nil, err
		}
		obj[key] = value
		p.skipWhitespace()
		if p.peek() == '}' {
			p.advance()
			break
		}
		if p.peek() != ',' {
			return nil, p.newError("se esperaba ',' entre pares clave-valor")
		}
		p.advance()
	}
	return obj, nil
}

func (p *JSONParser) parseArray() ([]interface{}, error) {
	var arr []interface{}
	p.advance()
	p.skipWhitespace()
	if p.peek() == ']' {
		p.advance()
		return arr, nil
	}
	for {
		p.skipWhitespace()
		val, err := p.parse()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
		p.skipWhitespace()
		if p.peek() == ']' {
			p.advance()
			break
		}
		if p.peek() != ',' {
			return nil, p.newError("se esperaba ',' entre elementos del arreglo")
		}
		p.advance()
	}
	return arr, nil
}

func (p *JSONParser) parseString() (string, error) {
	p.advance()
	var result []rune
	for {
		if p.pos >= len(p.input) {
			return "", p.newError("cadena no terminada")
		}
		r := p.advance()
		if r == '"' {
			break
		}
		if r == '\\' {
			esc := p.advance()
			switch esc {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case '/':
				result = append(result, '/')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			default:
				return "", p.newError("escape inválido en string")
			}
		} else {
			result = append(result, r)
		}
	}
	return string(result), nil
}

func (p *JSONParser) parseNumber() (float64, error) {
	start := p.pos
	hasDot := false
	hasExp := false

	r := p.peek()
	if r == '-' {
		p.advance()
	}
	if p.peek() == '0' {
		p.advance()
		if unicode.IsDigit(p.peek()) {
			return 0, p.newError("número inválido con cero inicial")
		}
	} else if unicode.IsDigit(p.peek()) {
		for unicode.IsDigit(p.peek()) {
			p.advance()
		}
	} else {
		return 0, p.newError("número mal formado")
	}

	if p.peek() == '.' {
		hasDot = true
		p.advance()
		if !unicode.IsDigit(p.peek()) {
			return 0, p.newError("se esperaba dígito después del punto decimal")
		}
		for unicode.IsDigit(p.peek()) {
			p.advance()
		}
	}

	if p.peek() == 'e' || p.peek() == 'E' {
		hasExp = true
		p.advance()
		if p.peek() == '+' || p.peek() == '-' {
			p.advance()
		}
		if !unicode.IsDigit(p.peek()) {
			return 0, p.newError("se esperaba dígito después del exponente")
		}
		for unicode.IsDigit(p.peek()) {
			p.advance()
		}
	}

	numStr := p.input[start:p.pos]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil || (!hasDot && !hasExp && numStr != fmt.Sprintf("%.0f", num)) {
		return 0, p.newError("número inválido")
	}
	return num, nil
}

func (p *JSONParser) parseBoolean() (bool, error) {
	if p.pos+4 <= len(p.input) && p.input[p.pos:p.pos+4] == "true" {
		p.pos += 4
		p.col += 4
		return true, nil
	} else if p.pos+5 <= len(p.input) && p.input[p.pos:p.pos+5] == "false" {
		p.pos += 5
		p.col += 5
		return false, nil
	}
	return false, p.newError("booleano inválido")
}

func (p *JSONParser) parseNull() (interface{}, error) {
	if p.pos+4 <= len(p.input) && p.input[p.pos:p.pos+4] == "null" {
		p.pos += 4
		p.col += 4
		return nil, nil
	}
	return nil, p.newError("se esperaba 'null'")
}

func (p *JSONParser) newError(msg string) error {
	return &JSONError{
		Message: msg,
		Line:    p.line,
		Column:  p.col,
	}
}
