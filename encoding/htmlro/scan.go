package htmlro

// The tokenizer states, after the tokenizer of golang.org/x/net/html
// (token.go), which this file follows closely: the same functions, the same
// cursor arithmetic on raw.end, so that the two tokenize alike. What
// differs is what they keep: this one indexes the current token in place
// and copies nothing.

import "io"

// skipWhiteSpace skips past any white space.
func (r *Reader) skipWhiteSpace() {
	if r.err != nil {
		return
	}
	for {
		c := r.readByte()
		if r.err != nil {
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f':
			// No-op.
		default:
			r.raw.end--
			return
		}
	}
}

// readRawOrRCDATA reads until the next "</foo>", where "foo" is the element
// of rawTag, typically script or textarea.
func (r *Reader) readRawOrRCDATA() {
	if r.rawTag == rawScript {
		r.readScript()
		r.textIsRaw = true
		r.rawTag = rawNone
		return
	}
loop:
	for {
		c := r.readByte()
		if r.err != nil {
			break loop
		}
		if c != '<' {
			continue loop
		}
		c = r.readByte()
		if r.err != nil {
			break loop
		}
		if c != '/' {
			r.raw.end--
			continue loop
		}
		if r.readRawEndTag() || r.err != nil {
			break loop
		}
	}
	r.data.end = r.raw.end
	// A textarea's or title's RCDATA can contain escaped entities.
	r.textIsRaw = r.rawTag != rawTextarea && r.rawTag != rawTitle
	r.rawTag = rawNone
}

// readRawEndTag attempts to read a tag like "</foo>", where "foo" is the
// element of rawTag. If it succeeds, it backs up the input position to
// reconsume the tag and returns true. Otherwise it returns false. The
// opening "</" has already been consumed.
func (r *Reader) readRawEndTag() bool {
	name := rawTagNames[r.rawTag]
	for i := 0; i < len(name); i++ {
		c := r.readByte()
		if r.err != nil {
			return false
		}
		if c != name[i] && c != name[i]-('a'-'A') {
			r.raw.end--
			return false
		}
	}
	c := r.readByte()
	if r.err != nil {
		return false
	}
	switch c {
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		// The 3 is 2 for the leading "</" plus 1 for the trailing character c.
		r.raw.end -= 3 + len(name)
		return true
	}
	r.raw.end--
	return false
}

// readScript reads until the next </script> tag, following the byzantine
// rules for escaping/hiding the closing tag.
func (r *Reader) readScript() {
	// Not a deferred closure: TinyGo would allocate it per call.
	r.scriptData()
	r.data.end = r.raw.end
}

func (r *Reader) scriptData() {
	var c byte

scriptData:
	c = r.readByte()
	if r.err != nil {
		return
	}
	if c == '<' {
		goto scriptDataLessThanSign
	}
	goto scriptData

scriptDataLessThanSign:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '/':
		goto scriptDataEndTagOpen
	case '!':
		goto scriptDataEscapeStart
	}
	r.raw.end--
	goto scriptData

scriptDataEndTagOpen:
	if r.readRawEndTag() || r.err != nil {
		return
	}
	goto scriptData

scriptDataEscapeStart:
	c = r.readByte()
	if r.err != nil {
		return
	}
	if c == '-' {
		goto scriptDataEscapeStartDash
	}
	r.raw.end--
	goto scriptData

scriptDataEscapeStartDash:
	c = r.readByte()
	if r.err != nil {
		return
	}
	if c == '-' {
		goto scriptDataEscapedDashDash
	}
	r.raw.end--
	goto scriptData

scriptDataEscaped:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDash
	case '<':
		goto scriptDataEscapedLessThanSign
	}
	goto scriptDataEscaped

scriptDataEscapedDash:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	}
	goto scriptDataEscaped

scriptDataEscapedDashDash:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	case '>':
		goto scriptData
	}
	goto scriptDataEscaped

