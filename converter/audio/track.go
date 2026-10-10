package audio

import (
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/xmp"
)

// Keys of the text fields a track's tags are normalized to, named as Vorbis
// comments name them. Each format's reader maps its own fields onto them.
const (
	keyTitle        = "TITLE"
	keySubtitle     = "SUBTITLE"
	keyArtist       = "ARTIST"
	keyAlbumArtist  = "ALBUMARTIST"
	keyAlbum        = "ALBUM"
	keyComposer     = "COMPOSER"
	keyLyricist     = "LYRICIST"
	keyConductor    = "CONDUCTOR"
	keyArranger     = "ARRANGER"
	keyRemixer      = "REMIXER"
	keyGenre        = "GENRE"
	keyDate         = "DATE"
	keyOriginalDate = "ORIGINALDATE"
	keyTrack        = "TRACKNUMBER"
	keyTrackTotal   = "TRACKTOTAL"
	keyDisc         = "DISCNUMBER"
	keyDiscTotal    = "DISCTOTAL"
	keyPublisher    = "PUBLISHER"
	keyCopyright    = "COPYRIGHT"
	keyLanguage     = "LANGUAGE"
	keyISRC         = "ISRC"
	keyComment      = "COMMENT"
	keyLyrics       = "LYRICS"
	keyGrouping     = "GROUPING"
	keyBPM          = "BPM"
)

// Bounds on what is kept of a file's tags: the lyrics, the other fields,
// and how many values a field may have.
const (
	maxLyrics = 256 << 10
	maxField  = 16 << 10
	maxValues = 64
)

// tags are the text fields of a track by key, in the order they were found.
type tags map[string][]string

// add appends the values that are not blank or repeated, cut to the bounds.
func (t tags) add(key string, values ...string) {
	for _, v := range values {
		v = clean(v)
		if key != keyLyrics {
			v = oneLine(v)
		}
		if limit := maxField; key == keyLyrics {
			v = cut(v, maxLyrics)
		} else {
			v = cut(v, limit)
		}
		if v == "" || len(t[key]) >= maxValues || slices.Contains(t[key], v) {
			continue
		}
		t[key] = append(t[key], v)
	}
}

// addNumber adds a number that may be written with its total ("3/12").
func (t tags) addNumber(key, totalKey, s string) {
	n, total, _ := strings.Cut(s, "/")
	t.add(key, n)
	t.add(totalKey, total)
}

// totalKey is the key of the total that goes with a number: the tracks of
// the album, the discs.
func totalKey(key string) string {
	if key == keyDisc {
		return keyDiscTotal
	}
	return keyTrackTotal
}

// has reports whether the key has a value.
func (t tags) has(key string) bool { return len(t[key]) > 0 }

