package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ConfigName := ".tp.toml"
	cfgFile := "/home/bmorriso/.config"
	ConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Print(err)
	}

	HomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(ConfigDir)
	fmt.Println(HomeDir)

	if ConfigDir == cfgFile {
		fmt.Println("ConfigDir == cfgFile")
	} else {
		fmt.Println("ConfigDir != cfgFile")
	}

	e := time.Now().Unix()

	n := time.Now().Local().Format("2006-01-02")

	fmt.Println(n)

	epoch := strconv.FormatInt(e, 10)

	s := strings.Split(ConfigName, ".")
	fmt.Println(s[1])

	// cfgFile, _ := os.UserHomeDir()
	bkupConfigFile := cfgFile + "/" + ConfigName + "-" + epoch

	fmt.Println(bkupConfigFile)

}
