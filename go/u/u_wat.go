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
	TermClearer
}

type FsEventer interface {
	OnFsEvent(Ctx, notify.EventInfo)
}

// Optional extension to `FsEventer`.
type FsEventSkipper interface {
	ShouldSkipFsEvent(notify.EventInfo) bool
}

func NotifyEventPath(eve notify.EventInfo) (_ string) {
	if eve != nil {
		return eve.Path()
	}
	return
}

/*
TODO: new FS events should kill the current call to `.Runner.OnFsEvent` by
canceling the context passed to it. There are common cases where multiple FS
events are generated almost simultaneously, for example when multiple files are
saved at once. Killing and restarting would prevent inconsistent states.
*/
type Watcher[A FsEventer] struct {
	WatcherCommon
	Runner A
	Filter notify.Event  // Allowed event types, default all.
	Delay  time.Duration // Time window for ignoring subsequent FS events.
	IsDir  bool          // `true` = dir, `false` = file.
	Create bool          // Auto-create if missing.
	events chan notify.EventInfo
}

/*
Note: defining this method on a value type, rather than a pointer type, makes it
safe for external callers to make multiple `.Run` calls on the same value,
without concurrency hazards.
*/
func (self Watcher[_]) Run(ctx Ctx) {
	self.NormIgnore()

	if self.Init {
		self.Runner.OnFsEvent(ctx, nil)
	}

	if self.Touched() {
		// File events are async. The event from creating a file or directory
		// can and does arrive after we've started watching. TODO recheck.
		self.wait(ctx)
	}

	defer self.unwatch()
	self.watch()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			if self.Verb {
				log.Println(`context canceled, watcher stopping`)
			}
			return

		case eve := <-self.validEventsChan():
			self.onFsEvent(ctx, eve)
		}
	}
}

func (self *Watcher[_]) unwatch() {
	tar := self.events
	if tar == nil {
		return
	}

	notify.Stop(tar)
	if self.Verb {
		log.Printf(`unwatching %q`, self.Path)
	}
	self.events = nil
}

func (self *Watcher[_]) watch() {
	self.unwatch()

	/**
	Replacing the channel, instead of reusing it, is one of our measures for
	avoiding self-triggering via events generated by `self.Runner.OnFsEvent`.
	If we stop and then watch on the same channel, the channel may receive FS
	events generated while we weren't watching. This is also why we use an
	unbuffered channel.
	*/
	self.events = make(chan notify.EventInfo)

	if self.Verb {
		log.Printf(`watching %q`, self.Path)
	}
	gg.Try(notify.Watch(self.Pattern(), self.events, gg.Or(self.Filter, notify.All)))
}

func (self *Watcher[_]) watchDelayed(ctx Ctx) {
	self.wait(ctx)
	self.watch()
}

func (self *Watcher[_]) onFsEvent(ctx Ctx, eve notify.EventInfo) {
	if self.ShouldSkipFsEvent(eve) {
		log.Println(`ignoring file event:`, eve)
		return
	}

	if self.Clear {
		gg.TermClearHard()
	}

	if self.Verb {
		log.Println(`file event:`, eve)
	}

	/**
	Ideally, this would prevent us from self-triggering FS events due to
	`self.Runner.OnFsEvent`. However, watch/unwatch and event processing
	seems asynchronous in the library we use, so in practice this doesn't
	seem to help. We still rely on delays to avoid self-triggering.
	*/
	self.unwatch()
	defer self.watch()
	// defer self.watchDelayed(ctx)

	self.Runner.OnFsEvent(ctx, eve)
}

/*
Supported for technical reasons. If the target doesn't exist, the watcher
library returns an error. Auto-creating the target avoids that.
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

func (self Watcher[_]) ShouldSkipFsEvent(eve notify.EventInfo) bool {
	if eve == nil || self.ShouldSkipPath(eve.Path()) {
		return true
	}
	impl := gg.AnyAs[FsEventSkipper](self.Runner)
	return impl != nil && impl.ShouldSkipFsEvent(eve)
}

func (self Watcher[_]) ShouldSkipPath(path string) bool {
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

// Internal sanity check. Note: in Go, listening on a nil chan blocks forever.
func (self Watcher[_]) validEventsChan() chan notify.EventInfo {
	tar := self.events
	if tar == nil {
		panic(gg.Errf(`unable to listen on nil %T`, tar))
	}
	return tar
}
