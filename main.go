package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
)

func healthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"ok":true}`))
}

func ready(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if os.Getenv("BUG") == "1" {
        w.WriteHeader(http.StatusServiceUnavailable)
        _, _ = w.Write([]byte(`{"ready":false}`))
        return
    }
    _, _ = w.Write([]byte(`{"ready":true}`))
}

// calcHandler demonstrates a subtle bug: out-of-range slice access when index is too large.
// Example: /calc?nums=1,2,3&index=10 will panic with index out of range.
func calcHandler(w http.ResponseWriter, r *http.Request) {
    numsParam := r.URL.Query().Get("nums")
    if numsParam == "" {
        numsParam = "1,2,3"
    }
    indexStr := r.URL.Query().Get("index")
    if indexStr == "" {
        indexStr = "0"
    }
    parts := strings.Split(numsParam, ",")
    idx, _ := strconv.Atoi(indexStr)
    // FIX: Add bounds check before accessing parts[idx] to prevent panic on invalid indices.
    // Returns 400 Bad Request instead of crashing the server when idx is out of range.
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }
    // Safe to access parts[idx] here because idx has been validated above.
    n, _ := strconv.Atoi(parts[idx])
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(fmt.Sprintf(`{"value":%d}`, n)))
}

// crashHandler demonstrates nil pointer dereference when BUG=1.
func crashHandler(w http.ResponseWriter, r *http.Request) {
    if os.Getenv("BUG") == "1" {
        var p *int
        // BUG: nil pointer dereference
        _ = *p
    }
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"ok":true}`))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthz)
    mux.HandleFunc("/ready", ready)
    mux.HandleFunc("/calc", calcHandler)
    mux.HandleFunc("/crash", crashHandler)
    addr := ":8080"
    log.Printf("faulty-app listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
