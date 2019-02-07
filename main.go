package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Task struct {
	Created     time.Time
	Description string
	Done        bool
}

var header = []string{"Created", "Description", "Done"}

var layout = "2006-01-02 15:04:05"

// var dateFormat2 = timeToString(dateFormat)
//

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}

func main() {

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "xxxxxx")
	q := datastore.NewQuery("Task")
	it := client.Run(ctx, q)

	//書き込みファイル作成
	file, err := os.Create("sample.csv")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(header)

	for {
		var task Task
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next Task: %v", err)
		}
		time := timeToString(task.Created)
		desc := task.Description
		done := fmt.Sprint(task.Done)

		data := []string{time, desc, done}
		writer.Write(data)

	}
	writer.Flush() // ファイル出力
	log.Println("finish!")

}
