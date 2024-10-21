package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnparseValue(t *testing.T) {
	var (
		err    error
		schema string

		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
		u   uint
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		i   int
		f32 float32
		f64 float64
		b   bool
		str string
	)

	schema, err = SchemaFor(&u8)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "u8", "Expected SchemaFor to return u8")

	schema, err = SchemaFor(&u16)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "u16", "Expected SchemaFor to return u16")

	schema, err = SchemaFor(&u32)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "u32", "Expected SchemaFor to return u32")

	schema, err = SchemaFor(&u64)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "u64", "Expected SchemaFor to return u64")

	schema, err = SchemaFor(&u)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "uint", "Expected SchemaFor to return u32")

	schema, err = SchemaFor(&i8)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "i8", "Expected SchemaFor to return i8")

	schema, err = SchemaFor(&i16)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "i16", "Expected SchemaFor to return i16")

	schema, err = SchemaFor(&i32)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "i32", "Expected SchemaFor to return i32")

	schema, err = SchemaFor(&i64)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "i64", "Expected SchemaFor to return i64")

	schema, err = SchemaFor(&i)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "int", "Expected SchemaFor to return i32")

	schema, err = SchemaFor(&f32)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "f32", "Expected SchemaFor to return f32")

	schema, err = SchemaFor(&f64)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "f64", "Expected SchemaFor to return f64")

	schema, err = SchemaFor(&b)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "bool", "Expected SchemaFor to return bool")

	schema, err = SchemaFor(&str)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "string", "Expected SchemaFor to return string")
}

func TestUnparseOptional(t *testing.T) {
	var val *string
	schema, err := SchemaFor(&val)
	assert.Nil(t, err, "Expected SchemaFor to return without error")
	assert.Equal(t, schema, "optional<string>",
		"Expected SchemaFor to return optional<string>")
}
