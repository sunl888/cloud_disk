package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	folderMap := make(map[int64]string, 5)
	folderMap[1] = "根目录"
	folderMap[2] = "目录2"
	folderMap[3] = "目录3"
	folderMap[4] = "目录4"
	folderMap[5] = "目录5"

	key := "1-2-3-4-5-"
	currentKey := int64(2)
	path := generatePath(folderMap, currentKey, key)
	fmt.Println(path)
}
func generatePath(folderMap map[int64]string, currentId int64, key string) (path string) {
	flag := false
	key2Arr := strings.Split(key, "-")
	for _, v := range key2Arr {
		id2Int64, _ := strconv.ParseInt(v, 10, 64)
		if id2Int64 == currentId {
			flag = true
		}
		if flag {
			path += folderMap[id2Int64] + "/"
		}
	}
	firstIndex := strings.Index(path, "/")
	lastIndex := strings.LastIndex(path, "/")
	return path[firstIndex+1 : lastIndex-1]
}
