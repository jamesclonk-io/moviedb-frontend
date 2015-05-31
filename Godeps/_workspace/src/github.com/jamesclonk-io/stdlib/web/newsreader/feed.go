package newsreader

import (
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Sirupsen/logrus"
	rss "github.com/jteeuwen/go-pkg-rss"
)

var redditRx = regexp.MustCompile(`<br\/>\s+<a href="(.*)">\[link\]`)

type Feeds []Feed
type Feed struct {
	Id    int
	Title string
	URL   string
	Items []FeedItem
}

type FeedItem struct {
	Title    string
	URL      string
	Comments string
}

type Work struct {
	Index int
	URL   string
}

func (n *NewsReader) InitializeFeeds() {
	n.UpdateFeeds()
}

func (n *NewsReader) GetFeeds() (feeds Feeds, err error) {
	// check if enough time has passed and a feed update is needed
	if n.timestamp.Add(time.Duration(n.Config.CacheDuration)).Before(time.Now()) {
		go n.UpdateFeeds()
	}

	return n.Feeds, nil
}

func (n *NewsReader) UpdateFeeds() {
	n.timestamp = time.Now()
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var workers = 3
	var feeds Feeds
	var workChan = make(chan Work, 1000)
	var feedChan = make(chan Feed, 1000)
	var doneChan = make(chan bool, 3)
	defer close(doneChan)

	channelHandler := func(feedId int, feedChan chan Feed) func(*rss.Feed, []*rss.Channel) {
		return func(feed *rss.Feed, channels []*rss.Channel) {
			for _, channel := range channels {
				reddit := strings.Contains(getFeedTitleUrl(channel), "www.reddit.com") // flag feed as reddit

				var feedItems []FeedItem
				for i, item := range channel.Items {
					if i < 30 { // max items per feed
						feedItem := FeedItem{
							Title:    item.Title,
							URL:      getFeedItemUrl(item),
							Comments: item.Comments,
						}
						// special handling for reddit (links & comments)
						if reddit {
							// comments url for reddit rss is actually the feed item link itself
							feedItem.Comments = getFeedItemUrl(item)
							feedItem.URL = getRedditFeedItemUrl(item)
						}

						feedItems = append(feedItems, feedItem)
					}
				}

				feed := Feed{
					Id:    feedId,
					Title: channel.Title,
					URL:   getFeedTitleUrl(channel),
					Items: feedItems,
				}
				feedChan <- feed
			}
		}
	}

	// create n worker goroutines
	for w := 1; w <= workers; w++ {
		go func(workChan chan Work, feedChan chan Feed, doneChan chan bool) {
			for work := range workChan {
				rssFeed := rss.New(5, true, channelHandler(work.Index, feedChan), nil)
				if err := rssFeed.Fetch(work.URL, nil); err != nil &&
					// ignore these encoding errors
					!strings.Contains(err.Error(), `encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil`) {
					n.log.WithFields(logrus.Fields{
						"error":    err,
						"feed_id":  work.Index,
						"feed_url": work.URL,
					}).Error("Could not fetch feed data")
				}
			}
			doneChan <- true
		}(workChan, feedChan, doneChan)
	}

	// distribute work to workers
	for idx, url := range n.Config.Feeds {
		workChan <- Work{idx, url}
	}
	close(workChan)

	// wait till all workers are done
	for w := 1; w <= workers; w++ {
		<-doneChan
	}
	close(feedChan)

	// read all feed results
	for feed := range feedChan {
		// cut off feed size at 30 lines
		for feed.Lines() > 30 {
			feed.Items = feed.Items[:len(feed.Items)-1]
		}

		feeds = append(feeds, feed)
	}

	// only replace feeds if we got enough of them back
	if len(feeds) >= 4 {
		feeds.sort()
		n.Feeds = feeds
	}
}

func (f Feed) Lines() int {
	// calculate number of lines a feed will occupy
	var lines float64
	for _, item := range f.Items {
		lines += math.Ceil(float64(utf8.RuneCountInString(item.Title)) / 82.0)
	}
	return len(f.Items) + ((int(lines) - len(f.Items)) / 2)
}

func getFeedTitleUrl(feed *rss.Channel) string {
	if len(feed.Links) > 0 {
		return feed.Links[0].Href
	}
	return "#"
}

func getFeedItemUrl(item *rss.Item) string {
	if len(item.Links) > 0 {
		return item.Links[0].Href
	}
	return "#"
}

func getRedditFeedItemUrl(item *rss.Item) string {
	matched := redditRx.FindStringSubmatch(item.Description)
	if len(matched) > 1 {
		return matched[1]
	}
	return "#"
}

func isFeedCommentUrl(comments string) bool {
	return strings.HasPrefix(comments, "http")
}
