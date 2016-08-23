package main

/*
TODO:
* Add ability to have multiple paths at the end schown <user> <file1> <file2>
* Add recursion
* Add whitelist based on config file
  * Read pathlist into an array
* Add bad char identification
*/

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	// Have to use the local user until go 1.7 comes out
	// This was copied directly from os/user in the 1.7 branch
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
	string_array := strings.Split(raw_string, ":")
	if len(string_array[1]) == 0 {
		return []string{string_array[0], "DEFAULTGROUP"}
	} else {
		return string_array
	}

}

// This takes a username and returns an int for use
// with os.chown()
// func Chown(name string, uid, gid int) error
func uidFromUsername(username string) int {
	user_obj := loadUser(username)
	uid, _ := strconv.Atoi(user_obj.Uid)
	return uid
}

// Same as uidFromUsername, but for groups. Allowing us to ensure
// the proper return for chownage
func gidFromGroupname(groupname string) int {
	group_struct, err := user.LookupGroup(groupname)
	if err != nil {
		fmt.Println("Error: group not found")
		os.Exit(1)
	}
	gid, _ := strconv.Atoi(group_struct.Gid)
	return gid
}

// Takes a username and returns their default group
func getUsersDefaultGroup(username string) int {
	user_obj := loadUser(username)
	default_gid, _ := strconv.Atoi(user_obj.Gid)
	return default_gid
}

func isSymlink(filename string) bool {
	fileinfo, _ := os.Lstat(filename)
	mode := fileinfo.Mode()

	if os.FileMode.IsRegular(mode) {
		return false
	} else {
		//fmt.Println("Symlink")
		return true
	}
}

func chown(filename string, uid int, gid int) {
	if isSymlink(filename) != true {
		ret_val := os.Chown(filename, uid, gid)
		if ret_val != nil {
			fmt.Println(ret_val)
		}
	}

}

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

	if len(allArgs) < 2 {
		usage()
		os.Exit(1)
	}
	//Get user/group string slice
	user_group := userGroupParser(allArgs[0])
	if len(user_group) > 2 {
		fmt.Println("Error: Invalid user/group given.")
		usage()
		os.Exit(1)
	}
	file_path := allArgs[1]
	user_name := user_group[0]
	group_name := user_group[1]

	// Get the GID to use
	var group_gid int

	if group_name == "DEFAULTGROUP" {
		group_gid = getUsersDefaultGroup(user_name)
		//fmt.Println("Default GID:", group_gid)
	} else if user_group[1] == "" {
		group_gid = -1
		//fmt.Println("GID Set to -1")
	} else {
		group_gid = gidFromGroupname(group_name)
		//fmt.Println("Given GID:", group_gid)
	}

	user_obj := loadUser(user_name)
	//fmt.Println("UID for", user_name, "=", user_obj.Uid)
	//fmt.Println("GID for", group_name, "=", group_gid)
	user_uid, _ := strconv.Atoi(user_obj.Uid)

	chown(file_path, user_uid, group_gid)
}
