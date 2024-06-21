package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"main/tools"
	"net/http"
)

func middleware(method string, handlerFunc http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		// ngecek apakah request method sesuai dengan di mux
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// kalau sesuai nanti diarahin ke function handler
		handlerFunc(w, r)
	}
}

func getMessageHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello World!")
}

func sendFileHandler(w http.ResponseWriter, r *http.Request){
	jsonData := r.FormValue("Person")

	var person data.Person
	err := json.Unmarshal([]byte(jsonData), &person)
	tools.ErrorHandler(err)


	fmt.Println("JSON : ", person)

	file, handler, err := r.FormFile("File")
	// handler: berisi informaasi ssperti nama file, size file, dsb.
	// file: baal berisi semua data
	tools.ErrorHandler(err)

	
	fmt.Printf("Received File : %s\n", handler.Filename)

	fileContent, err := io.ReadAll(file)
	tools.ErrorHandler(err)

	fmt.Printf("File Content : \n%s\n", fileContent)
	fmt.Fprint(w, "Succesfully Received Data")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware(http.MethodGet, getMessageHandler))
	mux.HandleFunc("/sendFile", middleware(http.MethodPost, sendFileHandler))

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServeTLS("../cert.pem", "../key.pem")
	tools.ErrorHandler(err)
}