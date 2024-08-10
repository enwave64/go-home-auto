package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type cmdresult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func homepage(write http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(write, "Go Home Simple REST API Server")
}

func getdate(write http.ResponseWriter, r *http.Request) {
	result := cmdresult{}

	out, err := exec.Command("date").Output()
	if err == nil {
		result.Success = true
		result.Message = "The date is " + string(out) + "and your request was: " + r.RequestURI
	}

	json.NewEncoder(write).Encode(result)
}

func play(write http.ResponseWriter, r *http.Request) {
	result := cmdresult{}
	fmt.Println("your URL query params: ", r.URL.Query())
	track := r.URL.Query().Get("track")
	if track == "" {
		track = "1"
	}
	out, err := exec.Command("mpc", "play", track).Output()
	result.Message = "MPD output: " + string(out)
	if err == nil {
		result.Success = true
	} else {
		result.Success = false
		errBytes := []byte(fmt.Sprintf("%v", err))
		result.Message += "Error msg: " + string(errBytes)
	}
	json.NewEncoder(write).Encode(result)
}

func stop(write http.ResponseWriter, r *http.Request) {
	result := cmdresult{}
	fmt.Println("your URL query params: ", r.URL.Query())
	out, err := exec.Command("mpc", "stop").Output()
	result.Message = "MPD output: " + string(out)
	if err == nil {
		result.Success = true
	} else {
		result.Success = false
		errBytes := []byte(fmt.Sprintf("%v", err))
		result.Message += "Error msg: " + string(errBytes)
	}
	json.NewEncoder(write).Encode(result)
}

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/api/v1/getdate", getdate)
	http.HandleFunc("/api/v1/play", play)
	http.HandleFunc("/api/v1/stop", stop)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
