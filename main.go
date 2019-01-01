package main

import (
	"fmt"
	"strings"
)

func main() {
	key := "1-2-3-4-5-"
	startKey := "3"
	newStartKey := "1-"

	idMap := make(map[string]string, 4)
	idMap["3"] = "6"
	idMap["4"] = "7"
	idMap["5"] = "8"

	//fmt.Println(updateKey(newStartKey, key, startKey))

	fmt.Println(replaceKey(idMap, newStartKey, key, startKey))

	sql := "INSERT INTO `folder_files` VALUES (1,2),(2,2),(3,2),"

	sql = strings.TrimRight(sql, ",")
	fmt.Println(sql)
}

func updateKey(parentKey, key, startId string) string {
	keys := strings.Split(key, "-")
	for index, key := range keys {
		if key == startId {
			return parentKey + strings.Join(keys[index:], "-")
		}
	}
	return ""
}

func replaceKey(idMap map[string]string, parentKey, key, startId string) string {
	newKey := updateKey(parentKey, key, startId)
	if newKey == "" {
		return ""
	}
	keys := strings.Split(key, "-")
	for index, key := range keys {
		fmt.Println(idMap[key])
		if newId, ok := idMap[key]; ok {
			keys[index] = newId
		}
	}
	return strings.Join(keys, "-")
}
