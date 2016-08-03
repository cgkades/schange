package main

import "./user"
import "fmt"

func main() {
	group, _ := user.LookupGroupId("2133516900")
	fmt.Println(group)
}
