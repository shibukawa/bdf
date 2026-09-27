package cgm

import "strconv"

// code identifies an element by its class and id (ISO/IEC 8632-1 §8).
type code struct{ class, id int }

// Delimiter elements (class 0).
var (
	eNoOp              = code{0, 0}
	eBeginMetafile     = code{0, 1}
	eEndMetafile       = code{0, 2}
	eBeginPicture      = code{0, 3}
	eBeginPictureBody  = code{0, 4}
	eEndPicture        = code{0, 5}
	eBeginSegment      = code{0, 6}
	eEndSegment        = code{0, 7}
	eBeginFigure       = code{0, 8}
	eEndFigure         = code{0, 9}
	eBeginProtection   = code{0, 13}
	eEndProtection     = code{0, 14}
	eBeginCompoundLine = code{0, 15}
	eEndCompoundLine   = code{0, 16}
	eBeginCompoundText = code{0, 17}
	eEndCompoundText   = code{0, 18}
	eBeginTileArray    = code{0, 19}
	eEndTileArray      = code{0, 20}
	eBeginAPS          = code{0, 21}
	eBeginAPSBody      = code{0, 22}
	eEndAPS            = code{0, 23}
)

// Metafile descriptor elements (class 1).
var (
	eMetafileVersion     = code{1, 1}
	eMetafileDescription = code{1, 2}
	eVDCType             = code{1, 3}
	eIntegerPrecision    = code{1, 4}
	eRealPrecision       = code{1, 5}
	eIndexPrecision      = code{1, 6}
	eColourPrecision     = code{1, 7}
	eColourIndexPrec     = code{1, 8}
	eMaxColourIndex      = code{1, 9}
	eColourValueExtent   = code{1, 10}
	eElementList         = code{1, 11}
	eDefaultsReplacement = code{1, 12}
	eFontList            = code{1, 13}
	eCharacterSetList    = code{1, 14}
	eCharacterCoding     = code{1, 15}
	eNamePrecision       = code{1, 16}
	eColourModel         = code{1, 19}
)

// Picture descriptor elements (class 2).
var (
	eScalingMode         = code{2, 1}
	eColourSelectionMode = code{2, 2}
	eLineWidthMode       = code{2, 3}
	eMarkerSizeMode      = code{2, 4}
	eEdgeWidthMode       = code{2, 5}
	eVDCExtent           = code{2, 6}
	eBackgroundColour    = code{2, 7}
	eDeviceViewport      = code{2, 8}
	eDeviceViewportMode  = code{2, 9}
	eLineRep             = code{2, 11}
	eMarkerRep           = code{2, 12}
	eTextRep             = code{2, 13}
	eFillRep             = code{2, 14}
	eEdgeRep             = code{2, 15}
	eInteriorStyleMode   = code{2, 16}
	eLineEdgeTypeDef     = code{2, 17}
	eHatchStyleDef       = code{2, 18}
	eGeoPatternDef       = code{2, 19}
)

// Control elements (class 3).
var (
	eVDCIntegerPrecision = code{3, 1}
	eVDCRealPrecision    = code{3, 2}
	eAuxiliaryColour     = code{3, 3}
	eTransparency        = code{3, 4}
	eClipRectangle       = code{3, 5}
	eClipIndicator       = code{3, 6}
	eNewRegion           = code{3, 10}
	eSaveContext         = code{3, 11}
	eRestoreContext      = code{3, 12}
	eMitreLimit          = code{3, 19}
	eTransparentCell     = code{3, 20}
)

