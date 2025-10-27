package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "net/http/httptest"
    "testing"
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
    // BUG: no bounds check
    if idx < 0 || idx >= len(parts) {
        http.Error(w, "index out of range", http.StatusBadRequest)
        return
    }
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

func TestCalcHandler_IndexOutOfRange(t *testing.T) {
    req := httptest.NewRequest("GET", "/calc?nums=1,2,3&index=10", nil)
    rr := httptest.NewRecorder()
    calcHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", rr.Code)
    }
}

func TestCalcHandler_NegativeIndex(t *testing.T) {
    req := httptest.NewRequest("GET", "/calc?nums=1,2,3&index=-1", nil)
    rr := httptest.NewRecorder()
    calcHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", rr.Code)
    }
}

func TestCalcHandler_ValidIndex(t *testing.T) {
    req := httptest.NewRequest("GET", "/calc?nums=1,2,3&index=1", nil)
    rr := httptest.NewRecorder()
    calcHandler(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
    want := `{"value":2}`
    got := strings.TrimSpace(rr.Body.String())
    if got != want {
        t.Fatalf("expected %s, got %s", want, got)
    }
}
