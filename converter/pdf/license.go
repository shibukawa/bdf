package pdf

import "github.com/shibukawa/bdf/converter/internal/sfnt"

// Embedding permissions of the OS/2 fsType field (OpenType spec). Bits 0–3
// hold one usage permission; when several are set the least restrictive one
// applies. Bits 8 and 9 add restrictions on top of it.
const (
	fsRestricted   = 0x0002 // must not be embedded without the owner's permission
	fsPreviewPrint = 0x0004 // may be embedded in documents that are only viewed and printed
	fsEditable     = 0x0008 // may be embedded in documents that are edited
	fsNoSubsetting = 0x0100 // must be embedded whole
	fsBitmapOnly   = 0x0200 // only bitmaps may be embedded
)

// fontLicense is what a font program says about embedding it.
type fontLicense struct {
	known   bool   // the program has an OS/2 table (TrueType, OpenType)
	fsType  uint16 // OS/2 fsType
	notices []sfnt.NameRecord
}

// restricted reports a Restricted License font: bit 1 set and neither of the
// less restrictive Preview & Print or Editable bits.
func (l fontLicense) restricted() bool {
	return l.known && l.fsType&(fsRestricted|fsPreviewPrint|fsEditable) == fsRestricted
}

func (l fontLicense) bitmapOnly() bool   { return l.known && l.fsType&fsBitmapOnly != 0 }
func (l fontLicense) noSubsetting() bool { return l.known && l.fsType&fsNoSubsetting != 0 }

// outFSType is the fsType written to the rebuilt font. A BDF document is
// only viewed, so the original permission carries over as is; a program that
// does not say (bare CFF, TrueType without OS/2) gets Preview & Print rather
// than the Installable of a zero field.
func (l fontLicense) outFSType() uint16 {
	if l.known {
		return l.fsType
	}
	return fsPreviewPrint
}

// license returns the embedding permission and notices of a font program.
func (p *fontProgram) license() fontLicense {
	var l fontLicense
	if p.sf != nil {
		l.known, l.fsType, l.notices = p.sf.HasFSType, p.sf.FSType, p.sf.Notices
	}
	if len(l.notices) == 0 && p.cff != nil {
		// A bare CFF program has only the Top DICT Copyright and Notice.
		if p.cff.copyright != "" {
			l.notices = append(l.notices, sfnt.NameRecord{ID: 0, Text: p.cff.copyright})
		}
		if p.cff.notice != "" && p.cff.notice != p.cff.copyright {
			id := uint16(7) // trademark notice, where Adobe fonts put it
			if p.cff.copyright == "" {
				id = 0
			}
			l.notices = append(l.notices, sfnt.NameRecord{ID: id, Text: p.cff.notice})
		}
	}
	return l
}
