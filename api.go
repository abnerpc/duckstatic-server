package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setAPIRoute(route string) string {
	return fmt.Sprintf("/api/%s/", route)
}

// SetupAPIHandlers set up the api handlers
func SetupAPIHandlers() {

	changePassword := http.HandlerFunc(changePasswordKeyHandler)
	addUser := http.HandlerFunc(addUserHandler)
	deleteUser := http.HandlerFunc(deleteUserHandler)

	http.Handle(setAPIRoute("change_password"), AdminSecureMiddleware(changePassword))
	http.Handle(setAPIRoute("add_user"), AdminSecureMiddleware(addUser))
	http.Handle(setAPIRoute("delete_user"), AdminSecureMiddleware(deleteUser))

}

func changePasswordKeyHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		UserName    string `json:"username"`
		NewPassword string `json:"new_password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.UserName == "" || user.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}
	msg, ok := CurrentConfig.changePassword(user.UserName, user.NewPassword)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		UserType byte   `json:"user_type"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.UserName == "" || user.Password == "" || (user.UserType != Admin && user.UserType != User) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := CurrentConfig.addUser(user.UserName, user.Password, user.UserType)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		UserName string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.UserName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := CurrentConfig.deleteUser(user.UserName)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}
