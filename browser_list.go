package devicedetector

type hintBrowser struct {
	Name    string
	Version string
}

type BrowserList struct {
	list []hintBrowser
}

// newBrowserList creates a new *BrowserList.
func newBrowserList(initialList ...hintBrowser) *BrowserList {
	return &BrowserList{
		list: initialList,
	}
}

// Reject passes each member of the given BrowserList through a function and returns
// a new *BrowserList with only the entries where the func returned false.
func (b *BrowserList) Reject(reject func(browser hintBrowser) bool) *BrowserList {
	var newList []hintBrowser
	for _, browser := range b.list {
		if reject(browser) {
			continue
		}
		newList = append(newList, browser)
	}
	return &BrowserList{
		list: newList,
	}
}

// Last returns the last browser in the list.
func (b *BrowserList) Last() hintBrowser {
	if len(b.list) > 0 {
		return b.list[len(b.list)-1]
	}
	return hintBrowser{}
}
