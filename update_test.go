package main

import "testing"

func TestExtractingPathPlainHTTP(t *testing.T) {
	path := extractPath("github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestExtractingPathHTTP(t *testing.T) {
	path := extractPath("http://github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestExtractingPathHTTPPort(t *testing.T) {
	path := extractPath("http://github.com:8080/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestExtractingPathHTTPS(t *testing.T) {
	path := extractPath("https://github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestExtractingPathPlainSSH(t *testing.T) {
	path := extractPath("git@github.com:okkur/gitor")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestExtractingPathSSH(t *testing.T) {
	path := extractPath("ssh://git@github.com:okkur/gitor")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}
