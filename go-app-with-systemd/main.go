package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ming-go/pkg/mhttp"
)

func main() {
	//var ret uintptr
	//var err uintptr
	//ret, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	//if err != 0 {
	//	log.Println("fork fail")
	//	os.Exit(-1)
	//}

	//if ret != 0 {
	// child process

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Hello, world!")
	})

	cfg := mhttp.NewConfig().SetMux(mux).SetAddr("", "8787")
	s := mhttp.NewDefaultServer(cfg)
	log.Println(s.Run())
	os.Exit(0)
	//}

	<-time.After(10 * time.Minute)

}

func monitor() {

}

func monkey() {

}