// Graphical primitives (class 4).
var (
	ePolyline          = code{4, 1}
	eDisjointPolyline  = code{4, 2}
	ePolymarker        = code{4, 3}
	eText              = code{4, 4}
	eRestrictedText    = code{4, 5}
	eAppendText        = code{4, 6}
	ePolygon           = code{4, 7}
	ePolygonSet        = code{4, 8}
	eCellArray         = code{4, 9}
	eGDP               = code{4, 10}
	eRectangle         = code{4, 11}
	eCircle            = code{4, 12}
	eArc3Point         = code{4, 13}
	eArc3PointClose    = code{4, 14}
	eArcCentre         = code{4, 15}
	eArcCentreClose    = code{4, 16}
	eEllipse           = code{4, 17}
	eEllipticalArc     = code{4, 18}
	eEllipticalArcClos = code{4, 19}
	eArcCentreReversed = code{4, 20}
	eConnectingEdge    = code{4, 21}
	eHyperbolicArc     = code{4, 22}
	eParabolicArc      = code{4, 23}
	eNUBSpline         = code{4, 24}
	eNURBSpline        = code{4, 25}
	ePolybezier        = code{4, 26}
	ePolysymbol        = code{4, 27}
	eBitonalTile       = code{4, 28}
	eTile              = code{4, 29}
)

// Attribute elements (class 5).
var (
	eLineBundleIndex    = code{5, 1}
	eLineType           = code{5, 2}
	eLineWidth          = code{5, 3}
	eLineColour         = code{5, 4}
	eMarkerBundleIndex  = code{5, 5}
	eMarkerType         = code{5, 6}
	eMarkerSize         = code{5, 7}
	eMarkerColour       = code{5, 8}
	eTextBundleIndex    = code{5, 9}
	eTextFontIndex      = code{5, 10}
	eTextPrecision      = code{5, 11}
	eCharExpansion      = code{5, 12}
	eCharSpacing        = code{5, 13}
	eTextColour         = code{5, 14}
	eCharHeight         = code{5, 15}
	eCharOrientation    = code{5, 16}
	eTextPath           = code{5, 17}
	eTextAlignment      = code{5, 18}
	eCharSetIndex       = code{5, 19}
	eAltCharSetIndex    = code{5, 20}
	eFillBundleIndex    = code{5, 21}
	eInteriorStyle      = code{5, 22}
	eFillColour         = code{5, 23}
	eHatchIndex         = code{5, 24}
	ePatternIndex       = code{5, 25}
	eEdgeBundleIndex    = code{5, 26}
	eEdgeType           = code{5, 27}
	eEdgeWidth          = code{5, 28}
	eEdgeColour         = code{5, 29}
	eEdgeVisibility     = code{5, 30}
	eFillReferencePoint = code{5, 31}
	ePatternTable       = code{5, 32}
	ePatternSize        = code{5, 33}
	eColourTable        = code{5, 34}
	eASF                = code{5, 35}
	eLineCap            = code{5, 37}
	eLineJoin           = code{5, 38}
	eLineTypeOffset     = code{5, 40}
	eRestrictedTextType = code{5, 42}
	eInterpolatedInt    = code{5, 43}
	eEdgeCap            = code{5, 44}
	eEdgeJoin           = code{5, 45}
	eEdgeTypeOffset     = code{5, 47}
)

// Segment elements (class 8) and application structures (class 9).
var (
	eCopySegment    = code{8, 1}
	eSegmentTransf  = code{8, 4}
	eAPSAttribute   = code{9, 1}
	eEndMFDefaults  = code{-1, 1} // clear text: the end of BEGMFDEFAULTS
	eUnknownKeyword = code{-1, 2} // clear text: a keyword not in the table
)

