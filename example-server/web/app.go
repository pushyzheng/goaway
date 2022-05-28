package main

import (
	"flag"
	"fmt"
	"net/http"
)

var port = flag.Int("p", 5000, "Input server port")

func echo(cont string) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.RequestURI)
		w.Write([]byte(cont))
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", echo("Hello World"))  // homepage
	http.HandleFunc("/public", echo("public")) // public path
	http.HandleFunc("/你好世界", echo("你好世界"))     // test for chinese path
	http.HandleFunc("/admin", echo("admin"))   // private path

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	fmt.Printf("Running on http://%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
