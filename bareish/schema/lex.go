package schema

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

// A scanner for reading lexographic tokens from a BARE schema language
// document.
type Scanner struct {
	// TODO: track lineno/colno information and attach it to the tokens
	// returned, for better error reporting
	br       *bufio.Reader
	pushback []Token
}

// Creates a new BARE schema language scanner for the given reader.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(reader), nil}
}

// Returns the next token from the reader. If the token has a string associated
// with it (e.g. UserTypeName, Name, and Integer), the second return value is
// set to that string.
func (sc *Scanner) Next() (Token, error) {
	if len(sc.pushback) != 0 {
		tok := sc.pushback[0]
		sc.pushback = sc.pushback[1:]
		return tok, nil
	}

	var (
		err error
		r   rune
	)

	for {
		r, _, err = sc.br.ReadRune()
		if err != nil {
			break
		}

		if unicode.IsSpace(r) {
			continue
		}
		if unicode.IsLetter(r) {
			sc.br.UnreadRune()
			return sc.scanWord()
		}
		if unicode.IsDigit(r) {
			sc.br.UnreadRune()
			return sc.scanInteger()
		}

		switch r {
		case '#':
			sc.br.ReadString('\n')
			continue
		case '<':
			return Token{TLANGLE, ""}, nil
		case '>':
			return Token{TRANGLE, ""}, nil
		case '{':
			return Token{TLBRACE, ""}, nil
		case '}':
			return Token{TRBRACE, ""}, nil
		case '[':
			return Token{TLBRACKET, ""}, nil
		case ']':
			return Token{TRBRACKET, ""}, nil
		case '(':
			return Token{TLPAREN, ""}, nil
		case ')':
			return Token{TRPAREN, ""}, nil
		case '|':
			return Token{TPIPE, ""}, nil
		case '=':
			return Token{TEQUAL, ""}, nil
		case ':':
			return Token{TCOLON, ""}, nil
		}

		return Token{}, &ErrUnknownToken{r}
	}

	return Token{}, err
}

// Pushes a token back to the scanner, causing it to be returned on the next
// call to Next.
func (sc *Scanner) PushBack(tok Token) {
	sc.pushback = append(sc.pushback, tok)
}

// Returned when the lexer encounters an unexpected character
type ErrUnknownToken struct {
	token rune
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("Unknown token '%c'", e.token)
}

func (sc *Scanner) scanWord() (Token, error) {
	var buf bytes.Buffer

	for {
		r, _, err := sc.br.ReadRune()
		if err != nil  {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			buf.WriteRune(r)
		} else {
			sc.br.UnreadRune()
			break
		}
	}

	tok := buf.String()
	switch tok {
	case "type":
		return Token{TTYPE, ""}, nil
	case "enum":
		return Token{TENUM, ""}, nil
	case "uint":
		return Token{TUINT, ""}, nil
	case "u8":
		return Token{TU8, ""}, nil
	case "u16":
		return Token{TU16, ""}, nil
	case "u32":
		return Token{TU32, ""}, nil
	case "u64":
		return Token{TU64, ""}, nil
	case "int":
		return Token{TINT, ""}, nil
	case "i8":
		return Token{TI8, ""}, nil
	case "i16":
		return Token{TI16, ""}, nil
	case "i32":
		return Token{TI32, ""}, nil
	case "i64":
		return Token{TI64, ""}, nil
	case "f32":
		return Token{TF32, ""}, nil
	case "f64":
		return Token{TF64, ""}, nil
	case "bool":
		return Token{TBOOL, ""}, nil
	case "string":
		return Token{TSTRING, ""}, nil
	case "data":
		return Token{TDATA, ""}, nil
	case "void":
		return Token{TVOID, ""}, nil
	case "optional":
		return Token{TOPTIONAL, ""}, nil
	case "map":
		return Token{TMAP, ""}, nil
	}

	return Token{TNAME, tok}, nil
}

func (sc *Scanner) scanInteger() (Token, error) {
	var buf bytes.Buffer

	for {
		r, _, err := sc.br.ReadRune()
		if err != nil  {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			sc.br.UnreadRune()
			break
		}
	}

	return Token{TINTEGER, buf.String()}, nil
}

// A single lexographic token from a schema language token stream
type Token struct {
	Token TokenKind
	Value string
}

type TokenKind int

const (
	TTYPE TokenKind = iota
	TENUM

	// NAME is used for name, user-type-name, and enum-value-name.
	// Distinguishing between these requires context.
	TNAME
	TINTEGER

	TUINT
	TU8
	TU16
	TU32
	TU64
	TINT
	TI8
	TI16
	TI32
	TI64
	TF32
	TF64
	TBOOL
	TSTRING
	TDATA
	TVOID
	TMAP
	TOPTIONAL

	// <
	TLANGLE
	// >
	TRANGLE
	// {
	TLBRACE
	// }
	TRBRACE
	// [
	TLBRACKET
	// ]
	TRBRACKET
	// (
	TLPAREN
	// )
	TRPAREN
	// |
	TPIPE
	// =
	TEQUAL
	// :
	TCOLON
)

func (t Token) String() string {
	switch t.Token {
	case TTYPE:
		return "type"
	case TENUM:
		return "enum"
	case TNAME:
		return "name"
	case TINTEGER:
		return "integer"
	case TUINT:
		return "uint"
	case TU8:
		return "u8"
	case TU16:
		return "u16"
	case TU32:
		return "u32"
	case TU64:
		return "u64"
	case TINT:
		return "int"
	case TI8:
		return "i8"
	case TI16:
		return "i16"
	case TI32:
		return "i32"
	case TI64:
		return "i64"
	case TF32:
		return "f32"
	case TF64:
		return "f64"
	case TBOOL:
		return "bool"
	case TSTRING:
		return "string"
	case TDATA:
		return "data"
	case TVOID:
		return "void"
	case TMAP:
		return "map"
	case TOPTIONAL:
		return "optional"
	case TLANGLE:
		return "<"
	case TRANGLE:
		return ">"
	case TLBRACE:
		return "{"
	case TRBRACE:
		return "}"
	case TLBRACKET:
		return "["
	case TRBRACKET:
		return "]"
	case TLPAREN:
		return "("
	case TRPAREN:
		return ")"
	case TPIPE:
		return "|"
	case TEQUAL:
		return "="
	case TCOLON:
		return ":"
	default:
		panic(errors.New("Invalid token value"))
	}
}
