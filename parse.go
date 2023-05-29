package devicedetector

import (
	"io/ioutil"
	"log"

	"github.com/robicode/device-detector/util"
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

func parseDevices(fileList *CacheFileList) (map[string]CachedDeviceList, error) {
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

func parseBots(fileList *CacheFileList) ([]CachedBot, error) {
	var bots []CachedBot

	for _, filename := range fileList.filenames {
		fnBots, err := parseEmbeddedBotsFile(filename)
		if err != nil {
			return nil, err
		}
		bots = append(bots, fnBots...)
	}

	bots = util.ReverseArray(bots)
	return bots, nil
}

func parseEmbeddedBotsFile(filename string) ([]CachedBot, error) {
	var bots []CachedBot
	bytes, err := embeddedData.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &bots)
	if err != nil {
		return nil, err
	}
	return bots, nil
}

func parseEmbeddedOSFile(filename string) ([]CachedOS, error) {
	var oss []CachedOS
	bytes, err := embeddedData.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &oss)
	if err != nil {
		return nil, err
	}
	return oss, nil
}

func parseOSs(fileList *CacheFileList) ([]CachedOS, error) {
	var oss []CachedOS

	for _, filename := range fileList.filenames {
		fnOSs, err := parseEmbeddedOSFile(filename)
		if err != nil {
			return nil, err
		}

		oss = append(oss, fnOSs...)
	}
	return oss, nil
}

func parseHintsFile(filename string) (map[string]string, error) {
	output := make(map[string]string)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}
	return output, err
}

func parseHints(fileList *CacheFileList) (map[string]string, error) {
	list := make(map[string]string)

	for _, filename := range fileList.filenames {
		m, err := parseHintsFile(filename)
		if err != nil {
			return nil, err
		}
		for name, item := range m {
			if _, ok := list[name]; ok {
				log.Println("Item already in list:", name)
			} else {
				list[name] = item
			}
		}
	}
	return list, nil
}

func parseClientFile(filename string) ([]CachedClient, error) {
	var clients []CachedClient

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func parseClients(fileList *CacheFileList) ([]CachedClient, error) {
	var clients []CachedClient

	for _, filename := range fileList.filenames {
		list, err := parseClientFile(filename)
		if err != nil {
			return clients, err
		}
		clients = append(clients, list...)
	}

	return clients, nil
}
