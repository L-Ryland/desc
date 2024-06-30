package gateway

import (
	"net/http"

	"github.com/gorilla/mux"

	"server/datasys"
	"server/usersys"
)

const pathPerfix = "/v1"

func NewService(router *mux.Router) {
	// user
	router.HandleFunc(pathPerfix+"/register", usersys.HandleRegister).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/login", usersys.HandleLogin).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/logout", usersys.HandleLogout).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/auth", usersys.HandleGetAuth).Methods(http.MethodGet)
	router.HandleFunc(pathPerfix+"/user", usersys.HandleGetUsers).Methods(http.MethodGet)
	router.HandleFunc(pathPerfix+"/user", usersys.HandleAddUser).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/user/{id}", usersys.HandleRemoveUser).Methods(http.MethodDelete)

	// web data
	router.HandleFunc(pathPerfix+"/web", datasys.HandleAddWeb).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/web/{tags}", datasys.HandleSearchWeb).Methods(http.MethodGet)
	router.HandleFunc(pathPerfix+"/web/{id}", datasys.HandleDeleteWeb).Methods(http.MethodDelete)
	router.HandleFunc(pathPerfix+"/web/{id}", datasys.HandlePatchWeb).Methods(http.MethodPatch)

	//tag data
	router.HandleFunc(pathPerfix+"/tag", datasys.HandleAddTag).Methods(http.MethodPost)
	router.HandleFunc(pathPerfix+"/tag", datasys.HandleGetAllTags).Methods(http.MethodGet)
	router.HandleFunc(pathPerfix+"/tag/{name}", datasys.HandleReorderTag).Methods(http.MethodPatch)
	router.HandleFunc(pathPerfix+"/tag/{name}", datasys.HandleDeleteTag).Methods(http.MethodDelete)

	// category data
	router.HandleFunc(pathPerfix+"/categories", datasys.HandleGetAllCategories).Methods(http.MethodGet)
	router.HandleFunc(pathPerfix+"/categories/{id}", datasys.HandleUpdateCategory).Methods(http.MethodPatch)
}
