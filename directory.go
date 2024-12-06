package main

import (
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)


var Directory map[string]string

func ReadDirectory() {
	csv, err := fs.ReadFile(os.DirFS("."), "directory.csv")
	if err != nil {
		log.Fatalln("failed to read directory file", err)
	}
	csvstr := strings.Split(string(csv), ", ")
	directory := make(map[string]string)

	for _, entry := range csvstr {
		data := strings.Split(entry, " <")
		name := data[0]
		email, f := strings.CutSuffix(data[1], "@[redacted]>")
		if !f {
			continue
		}
		directory[strings.ToLower(email)] = name
	}

	Directory = directory
}

func CoolMode() bool {
	return os.Getenv("COOL_MODE") == "yes"
}

func GetName(r *http.Request) (string, error) {
	if !CoolMode() {
		// revert to boring mode
		name := r.FormValue("author")

		if len(name) <= 1 {
			return "", errors.New("name too short")
		} else if len(name) > 70 {
			return "", errors.New("name too long")
		}

		return name, nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	hostname, err := net.LookupAddr(ip)

	if err != nil {
		return "", err
	}

	name := strings.Split(hostname[0], ".")[0]
	if name == "localhost" {
		return "Test User", nil
	}

	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		return "", errors.New("invalid hostname")
	}
	username := strings.ToLower(parts[1])
	full_name, found := Directory[username]

	if !found {
		return "", errors.New("username not in directory")
	}

	log.Println("authenticated", full_name)
	return full_name, nil
}
