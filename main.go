package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

const frame_rate int = 60
const password string = "TestPassword"
const httpPort int = 3000

var webcam *gocv.VideoCapture
var buffer = make(map[int][]byte)
var frame []byte
var mutex = &sync.Mutex{}
var err error
var token string

func handleVideoStream(w http.ResponseWriter, r *http.Request) {
	user_token := r.URL.Query().Get("token")
	if user_token == token {
		r.Cookies()
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		data := ""
		for {
			mutex.Lock()
			data = "--frame\r\n  Content-Type: image/jpeg\r\n\r\n" + string(frame) + "\r\n\r\n"
			mutex.Unlock()
			time.Sleep(time.Duration(1000/frame_rate) * time.Millisecond) //60fps (ish)
			w.Write([]byte(data))
		}
	} else {
		sendForbiddenResponse(w, r, "Invalid Token. You need to log in.")
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	user_password := r.URL.Query().Get("password")
	if user_password == password {
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
	fmt.Println("Server available on port 3000")
	mux := http.NewServeMux()

	mux.HandleFunc("/video", handleVideoStream)
	mux.HandleFunc("/authenticate", handleLogin)
	mux.HandleFunc("/validate", handleTokenValidation)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	mux.Handle("/", http.RedirectHandler("/static/html/login.html", http.StatusSeeOther))
	mux.Handle("/stream", http.RedirectHandler("/static/html/stream.html", http.StatusSeeOther))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux))
}

func getframes() {
	img := gocv.NewMat()
	defer img.Close()
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed\n")
			return
		}
		if img.Empty() {
			continue
		}
		frame_buffer, _ := gocv.IMEncode(".jpg", img)
		frame = frame_buffer.GetBytes()
	}
}

func main() {
	webcam, err = gocv.VideoCaptureDevice(0)
	token = uuid.NewString()

	if err != nil {
		fmt.Printf("Error opening capture device: \n")
		return
	}

	defer webcam.Close()
	go getframes()
	startServer()
}
