package u

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

type WatcherCommon struct {
	Pathed
	Verbose
	Inited
	Ignored
}

type FsEventer interface {
	OnFsEvent(Ctx, notify.EventInfo)
}

/*
TODO: new FS events should kill the current call to `.Runner.OnFsEvent` by
canceling the context passed to it. There are common cases where multiple FS
events are generated almost simultaneously, for example when multiple files are
saved at once. Killing and restarting prevents inconsistent states.
*/
type Watcher[A FsEventer] struct {
	WatcherCommon
	Runner A
	Filter notify.Event  // Allowed event types, default all.
	Delay  time.Duration // Time window for ignoring subsequent FS events.
	IsDir  bool          // `true` = dir, `false` = file.
	Create bool          // Auto-create if missing.
	Clear  bool          // Clear terminal on FS event.
}

func (self Watcher[_]) Run(ctx Ctx) {
	self.NormIgnore()

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
			if self.Skip(eve) {
				log.Println(`skipping file event:`, eve)
				continue
			}
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

/*
Normalizes relative paths to absolute for compatibility with
`notify.EventInfo.Path` which returns absolute paths.
*/
func (self *Watcher[_]) NormIgnore() {
	for ind, val := range self.Ignore {
		self.Ignore[ind] = gg.Try1(filepath.Abs(val))
	}
}

func (self Watcher[_]) Skip(eve notify.EventInfo) bool {
	return eve == nil || self.SkipPath(eve.Path())
}

func (self Watcher[_]) SkipPath(path string) bool {
	return gg.Some(self.Ignore, func(val string) bool {
		return IsPathAncestorOf(val, path)
	})
}

func (self Watcher[_]) wait(ctx Ctx) {
	if self.Delay >= 0 {
		Wait(ctx, self.getDelay())
	}
}

func (self Watcher[_]) getDelay() time.Duration {
	return gg.Or(self.Delay, time.Millisecond*100)
}

func NotifyEventPath(eve notify.EventInfo) (_ string) {
	if eve != nil {
		return eve.Path()
	}
	return
}
