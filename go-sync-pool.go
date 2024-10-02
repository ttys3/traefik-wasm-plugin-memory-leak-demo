package main

// #include <stdlib.h>
// #include <stdio.h>
//
// typedef struct {
//     int  id;
//     char data[1024*1024*2];
// } LargeBuffer;
//
// LargeBuffer* createLargeBuffer() {
// LargeBuffer* buf = (LargeBuffer*)malloc(sizeof(LargeBuffer));
// if (buf == NULL) {
// 	printf("Error: Failed to allocate memory for LargeBuffer\n");
//     return NULL;
// }
// buf->id = 123;
// return buf;
// }
//
// void freeLargeBuffer(LargeBuffer* buf) {
//     free(buf);
// }
import "C"
import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	_ "net/http/pprof"
	"runtime"
)

var totalCreated atomic.Int64

var lock sync.Mutex

// LargeBufferWrapper wraps the C LargeBuffer struct
type LargeBufferWrapper struct {
	cBuffer *C.LargeBuffer
}

// NewLargeBufferWrapper creates a new LargeBufferWrapper
func NewLargeBufferWrapper() *LargeBufferWrapper {
	totalCreated.Add(1)

	return &LargeBufferWrapper{
		cBuffer: C.createLargeBuffer(),
	}
}

// Close frees the allocated memory
func (lbw *LargeBufferWrapper) Close() {
	if lbw.cBuffer != nil {
		totalCreated.Add(-1)
		C.freeLargeBuffer(lbw.cBuffer)
		lbw.cBuffer = nil
	}
}

// Data returns a byte slice of the underlying buffer
func (lbw *LargeBufferWrapper) Data() []byte {
	return (*[1024 * 1024 * 2]byte)(unsafe.Pointer(lbw.cBuffer))[:]
}

var bytePool = sync.Pool{}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get a byte slice from the pool
	buf, ok := bytePool.Get().(*LargeBufferWrapper)

	if buf == nil {
		buf = NewLargeBufferWrapper()
		fmt.Printf("buffer created, current_total: %d\n", totalCreated.Load())

		// runtime.SetFinalizer(buf, func(buf *LargeBufferWrapper) {
		// 	fmt.Printf("finalizer called, current_total: %d\n", totalCreated.Load())
		// 	buf.Close()
		// })
	} else if !ok {
		fmt.Println("-------------- buffer invalid --------------")
	}

	// Ensure the buffer is put back into the pool at the end of this function
	defer func() {
		// Simulate some work
		time.Sleep(time.Millisecond * 10)
		bytePool.Put(buf)
	}()

	buf.cBuffer.id = C.int(rand.IntN(1000000))
	dummy := []byte(fmt.Sprintf("hello world %d", buf.cBuffer.id))
	// write dummy data to the buffer
	copy(buf.Data(), dummy)

	// Process the buffer (this example just echoes it back)
	w.Write(buf.Data()[:len(dummy)])
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
