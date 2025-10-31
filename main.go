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
// Example (before fix): /calc?nums=1,2,3&index=10 would panic with "index out of range".
// The bounds check below prevents that panic by validating the parsed index.
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

    // Convert the index string to an int. We intentionally ignore the Atoi error here
    // to preserve previous behavior (non-numeric values effectively become 0),
    // and rely on the bounds check below to ensure safety.
    idx, _ := strconv.Atoi(indexStr)

    // FIX: Validate idx is within [0, len(parts)) to avoid a runtime panic from
    // out-of-range slice access ("index out of range"). If invalid, return 400.
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }

    // Safe to access parts[idx] after the bounds check; still ignoring Atoi error
    // here to keep behavior consistent with prior implementation.
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
