package watch

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/glob"
	"golang.org/x/exp/fsnotify"
)

type Closer interface {
	Close() error
}

func Watch(c *slurp.C, task func(string), globs ...string) Closer {

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
		w.Watch(matchpair.Name)
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event != nil && event.IsModify() && !event.IsAttrib(){
					task(event.Name)
				}
			case err := <-w.Error:
				c.Println(err)
			}
		}
	}()

	return w
}
