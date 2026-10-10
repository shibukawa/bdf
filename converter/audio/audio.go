// Package audio converts audio files — MP3, AAC (in an MP4 container, as
// .m4a, or raw as ADTS), FLAC, Ogg Vorbis and Opus, WAV and AIFF — into BDF
// documents of one page: a card that shows the cover art the file carries
// (or a placeholder when it has none) across the top, and under it what
// the tags say: the title, the artists, the album, the other fields in a
// table, the lyrics and the chapters.
//
// The view plays (docs/spec.md §4.4): the file is stored as it is, for the
// browser's own player, and cues tie its time to the page, to the lines of
// synchronized lyrics (an ID3 SYLT frame, or LRC text in a lyrics field)
// or else to the chapters, so that a viewer shows the line being sung and
// plays from a line clicked on. Options.NoPlay leaves the audio out, and a
// file larger than Options.MaxAudio gets its card without it.
//
// The cover passes through as the images of converter/image do: stored as
// it is, in the format the file holds it in (JPEG or PNG, usually), for
// the browser to decode. The tags become the document's Dublin Core the
// way other formats' document properties do (docs/spec.md §4.3), from ID3v2
// and ID3v1 (MP3, and the ID3 chunks of WAV and AIFF), iTunes-style MP4
// metadata (with the XMP an MP4 file may carry), Vorbis comments (FLAC and
// Ogg), RIFF INFO and the text chunks of AIFF. The duration, bitrate,
// sample rate and channels come from the stream itself. See
// docs/design.md §3.31.
package audio

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/imgconv"

	// the formats of cover art a browser shows, for their sizes
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title (by default the track's title,
	// else the file name).
	Title string
	// FileName is the file's name: the title of a file without one, and
	// the alternative text of a cover without a title.
	FileName string
	// Images controls whether the cover is re-encoded (see imgconv); the
	// zero value keeps it as it is.
	Images imgconv.Options
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// NoPlay leaves the audio out: the document is the card alone, which
	// is all that a file list, a thumbnail or a search index needs.
	NoPlay bool
	// MaxAudio is the size in bytes of the largest file that is stored for
	// playing; a larger file gets its card without the audio, and a
	// warning. 0 is DefaultMaxAudio.
	MaxAudio int64
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// DefaultMaxAudio is the size of the largest file stored for playing when
// Options.MaxAudio is 0. The audio is held in memory while the document is
// built, and by the viewer that plays it.
const DefaultMaxAudio = 256 << 20

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Format is the file format: "mp3", "m4a", "aac", "flac", "ogg", "wav"
	// or "aiff".
	Format string
	// Codec names the audio coding ("MP3", "AAC", "Vorbis", "PCM"); "" when
	// it is not known.
	Codec string
	// Duration is the length in seconds (0 when it is not known), Bitrate
	// the bits per second, SampleRate the samples per second and Channels
	// the number of channels (each 0 when not known).
	Duration   float64
	Bitrate    int
	SampleRate int
	Channels   int
	// Cover is the format of the cover art the page shows ("jpeg", "png",
	// …), "" when the file has none it can show; CoverWidth and
	// CoverHeight are its size in pixels.
	Cover                   string
	CoverWidth, CoverHeight int
	// Audio is the size in bytes of the audio stored for playing; 0 when
	// the document does not play.
	Audio int64
	// Lyrics is the number of lines of lyrics on the card, and Synced
	// reports that they have the times they are sung at.
	Lyrics int
	Synced bool
	// Chapters is the number of chapters on the card.
	Chapters int
}

// maxSize bounds the files read.
const maxSize = 1 << 30

// maxCover bounds the cover art stored.
const maxCover = 64 << 20

// ConvertFile converts an audio file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.FileName == "" {
		o.FileName = filepath.Base(path)
	}
	return Convert(f, st.Size(), &o)
}

