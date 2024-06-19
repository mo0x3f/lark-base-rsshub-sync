package rsshub

type Feed struct {
	Title string
	Items []*Item
}

// IsOverrideMode 如果 RSS 都没有发布时间，则为全量覆盖模式
func (feed *Feed) IsOverrideMode() bool {
	for _, item := range feed.Items {
		if item.PubDate != 0 {
			return false
		}
	}
	return true
}

type Item struct {
	Guid         string
	Title        string
	Description  string
	Link         string
	Authors      []string
	PubDate      int64
	CategoryList []string
}