// first returns the first value of the key, or "".
func (t tags) first(key string) string {
	if v := t[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// fill adds the values of another source for the keys that have none: the
// source of lower precedence (an ID3v1 tag after an ID3v2 one).
func (t tags) fill(src tags) {
	for k, v := range src {
		if !t.has(k) {
			t.add(k, v...)
		}
	}
}

// userKey maps the names of user-defined fields (ID3 TXXX descriptions,
// iTunes "----" names, RIFF INFO's less usual chunks) to a key, or "".
func userKey(name string) string {
	switch strings.ToUpper(strings.Join(strings.FieldsFunc(name, func(r rune) bool { return r == ' ' || r == '_' || r == '-' }), "")) {
	case "TITLE":
		return keyTitle
	case "SUBTITLE", "VERSION":
		return keySubtitle
	case "ARTIST":
		return keyArtist
	case "ALBUMARTIST", "ALBUMARTISTS", "BAND":
		return keyAlbumArtist
	case "ALBUM":
		return keyAlbum
	case "COMPOSER":
		return keyComposer
	case "LYRICIST", "WRITER", "TEXT":
		return keyLyricist
	case "CONDUCTOR":
		return keyConductor
	case "ARRANGER":
		return keyArranger
	case "REMIXER", "MIXARTIST":
		return keyRemixer
	case "GENRE":
		return keyGenre
	case "DATE", "YEAR", "RELEASEDATE":
		return keyDate
	case "ORIGINALDATE", "ORIGINALYEAR", "ORIGINALRELEASEDATE":
		return keyOriginalDate
	case "TRACKNUMBER", "TRACK":
		return keyTrack
	case "TRACKTOTAL", "TOTALTRACKS":
		return keyTrackTotal
	case "DISCNUMBER", "DISC":
		return keyDisc
	case "DISCTOTAL", "TOTALDISCS":
		return keyDiscTotal
	case "PUBLISHER", "LABEL", "ORGANIZATION":
		return keyPublisher
	case "COPYRIGHT", "LICENSE":
		return keyCopyright
	case "LANGUAGE":
		return keyLanguage
	case "ISRC":
		return keyISRC
	case "COMMENT", "DESCRIPTION":
		return keyComment
	case "LYRICS", "UNSYNCEDLYRICS", "UNSYNCHRONISEDLYRICS", "USLT", "SYNCEDLYRICS", "SYNCHRONISEDLYRICS":
		return keyLyrics
	case "GROUPING", "CONTENTGROUP":
		return keyGrouping
	case "BPM", "TEMPO":
		return keyBPM
	}
	return ""
}

// picture is a picture a file carries: usually the cover of the album.
type picture struct {
	// typ is the ID3 picture type: 3 is the front cover, 0 other.
	typ  int
	mime string
	data []byte
}

// chapter is a chapter of an audiobook or podcast.
type chapter struct {
	start float64 // seconds
	title string
}

// maxChapters bounds the chapters kept.
const maxChapters = 1000

// track is what the reader of a format found in a file.
type track struct {
	// format is the file format: "mp3", "m4a", "aac", "flac", "ogg", "wav"
	// or "aiff".
	format string
	// codec names the audio coding for people ("MP3", "AAC", "Vorbis",
	// "PCM"); "" when it is not known.
	codec string
	// duration is in seconds, 0 when it is not known.
	duration float64
	// bitrate is in bits per second, 0 when it is not known.
	bitrate int
	// sampleRate is in hertz, channels the number of channels, bits the
	// bits per sample of PCM and lossless coding; each 0 when not known.
	sampleRate, channels, bits int
	// vbr reports a stream whose bitrate varies.
	vbr bool
	// protected reports audio that is encrypted (DRM).
	protected bool
	// preskip is the samples an Opus decoder drops at the start.
	preskip int64
	// scanned is the number of frames a scan of the stream counted.
	scanned int

	tags     tags
	pictures []picture
	chapters []chapter
	// synced holds the lines of an ID3 SYLT frame (synchronized lyrics).
	synced []lyricLine
	// lyrics are the lines the card shows, picked by setLyrics, and
	// lyricsSynced reports that they have their times.
	lyrics       []lyricLine
	lyricsSynced bool
	// xmp holds XMP packets (MP4's uuid box).
	xmp      [][]byte
	warnings []string
}

func newTrack(format string) *track { return &track{format: format, tags: tags{}} }

func (t *track) warn(msg string) { t.warnings = append(t.warnings, msg) }

func (t *track) addChapter(start float64, title string) {
	if len(t.chapters) < maxChapters {
		t.chapters = append(t.chapters, chapter{start, oneLine(clean(title))})
	}
}

// cover picks the picture to show: the front cover, else the first picture
// of no particular kind, else the first picture.
func (t *track) cover() *picture {
	best := -1
	for i, p := range t.pictures {
		switch {
		case p.typ == 3:
			return &t.pictures[i]
		case best < 0 || p.typ == 0 && t.pictures[best].typ != 0:
			best = i
		}
	}
	if best < 0 {
		return nil
	}
	return &t.pictures[best]
}

// mimeTypes are the media types of the formats (meta.dc.format).
var mimeTypes = map[string]string{
	"mp3": "audio/mpeg", "m4a": "audio/mp4", "aac": "audio/aac", "flac": "audio/flac", "ogg": "audio/ogg", "wav": "audio/wav", "aiff": "audio/aiff",
}

// formatNames name the formats for people.
var formatNames = map[string]string{
	"mp3": "MP3", "m4a": "M4A", "aac": "ADTS", "flac": "FLAC", "ogg": "Ogg", "wav": "WAV", "aiff": "AIFF",
}

// dublinCore maps the tags of a track to the document's metadata
// (docs/spec.md §4.3): XMP first when the file has it, then the tags.
func (t *track) dublinCore() bdf.DublinCore {
	var dc bdf.DublinCore
	for _, x := range t.xmp {
		xmp.Merge(&dc, xmp.DublinCore(x))
	}
	var native bdf.DublinCore
	tg := t.tags
	if s := tg.first(keyTitle); s != "" {
		native.Title = bdf.DCValues{s}
	}
	native.Creator = append(native.Creator, tg[keyArtist]...)
	for _, key := range []string{keyAlbumArtist, keyComposer, keyLyricist, keyConductor, keyArranger, keyRemixer} {
		for _, v := range tg[key] {
			if !slices.Contains(native.Creator, v) && !slices.Contains(native.Contributor, v) {
				native.Contributor = append(native.Contributor, v)
			}
		}
	}
	native.Subject = append(native.Subject, tg[keyGenre]...)
	if s := tg.first(keyComment); s != "" {
		native.Description = bdf.DCValues{s}
	}
	native.Publisher = append(native.Publisher, tg[keyPublisher]...)
	if d := isoDate(tg.first(keyDate)); d != "" {
		native.Date = bdf.DCValues{d}
	}
	native.Type = bdf.DCValues{"Sound"}
	if m := mimeTypes[t.format]; m != "" {
		native.Format = bdf.DCValues{m}
	}
	native.Identifier = append(native.Identifier, tg[keyISRC]...)
	if s := tg.first(keyAlbum); s != "" {
		native.Relation = bdf.DCValues{s}
	}
	for _, l := range tg[keyLanguage] {
		if b := bcp47(l); b != "" && !slices.Contains(native.Language, b) {
			native.Language = append(native.Language, b)
		}
	}
	native.Rights = append(native.Rights, tg[keyCopyright]...)
	xmp.Merge(&dc, native)
	return dc
}

// isoDate normalizes a date as tags write it (ISO 8601 as ID3v2.4, iTunes
// and Vorbis comments recommend, "2026/09/01", a year) to the W3C date and
// time format, or "" when it is not one.
func isoDate(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 10 && s[4] == '/' && s[7] == '/' {
		s = s[:4] + "-" + s[5:7] + "-" + s[8:]
	}
	return xmp.Date(s)
}

// bcp47 turns a language as tags write it (ISO 639-2 as ID3 and MP4 use,
// a BCP 47 tag as Vorbis comments may) into a BCP 47 tag; "" for an
// undetermined language.
func bcp47(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "", "und", "xxx", "mul", "zxx", "mis":
		return ""
	}
	if len(s) == 3 {
		if two, ok := iso639[s]; ok {
			return two
		}
	}
	if i := strings.IndexAny(s, "-_"); i > 0 {
		if two, ok := iso639[s[:i]]; ok {
			return two + "-" + strings.ToUpper(s[i+1:])
		}
		return s[:i] + "-" + strings.ToUpper(s[i+1:])
	}
	return s
}

// iso639 maps ISO 639-2 codes (bibliographic and terminological) to the
// ISO 639-1 codes BCP 47 uses.
var iso639 = map[string]string{
	"afr": "af", "alb": "sq", "sqi": "sq", "amh": "am", "ara": "ar", "arm": "hy", "hye": "hy", "aze": "az", "baq": "eu", "eus": "eu",
	"bel": "be", "ben": "bn", "bos": "bs", "bul": "bg", "bur": "my", "mya": "my", "cat": "ca", "chi": "zh", "zho": "zh", "cze": "cs", "ces": "cs",
	"dan": "da", "dut": "nl", "nld": "nl", "eng": "en", "epo": "eo", "est": "et", "fin": "fi", "fre": "fr", "fra": "fr", "geo": "ka", "kat": "ka",
	"ger": "de", "deu": "de", "gle": "ga", "glg": "gl", "gre": "el", "ell": "el", "guj": "gu", "heb": "he", "hin": "hi", "hrv": "hr", "hun": "hu",
	"ice": "is", "isl": "is", "ind": "id", "ita": "it", "jpn": "ja", "kan": "kn", "kaz": "kk", "khm": "km", "kor": "ko", "kur": "ku", "lao": "lo",
	"lat": "la", "lav": "lv", "lit": "lt", "mac": "mk", "mkd": "mk", "mal": "ml", "mar": "mr", "may": "ms", "msa": "ms", "mlt": "mt", "mon": "mn",
	"nep": "ne", "nor": "no", "nob": "nb", "nno": "nn", "pan": "pa", "per": "fa", "fas": "fa", "pol": "pl", "por": "pt", "rum": "ro", "ron": "ro",
	"rus": "ru", "sin": "si", "slo": "sk", "slk": "sk", "slv": "sl", "spa": "es", "srp": "sr", "swa": "sw", "swe": "sv", "tam": "ta", "tel": "te",
	"tgl": "tl", "tha": "th", "tur": "tr", "ukr": "uk", "urd": "ur", "uzb": "uz", "vie": "vi", "wel": "cy", "cym": "cy", "xho": "xh", "yid": "yi", "zul": "zu",
}

// genreName resolves the numbered genres of ID3v1 (which ID3v2 and MP4
// refer to) to their names: "17" or "(17)" is Rock, "(17)Pop" is both,
// "RX" a remix and "CR" a cover.
func genreNames(s string) []string {
	var out []string
	add := func(v string) {
		if v = strings.TrimSpace(v); v != "" && !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	for s != "" {
		switch {
		case strings.HasPrefix(s, "(("):
			add(s[1:])
			s = ""
		case s[0] == '(':
			end := strings.IndexByte(s, ')')
			if end < 0 {
				add(s)
				s = ""
				continue
			}
			add(genreRef(s[1:end]))
			s = s[end+1:]
		default:
			if n, err := strconv.Atoi(s); err == nil {
				add(genreRef(strconv.Itoa(n)))
			} else {
				add(s)
			}
			s = ""
		}
	}
	return out
}

func genreRef(ref string) string {
	switch ref {
	case "RX":
		return "Remix"
	case "CR":
		return "Cover"
	}
	if n, err := strconv.Atoi(ref); err == nil && n >= 0 && n < len(id3Genres) {
		return id3Genres[n]
	}
	return ref
}

// id3Genres are the genres of ID3v1 by number, with Winamp's additions.
var id3Genres = []string{
	"Blues", "Classic Rock", "Country", "Dance", "Disco", "Funk", "Grunge", "Hip-Hop", "Jazz", "Metal",
	"New Age", "Oldies", "Other", "Pop", "R&B", "Rap", "Reggae", "Rock", "Techno", "Industrial",
	"Alternative", "Ska", "Death Metal", "Pranks", "Soundtrack", "Euro-Techno", "Ambient", "Trip-Hop", "Vocal", "Jazz+Funk",
	"Fusion", "Trance", "Classical", "Instrumental", "Acid", "House", "Game", "Sound Clip", "Gospel", "Noise",
	"Alternative Rock", "Bass", "Soul", "Punk", "Space", "Meditative", "Instrumental Pop", "Instrumental Rock", "Ethnic", "Gothic",
	"Darkwave", "Techno-Industrial", "Electronic", "Pop-Folk", "Eurodance", "Dream", "Southern Rock", "Comedy", "Cult", "Gangsta",
	"Top 40", "Christian Rap", "Pop/Funk", "Jungle", "Native American", "Cabaret", "New Wave", "Psychedelic", "Rave", "Showtunes",
	"Trailer", "Lo-Fi", "Tribal", "Acid Punk", "Acid Jazz", "Polka", "Retro", "Musical", "Rock & Roll", "Hard Rock",
	"Folk", "Folk-Rock", "National Folk", "Swing", "Fast Fusion", "Bebop", "Latin", "Revival", "Celtic", "Bluegrass",
	"Avantgarde", "Gothic Rock", "Progressive Rock", "Psychedelic Rock", "Symphonic Rock", "Slow Rock", "Big Band", "Chorus", "Easy Listening", "Acoustic",
	"Humour", "Speech", "Chanson", "Opera", "Chamber Music", "Sonata", "Symphony", "Booty Bass", "Primus", "Porn Groove",
	"Satire", "Slow Jam", "Club", "Tango", "Samba", "Folklore", "Ballad", "Power Ballad", "Rhythmic Soul", "Freestyle",
	"Duet", "Punk Rock", "Drum Solo", "A cappella", "Euro-House", "Dance Hall", "Goa", "Drum & Bass", "Club-House", "Hardcore",
	"Terror", "Indie", "BritPop", "", "Polsk Punk", "Beat", "Christian Gangsta Rap", "Heavy Metal", "Black Metal", "Crossover",
	"Contemporary Christian", "Christian Rock", "Merengue", "Salsa", "Thrash Metal", "Anime", "JPop", "Synthpop", "Abstract", "Art Rock",
	"Baroque", "Bhangra", "Big Beat", "Breakbeat", "Chillout", "Downtempo", "Dub", "EBM", "Eclectic", "Electro",
	"Electroclash", "Emo", "Experimental", "Garage", "Global", "IDM", "Illbient", "Industro-Goth", "Jam Band", "Krautrock",
	"Leftfield", "Lounge", "Math Rock", "New Romantic", "Nu-Breakz", "Post-Punk", "Post-Rock", "Psytrance", "Shoegaze", "Space Rock",
	"Trop Rock", "World Music", "Neoclassical", "Audiobook", "Audio Theatre", "Neue Deutsche Welle", "Podcast", "Indie Rock", "G-Funk", "Dubstep",
	"Garage Rock", "Psybient",
}

// cut shortens a string to at most n bytes at a character boundary.
func cut(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return strings.TrimSpace(s[:n])
}
