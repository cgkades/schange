package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"strconv"

	"./user"
)

//GO has a good user struct that we can get from the os module
//We'll use this to grab it and return the struct
func loadUser(username string) *user.User {
	user_obj, err := user.Lookup(username)
	if err != nil {
		fmt.Println("Error: user not found")
		os.Exit(1)
	}
	return user_obj
}

//Need to change the usage output so it looks right
func usage() {
	fmt.Fprintf(os.Stderr, "usage: schown <options> owner[:group] <path>\n")
	flag.PrintDefaults()
}

func userGroupParser(raw_string string) []string {
	colon_location := strings.Index(raw_string, ":")
	if colon_location == -1 {
		return []string{raw_string, ""}
	}
	return strings.Split(raw_string, ":")
}

// This takes a username and returns an int for use
// with os.chown()
// func Chown(name string, uid, gid int) error
func uidFromUsername(username string) int {
	user_obj := loadUser(username)
	return strconv.Atoi(user_obj.Uid)
}

//func getUsersDefaultGroup(username group)

func main() {
	//Setup Flags
	flag.Usage = usage
	recursive := flag.Bool("R", false, "Recursive")
	flag.Parse()
	allArgs := flag.Args()

	//Check for recursive flag
	if *recursive {
		fmt.Println("Recursive was set")
	}

	fmt.Println("args:", allArgs)
	//Get user/group string slice
	user_group := userGroupParser(allArgs[0])
	if len(user_group) > 2 {
		fmt.Println("Error: Invalid user/group given.")
		usage()
		os.Exit(1)
	}
	fmt.Println("User/Group:", user_group)

	if user_group[1] == "" {
		fmt.Println("No group given.")
	} else {
		group, _ := user.LookupGroup(user_group[1])
		fmt.Println("Group:", group.Gid)
	}

	user_obj := loadUser("byoakum")
	fmt.Println("UID:", user_obj.Uid)
	fmt.Println("GID:", user_obj.Gid)
}
