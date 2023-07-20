package u

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

type FsEventer interface {
	OnFsEvent(Ctx, notify.EventInfo)
}

type Watcher[A FsEventer] struct {
	Runner A
	Path   string
	Filter notify.Event  // Allowed event types, default all.
	Delay  time.Duration // Time window for ignoring subsequent FS events.
	IsDir  bool          // `true` = dir, `false` = file.
	Create bool          // Auto-create if missing.
	Verb   bool          // Verbose logging.
	Init   bool          // Run once before watching.
	Clear  bool          // Clear terminal on FS event.
}

func (self Watcher[_]) Run(ctx Ctx) {
	if self.Init {
		self.Runner.OnFsEvent(ctx, nil)
	}

	if self.Touched() {
		// File events are async. The event from creating a file or directory
		// can and does arrive after we've started watching.
		self.wait(ctx)
	}

	// Note: unbuffered channel makes it easier to avoid self-triggering.
	events := make(chan notify.EventInfo)
	defer notify.Stop(events)

	gg.Try(notify.Watch(self.Pattern(), events, gg.Or(self.Filter, notify.All)))

	if self.Verb {
		log.Printf(`watching %q`, self.Path)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			if self.Verb {
				log.Println(`context canceled, watcher stopping`)
			}
			return

		case eve := <-events:
			if self.Clear {
				gg.TermClearHard()
			}
			if self.Verb {
				log.Println(`file event:`, eve)
			}
			self.Runner.OnFsEvent(ctx, eve)
			self.wait(ctx)
		}
	}
}

/*
Supported for technical reasons. If the target doesn't exist, the watcher
library returns an error.
*/
func (self Watcher[_]) Touched() bool {
	if !self.Create {
		return false
	}
	if self.IsDir {
		return TouchedDirRec(self.Path)
	}
	return TouchedFileRec(self.Path)
}

func (self Watcher[_]) Pattern() string {
	if gg.IsNotZero(self.Path) && self.IsDir {
		return filepath.Join(self.Path, `...`)
	}
	return self.Path
}

func (self Watcher[_]) wait(ctx Ctx) {
	if self.Delay >= 0 {
		Wait(ctx, self.getDelay())
	}
}

func (self Watcher[_]) getDelay() time.Duration {
	return gg.Or(self.Delay, time.Millisecond*100)
}
