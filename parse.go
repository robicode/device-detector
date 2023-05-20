package devicedetector

import (
	"gopkg.in/yaml.v3"
)

func parseEmbeddedDeviceFile(filename string) (*CachedDeviceList, error) {
	devices := make(map[string]CachedDevice)
	bytes, err := embeddedData.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &devices)
	if err != nil {
		return nil, err
	}

	deviceList := NewDeviceList()
	deviceList.list = devices

	return deviceList, nil
}

func parse(fileList *CacheFileList) (map[string]CachedDeviceList, error) {
	deviceTree := make(map[string]CachedDeviceList)

	for _, filename := range fileList.filenames {
		list, err := parseEmbeddedDeviceFile(filename)
		if err != nil {
			return deviceTree, err
		}

		deviceTree[filename] = *list
	}

	return deviceTree, nil
}
