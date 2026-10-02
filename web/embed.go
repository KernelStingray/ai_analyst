package web

import (
	"embed"
	"net/http"
)

//go:embed index.html static/*
var embeddedFiles embed.FS

// GetFileSystem returns an http.FileSystem serving the embedded UI assets
func GetFileSystem() http.FileSystem {
	return http.FS(embeddedFiles)
}

// ReadIndexHTML returns the raw bytes of embedded index.html
func ReadIndexHTML() ([]byte, error) {
	return embeddedFiles.ReadFile("index.html")
}

