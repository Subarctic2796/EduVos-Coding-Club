package urlinfo

import "sync"

type ShortUrl string
type FullUrl string

var (
	URL_CACHE = make(map[ShortUrl]FullUrl)
	LOCK      sync.Mutex
)
