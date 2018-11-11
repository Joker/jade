package main

import (
	"fmt"
	"net/http"

	pool "github.com/valyala/bytebufferpool"
)

func main() {
	buffer := pool.Get()

	Index("Jade.go", true, buffer)

	fmt.Printf("\nOutput:\n\n%s", buffer)

	http.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
		wr.Write(buffer.Bytes())
	})
	http.ListenAndServe(":8080", nil)
}
