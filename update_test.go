package main

import "testing"

func TestPlainInit(t *testing.T) {
	path := parseURL("github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}
func TestPlainInitWithHTTP(t *testing.T) {
	path := parseURL("http://github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestPlainInitWithHTTPS(t *testing.T) {
	path := parseURL("https://github.com/okkur/gitor.git")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}

func TestPlainInitExtractingHostAndPath(t *testing.T) {
	path := parseURL("git@github.com:okkur/gitor")
	if path != "github.com/okkur/gitor" {
		t.Error("Expected github.com/okkur/gitor, got ", path)
	}
}
