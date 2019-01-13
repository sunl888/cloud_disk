package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	dir = "/home/sunlong/download/cloud/"
)

func test(file *os.File) {
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, file); err != nil {
		log.Println(err.Error())
		return
	}
	md5sum := md5hash.Sum(nil)
	fmt.Println(md5sum, fmt.Sprintf("%x", md5sum))
	fileStat, err := file.Stat()
	if err != nil {
		log.Print(err)
		return
	}
	//&{name:1198055198.jpg size:18086281 mode:484 modTime:{wall:836041121 ext:63682999419 loc:0x2680380} sys:{Dev:2066 Ino:2502656 Nlink:1 Mode:33252 Uid:1000 Gid:1000 X__pad0:0 Rdev:0 Size:18086281 Blksize:4096 Blocks:35328 Atim:{Sec:1547402619 Nsec:788039608} Mtim:{Sec:1547402619 Nsec:836041121} Ctim:{Sec:1547402619 Nsec:836041121} X__unused:[0 0 0]}}
	fmt.Printf("%+v\n", fileStat)
}

//go1.11.linux-amd64.tar.gz
func main() {
	file, err := os.OpenFile(dir+"1198055198.jpg", os.O_RDWR, 0644)
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()
	test(file)
	fmt.Println("\n\n")
	file1, err := os.OpenFile(dir+"2.jpg", os.O_RDONLY, 0644)
	if err != nil {
		log.Print(err)
		return
	}
	defer file1.Close()
	test(file1)

	return
	//size := int64(16 << 20) // 16*2^20
	//blocks := make([][]byte, 0, 10)
	//block := make([]byte, 0, size)
	//offset := int64(0)
	//n := 0
	//for i := int64(1); i <= fileStat.Size()/size; i++ {
	//	_, err := file.ReadAt(block, offset)
	//	if err != nil {
	//		log.Print(err)
	//		return
	//	}
	//	blocks = append(blocks, block)
	//	offset += size + 1
	//	n++
	//}
	//writeFile, err := os.OpenFile(dir+"test.tar.gz", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0766)
	//for i := 0; i < n; i++ {
	//	_, err = writeFile.Write(blocks[i])
	//	if err != nil {
	//		log.Print(err)
	//		return
	//	}
	//}
	//defer writeFile.Close()
	//newFileStat, err := file.Stat()
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//fmt.Println(fileStat.Size(), fileStat.Name(), newFileStat.Name(), newFileStat.Size())
}
