//中間ファイルを作成せずに圧縮ファイルを生成する手法
package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
)

func before() {
	f, err := os.Create("sample.txt")
	if err != nil {
		panic(err)
	}
	f.Write([]byte("sample"))
	f.Close()
}

func main() {
	before()

	//圧縮対象ファイル取得
	files := find()

	//圧縮ファイルデータ生成
	//この時点では、メモリ上に存在する
	b := compress(files)

	//例としてこのタイミングで、ファイル化
	//この時点では、あまり意味は無い
	//ローカルに残さず、s3にuploadするなどが可能となる
	//file ioがかからないので、高速化
	if err := save(b); err != nil {
		panic(err)
	}
}

func find() []string {
	return []string{"sample.txt"}
}

func save(b *bytes.Buffer) error {
	zf, err := os.Create("sample.zip")
	if err != nil {
		return err
	}
	zf.Write(b.Bytes())
	zf.Close()
	return nil
}

func compress(files []string) *bytes.Buffer {
	b := new(bytes.Buffer)
	w := zip.NewWriter(b)

	for _, file := range files {
		info, _ := os.Stat(file)

		hdr, _ := zip.FileInfoHeader(info)
		hdr.Name = "files/" + file
		f, err := w.CreateHeader(hdr)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		f.Write(body)
	}

	w.Close()

	return b
}
