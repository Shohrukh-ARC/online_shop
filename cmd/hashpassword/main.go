package main

import (
	"fmt"
	"log"

	"example.com/online-shop/admin"
)

func main() {
	hash, err := admin.HashPassword("Password!!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hash)
}
