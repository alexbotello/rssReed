package main

// byResult type is used for sorting Result items from newest to oldest
type byResult []*Result

func (t byResult) Len() int {
	return len(t)
}

func (t byResult) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byResult) Less(i, j int) bool {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	a, b := t[i].Item, t[j].Item
	if a.PublishedParsed == nil || b.PublishedParsed == nil {
		return false
	}
	return a.PublishedParsed.Before(*b.PublishedParsed)
}

// byItem type is used for sorting Item items from newest to oldest
type byItem []Item

func (t byItem) Len() int {
	return len(t)
}

func (t byItem) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byItem) Less(i, j int) bool {
	return t[i].Date.Before(*t[j].Date)
}
