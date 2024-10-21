package schema

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var (
	userTypeNameRE = regexp.MustCompile(`[A-Z][A-Za-z0-9]*`)
	userEnumNameRE = regexp.MustCompile(`[A-Z][A-Za-z0-9]*`)
	fieldNameRE = regexp.MustCompile(`[a-z][A-Za-z0-9]*`)
	enumValueRE = regexp.MustCompile(`[A-Z][A-Z0-9_]*`)
)

// Returned when the lexer encounters an unexpected token
type ErrUnexpectedToken struct {
	token    Token
	expected string
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("Unexpected token '%s'; expected %s",
		e.token.String(), e.expected)
}

// Parses a BARE schema definition language document from the given reader and
// returns a list of the user-defined types it specifies.
func Parse(reader io.Reader) ([]SchemaType, error) {
	scanner := NewScanner(reader)
	var stypes []SchemaType
	for {
		st, err := parseSchemaType(scanner)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		stypes = append(stypes, st)
	}
	return stypes, nil
}

func parseSchemaType(scanner *Scanner) (SchemaType, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Token {
	case TTYPE:
		scanner.PushBack(tok)
		return parseUserType(scanner)
	case TENUM:
		scanner.PushBack(tok)
		return parseUserEnum(scanner)
	}

	return nil, &ErrUnexpectedToken{tok, "'type' or 'enum'"}
}

func parseUserType(scanner *Scanner) (SchemaType, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TTYPE {
		return nil, &ErrUnexpectedToken{tok, "type"}
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TNAME {
		return nil, &ErrUnexpectedToken{tok, "type name"}
	}

	udt := &UserDefinedType{name: tok.Value}
	udt.type_, err = parseType(scanner)
	if err != nil {
		return nil, err
	}

	if !userTypeNameRE.MatchString(udt.Name()) {
		return nil, fmt.Errorf("Invalid name for user type %s", udt.Name())
	}

	return udt, nil
}

func parseUserEnum(scanner *Scanner) (SchemaType, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TENUM {
		return nil, &ErrUnexpectedToken{tok, "enum"}
	}

	var name string
	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TNAME {
		return nil, &ErrUnexpectedToken{tok, "enum name"}
	}
	name = tok.Value

	var kind TypeKind
	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	switch tok.Token {
	case TU8:
		kind = U8
	case TU16:
		kind = U16
	case TU32:
		kind = U32
	case TU64:
		kind = U64
	default:
		kind = UINT
		scanner.PushBack(tok)
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLBRACE {
		return nil, &ErrUnexpectedToken{tok, "{"}
	}

	var value uint
	var evs []EnumValue
	for {
		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token != TNAME {
			return nil, &ErrUnexpectedToken{tok, "value name"}
		}

		var ev EnumValue
		ev.name = tok.Value
		if !enumValueRE.MatchString(ev.name) {
			return nil, fmt.Errorf("Invalid name for enum value %s", ev.name)
		}

		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token == TEQUAL {
			tok, err = scanner.Next()
			if err != nil {
				return nil, err
			}
			if tok.Token != TINTEGER {
				return nil, &ErrUnexpectedToken{tok, "integer"}
			}

			v, _ := strconv.ParseUint(tok.Value, 10, 32)
			value = uint(v)
			ev.value = value
		} else {
			ev.value = value
			value += 1
			scanner.PushBack(tok)
		}

		evs = append(evs, ev)

		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}

		if tok.Token == TRBRACE {
			break
		} else if tok.Token == TNAME {
			scanner.PushBack(tok)
		} else {
			return nil, &ErrUnexpectedToken{tok, "value name"}
		}
	}

	if !userEnumNameRE.MatchString(name) {
		return nil, fmt.Errorf("Invalid name for user enum %s", name)
	}

	return &UserDefinedEnum{name, kind, evs}, nil
}

func parseType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Token {
	case TUINT:
		return &PrimitiveType{UINT}, nil
	case TU8:
		return &PrimitiveType{U8}, nil
	case TU16:
		return &PrimitiveType{U16}, nil
	case TU32:
		return &PrimitiveType{U32}, nil
	case TU64:
		return &PrimitiveType{U64}, nil
	case TINT:
		return &PrimitiveType{INT}, nil
	case TI8:
		return &PrimitiveType{I8}, nil
	case TI16:
		return &PrimitiveType{I16}, nil
	case TI32:
		return &PrimitiveType{I32}, nil
	case TI64:
		return &PrimitiveType{I64}, nil
	case TF32:
		return &PrimitiveType{F32}, nil
	case TF64:
		return &PrimitiveType{F64}, nil
	case TBOOL:
		return &PrimitiveType{Bool}, nil
	case TSTRING:
		return &PrimitiveType{String}, nil
	case TVOID:
		return &PrimitiveType{Void}, nil
	case TOPTIONAL:
		scanner.PushBack(tok)
		return parseOptionalType(scanner)
	case TDATA:
		scanner.PushBack(tok)
		return parseDataType(scanner)
	case TMAP:
		scanner.PushBack(tok)
		return parseMapType(scanner)
	case TLBRACKET:
		scanner.PushBack(tok)
		return parseArrayType(scanner)
	case TLPAREN:
		scanner.PushBack(tok)
		return parseUnionType(scanner)
	case TLBRACE:
		scanner.PushBack(tok)
		return parseStructType(scanner)
	case TNAME:
		return &NamedUserType{name: tok.Value}, nil
	}

	return nil, &ErrUnexpectedToken{tok, "type"}
}