scriptDataEscapedLessThanSign:
	c = r.readByte()
	if r.err != nil {
		return
	}
	if c == '/' {
		goto scriptDataEscapedEndTagOpen
	}
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		goto scriptDataDoubleEscapeStart
	}
	r.raw.end--
	goto scriptData

scriptDataEscapedEndTagOpen:
	if r.readRawEndTag() || r.err != nil {
		return
	}
	goto scriptDataEscaped

scriptDataDoubleEscapeStart:
	r.raw.end--
	for i := 0; i < len("script"); i++ {
		c = r.readByte()
		if r.err != nil {
			return
		}
		if c != "script"[i] && c != "SCRIPT"[i] {
			r.raw.end--
			goto scriptDataEscaped
		}
	}
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		goto scriptDataDoubleEscaped
	}
	r.raw.end--
	goto scriptDataEscaped

scriptDataDoubleEscaped:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDash:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDashDash:
	c = r.readByte()
	if r.err != nil {
		return
	}
	switch c {
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	case '>':
		goto scriptData
	}
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedLessThanSign:
	c = r.readByte()
	if r.err != nil {
		return
	}
	if c == '/' {
		goto scriptDataDoubleEscapeEnd
	}
	r.raw.end--
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapeEnd:
	if r.readRawEndTag() {
		r.raw.end += len("</script>")
		goto scriptDataEscaped
	}
	if r.err != nil {
		return
	}
	goto scriptDataDoubleEscaped
}

// readComment reads the next comment token starting with "<!--". The
// opening "<!--" has already been consumed.
func (r *Reader) readComment() {
	r.data.start = r.raw.end
	r.commentData()
	if r.data.end < r.data.start {
		// It's a comment with no data, like <!-->.
		r.data.end = r.data.start
	}
}

func (r *Reader) commentData() {
	var dashCount int
	beginning := true
	for {
		c := r.readByte()
		if r.err != nil {
			r.data.end = r.calculateAbruptCommentDataEnd()
			return
		}
		switch c {
		case '-':
			dashCount++
			continue
		case '>':
			if dashCount >= 2 || beginning {
				r.data.end = r.raw.end - len("-->")
				return
			}
		case '!':
			if dashCount >= 2 {
				c = r.readByte()
				if r.err != nil {
					r.data.end = r.calculateAbruptCommentDataEnd()
					return
				} else if c == '>' {
					r.data.end = r.raw.end - len("--!>")
					return
				} else if c == '-' {
					dashCount = 1
					beginning = false
					continue
				}
			}
		}
		dashCount = 0
		beginning = false
	}
}

func (r *Reader) calculateAbruptCommentDataEnd() int {
	raw := r.buf[r.raw.start:r.raw.end]
	const prefixLen = len("<!--")
	if len(raw) >= prefixLen {
		raw = raw[prefixLen:]
		if hasSuffix(raw, "--!") {
			return r.raw.end - 3
		} else if hasSuffix(raw, "--") {
			return r.raw.end - 2
		} else if hasSuffix(raw, "-") {
			return r.raw.end - 1
		}
	}
	return r.raw.end
}

func hasSuffix(b []byte, suffix string) bool {
	if len(b) < len(suffix) {
		return false
	}
	b = b[len(b)-len(suffix):]
	for i := range b {
		if b[i] != suffix[i] {
			return false
		}
	}
	return true
}

// readUntilCloseAngle reads until the next ">".
func (r *Reader) readUntilCloseAngle() {
	r.data.start = r.raw.end
	for {
		c := r.readByte()
		if r.err != nil {
			r.data.end = r.raw.end
			return
		}
		if c == '>' {
			r.data.end = r.raw.end - len(">")
			return
		}
	}
}

