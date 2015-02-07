package watch

import (
	"io"
	"log"

	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/glob"
	"gopkg.in/fsnotify.v1"
)

func Watch(c *slurp.C, task func(), globs ...string) io.Closer {

	files, err := glob.Glob(globs...)

	if err != nil {
		c.Println(err)
		return nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		c.Println(err)
		return nil
	}

	for matchpair := range files {
		w.Add(matchpair.Name)
	}

	go func() {
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					task()
				}
			case err := <-w.Errors:
				c.Println(err)
			}
		}
	}()

	return w
}
