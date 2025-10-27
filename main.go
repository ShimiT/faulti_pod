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
    n, err := strconv.Atoi(parts[idx])
    if err != nil {
        http.Error(w, "invalid number format", http.StatusBadRequest)
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
