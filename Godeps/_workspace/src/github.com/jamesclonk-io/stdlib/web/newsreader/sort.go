package newsreader

import (
	"math"
	"sort"
)

type feedSort struct {
	feeds Feeds
	by    func(f1, f2 *Feed) bool
}

func (fs *feedSort) Len() int {
	return len(fs.feeds)
}

func (fs *feedSort) Swap(l, r int) {
	fs.feeds[l], fs.feeds[r] = fs.feeds[r], fs.feeds[l]
}

func (fs *feedSort) Less(l, r int) bool {
	return fs.by(&fs.feeds[l], &fs.feeds[r])
}

func (feeds *Feeds) sortBy(by func(f1, f2 *Feed) bool) *Feeds {
	fs := &feedSort{
		feeds: *feeds,
		by:    by,
	}
	sort.Sort(fs)
	return feeds
}

func (feeds *Feeds) sortByLines() *Feeds {
	feeds.sortBy(func(f1, f2 *Feed) bool {
		return f1.Lines() > f2.Lines()
	})
	return feeds
}

func (feeds *Feeds) sortById() *Feeds {
	feeds.sortBy(func(f1, f2 *Feed) bool {
		return f1.Id < f2.Id
	})
	return feeds
}

func (feeds *Feeds) sort() *Feeds {
	feeds.sortBy(func(f1, f2 *Feed) bool {
		if math.Abs(float64(f1.Lines()-f2.Lines())) < 2 {
			return f1.Id < f2.Id
		}
		return f1.Lines() > f2.Lines()
	})
	return feeds
}
