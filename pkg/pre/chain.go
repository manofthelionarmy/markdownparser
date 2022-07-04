package pre

import "io"

// ChainProcessors chains together processors, and what is returned is the source for the last preproccers pipe writer
func ChainProcessors(src io.Reader, processors ...Processor) (prevPr io.Reader) {
	var currPw io.Writer
	for i, p := range processors {
		if i == 0 {
			prevPr, currPw = io.Pipe()
			p.SetSource(src)
			p.SetTarget(currPw)
		} else {
			p.SetSource(prevPr)
			prevPr, currPw = io.Pipe()
			p.SetTarget(currPw)
		}
	}
	return
}
