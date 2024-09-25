package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

var bytePool = sync.Pool{
	New: func() interface{} {
		// Create a new byte slice when a new one is needed
		return make([]byte, 1024*1024*2) // allocate 2 MB
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get a byte slice from the pool
	buf := bytePool.Get().([]byte)

	// Ensure the buffer is put back into the pool at the end of this function
	defer bytePool.Put(buf)

	// Read the request body into the buffer (assuming a small request body for simplicity)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the buffer (this example just echoes it back)
	w.Write(buf[:n])
}

func main() {
	http.HandleFunc("/gc", func(w http.ResponseWriter, r *http.Request) {
		// Print the memory statistics after allocating memory
		printMemStats()

		// Force a garbage collection cycle
		runtime.GC()

		// Print the memory statistics after garbage collection
		printMemStats()
	})

	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8079")
	err := http.ListenAndServe(":8079", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func printMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
		bToMb(mem.Alloc), bToMb(mem.TotalAlloc), bToMb(mem.Sys), mem.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
