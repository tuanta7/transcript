package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := NewRecorder()
	err := r.Record()
	if err != nil {
		fmt.Println(err)
		return
	}
}
