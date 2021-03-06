package main

import (
	. "github.com/xitongsys/parquet-go/Marshal"
	. "github.com/xitongsys/parquet-go/ParquetHandler"
	. "github.com/xitongsys/parquet-go/ParquetType"
	"log"
	"os"
)

type Student struct {
	Name   UTF8
	Age    INT32
	Id     INT64
	Weight DOUBLE
	Sex    BOOLEAN
	School UTF8
}

type MyFile struct {
	FilePath string
	File     *os.File
}

func (self *MyFile) Create(name string) (ParquetFile, error) {
	file, err := os.Create(name)
	myFile := new(MyFile)
	myFile.File = file
	return myFile, err

}
func (self *MyFile) Open(name string) (ParquetFile, error) {
	var (
		err error
	)
	if name == "" {
		name = self.FilePath
	}

	myFile := new(MyFile)
	myFile.FilePath = name
	myFile.File, err = os.Open(name)
	return myFile, err
}
func (self *MyFile) Seek(offset int, pos int) (int64, error) {
	return self.File.Seek(int64(offset), pos)
}

func (self *MyFile) Read(b []byte) (n int, err error) {
	return self.File.Read(b)
}

func (self *MyFile) Write(b []byte) (n int, err error) {
	return self.File.Write(b)
}

func (self *MyFile) Close() {
	self.File.Close()
}

func main() {
	fname := os.Args[1]
	var f ParquetFile
	f = &MyFile{}
	f, _ = f.Open(fname)
	ph := NewParquetHandler()
	np := 20
	rowGroupNum := ph.ReadInit(f, int64(np))
	for i := 0; i < rowGroupNum; i++ {

		doneChan := make(chan int)
		tmap, num := ph.ReadOneRowGroup()
		stus := make([]Student, num)

		delta := (num + np - 1) / np

		for c := 0; c < np; c++ {
			bgn := c * delta
			end := bgn + delta
			if end > num {
				end = num
			}
			if bgn >= num {
				bgn, end = num, num
			}
			go func(b, e, index int) {
				tmp := stus[b:b]
				Unmarshal(tmap, b, e, &tmp, ph.SchemaHandler)
				doneChan <- 0
			}(bgn, end, c)
		}
		for c := 0; c < np; c++ {
			<-doneChan
		}
		log.Println(i)
		//log.Println(stus)
	}

	f.Close()
}
