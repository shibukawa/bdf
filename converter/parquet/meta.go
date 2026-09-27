package parquet

// The structures of parquet.thrift that the reader needs, with the fields
// it reads (see https://github.com/apache/parquet-format).

// Physical types.
const (
	typeBoolean = 0
	typeInt32   = 1
	typeInt64   = 2
	typeInt96   = 3
	typeFloat   = 4
	typeDouble  = 5
	typeBinary  = 6 // BYTE_ARRAY
	typeFixed   = 7 // FIXED_LEN_BYTE_ARRAY
)

// Field repetitions.
const (
	repRequired = 0
	repOptional = 1
	repRepeated = 2
)

// Converted types (the annotations before logical types).
const (
	convNone        = -1
	convUTF8        = 0
	convMap         = 1
	convMapKeyValue = 2
	convList        = 3
	convEnum        = 4
	convDecimal     = 5
	convDate        = 6
	convTimeMillis  = 7
	convTimeMicros  = 8
	convTSMillis    = 9
	convTSMicros    = 10
	convUint8       = 11
	convUint16      = 12
	convUint32      = 13
	convUint64      = 14
	convInt8        = 15
	convInt16       = 16
	convInt32       = 17
	convInt64       = 18
	convJSON        = 19
	convBSON        = 20
	convInterval    = 21
)

// Logical types: the field ids of the LogicalType union.
const (
	logNone      = 0
	logString    = 1
	logMap       = 2
	logList      = 3
	logEnum      = 4
	logDecimal   = 5
	logDate      = 6
	logTime      = 7
	logTimestamp = 8
	logInteger   = 10
	logUnknown   = 11 // always null
	logJSON      = 12
	logBSON      = 13
	logUUID      = 14
	logFloat16   = 15
	logVariant   = 16
	logGeometry  = 17
	logGeography = 18
	logFile      = 19
)

// Time units.
const (
	unitMillis = 1
	unitMicros = 2
	unitNanos  = 3
)

// Compression codecs.
const (
	codecNone   = 0
	codecSnappy = 1
	codecGzip   = 2
	codecLZO    = 3
	codecBrotli = 4
	codecLZ4    = 5 // Hadoop's framing, or a raw block
	codecZstd   = 6
	codecLZ4Raw = 7
)

var codecNames = map[int32]string{codecNone: "uncompressed", codecSnappy: "Snappy", codecGzip: "gzip", codecLZO: "LZO",
	codecBrotli: "Brotli", codecLZ4: "LZ4", codecZstd: "Zstandard", codecLZ4Raw: "LZ4"}

// Encodings.
const (
	encPlain          = 0
	encPlainDict      = 2
	encRLE            = 3
	encBitPacked      = 4
	encDeltaBinary    = 5
	encDeltaLength    = 6
	encDeltaByteArray = 7
	encRLEDict        = 8
	encByteStreamSplt = 9
)

// Page types.
const (
	pageData       = 0
	pageIndex      = 1
	pageDictionary = 2
	pageDataV2     = 3
)

type fileMeta struct {
	schema    []schemaElement
	numRows   int64
	rowGroups []rowGroup
	keyValues map[string]string
	createdBy string
}

type schemaElement struct {
	typ         int32 // -1 for groups
	typeLength  int32
	repetition  int32
	name        string
	numChildren int32
	converted   int32
	scale       int32
	precision   int32
	logical     logicalType
}

type logicalType struct {
	kind      int16 // log*
	unit      int16 // of times and timestamps
	utc       bool  // isAdjustedToUTC
	scale     int32 // of decimals
	precision int32
	bitWidth  int8 // of integers
	signed    bool
}

type rowGroup struct {
	columns []columnChunk
	numRows int64
}

type columnChunk struct {
	filePath  string
	meta      *columnMeta
	encrypted bool
}

type columnMeta struct {
	typ             int32
	codec           int32
	numValues       int64
	totalCompressed int64
	dataPageOffset  int64
	dictPageOffset  int64 // 0 when there is none
}

type pageHeader struct {
	typ          int32
	uncompressed int32
	compressed   int32
	// data pages
	numValues   int32
	encoding    int32
	defEncoding int32
	repEncoding int32
	// data pages v2
	defLength    int32
	repLength    int32
	uncompressV2 bool // is_compressed = false
	// dictionary pages
	dictValues int32
}

func readFileMeta(b []byte) (*fileMeta, error) {
	t := &thrift{b: b}
	m := &fileMeta{}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 2 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { m.schema = append(m.schema, readSchemaElement(t)) })
		case id == 3 && typ == tI64:
			m.numRows = t.i64()
		case id == 4 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { m.rowGroups = append(m.rowGroups, readRowGroup(t)) })
		case id == 5 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() {
				var k, v string
				t.fields(func(id int16, typ byte) {
					switch {
					case id == 1 && typ == tBinary:
						k = t.string()
					case id == 2 && typ == tBinary:
						v = t.string()
					default:
						t.skip(typ)
					}
				})
				if m.keyValues == nil {
					m.keyValues = map[string]string{}
				}
				m.keyValues[k] = v
			})
		case id == 6 && typ == tBinary:
			m.createdBy = t.string()
		default:
			t.skip(typ)
		}
	})
	if t.err != nil {
		return nil, t.err
	}
	return m, nil
}

