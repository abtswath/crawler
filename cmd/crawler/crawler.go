package main

import (
	"crawler/pkg/task"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	t, err := task.New("http://192.168.99.241/", task.Option{
		Timeout: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer t.Close()
	t.Run()
	fmt.Println(t.Collection.GetAll())
	fs, _ := os.OpenFile("./result.json", os.O_CREATE|os.O_RDWR, os.ModePerm)
	t.Collection.ToJSON(fs)
	fs.Close()
}
