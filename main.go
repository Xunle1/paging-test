package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need username and password of MySQL!")
		os.Exit(0)
	}
	username := os.Args[1]
	password := os.Args[2]
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", username, password)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot open MySQL connection!")
		os.Exit(0)
	}

	// Insert(db)
	page := 345543
	lastId := page
	pageSize := 1000

	fmt.Println("-------------Paging by offset-------------")
	t1 := time.Now()
	PagingByOffset(page, pageSize, db)
	t2 := time.Since(t1)
	fmt.Printf("Paging by offset took: \n%s", t2.String())

	fmt.Println()
	fmt.Println("-------------Paging by lastId-------------")
	t3 := time.Now()
	PagingById(lastId, pageSize, db)
	t4 := time.Since(t3)
	fmt.Printf("Paging by lastId took: \n%s", t4.String())
	fmt.Println()
}

type User struct {
	ID       uint `gorm:"column:id;primaryKey"`
	Username string
	Email    string
}

func Insert(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Println("Cannot creat table `user`.")
		return
	}
	fmt.Println("Create table `user` successfully!")
	db.Where("1 = 1").Delete(&User{})
	db.Exec("truncate table users")
	users := make([]User, 0)

	fmt.Println("Start insertion...")
	for i := 1; i <= 1_000_000; i++ {
		user := User{
			Username: "user-" + strconv.Itoa(i),
			Email:    fmt.Sprintf("foo-%s@bar.com", strconv.Itoa(i)),
		}
		users = append(users, user)
	}
	for i := 0; i < len(users); i += 1000 {
		end := i + 1000
		if err := db.Model(&User{}).Save(users[i:end]).Error; err != nil {
			fmt.Println("Failed to insert users:", err.Error())
			return
		}
	}

	fmt.Println("Inserted 1000000 users!")
}

func PagingByOffset(page, pageSize int, db *gorm.DB) {
	var users []User
	db.Model(&User{}).Scopes(func(d *gorm.DB) *gorm.DB {
		return d.Limit(pageSize).Offset(page)
	}).Find(&users)
}

func PagingById(page, pageSize int, db *gorm.DB) {
	var users []User
	db.Model(&User{}).Scopes(func(d *gorm.DB) *gorm.DB {
		return d.Limit(pageSize).Where("id > ?", page)
	}).Find(&users)
}
