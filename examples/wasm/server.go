//go:build !js

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	port := flag.Int("port", 8080, "port to serve on")
	build := flag.String("build", "", "path to example to build (e.g., ./examples/shapes)")
	flag.Parse()

	// Get the directory where this server.go is located
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// Fallback to current directory
		dir, _ = os.Getwd()
	}

	// If running via go run, use the source file location
	if _, err := os.Stat(filepath.Join(dir, "index.html")); os.IsNotExist(err) {
		dir, _ = os.Getwd()
	}

	// Build WASM if requested
	if *build != "" {
		fmt.Printf("Building %s to main.wasm...\n", *build)
		cmd := exec.Command("go", "build", "-o", filepath.Join(dir, "main.wasm"), *build)
		cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Build failed: %v", err)
		}
		fmt.Println("Build complete!")
	}

	// Check if main.wasm exists
	wasmPath := filepath.Join(dir, "main.wasm")
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		fmt.Println("Warning: main.wasm not found. Build an example first:")
		fmt.Println("  go run server.go -build ./examples/shapes")
		fmt.Println("  or")
		fmt.Println("  GOOS=js GOARCH=wasm go build -o main.wasm ./examples/shapes")
		fmt.Println()
	}

	// Set up file server with correct MIME types
	fs := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set correct MIME type for WASM files
		if filepath.Ext(r.URL.Path) == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fs.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Serving at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