// Convert converts an audio file read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("audio: the file is larger than %d bytes", maxSize)
	}
	t, err := read(r, size)
	if err != nil {
		return nil, fmt.Errorf("audio: %w", err)
	}
	t.setLyrics()
	res := &Result{Doc: bdf.NewDocument(), Format: t.format, Codec: t.codec, Duration: t.duration, Bitrate: t.bitrate, SampleRate: t.sampleRate, Channels: t.channels,
		Lyrics: len(t.lyrics), Synced: t.lyricsSynced, Chapters: len(t.chapters)}
	warn := func(msg string) {
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			res.Warnings = append(res.Warnings, msg)
		}
	}
	for _, w := range t.warnings {
		warn(w)
	}
	if t.protected {
		warn("the audio is protected (DRM): its tags are read, but a player cannot play it")
	}
	doc := res.Doc
	doc.Meta.Source = t.format
	doc.Meta.DC = t.dublinCore()
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	title := doc.Meta.DC.Title.First()
	if title == "" {
		title = strings.TrimSuffix(opts.FileName, filepath.Ext(opts.FileName))
	}
	if title == "" {
		title = "Audio"
	}
	cover := pickCover(t, opts, warn)
	if cover != nil {
		res.Cover, res.CoverWidth, res.CoverHeight = cover.format, cover.w, cover.h
	}
	alt := "Cover of " + title
	if s := t.tags.first(keyAlbum); s != "" {
		alt = "Cover of " + s
	}
	w, h, obj, cues := buildCard(doc, t, cover, title, alt)
	view := doc.NewView("card", bdf.ViewFixed, title)
	view.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: obj})
	if data := storedAudio(r, size, t, opts, warn); data != nil {
		res.Audio = int64(len(data))
		view.Play = &bdf.Play{Audio: doc.AddPart(bdf.PartAudio, data).String(), Type: mimeTypes[t.format]}
		if cues != nil {
			view.Play.Cues = doc.AddPart(bdf.PartIndex, bdf.EncodeCues(cues)).String()
		}
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			return nil, fmt.Errorf("audio: %w", err)
		}
	}
	return res, nil
}