// readMarkupDeclaration reads the next token starting with "<!". It might
// be a "<!--comment-->", a "<!DOCTYPE foo>", a "<![CDATA[section]]>" or
// "<!a bogus comment". The opening "<!" has already been consumed.
func (r *Reader) readMarkupDeclaration() Kind {
	r.data.start = r.raw.end
	var c [2]byte
	for i := 0; i < 2; i++ {
		c[i] = r.readByte()
		if r.err != nil {
			// bogus comment
			r.data.end = r.raw.end
			if i == 1 && c[0] == '>' {
				r.data.end--
			}
			return Comment
		}
	}
	if c[0] == '-' && c[1] == '-' {
		r.readComment()
		return Comment
	}
	r.raw.end -= 2
	if r.readDoctype() {
		return Doctype
	}
	if r.allowCDATA && r.readCDATA() {
		r.convertNUL = true
		return Text
	}
	// It's a bogus comment.
	r.readUntilCloseAngle()
	return Comment
}

// readDoctype attempts to read a doctype declaration and returns true if
// successful. The opening "<!" has already been consumed.
func (r *Reader) readDoctype() bool {
	const s = "DOCTYPE"
	for i := 0; i < len(s); i++ {
		c := r.readByte()
		if r.err != nil {
			if r.err == io.EOF {
				// Back up to read the fragment of "DOCTYPE" again, reset
				// err to signal EOF on the next call
				r.raw.end = r.data.start
				r.err = nil
				return false
			}
			r.data.end = r.raw.end
			return false
		}
		if c != s[i] && c != s[i]+('a'-'A') {
			// Back up to read the fragment of "DOCTYPE" again.
			r.raw.end = r.data.start
			return false
		}
	}
	if r.skipWhiteSpace(); r.err != nil {
		r.data.start = r.raw.end
		r.data.end = r.raw.end
		return true
	}
	r.readUntilCloseAngle()
	return true
}

// readCDATA attempts to read a CDATA section and returns true if
// successful. The opening "<!" has already been consumed.
func (r *Reader) readCDATA() bool {
	const s = "[CDATA["
	for i := 0; i < len(s); i++ {
		c := r.readByte()
		if r.err != nil {
			if r.err == io.EOF {
				// Back up to read the fragment of "[CDATA[" again, reset
				// err to signal EOF on the next call
				r.raw.end = r.data.start
				r.err = nil
				return false
			}
			r.data.end = r.raw.end
			return false
		}
		if c != s[i] {
			// Back up to read the fragment of "[CDATA[" again.
			r.raw.end = r.data.start
			return false
		}
	}
	r.data.start = r.raw.end
	brackets := 0
	for {
		c := r.readByte()
		if r.err != nil {
			r.data.end = r.raw.end
			return true
		}
		switch c {
		case ']':
			brackets++
		case '>':
			if brackets >= 2 {
				r.data.end = r.raw.end - len("]]>")
				return true
			}
			brackets = 0
		default:
			brackets = 0
		}
	}
}

// startTagIs reports whether the start tag in buf[data.start:data.end]
// case-insensitively matches s, which is lower case.
func (r *Reader) startTagIs(s string) bool {
	return EqualFold(r.buf[r.data.start:r.data.end], s)
}

// readStartTag reads the next start tag token. The opening "<a" has
// already been consumed, where 'a' means anything in [A-Za-z].
func (r *Reader) readStartTag() Kind {
	r.readTag(true)
	if r.err != nil {
		return None
	}
	// Several tags flag the next token as raw text.
	c := r.buf[r.data.start]
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	switch c {
	case 'i':
		if r.startTagIs("iframe") {
			r.rawTag = rawIframe
		}
	case 'n':
		switch {
		case r.startTagIs("noembed"):
			r.rawTag = rawNoembed
		case r.startTagIs("noframes"):
			r.rawTag = rawNoframes
		case r.startTagIs("noscript"):
			r.rawTag = rawNoscript
		}
	case 'p':
		if r.startTagIs("plaintext") {
			r.rawTag = rawPlaintext
		}
	case 's':
		switch {
		case r.startTagIs("script"):
			r.rawTag = rawScript
		case r.startTagIs("style"):
			r.rawTag = rawStyle
		}
	case 't':
		switch {
		case r.startTagIs("textarea"):
			r.rawTag = rawTextarea
		case r.startTagIs("title"):
			r.rawTag = rawTitle
		}
	case 'x':
		if r.startTagIs("xmp") {
			r.rawTag = rawXmp
		}
	}
	// Look for a self-closing token (e.g. <br/>).
	//
	// The last character of the tag (ignoring the closing bracket) being a
	// solidus is not enough: the tag may have an unquoted attribute value
	// that ends in one (<p a=/>). So the last non-bracket character of the
	// tag (raw.end-2) must not be the last character of the last attribute
	// value (attr[n-1][1].end-1), if the tag has attributes.
	nAttrs := len(r.attr)
	if r.err == nil && r.buf[r.raw.end-2] == '/' && (nAttrs == 0 || r.raw.end-2 != r.attr[nAttrs-1][1].end-1) {
		return SelfClosingTag
	}
	return StartTag
}