// clearText maps the element names of the clear text encoding (ISO/IEC
// 8632-4) to their elements. The INCR forms of the point lists give each
// point after the first relative to the one before.
var clearText = map[string]code{
	"BEGMF": eBeginMetafile, "ENDMF": eEndMetafile, "BEGPIC": eBeginPicture, "BEGPICBODY": eBeginPictureBody,
	"ENDPIC": eEndPicture, "BEGSEG": eBeginSegment, "ENDSEG": eEndSegment, "BEGFIGURE": eBeginFigure,
	"ENDFIGURE": eEndFigure, "BEGPROTREGION": eBeginProtection, "ENDPROTREGION": eEndProtection,
	"BEGCOMPOLINE": eBeginCompoundLine, "ENDCOMPOLINE": eEndCompoundLine, "BEGCOMPOTEXTPATH": eBeginCompoundText,
	"ENDCOMPOTEXTPATH": eEndCompoundText, "BEGTILEARRAY": eBeginTileArray, "ENDTILEARRAY": eEndTileArray,
	"BEGAPS": eBeginAPS, "BEGAPSBODY": eBeginAPSBody, "ENDAPS": eEndAPS,

	"MFVERSION": eMetafileVersion, "MFDESC": eMetafileDescription, "VDCTYPE": eVDCType,
	"INTEGERPREC": eIntegerPrecision, "REALPREC": eRealPrecision, "INDEXPREC": eIndexPrecision,
	"COLRPREC": eColourPrecision, "COLRINDEXPREC": eColourIndexPrec, "MAXCOLRINDEX": eMaxColourIndex,
	"COLRVALUEEXT": eColourValueExtent, "MFELEMLIST": eElementList, "BEGMFDEFAULTS": eDefaultsReplacement,
	"ENDMFDEFAULTS": eEndMFDefaults, "FONTLIST": eFontList, "CHARSETLIST": eCharacterSetList,
	"CHARCODING": eCharacterCoding, "NAMEPREC": eNamePrecision, "MAXVDCEXT": {1, 17}, "SEGPRIEXT": {1, 18},
	"COLRMODEL": eColourModel, "COLRCALIB": {1, 20}, "FONTPROP": {1, 21}, "GLYPHMAP": {1, 22},
	"SYMBOLLIBLIST": {1, 23}, "PICDIR": {1, 24},

	"SCALEMODE": eScalingMode, "COLRMODE": eColourSelectionMode, "LINEWIDTHMODE": eLineWidthMode,
	"MARKERSIZEMODE": eMarkerSizeMode, "EDGEWIDTHMODE": eEdgeWidthMode, "VDCEXT": eVDCExtent,
	"BACKCOLR": eBackgroundColour, "DEVVP": eDeviceViewport, "DEVVPMODE": eDeviceViewportMode, "DEVVPMAP": {2, 10},
	"LINEREP": eLineRep, "MARKERREP": eMarkerRep, "TEXTREP": eTextRep, "FILLREP": eFillRep, "EDGEREP": eEdgeRep,
	"INTSTYLEMODE": eInteriorStyleMode, "LINEEDGETYPEDEF": eLineEdgeTypeDef, "HATCHSTYLEDEF": eHatchStyleDef,
	"GEOPATDEF": eGeoPatternDef, "APSDIR": {2, 20},

	"VDCINTEGERPREC": eVDCIntegerPrecision, "VDCREALPREC": eVDCRealPrecision, "AUXCOLR": eAuxiliaryColour,
	"TRANSPARENCY": eTransparency, "CLIPRECT": eClipRectangle, "CLIP": eClipIndicator, "LINECLIPMODE": {3, 7},
	"MARKERCLIPMODE": {3, 8}, "EDGECLIPMODE": {3, 9}, "NEWREGION": eNewRegion, "SAVEPRIMCONT": eSaveContext,
	"RESPRIMCONT": eRestoreContext, "PROTREGION": {3, 17}, "GENTEXTPATHMODE": {3, 18}, "MITRELIMIT": eMitreLimit,
	"TRANSPCELLCOLR": eTransparentCell,

	"LINE": ePolyline, "DISJTLINE": eDisjointPolyline, "MARKER": ePolymarker, "TEXT": eText,
	"RESTRTEXT": eRestrictedText, "APNDTEXT": eAppendText, "POLYGON": ePolygon, "POLYGONSET": ePolygonSet,
	"CELLARRAY": eCellArray, "GDP": eGDP, "RECT": eRectangle, "CIRCLE": eCircle, "ARC3PT": eArc3Point,
	"ARC3PTCLOSE": eArc3PointClose, "ARCCTR": eArcCentre, "ARCCTRCLOSE": eArcCentreClose, "ELLIPSE": eEllipse,
	"ELLIPARC": eEllipticalArc, "ELLIPARCCLOSE": eEllipticalArcClos, "ARCCTRREV": eArcCentreReversed,
	"CONNEDGE": eConnectingEdge, "HYPERBARC": eHyperbolicArc, "PARABARC": eParabolicArc, "NUB": eNUBSpline,
	"NURB": eNURBSpline, "POLYBEZIER": ePolybezier, "SYMBOL": ePolysymbol, "BITONALTILE": eBitonalTile,
	"TILE": eTile,

	"LINEINDEX": eLineBundleIndex, "LINETYPE": eLineType, "LINEWIDTH": eLineWidth, "LINECOLR": eLineColour,
	"MARKERINDEX": eMarkerBundleIndex, "MARKERTYPE": eMarkerType, "MARKERSIZE": eMarkerSize,
	"MARKERCOLR": eMarkerColour, "TEXTINDEX": eTextBundleIndex, "TEXTFONTINDEX": eTextFontIndex,
	"TEXTPREC": eTextPrecision, "CHAREXPAN": eCharExpansion, "CHARSPACE": eCharSpacing, "TEXTCOLR": eTextColour,
	"CHARHEIGHT": eCharHeight, "CHARORI": eCharOrientation, "TEXTPATH": eTextPath, "TEXTALIGN": eTextAlignment,
	"CHARSETINDEX": eCharSetIndex, "ALTCHARSETINDEX": eAltCharSetIndex, "FILLINDEX": eFillBundleIndex,
	"INTSTYLE": eInteriorStyle, "FILLCOLR": eFillColour, "HATCHINDEX": eHatchIndex, "PATINDEX": ePatternIndex,
	"EDGEINDEX": eEdgeBundleIndex, "EDGETYPE": eEdgeType, "EDGEWIDTH": eEdgeWidth, "EDGECOLR": eEdgeColour,
	"EDGEVIS": eEdgeVisibility, "FILLREFPT": eFillReferencePoint, "PATTABLE": ePatternTable,
	"PATSIZE": ePatternSize, "COLRTABLE": eColourTable, "ASF": eASF, "PICKID": {5, 36}, "LINECAP": eLineCap,
	"LINEJOIN": eLineJoin, "LINETYPECONT": {5, 39}, "LINETYPEINITOFFSET": eLineTypeOffset,
	"TEXTSCORETYPE": {5, 41}, "RESTRTEXTTYPE": eRestrictedTextType, "INTERPINT": eInterpolatedInt,
	"EDGECAP": eEdgeCap, "EDGEJOIN": eEdgeJoin, "EDGETYPECONT": {5, 46}, "EDGETYPEINITOFFSET": eEdgeTypeOffset,
	"SYMBOLLIBINDEX": {5, 48}, "SYMBOLCOLR": {5, 49}, "SYMBOLSIZE": {5, 50}, "SYMBOLORI": {5, 51},

	"ESCAPE": {6, 1}, "MESSAGE": {7, 1}, "APPLDATA": {7, 2},
	"COPYSEG": eCopySegment, "INHFILTER": {8, 2}, "CLIPINH": {8, 3}, "SEGTRAN": eSegmentTransf,
	"SEGHIGHL": {8, 5}, "SEGDISPPRI": {8, 6}, "SEGPICKPRI": {8, 7},
	"APSATTR": eAPSAttribute,
}

// incremental names the clear text elements whose points after the first
// are relative.
var incremental = map[string]code{
	"INCRLINE": ePolyline, "INCRDISJTLINE": eDisjointPolyline, "INCRMARKER": ePolymarker,
	"INCRPOLYGON": ePolygon, "INCRPOLYGONSET": ePolygonSet,
}

// names are the element names of warnings.
var names = map[code]string{}

func init() {
	for k, c := range clearText {
		if c.class >= 0 {
			names[c] = k
		}
	}
}

// name returns the clear text name of an element, or its class and id.
func (c code) name() string {
	if n, ok := names[c]; ok {
		return n
	}
	return "class " + strconv.Itoa(c.class) + " element " + strconv.Itoa(c.id)
}
