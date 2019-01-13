package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	dir = "/home/sunlong/图片/"
)

func test(file *os.File) {
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, file); err != nil {
		log.Println(err.Error())
		return
	}
	md5sum := md5hash.Sum(nil)
	fmt.Println(md5sum, fmt.Sprintf("%x", md5sum))
}

//go1.11.linux-amd64.tar.gz
func main() {
	file, err := os.Open(dir + "1527580104.jpg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	var b []byte
	a, c := file.Read(b)
	fmt.Println(a, b, c)
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	size := int64(16 << 20) // 16*2^20
	blocks := make([][]byte, 0, 10)
	offset := int64(0)
	n := 0
	fmt.Println(size)
	for i := int64(0); i <= fileStat.Size()/size; i++ {
		block := make([]byte, 0, size)
		_, err := file.ReadAt(block, 0)
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Println(block)
		blocks = append(blocks, block)
		offset += size + 1
		n++
	}
	writeFile, err := os.OpenFile(dir+"test.jpg", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	for i := 0; i < n; i++ {
		_, err = writeFile.WriteAt(blocks[i], 0)
		if err != nil {
			log.Print(err)
			return
		}
	}
	defer writeFile.Close()
	newFileStat, err := writeFile.Stat()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(fileStat.Size(), fileStat.Name(), newFileStat.Name(), newFileStat.Size())
}