func readSchemaElement(t *thrift) schemaElement {
	e := schemaElement{typ: -1, repetition: -1, converted: convNone}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tI32:
			e.typ = t.i32()
		case id == 2 && typ == tI32:
			e.typeLength = t.i32()
		case id == 3 && typ == tI32:
			e.repetition = t.i32()
		case id == 4 && typ == tBinary:
			e.name = t.string()
		case id == 5 && typ == tI32:
			e.numChildren = t.i32()
		case id == 6 && typ == tI32:
			e.converted = t.i32()
		case id == 7 && typ == tI32:
			e.scale = t.i32()
		case id == 8 && typ == tI32:
			e.precision = t.i32()
		case id == 10 && typ == tStruct:
			e.logical = readLogicalType(t)
		default:
			t.skip(typ)
		}
	})
	return e
}

func readLogicalType(t *thrift) logicalType {
	var l logicalType
	t.fields(func(id int16, typ byte) {
		if typ != tStruct {
			t.skip(typ)
			return
		}
		l.kind = id
		t.fields(func(fid int16, typ byte) {
			switch {
			case id == logDecimal && fid == 1 && typ == tI32:
				l.scale = t.i32()
			case id == logDecimal && fid == 2 && typ == tI32:
				l.precision = t.i32()
			case (id == logTime || id == logTimestamp) && fid == 1 && (typ == tTrue || typ == tFalse):
				l.utc = typ == tTrue
			case (id == logTime || id == logTimestamp) && fid == 2 && typ == tStruct:
				t.fields(func(uid int16, typ byte) {
					l.unit = uid
					t.skip(typ)
				})
			case id == logInteger && fid == 1 && typ == tByte:
				l.bitWidth = int8(t.byte())
			case id == logInteger && fid == 2 && (typ == tTrue || typ == tFalse):
				l.signed = typ == tTrue
			default:
				t.skip(typ)
			}
		})
	})
	return l
}

func readRowGroup(t *thrift) rowGroup {
	var g rowGroup
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { g.columns = append(g.columns, readColumnChunk(t)) })
		case id == 3 && typ == tI64:
			g.numRows = t.i64()
		default:
			t.skip(typ)
		}
	})
	return g
}

func readColumnChunk(t *thrift) columnChunk {
	var c columnChunk
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tBinary:
			c.filePath = t.string()
		case id == 3 && typ == tStruct:
			c.meta = readColumnMeta(t)
		case id == 8 || id == 9:
			c.encrypted = true
			t.skip(typ)
		default:
			t.skip(typ)
		}
	})
	return c
}

func readColumnMeta(t *thrift) *columnMeta {
	m := &columnMeta{}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tI32:
			m.typ = t.i32()
		case id == 4 && typ == tI32:
			m.codec = t.i32()
		case id == 5 && typ == tI64:
			m.numValues = t.i64()
		case id == 7 && typ == tI64:
			m.totalCompressed = t.i64()
		case id == 9 && typ == tI64:
			m.dataPageOffset = t.i64()
		case id == 11 && typ == tI64:
			m.dictPageOffset = t.i64()
		default:
			t.skip(typ)
		}
	})
	return m
}

// readPageHeader reads a page header from the start of b and returns its
// length.
func readPageHeader(b []byte) (*pageHeader, int, error) {
	t := &thrift{b: b}
	h := &pageHeader{typ: -1}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tI32:
			h.typ = t.i32()
		case id == 2 && typ == tI32:
			h.uncompressed = t.i32()
		case id == 3 && typ == tI32:
			h.compressed = t.i32()
		case id == 5 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.numValues = t.i32()
				case id == 2 && typ == tI32:
					h.encoding = t.i32()
				case id == 3 && typ == tI32:
					h.defEncoding = t.i32()
				case id == 4 && typ == tI32:
					h.repEncoding = t.i32()
				default:
					t.skip(typ)
				}
			})
		case id == 7 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.dictValues = t.i32()
				case id == 2 && typ == tI32:
					h.encoding = t.i32()
				default:
					t.skip(typ)
				}
			})
		case id == 8 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.numValues = t.i32()
				case id == 4 && typ == tI32:
					h.encoding = t.i32()
				case id == 5 && typ == tI32:
					h.defLength = t.i32()
				case id == 6 && typ == tI32:
					h.repLength = t.i32()
				case id == 7 && (typ == tTrue || typ == tFalse):
					h.uncompressV2 = typ == tFalse
				default:
					t.skip(typ)
				}
			})
		default:
			t.skip(typ)
		}
	})
	if t.err != nil {
		return nil, 0, t.err
	}
	return h, t.pos, nil
}
