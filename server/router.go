package main

import (
	"net/http"
	"path"

	"github.com/ageapps/Peerster/pkg/utils"

	"github.com/ageapps/Peerster/pkg/file"
	"github.com/gorilla/mux"
)

// Routes arr
type Routes []Route

// MAX_UPLOAD_SIZE contant
const MAX_UPLOAD_SIZE = 1000 * 1024 // 10 MB
var uploadPath = path.Join(utils.GetRootPath(), file.SharedBlobsDir)

// NewRouter func
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	fsWeb := http.FileServer(http.Dir("../web/"))
	fsFiles := http.FileServer(http.Dir("../" + file.DownloadsDir + "/"))
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", fsFiles))

	router.PathPrefix("/").Handler(fsWeb)

	return router
}
