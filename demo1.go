package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	folderMap := make(map[int64]string, 5)
	folderMap[1] = "根目录"
	folderMap[2] = "1.1"
	folderMap[3] = "1.1.1"
	folderMap[4] = "1.1.1.1"
	folderMap[5] = "1.1.2"
	folderMap[6] = "1.1.2.1"

	key := "1-2-3-"
	path := mergePath(folderMap, 1, key, 4)
	fmt.Println(path)

	tmp := "bytes=0-"
	prefixIndex := strings.Index(tmp, "-")
	fmt.Println(prefixIndex)
	start, err := strconv.ParseInt(tmp[6:prefixIndex], 10, 64)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println(tmp[prefixIndex+1:] == "")
	end, err := strconv.ParseInt(tmp[prefixIndex+1:], 10, 64)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println(start, end)
}

func mergePath(folderMap map[int64]string, currentId int64, key string, withId int64) (path string) {
	if currentId == withId {
		return ""
	}
	key2Arr := strings.Split(key, "-")
	for _, v := range key2Arr {
		id2Int64, _ := strconv.ParseInt(v, 10, 64)
		if id2Int64 > currentId {
			path += folderMap[id2Int64] + "/"
		}
	}
	path += folderMap[withId]
	return path
}
