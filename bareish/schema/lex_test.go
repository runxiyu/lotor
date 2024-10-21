package schema

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanWords(t *testing.T) {
	cases := map[string]TokenKind{
		"uint": TUINT,
		"u8": TU8,
		"u16": TU16,
		"u32": TU32,
		"u64": TU64,
		"int": TINT,
		"i8": TI8,
		"i16": TI16,
		"i32": TI32,
		"i64": TI64,
		"f32": TF32,
		"f64": TF64,
		"bool": TBOOL,
		"string": TSTRING,
		"data": TDATA,
		"void": TVOID,
		"map": TMAP,
		"optional": TOPTIONAL,
	}

	for input, reference := range cases {
		scanner := NewScanner(strings.NewReader(input))
		tok, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error")
		assert.Empty(t, tok.Value, "Expected Scan to return no value")
		assert.Equal(t, reference, tok.Token,
			"Expected Scan to return reference value for %s", input)
		_, err = scanner.Next()
		assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
	}

	scanner := NewScanner(strings.NewReader("hello"))
	tok, err := scanner.Next()
	assert.NoError(t, err, "Expected Scan to return without error")
	assert.Equal(t, tok.Value, "hello", "Expected Scan to return value 'hello'")
	assert.Equal(t, TNAME, tok.Token, "Expected Scan to return TNAME")
	_, err = scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}

func TestScanInteger(t *testing.T) {
	scanner := NewScanner(strings.NewReader("12345"))
	tok, err := scanner.Next()
	assert.NoError(t, err, "Expected Scan to return without error")
	assert.Equal(t, tok.Value, "12345", "Expected Scan to return value '12345'")
	assert.Equal(t, TINTEGER, tok.Token, "Expected Scan to return TINTEGER")
	_, err = scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}

func TestScanSymbols(t *testing.T) {
	cases := map[string]TokenKind{
		"<": TLANGLE,
		">": TRANGLE,
		"{": TLBRACE,
		"}": TRBRACE,
		"[": TLBRACKET,
		"]": TRBRACKET,
		"(": TLPAREN,
		")": TRPAREN,
	}

	for input, reference := range cases {
		scanner := NewScanner(strings.NewReader(input))
		tok, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error")
		assert.Empty(t, tok.Value, "Expected Scan to return no value")
		assert.Equal(t, reference, tok.Token,
			"Expected Scan to return reference value for %s", input)
		_, err = scanner.Next()
		assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
	}
}

func TestScanSample(t *testing.T) {
	sample := `
	type PublicKey data<128>
	type Time string # ISO T8601

	enum Department u8 {
		ACCOUNTING
		ADMINISTRATION
		CUSTOMER_SERVICE
		DEVELOPMENT

		# Reserved for the CEO
		JSMITH = 99
	}

	type Customer {
		name: string
		email: string
		address: Address
		orders: []{
			orderId: i64
			quantity: i32
		}
		metadata: map[string]data
	}

	type Person (Customer | Employee)`
	reference := []Token{
		{TTYPE, ""}, {TNAME, "PublicKey"}, {TDATA, ""},
			{TLANGLE, ""}, {TINTEGER, "128"}, {TRANGLE, ""},
		{TTYPE, ""}, {TNAME, "Time"}, {TSTRING, ""},
		{TENUM, ""}, {TNAME, "Department"}, {TU8, ""}, {TLBRACE, ""},
			{TNAME, "ACCOUNTING"},
			{TNAME, "ADMINISTRATION"},
			{TNAME, "CUSTOMER_SERVICE"},
			{TNAME, "DEVELOPMENT"},
			{TNAME, "JSMITH"}, {TEQUAL, ""}, {TINTEGER, "99"},
		{TRBRACE, ""},
		{TTYPE, ""}, {TNAME, "Customer"}, {TLBRACE, ""},
		{TNAME, "name"}, {TCOLON, ""}, {TSTRING, ""},
		{TNAME, "email"}, {TCOLON, ""}, {TSTRING, ""},
		{TNAME, "address"}, {TCOLON, ""}, {TNAME, "Address"},
		{TNAME, "orders"}, {TCOLON, ""}, {TLBRACKET, ""}, {TRBRACKET, ""}, {TLBRACE, ""},
			{TNAME, "orderId"}, {TCOLON, ""}, {TI64, ""},
			{TNAME, "quantity"}, {TCOLON, ""}, {TI32, ""},
		{TRBRACE, ""},
		{TNAME, "metadata"}, {TCOLON, ""},
			{TMAP, ""}, {TLBRACKET, ""}, {TSTRING, ""}, {TRBRACKET, ""},
				{TDATA, ""},
		{TRBRACE, ""},
		{TTYPE, ""}, {TNAME, "Person"},
			{TLPAREN, ""}, {TNAME, "Customer"},
			{TPIPE, ""}, {TNAME, "Employee"},
			{TRPAREN, ""},
	}
	scanner := NewScanner(strings.NewReader(sample))
	for i, ref := range reference {
		tok, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error for reference %d", i)
		assert.Equal(t, ref.Token, tok.Token, "Expected Scan to return correct token for reference %d", i)
		assert.Equal(t, ref.Value, tok.Value, "Expected Scan to return correct value for reference %d", i)
	}

	_, err := scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}
