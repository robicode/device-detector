package devicedetector

import (
	"path/filepath"
	"testing"
)

func Test_NewCacheFileList(t *testing.T) {
	cf := NewCacheFileList()

	if len(cf.filenames) != len(cacheFiles) {
		t.Errorf("expected length to be %d but was %d", len(cacheFiles), len(cf.filenames))
		t.Fail()
		return
	}
}

func TestCacheFileList_Get(t *testing.T) {
	cf := NewCacheFileList()

	for i, file := range cacheFiles {
		if cf.Get(i) != filepath.Join("regexes", file) {
			t.Errorf("expected index %d to be '%s' but was '%s'", i, filepath.Join("regexes", file), cf.Get(i))
			t.Fail()
		}
	}
}

func TestCacheFileList_Exclude(t *testing.T) {
	cf := NewCacheFileList()

	newFiles := cf.Exclude("regexes/device/televisions.yml")

	if len(newFiles.filenames) != len(cacheFiles)-1 {
		t.Errorf("expected count to be %d but was %d\nCount of original files: %d", len(cacheFiles)-1, len(newFiles.filenames),
			len(cf.filenames))
		t.Fail()
		return
	}

	cf = NewCacheFileList(deviceFilenames...)

	if len(cf.filenames) != len(deviceFilenames) {
		t.Errorf("expected count to be %d but was %d", len(deviceFilenames), len(cf.filenames))
		t.Fail()
		return
	}

	if len(cf.originalList) != len(deviceFilenames) {
		t.Errorf("expected original count to be %d but was %d", len(deviceFilenames), len(cf.originalList))
		t.Fail()
		return
	}
}

func TestCacheFileList_Exclusive(t *testing.T) {
	cf := NewCacheFileList()

	newCacheFiles := cf.Exclusive("regexes/device/shell_tv.yml")

	if len(newCacheFiles.filenames) != 1 {
		t.Error("expected list to only have 1 item but had", len(newCacheFiles.filenames))
		t.Fail()
		return
	}
}
