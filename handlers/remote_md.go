package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// GET /api/md/{slug} — fetches from GitHub raw
func GetMDFromGitHub(w http.ResponseWriter, r *http.Request) {
    slug := r.PathValue("slug")

    // Build the GitHub raw URL
    rawURL := "https://raw.githubusercontent.com/gabrielatwell1987/portfolio/main/src/content/posts/" + slug + ".md"

    // (rest is the same as ProxyMDFile — fetch, convert, serve)
    resp, err := http.Get(rawURL)
    if err != nil {
        http.Error(w, "failed to fetch post: "+err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        http.Error(w, "post not found", http.StatusNotFound)
        return
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    accept := r.Header.Get("Accept")
    if strings.Contains(accept, "text/html") {
        md := goldmark.New(goldmark.WithExtensions(extension.GFM))
        var buf strings.Builder
        if err := md.Convert(body, &buf); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte(buf.String()))
    } else {
        w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
        w.Write(body)
    }
}