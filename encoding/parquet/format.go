package parquet

// The constants of parquet.thrift that the writer and the reader use.

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
	convNone     = -1
	convUTF8     = 0
	convEnum     = 4
	convDecimal  = 5
	convDate     = 6
	convTSMillis = 9
	convTSMicros = 10
	convJSON     = 19
)

// Logical types: the field ids of the LogicalType union.
const (
	logNone      = 0
	logString    = 1
	logEnum      = 4
	logDecimal   = 5
	logDate      = 6
	logTimestamp = 8
	logJSON      = 12
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
	codecLZ4    = 5
	codecZstd   = 6
	codecLZ4Raw = 7
)

// Encodings.
const (
	encPlain     = 0
	encPlainDict = 2
	encRLE       = 3
	encBitPacked = 4
	encRLEDict   = 8
)

// Page types.
const (
	pageData       = 0
	pageIndex      = 1
	pageDictionary = 2
	pageDataV2     = 3
)
