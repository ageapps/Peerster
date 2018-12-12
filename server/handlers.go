package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// Health message
func Health(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	status := "HEALTHY"
	send(&w, &status)
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Peerster!")
}

// GetMessages func
func GetMessages(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for messages"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// GetPrivateMessages func
func GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for private messages"))
		return
	}
	send(&w, getGossiperPrivateMessages(name))
}

// GetRoutes func
func GetRoutes(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for routes"))
		return
	}
	send(&w, getGossiperRoutes(name))
}

// GetFiles func
func GetFiles(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for files"))
		return
	}
	send(&w, getGossiperFiles(name))
}

// GetNodes func
func GetNodes(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for nodes"))
		return
	}
	send(&w, getGossiperPeers(name))
}

// PostMessage func
func PostMessage(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for new message"))
		return
	}
	msg, ok := params["msg"].(string)
	if !ok || !sendMessage(name, msg) {
		sendError(&w, errors.New("Error while sending new message"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// PostPrivateMessage func
func PostPrivateMessage(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested for private message"))
		return
	}
	dest, ok := params["destination"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no destination requested"))
		return
	}
	msg, ok := params["msg"].(string)
	if !ok || !sendPrivateMessage(name, dest, msg) {
		sendError(&w, errors.New("Error while sending new message"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// PostSearch func
func PostSearch(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	search, ok := params["search"].(string)
	if !ok || !sendSearchMessage(name, search) {
		sendError(&w, errors.New("Error while searching"))
		return
	}
	sendOk(&w)
}

// PostRequest func
func PostRequest(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	dest, ok := params["destination"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no destination requested"))
		return
	}
	hash, ok := params["hash"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no hash requested"))
		return
	}
	file, ok := params["file"].(string)
	if !ok || !sendFileRequest(name, dest, file, hash) {
		sendError(&w, errors.New("Error while sending new request"))
		return
	}
	send(&w, getGossiperFiles(name))
}

// PostNode func
func PostNode(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	peer, ok := params["node"].(string)
	if !ok || !addPeer(name, peer) {
		sendError(&w, errors.New("Error while adding new peer"))
		return
	}
	send(&w, getGossiperPeers(name))
}

// GetID func
func GetID(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}

	send(&w, getStatusResponse(name))
}

// Delete gossiper
func Delete(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	deleteGossiper(name)
	sendOk(&w)
}

// Upload file
func Upload(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	if path := downloadFile(w, r); path != "" {
		send(&w, indexFileInGossiper(name, path))
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) string {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendError(&w, fmt.Errorf("file to big, %v", err))
		return ""
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		sendError(&w, fmt.Errorf("invalid file, %v", err))
		return ""
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		sendError(&w, fmt.Errorf("invalid file, %v", err))
		return ""
	}
	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" &&
		filetype != "application/pdf" {
		sendError(&w, fmt.Errorf("invalid file type, %v", err))
		return ""
	}
	fileName := strings.Split(handler.Filename, ".")[0]
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		sendError(&w, fmt.Errorf("can't read file type, %v", err))
		return ""
	}
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)
	logger.Logf("Saving file in path %v", newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		sendError(&w, fmt.Errorf("can't write file type, %v", err))
		return ""
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		sendError(&w, fmt.Errorf("can't write file type, %v", err))
		return ""
	}
	return fileName + fileEndings[0]
}

// Start gossiper
func Start(w http.ResponseWriter, r *http.Request) {

	params := *readBody(&w, r)

	name := params["name"].(string)
	var peers = utils.PeerAddresses{}
	var gossipAddr = utils.PeerAddress{}

	if address, ok := params["address"]; ok {
		if err := gossipAddr.Set(address.(string)); err != nil {
			sendError(&w, err)
			return
		}
	}
	if peerParams, ok := params["peers"]; ok {
		if err := peers.Set(peerParams.(string)); err != nil {
			sendError(&w, err)
			return
		}
	}
	gossiperName := startGossiper(name, gossipAddr.String(), &peers)
	if gossiperName == "" {
		sendError(&w, errors.New("Error starting gossiper"))
		return
	}
	send(&w, getStatusResponse(gossiperName))
}

func send(w *http.ResponseWriter, v interface{}) {
	if reflect.ValueOf(v).IsNil() {
		sendError(w, errors.New("Error sending response"))
		return
	}
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(http.StatusOK)
	if err := json.NewEncoder(*w).Encode(v); err != nil {
		panic(err)
	}
}
func sendOk(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusOK)
	fmt.Fprintf(*w, "OK")
}
func sendError(w *http.ResponseWriter, msg error) {
	(*w).WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(*w, "There was an error processing the request: %v\n", msg.Error())
	fmt.Printf("There was an error processing the request: %v\n", msg.Error())

}
func readBody(w *http.ResponseWriter, r *http.Request) *map[string]interface{} {
	var params map[string]interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		sendError(w, err)
		return nil
	}
	if err := r.Body.Close(); err != nil {
		sendError(w, err)
		return nil
	}
	if err := json.Unmarshal(body, &params); err != nil {
		sendError(w, err)
		return nil
	}
	return &params
}

func getNameFromRequest(r *http.Request) (string, bool) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		return "", false
	}
	return name[0], true
}
