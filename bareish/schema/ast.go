package schema

type SchemaTypeKind int

type SchemaType interface {
	Name() string
}

type UserDefinedType struct {
	name string
	type_ Type
}

func (udt *UserDefinedType) Name() string {
	return udt.name
}

func (udt *UserDefinedType) Type() Type {
	return udt.type_
}

type UserDefinedEnum struct {
	name   string
	kind   TypeKind
	values []EnumValue
}

func (ude *UserDefinedEnum) Name() string {
	return ude.name
}

func (ude *UserDefinedEnum) Kind() TypeKind {
	return ude.kind
}

func (ude *UserDefinedEnum) Values() []EnumValue {
	return ude.values
}

type EnumValue struct {
	name  string
	value uint
}

func (ev *EnumValue) Name() string {
	return ev.name
}

func (ev *EnumValue) Value() uint {
	return ev.value
}

type TypeKind int

const (
	UINT TypeKind = iota
	U8
	U16
	U32
	U64
	INT
	I8
	I16
	I32
	I64
	F32
	F64
	Bool
	String
	Void
	// data
	Data
	// data<length>
	DataFixed
	// [len]type
	Array
	// []type
	Slice
	// optional<type>
	Optional
	// data<len>
	DataArray
	// data
	DataSlice
	// map[type]type
	Map
	// (type | type | ...)
	Union
	// { fields... }
	Struct
	// Named user type
	UserType
)

func (tk TypeKind) String() string {
	switch tk {
	case UINT:
		return "UINT";
	case U8:
		return "U8";
	case U16:
		return "U16";
	case U32:
		return "U32";
	case U64:
		return "U64";
	case INT:
		return "INT";
	case I8:
		return "I8";
	case I16:
		return "I16";
	case I32:
		return "I32";
	case I64:
		return "I64";
	case F32:
		return "F32";
	case F64:
		return "F64";
	case Bool:
		return "Bool";
	case String:
		return "String";
	case Void:
		return "Void";
	case Data:
		return "Data";
	case DataFixed:
		return "DataFixed";
	case Array:
		return "Array";
	case Slice:
		return "Slice";
	case Optional:
		return "Optional";
	case DataArray:
		return "DataArray";
	case DataSlice:
		return "DataSlice";
	case Map:
		return "Map";
	case Union:
		return "Union";
	case Struct:
		return "Struct";
	case UserType:
		return "UserType";
	default:
		return "?"
	}
}

type Type interface {
	Kind() TypeKind
}

type PrimitiveType struct {
	kind TypeKind
}

func (pt *PrimitiveType) Kind() TypeKind {
	return pt.kind
}

type OptionalType struct {
	subtype Type
}

func (ot *OptionalType) Kind() TypeKind {
	return Optional
}

func (ot *OptionalType) Subtype() Type {
	return ot.subtype
}

type DataType struct {
	length uint
}

func (dt *DataType) Kind() TypeKind {
	if dt.length == 0 {
		return DataSlice
	}
	return DataArray
}

func (dt *DataType) Length() uint {
	return dt.length
}

type MapType struct {
	key   Type
	value Type
}

func (mt *MapType) Kind() TypeKind {
	return Map
}

func (mt *MapType) Key() Type {
	return mt.key
}

func (mt *MapType) Value() Type {
	return mt.value
}

type ArrayType struct {
	member Type
	length uint
}

func (at *ArrayType) Kind() TypeKind {
	if at.length == 0 {
		return Slice
	}
	return Array
}

func (at *ArrayType) Member() Type {
	return at.member
}

func (at *ArrayType) Length() uint {
	return at.length
}

type UnionType struct {
	types []UnionSubtype
}

func (ut *UnionType) Kind() TypeKind {
	return Union
}

func (ut *UnionType) Types() []UnionSubtype {
	return ut.types
}

type UnionSubtype struct {
	subtype Type
	tag     uint64
}

func (ust *UnionSubtype) Type() Type {
	return ust.subtype
}

func (ust *UnionSubtype) Tag() uint64 {
	return ust.tag
}

type StructType struct {
	fields []StructField
}

func (st *StructType) Kind() TypeKind {
	return Struct
}

func (st *StructType) Fields() []StructField {
	return st.fields
}

type StructField struct {
	name string
	type_ Type
}

func (sf *StructField) Name() string {
	return sf.name
}

func (sf *StructField) Type() Type {
	return sf.type_
}

// This has not been compared with the list of user-defined types and is not
// guaranteed to actually exist; the consumer of this type must perform this
// lookup itself.
type NamedUserType struct {
	name string
}

func (nut *NamedUserType) Kind() TypeKind {
	return UserType
}

func (nut *NamedUserType) Name() string {
	return nut.name
}
