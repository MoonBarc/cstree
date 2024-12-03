package main

import (
	"io"
	"net/http"
	"strings"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	author, err := GetName(r)
	if err != nil {
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB File upload limit
	if err != nil {
		http.Error(w, "bad request (maybe too large)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("pyfile")
	if err != nil {
		http.Error(w, "bad request (bad file)", http.StatusBadRequest)
		return
	}

	_, ends_in_py := strings.CutSuffix(header.Filename, ".py")
	if !ends_in_py {
		http.Error(w, "must be a python file", http.StatusBadRequest)
		return
	}

	file_str, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "bad file", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	id, err := TrDB.AddProgram(string(file_str), name, author)
	if err != nil {
		http.Error(w, "failed to add program to the DB", http.StatusInternalServerError)
		return
	}

	// add it to the start of the queue
	QueueMutex.Lock()
	defer QueueMutex.Unlock()
	Queue = append([]int64{id}, Queue...)

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}
