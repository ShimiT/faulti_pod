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

// calcHandler previously had a subtle bug: out-of-range slice access when index was too large.
// The explicit bounds check below prevents a panic by validating the index and returning 400 instead.
// Example: /calc?nums=1,2,3&index=10 now returns HTTP 400 rather than crashing the server.
func calcHandler(w http.ResponseWriter, r *http.Request) {
    // Read input numbers from query; default to "1,2,3" if not provided.
    numsParam := r.URL.Query().Get("nums")
    if numsParam == "" {
        numsParam = "1,2,3"
    }

    // Read index from query; default to "0" if not provided.
    indexStr := r.URL.Query().Get("index")
    if indexStr == "" {
        indexStr = "0"
    }

    parts := strings.Split(numsParam, ",")

    // Convert index string to int. We intentionally ignore the conversion error here;
    // on failure Atoi returns 0, and the subsequent bounds check ensures safety.
    idx, _ := strconv.Atoi(indexStr)

    // Bounds check to avoid out-of-range slice access. If idx is invalid, return 400.
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }

    // Safe to access parts[idx] after the bounds check.
    n, _ := strconv.Atoi(parts[idx])

    // Respond with the selected value as JSON.
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
