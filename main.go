package main

import (
	"encoding/json"
	"flag"
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

// uploadDir is where uploaded files are stored.
const uploadDir = "./static/servfile"

// maxUploadSize limits a single upload (default 1 GiB, override with -max-size).
var maxUploadSize int64 = 1 << 30

// safeName strips any directory components from a user-supplied file name
// so requests cannot read, write, or delete files outside uploadDir.
func safeName(name string) (string, bool) {
	name = filepath.Base(filepath.Clean("/" + name))
	if name == "/" || name == "." || name == ".." || name == "" || name[0] == '.' {
		return "", false
	}
	return name, true
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("Invalid request method")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		log.Println("Error retrieving the file:", err)
		return
	}
	defer file.Close()

	name, ok := safeName(handler.Filename)
	if !ok {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	f, err := os.Create(filepath.Join(uploadDir, name))
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

	log.Println("File uploaded successfully:", name)
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the main page
}

// listFiles returns the regular, non-hidden files in uploadDir.
func listFiles() ([]UploadedFile, error) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return nil, err
	}

	fileInfos := []UploadedFile{}
	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, UploadedFile{Name: file.Name(), Date: info.ModTime()})
	}
	return fileInfos, nil
}

func getLastUploadedFiles(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := listFiles()
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		log.Println("Unable to read directory:", err)
		return
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Date.After(fileInfos[j].Date)
	})

	if len(fileInfos) > 5 {
		fileInfos = fileInfos[:5]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfos)
}

func getServerFiles(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := listFiles()
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		log.Println("Unable to read directory:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfos)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fileName, ok := safeName(r.URL.Query().Get("file"))
	if !ok {
		http.Error(w, "Valid file name is required", http.StatusBadRequest)
		return
	}

	err := os.Remove(filepath.Join(uploadDir, fileName))
	if err != nil {
		http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		log.Println("Unable to delete file:", err)
		return
	}

	log.Println("File deleted successfully:", fileName)
	w.WriteHeader(http.StatusOK)
}
func downloadFile(w http.ResponseWriter, r *http.Request) {
	// Prevent path traversal
	fileName, ok := safeName(r.URL.Query().Get("file"))
	if !ok {
		http.Error(w, "Valid file name is required", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(uploadDir, fileName)

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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Write the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
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
	host := flag.String("host", "", "address to listen on (default: this machine's local network IP)")
	port := flag.Int("port", 8080, "port to listen on")
	maxSizeMB := flag.Int64("max-size", 1024, "maximum upload size in MB")
	flag.Parse()
	maxUploadSize = *maxSizeMB << 20

	// Create a log file
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	// Set log output to the log file
	log.SetOutput(logFile)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Println("Error creating upload directory:", err)
		return
	}

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/last-uploaded", getLastUploadedFiles)
	http.HandleFunc("/server-files", getServerFiles)
	http.HandleFunc("/delete", deleteFile)
	http.HandleFunc("/download", downloadFile)

	if *host == "" {
		*host = GetLocalIP()
	}
	addr := net.JoinHostPort(*host, fmt.Sprint(*port))

	log.Printf("Starting server on %s", addr)
	fmt.Printf("http://%s\n", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Error starting server:", err)
		fmt.Println("Error starting server:", err)
	}
}
