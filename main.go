package main

import (
	"fmt"
	"net/http"
	"net"
	"Listener"
	"os/signal"
	"syscall"
	"sync"
	"os"
)

func main(){
	muxServe := http.NewServeMux()
	muxServe.HandleFunc("/root",HandleRoot)
	fmt.Println("Inicio")
	srv := http.Server{
		Addr: "localhost:8080",
		Handler: muxServe,
	}
	srv.ListenAndServe()
	fmt.Println("Fim")
}

func HandleRoot(w http.ResponseWriter, r *http.Request){
	originalListener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}

	sl, err := customListener.New(originalListener)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", helloHttp)
	server := http.Server{}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
	}()

	fmt.Printf("Serving HTTP\n")
	select {
	case signal := <-stop:
		fmt.Printf("Got signal:%v\n", signal)
	}
	fmt.Printf("Stopping listener\n")
	sl.Stop()
	fmt.Printf("Waiting on server\n")
	wg.Wait()
}


func helloHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
}