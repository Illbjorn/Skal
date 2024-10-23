package skal

type srcFile struct {
	Path    string
	Content string
	Import  bool
}

type job struct {
	Main       *srcFile
	OutputPath string
	Imports    []*srcFile
}

func (j *job) Gen() chan *srcFile {
	ch := make(chan *srcFile)

	go func() {
		defer close(ch)
		// Send the imports first.
		for _, f := range j.Imports {
			ch <- f
		}

		// Finish with the main entrypoint.
		ch <- j.Main
	}()

	return ch
}
