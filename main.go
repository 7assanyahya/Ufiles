package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type UploadedFile struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

var uploadedFiles []UploadedFile

// func createFolders() {
// 	folders := []string{"folder1", "folder2", "folder3"}
// 	for _, folder := range folders {
// 		path := "./static/servfile/" + folder
// 		if _, err := os.Stat(path); os.IsNotExist(err) {
// 			err := os.Mkdir(path, os.ModePerm)
// 			if err != nil {
// 				log.Println("Error creating folder:", err)
// 			} else {
// 				log.Println("Folder created:", path)
// 			}
// 		}
// 	}
// }

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("Invalid request method")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		log.Println("Error retrieving the file:", err)
		return
	}
	defer file.Close()

	f, err := os.Create("./static/servfile/" + handler.Filename)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		log.Println("Error creating the file:", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		log.Println("Error saving the file:", err)
		return
	}

	uploadedFiles = append(uploadedFiles, UploadedFile{Name: handler.Filename, Date: time.Now()})
	log.Println("File uploaded successfully:", handler.Filename)
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the main page
}

func getLastUploadedFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./static/servfile")
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		log.Println("Unable to read directory:", err)
		return
	}

	var fileInfos []UploadedFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, UploadedFile{Name: file.Name(), Date: info.ModTime()})
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Date.After(fileInfos[j].Date)
	})

	if len(fileInfos) > 5 {
		fileInfos = fileInfos[:5]
	}

	json.NewEncoder(w).Encode(fileInfos)
}

func getServerFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./static/servfile")
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		log.Println("Unable to read directory:", err)
		return
	}

	var fileInfos []UploadedFile
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, UploadedFile{Name: file.Name(), Date: info.ModTime()})
	}

	json.NewEncoder(w).Encode(fileInfos)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	err := os.Remove("./static/servfile/" + fileName)
	if err != nil {
		http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		log.Println("Unable to delete file:", err)
		return
	}

	log.Println("File deleted successfully:", fileName)
	w.WriteHeader(http.StatusOK)
}
func downloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Prevent path traversal
	fileName = filepath.Clean(fileName)
	filePath := "./static/servfile/" + fileName

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		log.Println("Unable to open file:", err)
		return
	}
	defer file.Close()

	// Set the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Write the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Unable to write file to response", http.StatusInternalServerError)
		log.Println("Unable to write file to response:", err)
		return
	}

	log.Println("File downloaded successfully:", fileName)

}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	// Create a log file
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	// Set log output to the log file
	log.SetOutput(logFile)

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/last-uploaded", getLastUploadedFiles)
	http.HandleFunc("/server-files", getServerFiles)
	http.HandleFunc("/delete", deleteFile)
	http.HandleFunc("/download", downloadFile)

	log.Printf("Starting server on %s:8080", GetLocalIP())
	fmt.Printf("http://%s:8080\n", GetLocalIP())

	err = http.ListenAndServe(string(GetLocalIP())+":8080", nil)

	if err != nil {
		log.Println("Error starting server:", err)
		fmt.Println("Error starting server:", err)
	}
}
