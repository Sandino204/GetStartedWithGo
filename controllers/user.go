package controllers

import (
	"encoding/json"
	"net/http"
	"project/models"
	"regexp"
	"strconv"
)

type userController struct{
	userIDPattern *regexp.Regexp
}

func (uc *userController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetUsers(), w)
}

func (uc *userController) get(id int, w http.ResponseWriter) {
	u, err := models.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	encodeResponseAsJSON(u, w)
}

func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse the user"))
		return 
	}

	u, err = models.AddUser(u)
}

func (uc *userController) patch(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse the object"))
		return 
	}

	if id != u.ID {
		w.WriteHeader((http.StatusBadRequest))
		w.Write([]byte("ID of Must match the id in url"))
		return
	}

	u, err = models.UpdateUser(u)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return 
	}

	encodeResponseAsJSON(u, w)
}

func (uc *userController) delete(id int, w http.ResponseWriter) {
	err := models.DeleteUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *userController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet: 
			uc.getAll(w, r)
		case http.MethodPost: 
			uc.post(w, r)
		default: 
			w.WriteHeader(http.StatusNotImplemented)
		}

	} else {
		matches := uc.userIDPattern.FindStringSubmatch(r.URL.Path)
        if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		switch r.Method {
		case http.MethodGet:
			uc.get(id, w)
		case http.MethodPatch:
			uc.patch(id, w, r)
		case http.MethodDelete:
			uc.delete(id, w)
		default: 
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func newUserController() *userController {
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}