// readTag reads the next tag token and its attributes. If saveAttr, those
// attributes are saved in attr, otherwise attr is set to an empty slice.
// The opening "<a" or "</a" has already been consumed, where 'a' means
// anything in [A-Za-z].
func (r *Reader) readTag(saveAttr bool) {
	r.attr = r.attr[:0]
	r.nAttr = 0
	// Read the tag name and attribute key/value pairs.
	r.readTagName()
	if r.skipWhiteSpace(); r.err != nil {
		return
	}
	for {
		c := r.readByte()
		if r.err != nil || c == '>' {
			break
		}
		r.raw.end--
		r.readTagAttrKey()
		r.readTagAttrVal()
		// Save pendingAttr if saveAttr and that attribute has a non-empty
		// key, and the key hasn't been seen before.
		if saveAttr && r.pendingAttr[0].start != r.pendingAttr[0].end && !r.seenAttr(r.buf[r.pendingAttr[0].start:r.pendingAttr[0].end]) {
			r.attr = append(r.attr, r.pendingAttr)
		}
		if r.skipWhiteSpace(); r.err != nil {
			break
		}
	}
}

// seenAttr reports whether an attribute named key, compared ASCII
// case-insensitively, was already read for the current tag.
func (r *Reader) seenAttr(key []byte) bool {
	for _, a := range r.attr {
		if equalASCIIFold(r.buf[a[0].start:a[0].end], key) {
			return true
		}
	}
	return false
}

func equalASCIIFold(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if 'A' <= x && x <= 'Z' {
			x += 'a' - 'A'
		}
		if 'A' <= y && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}

// readTagName sets data to the "div" in "<div k=v>". The cursor (raw.end)
// is positioned such that the first byte of the tag name (the "d" in
// "<div") has already been consumed.
func (r *Reader) readTagName() {
	r.data.start = r.raw.end - 1
	for {
		c := r.readByte()
		if r.err != nil {
			r.data.end = r.raw.end
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f':
			r.data.end = r.raw.end - 1
			return
		case '/', '>':
			r.raw.end--
			r.data.end = r.raw.end
			return
		}
	}
}

// readTagAttrKey sets pendingAttr[0] to the "k" in "<div k=v>".
// Precondition: err == nil.
func (r *Reader) readTagAttrKey() {
	r.pendingAttr[0].start = r.raw.end
	for {
		c := r.readByte()
		if r.err != nil {
			r.pendingAttr[0].end = r.raw.end
			return
		}
		switch c {
		case '=':
			if r.pendingAttr[0].start+1 == r.raw.end {
				// WHATWG 13.2.5.32, if we see an equals sign before the attribute name
				// begins, we treat it as a character in the attribute name and continue.
				continue
			}
			fallthrough
		case ' ', '\n', '\r', '\t', '\f', '/', '>':
			// WHATWG 13.2.5.33 Attribute name state
			// We need to reconsume the char in the after attribute name state to support the / character
			r.raw.end--
			r.pendingAttr[0].end = r.raw.end
			return
		}
	}
}

