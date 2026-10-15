package imagebdf

import (
	"runtime"
	"sync"
)

// maxThreads bounds the goroutines one drawing splits the passes over a
// large canvas among: the conversion of the image drawn, the background,
// fills of large areas and the products of large clips. Pages drawn in
// parallel (Pages) divide them among their workers.
const maxThreads = 8

// parallelMinPixels is the least pixels of a pass that is split among
// goroutines: a pass over fewer is done before another goroutine would
// start on it. A variable, so that tests may split small passes.
var parallelMinPixels = 1 << 18

// defaultThreads is the number of threads of a new renderer.
func defaultThreads() int {
	return max(1, min(maxThreads, runtime.GOMAXPROCS(0)))
}

// splitRows calls fn for the bands of the rows [y0, y1) of a pass over
// width pixels a row, in up to threads goroutines at a time, and waits
// for them; k tells a band from the others (0 for the first), for the
// memory each uses. A pass over few pixels is done in the calling
// goroutine. A panic of a band is raised again in the caller, as a panic
// of the drawing.
func splitRows(threads, y0, y1, width int, fn func(y0, y1, k int)) {
	n := y1 - y0
	if n <= 0 {
		return
	}
	if threads <= 1 || plain || n*width < parallelMinPixels {
		fn(y0, y1, 0)
		return
	}
	threads = min(threads, n)
	var wg sync.WaitGroup
	panics := make([]any, threads)
	wg.Add(threads)
	for k := 0; k < threads; k++ {
		a, b := y0+n*k/threads, y0+n*(k+1)/threads
		go func() {
			defer wg.Done()
			defer func() {
				if p := recover(); p != nil {
					panics[k] = p
				}
			}()
			fn(a, b, k)
		}()
	}
	wg.Wait()
	for _, p := range panics {
		if p != nil {
			panic(p)
		}
	}
}