func parseOptionalType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TOPTIONAL {
		return nil, &ErrUnexpectedToken{tok, "optional"}
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLANGLE {
		return nil, &ErrUnexpectedToken{tok, "<"}
	}

	st, err := parseType(scanner)
	if err != nil {
		return nil, err
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TRANGLE {
		return nil, &ErrUnexpectedToken{tok, ">"}
	}
	return &OptionalType{subtype: st}, nil
}

func parseDataType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TDATA {
		return nil, &ErrUnexpectedToken{tok, "data"}
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLANGLE {
		scanner.PushBack(tok)
		return &DataType{0}, nil
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TINTEGER {
		return nil, &ErrUnexpectedToken{tok, "integer"}
	}
	length, _ := strconv.ParseUint(tok.Value, 10, 32)

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TRANGLE {
		return nil, &ErrUnexpectedToken{tok, ">"}
	}

	return &DataType{uint(length)}, nil
}

func parseMapType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TMAP {
		return nil, &ErrUnexpectedToken{tok, "map"}
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLBRACKET {
		return nil, &ErrUnexpectedToken{tok, "["}
	}

	key, err := parseType(scanner)
	if err != nil {
		return nil, err
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TRBRACKET {
		return nil, &ErrUnexpectedToken{tok, "]"}
	}

	value, err := parseType(scanner)
	if err != nil {
		return nil, err
	}

	return &MapType{key, value}, nil
}

func parseArrayType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLBRACKET {
		return nil, &ErrUnexpectedToken{tok, "["}
	}

	tok, err = scanner.Next()
	if err != nil {
		return nil, err
	}

	var length uint
	switch tok.Token {
	case TINTEGER:
		l, _ := strconv.ParseUint(tok.Value, 10, 32)
		length = uint(l)

		tok, err := scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token != TRBRACKET {
			return nil, &ErrUnexpectedToken{tok, "]"}
		}
		break
	case TRBRACKET:
		break
	default:
		return nil, &ErrUnexpectedToken{tok, "]"}
	}

	member, err := parseType(scanner)
	if err != nil {
		return nil, err
	}

	return &ArrayType{member, length}, nil
}

func parseUnionType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLPAREN {
		return nil, &ErrUnexpectedToken{tok, "("}
	}

	var (
		types []UnionSubtype
		tag   uint64
	)
	for {
		ty, err := parseType(scanner)
		if err != nil {
			return nil, err
		}

		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token == TEQUAL {
			tok, err := scanner.Next()
			if err != nil {
				return nil, err
			}
			if tok.Token != TINTEGER {
				return nil, &ErrUnexpectedToken{tok, "integer"}
			}
			tag, _ = strconv.ParseUint(tok.Value, 10, 64)
		} else {
			scanner.PushBack(tok)
		}

		types = append(types, UnionSubtype{
			subtype: ty,
			tag:     tag,
		})
		tag++

		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}

		if tok.Token == TPIPE {
			continue
		} else if tok.Token == TRPAREN {
			break
		} else {
			return nil, &ErrUnexpectedToken{tok, "'|' or ')'"}
		}
	}

	return &UnionType{types}, nil
}

func parseStructType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TLBRACE {
		return nil, &ErrUnexpectedToken{tok, "["}
	}

	var fields []StructField
	for {
		var sf StructField

		tok, err := scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token == TRBRACE {
			break
		}
		if tok.Token != TNAME {
			return nil, &ErrUnexpectedToken{tok, "field name"}
		}

		sf.name = tok.Value
		if !fieldNameRE.MatchString(sf.name) {
			return nil, fmt.Errorf("Invalid name for field %s", sf.name)
		}

		tok, err = scanner.Next()
		if err != nil {
			return nil, err
		}
		if tok.Token != TCOLON {
			return nil, &ErrUnexpectedToken{tok, ":"}
		}

		sf.type_, err = parseType(scanner)
		if err != nil {
			return nil, err
		}

		fields = append(fields, sf)
	}

	return &StructType{fields}, nil
}
