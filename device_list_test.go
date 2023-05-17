package devicedetector

import (
	"testing"
)

func TestDeviceList_Append(t *testing.T) {
	list := NewDeviceList()

	entry := CachedDevice{
		Regex: `(?:BUZZ [123]|CLEVER 1|URBAN [13](?: Pro)?)(?:[);/ ]|$)`,
		Type:  "smartphone",
	}

	nlist := list.Append("Ace", entry)

	if len(nlist.list) != 1 {
		t.Errorf("expected nlist to contain %d items but had %d instead.\nnList: %v\nlist:%v", len(list.list)+1, len(nlist.list),
			nlist, list)
		t.Fail()
		return
	}

	// Make sure the original is unmodified

	if len(list.list) > 0 {
		t.Error("expected original list to be empty but had", len(list.list), "items")
		t.Fail()
		return
	}

	// Don't allow appending the same item twice

	nnlist := nlist.Append("Ace", entry)
	if len(nnlist.list) != len(nlist.list) {
		t.Errorf("should not be able to add duplicate named items")
		t.Fail()
		return
	}

	// Don't allow appending two items w/the same Regex, even w/diff names

	nnnlist := nlist.Append("Another Ace", entry)
	if len(nnnlist.list) != len(nlist.list) {
		t.Error("should not be able to add two entries with the same regex")
		t.Fail()
		return
	}
}

func TestDeviceList_Delete(t *testing.T) {
	list := NewDeviceList()

	entry := CachedDevice{
		Regex: `(?:BUZZ [123]|CLEVER 1|URBAN [13](?: Pro)?)(?:[);/ ]|$)`,
		Type:  "smartphone",
	}

	list = list.Append("Ace", entry)

	newList := list.Delete("Ace")

	if len(newList.list) != 0 {
		t.Error("expected newList to be empty but length was", len(newList.list))
		t.Fail()
		return
	}

	// Make sure the original is unmodified

	if len(list.list) != 1 {
		t.Error("expected original not to be modified but new length is", len(list.list))
		t.Fail()
		return
	}
}
