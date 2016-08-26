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

// GO has a good user struct that we can get from the os module
// We'll use this to grab it and return the struct
func loadUser(username string) *user.User {
	user_obj, err := user.Lookup(username)
	if err != nil {
		fmt.Println("Error: user not found")
		os.Exit(1)
	}
	return user_obj
}

// Need to change the usage output so it looks right
func usage() {
	fmt.Fprintf(os.Stderr, "usage: schown <options> owner[:group] <path>\n")
	flag.PrintDefaults()
}

// Chown can take in a string like user:group or user: or user, we figure
// all of that out here and return a string array [0] is user [1] is group
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

// Returns UID, GID
func getFinalUIDandGID(username, groupname string) (int, int) {
	var gid int
	if groupname == "DEFAULTGROUP" {
		gid = getUsersDefaultGroup(user_name)
	} else if user_group[1] == "" { // If no ':' was passed then we will assume no group change
		// GID -1 (though not documented in the module) makes no change to group
		gid = -1
	} else {
		gid = gidFromGroupname(group_name)
	}
	user_obj := loadUser(user_name)
	uid, _ := strconv.Atoi(user_obj.Uid)

	return uid, gid

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

	uid, gid := getFinalUIDandGID(user_name, group_name)

	chown(file_path, uid, gid)
}
