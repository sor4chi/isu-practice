package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/semaphore"
)

var (
	db    *sqlx.DB
	total int
	cnt   int
)

type User struct {
	ID          int       `db:"id"`
	AccountName string    `db:"account_name"`
	Passhash    string    `db:"passhash"`
	Authority   int       `db:"authority"`
	DelFlg      int       `db:"del_flg"`
	CreatedAt   time.Time `db:"created_at"`
}

type Post struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Body         string    `db:"body"`
	Imgdata      []byte    `db:"imgdata"`
	Mime         string    `db:"mime"`
	CreatedAt    time.Time `db:"created_at"`
	CommentCount int
	Comments     []Comment
	User         User
	CSRFToken    string
}

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	UserID    int       `db:"user_id"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
	User      User
}

func main() {
	host := os.Getenv("ISUCONP_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ISUCONP_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCONP_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCONP_DB_PASSWORD")
	dbname := os.Getenv("ISUCONP_DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	defer db.Close()

	ids := []int{}
	err = db.Select(&ids, "SELECT `id` FROM `posts`")
	if err != nil {
		log.Print(err)
	}
	total = len(ids)
	downloadParallel(ids)
}

func downloadParallel(ids []int) {
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(20)
	for _, u := range ids {
		wg.Add(1)
		go downloadFromQuery(u, &wg, s)
	}
	wg.Wait()
}

func downloadFromQuery(_id int, wg *sync.WaitGroup, s *semaphore.Weighted) {
	defer wg.Done()
	if err := s.Acquire(context.Background(), 1); err != nil {
		return
	}
	defer s.Release(1)

	post := Post{}

	err := db.Get(&post, "SELECT * FROM `posts` WHERE `id` = ?", _id)
	if err != nil {
		log.Fatal(err)
	}
	suffix := strings.Split(post.Mime, "/")[1]
	if suffix == "jpeg" {
		suffix = "jpg"
	}
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.Mkdir("public", 0777)
	}
	fileName := fmt.Sprintf("../public/%d.%s", _id, suffix)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(post.Imgdata)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	cnt++
	fmt.Println("downloaded", cnt, "/", total)
}
