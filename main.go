package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
)

// healthz is a simple liveness probe endpoint returning a static JSON payload.
func healthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"ok":true}`))
}

// ready reports readiness. When BUG=1 it intentionally reports not-ready (503) to simulate a failing dependency.
func ready(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if os.Getenv("BUG") == "1" {
        w.WriteHeader(http.StatusServiceUnavailable)
        _, _ = w.Write([]byte(`{"ready":false}`))
        return
    }
    _, _ = w.Write([]byte(`{"ready":true}`))
}

// calcHandler parses a comma-separated list of numbers (nums) and returns the value at the given index.
// It includes a bounds check to avoid out-of-range slice access that would otherwise panic.
// Example: /calc?nums=1,2,3&index=10 now returns 400 instead of panicking.
func calcHandler(w http.ResponseWriter, r *http.Request) {
    // Use default list if "nums" is not provided.
    numsParam := r.URL.Query().Get("nums")
    if numsParam == "" {
        numsParam = "1,2,3"
    }
    // Default index to 0 when not provided.
    indexStr := r.URL.Query().Get("index")
    if indexStr == "" {
        indexStr = "0"
    }
    // Split the list and parse the desired index.
    parts := strings.Split(numsParam, ",")
    idx, _ := strconv.Atoi(indexStr)
    // Guard against negative or too-large index to prevent panic from parts[idx].
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }
    // Convert the selected element to an integer and return it as JSON.
    n, _ := strconv.Atoi(parts[idx])
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(fmt.Sprintf(`{"value":%d}`, n)))
}

// crashHandler demonstrates nil pointer dereference when BUG=1 to simulate a crash.
func crashHandler(w http.ResponseWriter, r *http.Request) {
    if os.Getenv("BUG") == "1" {
        var p *int
        // Intentional nil pointer dereference when BUG=1.
        _ = *p
    }
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"ok":true}`))
}

func main() {
    // Wire up HTTP handlers and start the HTTP server.
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthz)
    mux.HandleFunc("/ready", ready)
    mux.HandleFunc("/calc", calcHandler)
    mux.HandleFunc("/crash", crashHandler)
    addr := ":8080"
    log.Printf("faulty-app listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
