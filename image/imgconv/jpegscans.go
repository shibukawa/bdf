package imgconv

// MaxJPEGScans bounds the scans of a JPEG image that are read. A decoder
// goes through every block of an image once a scan, however little the
// scan says of them, so a file that is nothing but scans takes as long as
// thousands of images do: 65 kilobytes with 500 scans of an image of
// 8192 × 8192 pixels take twelve seconds. Encoders write about ten scans
// for a progressive image, and four at most for another.
const MaxJPEGScans = 100

// LimitJPEGScans returns a JPEG image without the scans after its first
// MaxJPEGScans, and whether it had more. What is left is the picture as
// far as those scans draw it, as a browser shows an image that is still
// arriving. Decode reads no more of an image than this; other data is
// returned as it is.
func LimitJPEGScans(data []byte) ([]byte, bool) {
	end := jpegScansEnd(data, MaxJPEGScans)
	if end < 0 {
		return data, false
	}
	cut := make([]byte, end+2)
	copy(cut, data[:end])
	cut[end], cut[end+1] = 0xff, 0xd9 // the end of the image
	return cut, true
}

// jpegScansEnd returns where the marker of the scan after the first max
// scans of a JPEG image is, or -1 for an image that has no such scan.
//
// It passes over the image as image/jpeg does between its segments: a
// marker is 0xff, any number of them, and a byte that is neither zero (a
// 0xff of the data of a scan is followed by one) nor a restart marker;
// what is not a marker is passed over, which is how the data of a scan is.
// A segment is as long as it says. Where the decoder would give up, this
// may read on: it only has to find every scan the decoder does.
func jpegScansEnd(data []byte, max int) int {
	if Sniff(data) != "jpeg" {
		return -1
	}
	scans := 0
	for i := 2; ; {
		for i < len(data) && data[i] != 0xff {
			i++
		}
		at := i
		for i < len(data) && data[i] == 0xff {
			i++
		}
		if i >= len(data) {
			return -1
		}
		marker := data[i]
		i++
		switch {
		case marker == 0x00, 0xd0 <= marker && marker <= 0xd7:
			continue // a 0xff of the data of a scan, or a restart marker
		case marker == 0xd9:
			return -1 // the end of the image
		}
		if i+2 > len(data) {
			return -1
		}
		n := int(data[i])<<8 | int(data[i+1])
		if n < 2 {
			return -1 // no segment: the decoder gives up here
		}
		if marker == 0xda { // the start of a scan
			if scans == max {
				return at
			}
			scans++
		}
		i += n
	}
}
