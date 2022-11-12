package fswatcher

import (
	"os"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

func New() *watcher {
	return &watcher{
		lastmap: cmap.New[int64](),
		fnmap:   cmap.New[func()](),
	}
}

type watcher struct {
	exit    chan struct{}
	ticker  *time.Ticker
	lastmap cmap.ConcurrentMap[string, int64]
	fnmap   cmap.ConcurrentMap[string, func()]
	once    sync.Once
}

func (w *watcher) Close() {
	w.lastmap.Clear()
	w.fnmap.Clear()
	w.ticker.Stop()
	w.exit <- struct{}{}
	close(w.exit)
}

func (w *watcher) Add(fpath string, fn func()) error {
	s, err := os.Stat(fpath)
	if err != nil {
		return err
	}
	w.lastmap.Set(fpath, s.ModTime().Unix())
	w.fnmap.Set(fpath, fn)
	return nil
}

func (w *watcher) Start() {
	w.once.Do(func() {
		w.exit = make(chan struct{})
		w.ticker = time.NewTicker(time.Minute)
		go w.timer()
	})
}

func (w *watcher) do() {
	for item := range w.lastmap.IterBuffered() {
		fpath, lastime := item.Key, item.Val
		s, err := os.Stat(fpath)
		if err != nil {
			continue
		}
		last := s.ModTime().Unix()
		if last != lastime {
			fn, ok := w.fnmap.Get(fpath)
			if !ok {
				continue
			}
			fn()
			w.lastmap.Set(fpath, last)
		}
	}
}

func (w *watcher) timer() {
	for {
		select {
		case <-w.ticker.C:
			w.do()
		case <-w.exit:
			return
		}
	}
}
