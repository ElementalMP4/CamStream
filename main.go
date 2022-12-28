package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

const frame_rate = 30

type CamstreamConfig struct {
	Password string
	Port     int
}

var config CamstreamConfig
var webcam *gocv.VideoCapture
var mutex = &sync.Mutex{}
var token string

func streamHandler(w http.ResponseWriter, r *http.Request) {
	user_token := r.URL.Query().Get("token")
	if user_token != token {
		sendForbiddenResponse(w, r, "Invalid Token. You need to log in.")
		return
	}

	if webcam == nil {
		mutex.Lock()
		defer mutex.Unlock()

		if webcam == nil {
			webcam, _ = gocv.VideoCaptureDevice(0)
		}
	}

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	for {
		frame := gocv.NewMat()
		if ok := webcam.Read(&frame); !ok {
			http.Error(w, "Failed to read frame", http.StatusInternalServerError)
			return
		}
		if frame.Empty() {
			http.Error(w, "Empty frame", http.StatusInternalServerError)
			return
		}

		buf, err := gocv.IMEncode(".jpg", frame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("--frame\n"))
		w.Write([]byte("Content-Type: image/jpeg\n"))
		w.Write([]byte("Content-Length: " + fmt.Sprint(len(buf.GetBytes())) + "\n"))
		w.Write([]byte("\n"))
		w.Write(buf.GetBytes())
		w.Write([]byte("\n"))
		if fw, ok := w.(http.Flusher); ok {
			fw.Flush()
			buf.Close()
			frame.Close()
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	user_password := r.URL.Query().Get("password")
	if user_password == config.Password {
		sendAcceptedResponse(w, r, token)
	} else {
		sendDeniedResponse(w, r, "Invalid Password")
	}
}

func handleTokenValidation(w http.ResponseWriter, r *http.Request) {
	user_token := r.URL.Query().Get("token")
	if user_token == token {
		sendAcceptedResponse(w, r, "Invalid Token")
	} else {
		sendDeniedResponse(w, r, "Invalid Token")
	}
}

func sendForbiddenResponse(w http.ResponseWriter, r *http.Request, response string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(response))
}

func sendDeniedResponse(w http.ResponseWriter, r *http.Request, response string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(response))
}

func sendAcceptedResponse(w http.ResponseWriter, r *http.Request, response string) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(response))
}

func startServer() {
	fmt.Printf("Server available on port %d\n", config.Port)

	http.HandleFunc("/video", streamHandler)
	http.HandleFunc("/authenticate", handleLogin)
	http.HandleFunc("/validate", handleTokenValidation)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	http.Handle("/", http.RedirectHandler("/static/html/login.html", http.StatusSeeOther))
	http.Handle("/stream", http.RedirectHandler("/static/html/stream.html", http.StatusSeeOther))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Couldn't find config!")
		return
	}

	config = CamstreamConfig{}
	err = json.Unmarshal([]byte(file), &config)
	token = uuid.NewString()

	if err != nil {
		fmt.Printf("Error opening capture device: \n")
		return
	}

	startServer()
}
