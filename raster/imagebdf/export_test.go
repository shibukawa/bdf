package imagebdf

// SetPlain makes the renderers draw without the shortcuts that change no
// pixel, or with them again, for the tests that compare both.
func SetPlain(on bool) { plain = on }
