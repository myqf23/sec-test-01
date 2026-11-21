package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func probe(w http.ResponseWriter, target string) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		fmt.Fprintf(w, "ERR: %v", err)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Fprintf(w, "code=%d\nlen=%d\nbody=%s", resp.StatusCode, len(b), string(b))
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	if target == "" {
		// 默认探针：官方 SSRF 靶场（ByteSRC 指定验证通道）
		probe(w, "https://src-ssrf.bytedance.net/ssrf?host=139.224.40.24.nip.io")
		return
	}
	if target != "https://src-ssrf.bytedance.net/ssrf" && target != "https://src-ssrf.bytedance.net/ssrf?host=139.224.40.24.nip.io" {
		fmt.Fprint(w, "only bytedance ssrf range allowed")
		return
	}
	probe(w, target)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthz)
	fmt.Println("listening on", port)
	http.ListenAndServe(":"+port, nil)
}
