package main

import "fmt"

func main() {
	config := loadConfig()
	fmt.Println(config.mattermostToken)
}
