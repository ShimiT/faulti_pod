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

    // Parse the provided index safely. If it's not a valid integer, return 400 instead of proceeding.
    idx, err := strconv.Atoi(indexStr)
    if err != nil {
        http.Error(w, "invalid index", http.StatusBadRequest)
        return
    }

    // Bounds check to prevent out-of-range slice access (preserves the core bug fix).
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }

    // Parse the selected element safely. If the chosen nums value isn't an integer, return 400
    // rather than using the zero value from a failed conversion.
    n, err := strconv.Atoi(parts[idx])
    if err != nil {
        http.Error(w, "invalid nums value", http.StatusBadRequest)
        return
    }

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
