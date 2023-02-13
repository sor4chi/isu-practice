package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	app "github.com/catatsuy/private-isu/webapp/golang"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/semaphore"
)

var (
	db *sqlx.DB
)

func main() {
	db = app.ConnectDb()
	ids := []int{}
	err := db.Select(&ids, "SELECT `id` FROM `posts`")
	if err != nil {
		log.Print(err)
	}
	downloadParallel(ids)
}

func downloadParallel(ids []int) {
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(5)
	for _, u := range ids {
		wg.Add(1)
		go downloadFromQuery(u, &wg, s)
	}
	wg.Wait()
}

func downloadFromQuery(_id int, wg *sync.WaitGroup, s *semaphore.Weighted) {
	fmt.Println("download start: ", _id)
	defer fmt.Println("download end: ", _id)
	defer wg.Done()
	if err := s.Acquire(context.Background(), 1); err != nil {
		return
	}
	defer s.Release(1)

	post := app.Post{}

	err := db.Get(&post, "SELECT * FROM `posts` WHERE `id` = ?", _id)
	if err != nil {
		log.Fatal(err)
	}
	suffix := strings.Split(post.Mime, "/")[1]
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.Mkdir("public", 0777)
	}
	fileName := fmt.Sprintf("public/%d.%s", _id, suffix)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(post.Imgdata)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
}
