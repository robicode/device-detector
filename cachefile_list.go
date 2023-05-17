package devicedetector

import (
	"log"
	"path/filepath"
	"strings"
)

// CacheFileList is a helper struct used to pass a list of caches to search into
// the cache's findRegex function.
type CacheFileList struct {
	originalList []string
	filenames    []string
}

// NewCacheFileList creates a new *CacheFileList, given the base list as an optional parameter.
// If no base is provided, uses the entires cacheFilenames array as a starting point.
func NewCacheFileList(baseList ...string) *CacheFileList {
	var files []string
	var originalFiles []string

	if baseList == nil || len(baseList) == 0 {
		for _, file := range cacheFiles {
			if strings.HasPrefix(file, "regexes/") {
				files = append(files, file)
			} else {
				files = append(files, filepath.Join("regexes", file))
			}
		}
		originalFiles = files
	} else {
		for _, file := range baseList {
			if strings.HasPrefix(file, "regexes/") {
				files = append(files, file)
				originalFiles = append(originalFiles, file)
			} else {
				files = append(files, filepath.Join("regexes", file))
				originalFiles = append(originalFiles, filepath.Join("regexes", file))
			}
		}
	}

	return &CacheFileList{
		filenames:    files,
		originalList: originalFiles,
	}
}

// Exclude returns a new *CacheFileList with the specified files removed.
func (c *CacheFileList) Exclude(excludes ...string) *CacheFileList {
	if len(excludes) == 0 {
		log.Println("Nothing to exclude.")
		return c
	}

	for _, exclude := range excludes {
		if !inStrArray(exclude, c.originalList) {
			return c
		}
	}

	newList := []string{}

	for _, file := range c.originalList {
		if !inStrArray(file, excludes) {
			newList = append(newList, file)
		}
	}

	newCacheFiles := NewCacheFileList()
	newCacheFiles.filenames = newList
	return newCacheFiles
}

// Includes returns true if the given filename is in the *CacheFileList.
func (c *CacheFileList) Includes(filename string) bool {
	for _, file := range c.filenames {
		if file == filename {
			return true
		}
	}
	return false
}

// Get returns the filename at the given index.
func (c *CacheFileList) Get(index int) string {
	if index > len(c.filenames)-1 {
		return ""
	}
	return c.filenames[index]
}

// Exclusive returns a copy of the *CacheFileList containing only the given filenames.
func (c *CacheFileList) Exclusive(filenames ...string) *CacheFileList {
	if len(filenames) == 0 || len(filenames) == len(deviceFilenames) {
		return c
	}

	newArr := []string{}

	for _, file := range filenames {
		if !inStrArray(file, c.originalList) {
			return c
		}
		newArr = append(newArr, file)
	}
	newFileList := NewCacheFileList()
	newFileList.filenames = newArr

	return newFileList
}