// storedAudio reads the file to store it for playing; nil when it is not
// stored: the options leave it out, it is larger than they allow, a player
// could not play it, or it could not be read whole.
func storedAudio(r io.ReaderAt, size int64, t *track, opts *Options, warn func(string)) []byte {
	if opts.NoPlay || t.protected || size <= 0 {
		return nil
	}
	limit := opts.MaxAudio
	if limit == 0 {
		limit = DefaultMaxAudio
	}
	if size > limit {
		warn(fmt.Sprintf("the file is larger than %s and its audio is not stored: the card does not play", byteSize(limit)))
		return nil
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(io.NewSectionReader(r, 0, size), data); err != nil {
		warn(fmt.Sprintf("the audio could not be read and is not stored: %v", err))
		return nil
	}
	return data
}

// byteSize writes a number of bytes in the unit that suits it ("256 MiB").
func byteSize(n int64) string {
	switch {
	case n >= 1<<20 && n%(1<<20) == 0:
		return fmt.Sprintf("%d MiB", n>>20)
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%d KiB", n>>10)
	}
	return fmt.Sprintf("%d bytes", n)
}

// pickCover picks the cover art to show and reads its size: a picture a
// browser can show, decodable (its header, at least), no larger than the
// bound.
func pickCover(t *track, opts *Options, warn func(string)) *coverImage {
	p := t.cover()
	if p == nil {
		return nil
	}
	if len(p.data) > maxCover {
		warn(fmt.Sprintf("the cover art is larger than %d bytes and is not shown", maxCover))
		return nil
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(p.data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		what := p.mime
		if what == "" {
			what = "of an unknown format"
		}
		warn(fmt.Sprintf("the cover art (%s) is not an image a browser shows, or is damaged, and is not shown", what))
		return nil
	}
	data := p.data
	if opts.Images.Mode == imgconv.Convert {
		if r, err := imgconv.Optimize(data, opts.Images); err == nil && r.Converted {
			data, format = r.Data, r.Format
		}
	}
	return &coverImage{data: data, format: format, w: cfg.Width, h: cfg.Height}
}

// Summary describes a result in a line ("1 page (cover 600 × 600 px, JPEG;
// MP3, 3:45, 320 kbps, 44.1 kHz, stereo; plays, 24 lines of synchronized
// lyrics)").
func (r *Result) Summary() string {
	s := "1 page ("
	if r.Cover != "" {
		s += fmt.Sprintf("cover %d × %d px, %s; ", r.CoverWidth, r.CoverHeight, strings.ToUpper(r.Cover))
	} else {
		s += "no cover; "
	}
	parts := []string{formatNames[r.Format]}
	if r.Codec != "" && !strings.EqualFold(r.Codec, formatNames[r.Format]) {
		parts[0] += " " + r.Codec
	}
	if r.Duration > 0 {
		parts = append(parts, clock(r.Duration))
	}
	if r.Bitrate > 0 {
		parts = append(parts, strconv.Itoa((r.Bitrate+500)/1000)+" kbps")
	}
	if r.SampleRate > 0 {
		parts = append(parts, kHz(r.SampleRate))
	}
	switch r.Channels {
	case 0:
	case 1:
		parts = append(parts, "mono")
	case 2:
		parts = append(parts, "stereo")
	default:
		parts = append(parts, fmt.Sprintf("%d channels", r.Channels))
	}
	s += strings.Join(parts, ", ")
	var more []string
	if r.Audio > 0 {
		more = append(more, "plays")
	}
	switch {
	case r.Lyrics > 0 && r.Synced:
		more = append(more, plural(r.Lyrics, "line")+" of synchronized lyrics")
	case r.Lyrics > 0:
		more = append(more, plural(r.Lyrics, "line")+" of lyrics")
	}
	if r.Chapters > 0 {
		more = append(more, plural(r.Chapters, "chapter"))
	}
	if len(more) > 0 {
		s += "; " + strings.Join(more, ", ")
	}
	return s + ")"
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// readAt returns up to n bytes at off.
func readAt(r io.ReaderAt, off, n int64) []byte {
	if n <= 0 {
		return nil
	}
	b := make([]byte, n)
	m, _ := r.ReadAt(b, off)
	return b[:m]
}

// read reads a file into a track: the ID3v2 tag it may start with, then
// the format the signature after it says.
func read(r io.ReaderAt, size int64) (*track, error) {
	head := readAt(r, 0, 32)
	start := int64(0)
	var id3 *id3v2
	if n := id3v2Size(head); n > 0 {
		b := readAt(r, 0, min(n, maxID3, size))
		if tag, err := parseID3v2(b); err == nil {
			id3 = tag
			if n > maxID3 {
				tag.truncated = true
			}
		}
		start = min(n, size)
		head = readAt(r, start, 32)
	}
	var t *track
	var err error
	switch {
	case bytes.HasPrefix(head, []byte("fLaC")):
		t = newTrack("flac")
		err = readFLAC(r, start, size, t)
	case bytes.HasPrefix(head, []byte("OggS")):
		t = newTrack("ogg")
		err = readOgg(r, size, t)
	case len(head) >= 12 && (string(head[:4]) == "RIFF" || string(head[:4]) == "RF64") && string(head[8:12]) == "WAVE":
		t, err = readWAV(r, size)
	case len(head) >= 12 && string(head[:4]) == "FORM" && (string(head[8:12]) == "AIFF" || string(head[8:12]) == "AIFC"):
		t, err = readAIFF(r, size)
	case len(head) >= 12 && string(head[4:8]) == "ftyp":
		switch mp4Kind(head, r, size) {
		case "audio":
			t, err = readMP4(r, size)
		case "video":
			return nil, fmt.Errorf("the file has a video track")
		default:
			return nil, fmt.Errorf("the MP4 file has no audio track")
		}
	default:
		// MPEG audio or ADTS frames, after junk at most
		end := size
		tail := readAt(r, size-128, 128)
		var v1 tags
		if v, ok := readID3v1(tail); ok {
			v1 = v
			end -= 128
		}
		if ape := apeTagSize(readAt(r, end-32, 32)); ape > 0 && ape < end {
			end -= ape
		}
		switch {
		case isMPEG(r, start, end):
			t = newTrack("mp3")
			t.setStream(scanMPEG(r, start, end))
		case isADTS(r, start, end):
			t = newTrack("aac")
			t.setStream(scanADTS(r, start, end))
		default:
			if info, ok := scanMPEG(r, start, end); ok {
				t = newTrack("mp3")
				t.setStream(info, true)
			} else if info, ok := scanADTS(r, start, end); ok {
				t = newTrack("aac")
				t.setStream(info, true)
			} else if id3 != nil {
				t = newTrack("mp3")
				t.warn("the audio after the ID3 tag was not recognized")
			} else {
				return nil, fmt.Errorf("not an MP3, AAC, FLAC, Ogg, WAV, AIFF or M4A file")
			}
		}
		if id3 != nil {
			t.applyID3v2(id3)
		}
		if v1 != nil {
			t.tags.fill(v1)
		}
		return t, nil
	}
	if err != nil {
		return nil, err
	}
	if id3 != nil {
		// an ID3 tag before another format's own tags: of lower precedence
		other := newTrack(t.format)
		other.applyID3v2(id3)
		t.tags.fill(other.tags)
		t.pictures = append(t.pictures, other.pictures...)
		if len(t.chapters) == 0 {
			t.chapters = other.chapters
		}
		if len(t.synced) == 0 {
			t.synced = other.synced
		}
		t.warnings = append(t.warnings, other.warnings...)
	}
	return t, nil
}

// setStream takes what a scan of the frames found.
func (t *track) setStream(info *streamInfo, ok bool) {
	if !ok || info == nil {
		return
	}
	t.codec, t.duration, t.bitrate, t.sampleRate, t.channels, t.vbr, t.scanned = info.codec, info.duration, info.bitrate, info.sampleRate, info.channels, info.vbr, info.frames
	if math.IsInf(t.duration, 0) || math.IsNaN(t.duration) {
		t.duration = 0
	}
}
