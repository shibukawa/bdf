package music

import "math"

// Krumhansl–Kessler key profiles: how well each pitch class fits a major
// and a minor key on its tonic.
var (
	majorProfile = [12]float64{6.35, 2.23, 3.48, 2.33, 4.38, 4.09, 2.52, 5.19, 2.39, 3.66, 2.29, 2.88}
	minorProfile = [12]float64{6.33, 2.68, 3.52, 5.38, 2.60, 3.53, 2.54, 4.75, 3.98, 2.69, 3.34, 3.17}
)

// estimateKey finds the key whose profile correlates best with the
// durations of the pitch classes played (Krumhansl–Schmuckler), as a key
// signature of at most six sharps or flats. flats counts the notes the
// input spelled with flats minus those spelled with sharps, to choose
// between F♯ and G♭ major.
func estimateKey(hist [12]float64, flats int) KeySig {
	total := 0.0
	for _, v := range hist {
		total += v
	}
	if total == 0 {
		return KeySig{}
	}
	best, bestR := KeySig{}, math.Inf(-1)
	for tonic := 0; tonic < 12; tonic++ {
		for _, minor := range []bool{false, true} {
			prof := majorProfile
			if minor {
				prof = minorProfile
			}
			var x, y [12]float64
			for i := 0; i < 12; i++ {
				x[i] = hist[(tonic+i)%12]
				y[i] = prof[i]
			}
			r := correlation(x[:], y[:])
			if r > bestR+1e-9 {
				bestR = r
				best = keyOf(tonic, minor, flats)
			}
		}
	}
	return best
}

// keyOf returns the key signature of a key by its tonic pitch class.
func keyOf(tonic int, minor bool, flats int) KeySig {
	major := tonic
	if minor {
		major = (tonic + 3) % 12 // the relative major
	}
	// fifths of each major tonic, preferring the signature with fewer
	// accidentals and flats for the ambiguous F♯/G♭
	fifths := [12]int{0, -5, 2, -3, 4, -1, 6, 1, -4, 3, -2, 5}[major]
	if fifths == 6 && flats > 0 {
		fifths = -6
	}
	if minor && fifths == 6 && flats > 0 {
		fifths = -6
	}
	return KeySig{Fifths: fifths, Minor: minor}
}

func correlation(x, y []float64) float64 {
	n := float64(len(x))
	var mx, my float64
	for i := range x {
		mx += x[i]
		my += y[i]
	}
	mx /= n
	my /= n
	var sxy, sxx, syy float64
	for i := range x {
		dx, dy := x[i]-mx, y[i]-my
		sxy += dx * dy
		sxx += dx * dx
		syy += dy * dy
	}
	if sxx == 0 || syy == 0 {
		return 0
	}
	return sxy / math.Sqrt(sxx*syy)
}

// lineOfFifths is the position of each step natural on the line of fifths
// (F -1, C 0, G 1, …, B 5).
var lineOfFifths = [7]int{0, 2, 4, -1, 1, 3, 5}

// spell writes a MIDI note number as a pitch in a key: the notes of the
// key as it spells them, the others with the accidental that keeps them
// nearest the key on the line of fifths (minor keys lean to sharps for
// their raised sixth and seventh). hint is the input's spelling (1 sharp,
// -1 flat), used for notes outside the key.
func spell(key int, k KeySig, hint int) Pitch {
	pc := ((key % 12) + 12) % 12
	oct := key/12 - 1
	if key < 0 {
		oct = (key-11)/12 - 1
	}
	center := float64(k.Fifths) + 2
	if k.Minor {
		center += 1.5
	}
	best := Pitch{}
	bestD := math.Inf(1)
	for step := 0; step < 7; step++ {
		for alter := -2; alter <= 2; alter++ {
			if (stepSemitones[step]+alter+12)%12 != pc {
				continue
			}
			lof := lineOfFifths[step] + 7*alter
			d := math.Abs(float64(lof) - center)
			diatonic := lof >= k.Fifths-1 && lof <= k.Fifths+5
			switch {
			case diatonic:
				d -= 100
			case hint > 0 && alter == 1, hint < 0 && alter == -1:
				d -= 10
			case alter == 2 || alter == -2:
				d += 20
			}
			if d < bestD-1e-9 || d < bestD+1e-9 && alter > best.Alter && k.Fifths >= 0 {
				bestD = d
				o := oct
				// B♯ belongs to the octave below its sounding C, C♭ above
				if step == 6 && stepSemitones[step]+alter >= 12 {
					o--
				}
				if step == 0 && alter < 0 {
					o++
				}
				best = Pitch{Step: step, Alter: alter, Octave: o}
			}
		}
	}
	return best
}
