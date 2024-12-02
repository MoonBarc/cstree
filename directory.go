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



func ReadDirectory() map[string]string {
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

	return directory
}

func GetName(directory map[string]string, r *http.Request) (string, error) {
	hostname, err := net.LookupAddr(strings.Split(r.RemoteAddr, ":")[0])

	if err != nil {
		return "", err
	}

	name := strings.Split(hostname[0], ".")[0]
	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		return "", errors.New("invalid hostname")
	}
	username := strings.ToLower(parts[1])
	full_name, found := directory[username]

	if !found {
		return "", errors.New("username not in directory")
	}

	log.Println("authenticated", full_name)
	return full_name, nil
}
