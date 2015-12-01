package main

import (
	"fmt"
	"github.com/lukamicoder/ini-parser"
)

func main() {
	path := "./config.ini"

	var config iniparser.Config

	err := config.LoadFile(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Sections:\n")
	sections := config.GetSections()
	for _, section := range sections {
		fmt.Printf(" - %v\n", section.Name)
	}
	fmt.Printf("\n")

	var port int
	port, err = config.GetInt("database", "port")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Port: %v\n", port)

	var dbfile string
	dbfile, err = config.GetString("database", "dbfile")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("DB file: %v\n\n", dbfile)

	var section *iniparser.Section
	section, err = config.GetSection("users")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Section '%v':\n", section.Name)
	for key, value := range section.Keys {
		fmt.Printf("%v: %v\n", key, value)
	}
}