// readTagAttrVal sets pendingAttr[1] to the "v" in "<div k=v>".
func (r *Reader) readTagAttrVal() {
	r.pendingAttr[1].start = r.raw.end
	r.pendingAttr[1].end = r.raw.end
	if r.skipWhiteSpace(); r.err != nil {
		return
	}
	c := r.readByte()
	if r.err != nil {
		return
	}
	if c == '/' {
		// WHATWG 13.2.5.34 After attribute name state
		// U+002F SOLIDUS (/) - Switch to the self-closing start tag state.
		return
	}
	if c != '=' {
		r.raw.end--
		return
	}
	if r.skipWhiteSpace(); r.err != nil {
		return
	}
	quote := r.readByte()
	if r.err != nil {
		return
	}
	switch quote {
	case '>':
		r.raw.end--
		return

	case '\'', '"':
		r.pendingAttr[1].start = r.raw.end
		for {
			c := r.readByte()
			if r.err != nil {
				r.pendingAttr[1].end = r.raw.end
				return
			}
			if c == quote {
				r.pendingAttr[1].end = r.raw.end - 1
				return
			}
		}

	default:
		r.pendingAttr[1].start = r.raw.end - 1
		for {
			c := r.readByte()
			if r.err != nil {
				r.pendingAttr[1].end = r.raw.end
				return
			}
			switch c {
			case ' ', '\n', '\r', '\t', '\f':
				r.pendingAttr[1].end = r.raw.end - 1
				return
			case '>':
				r.raw.end--
				r.pendingAttr[1].end = r.raw.end
				return
			}
		}
	}
}

// next scans the next token and returns its kind, None at the end of the
// input or on an error, which err then holds.
func (r *Reader) next() Kind {
	r.raw.start = r.raw.end
	r.data.start = r.raw.end
	r.data.end = r.raw.end
	r.alt = r.alt[:0]
	if r.err != nil {
		return None
	}
	if r.rawTag != rawNone {
		if r.rawTag == rawPlaintext {
			// Read everything up to EOF.
			for r.err == nil {
				r.readByte()
			}
			r.data.end = r.raw.end
			r.textIsRaw = true
		} else {
			r.readRawOrRCDATA()
		}
		if r.data.end > r.data.start {
			r.convertNUL = true
			return Text
		}
	}
	r.textIsRaw = false
	r.convertNUL = false

loop:
	for {
		c := r.readByte()
		if r.err != nil {
			break loop
		}
		if c != '<' {
			continue loop
		}

		// Check if the '<' we have just read is part of a tag, comment
		// or doctype. If not, it's part of the accumulated text token.
		c = r.readByte()
		if r.err != nil {
			break loop
		}
		var kind Kind
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			kind = StartTag
		case c == '/':
			kind = EndTag
		case c == '!' || c == '?':
			// Comment stands for any of "<!--actual comments-->",
			// "<!DOCTYPE declarations>" and "<?xml processing instructions?>".
			kind = Comment
		default:
			// Reconsume the current character.
			r.raw.end--
			continue
		}

		// We have a non-text token, but we might have accumulated some text
		// before that. If so, we return the text first, and return the non-
		// text token on the subsequent call to next.
		if x := r.raw.end - len("<a"); r.raw.start < x {
			r.raw.end = x
			r.data.end = x
			return Text
		}
		switch kind {
		case StartTag:
			return r.readStartTag()
		case EndTag:
			c = r.readByte()
			if r.err != nil {
				break loop
			}
			if c == '>' {
				// "</>" does not generate a token at all. Generate an empty
				// comment to allow passthrough clients to pick up the data.
				return Comment
			}
			if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
				r.readTag(false)
				if r.err != nil {
					return None
				}
				return EndTag
			}
			r.raw.end--
			r.readUntilCloseAngle()
			return Comment
		case Comment:
			if c == '!' {
				return r.readMarkupDeclaration()
			}
			r.raw.end--
			r.readUntilCloseAngle()
			return Comment
		}
	}
	if r.raw.start < r.raw.end {
		r.data.end = r.raw.end
		return Text
	}
	return None
}
