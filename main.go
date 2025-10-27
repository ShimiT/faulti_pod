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
    idx, err := strconv.Atoi(indexStr)
    if err != nil || idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }
    n, err := strconv.Atoi(parts[idx])
    if err != nil {
        http.Error(w, "invalid number", http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(fmt.Sprintf(`{"value":%d}`, n)))
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
    if os.Getenv("BUG") == "1" {
        var p *int
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
