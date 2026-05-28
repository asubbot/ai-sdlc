package main

import (
	"path/filepath"
	"testing"
)

func TestCheckEARSPattern(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantValid bool
	}{
		{"ubiquitous", "THE system SHALL log all requests.", true},
		{"event-driven", "WHEN a request arrives, THE system SHALL respond.", true},
		{"state-driven", "WHILE idle, THE system SHALL release connections.", true},
		{"unwanted", "IF connection fails, THEN THE system SHALL retry.", true},
		{"optional", "WHERE CSV is enabled, THE system SHALL add headers.", true},
		{"complex", "WHILE in maintenance, WHEN alert occurs, THE system SHALL notify.", true},
		{"no shall", "The system needs to be fast.", false},
		{"no ears pattern", "The system efficiently handles data as needed.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := reqBlock{Code: "REQ-00.001", Body: tt.body, Line: 1}
			findings := checkEARSPattern(block)
			isValid := len(findings) == 0
			if isValid != tt.wantValid {
				t.Errorf("checkEARSPattern(%q) valid=%v, want %v; findings=%v", tt.body, isValid, tt.wantValid, findings)
			}
		})
	}
}

func TestLintEARSFile_ValidFixtures(t *testing.T) {
	validFiles, _ := filepath.Glob("testdata/ears/valid-*.md")
	if len(validFiles) == 0 {
		t.Fatal("no valid fixture files found")
	}
	for _, f := range validFiles {
		t.Run(filepath.Base(f), func(t *testing.T) {
			result := lintEARSFile(f)
			if result.Errors > 0 {
				t.Errorf("expected no errors for valid fixture, got %d: %v", result.Errors, result.Findings)
			}
		})
	}
}

func TestLintEARSFile_InvalidFixtures(t *testing.T) {
	invalidFiles, _ := filepath.Glob("testdata/ears/invalid-*.md")
	if len(invalidFiles) == 0 {
		t.Fatal("no invalid fixture files found")
	}
	for _, f := range invalidFiles {
		t.Run(filepath.Base(f), func(t *testing.T) {
			result := lintEARSFile(f)
			if result.Errors+result.Warnings == 0 {
				t.Error("expected findings for invalid fixture")
			}
		})
	}
}

func TestBannedWords(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"efficiently", "System should efficiently handle data.", true},
		{"if possible", "Process data if possible.", true},
		{"etc", "Handle errors, etc.", true},
		{"clean", "THE system SHALL process requests within 100ms.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := reqBlock{Code: "REQ-00.001", Body: tt.body, Line: 1}
			findings := checkBannedWords(block)
			hasErr := len(findings) > 0
			if hasErr != tt.wantErr {
				t.Errorf("checkBannedWords(%q) hasErr=%v, want %v; findings=%v", tt.body, hasErr, tt.wantErr, findings)
			}
		})
	}
}

func TestPassiveVoice(t *testing.T) {
	tests := []struct {
		body     string
		wantWarn bool
	}{
		{"Data is processed by the system.", true},
		{"THE system SHALL process data.", false},
	}
	for _, tt := range tests {
		t.Run(tt.body, func(t *testing.T) {
			block := reqBlock{Code: "REQ-00.001", Body: tt.body, Line: 1}
			findings := checkPassiveVoice(block)
			hasWarn := len(findings) > 0
			if hasWarn != tt.wantWarn {
				t.Errorf("checkPassiveVoice(%q) hasWarn=%v, want %v", tt.body, hasWarn, tt.wantWarn)
			}
		})
	}
}

func TestMultipleThoughts(t *testing.T) {
	tests := []struct {
		body     string
		wantWarn bool
	}{
		{"THE system SHALL validate input and THE system SHALL sanitize output.", true},
		{"THE system SHALL validate and sanitize input.", false},
	}
	for _, tt := range tests {
		t.Run(tt.body, func(t *testing.T) {
			block := reqBlock{Code: "REQ-00.001", Body: tt.body, Line: 1}
			findings := checkMultipleThoughts(block)
			hasWarn := len(findings) > 0
			if hasWarn != tt.wantWarn {
				t.Errorf("checkMultipleThoughts(%q) hasWarn=%v, want %v", tt.body, hasWarn, tt.wantWarn)
			}
		})
	}
}

func TestLintEARSFile_EP099(t *testing.T) {
	result := lintEARSFile("testdata/EP-099/ep-requirements.md")
	if result.Total != 6 {
		t.Errorf("total reqs = %d, want 6", result.Total)
	}
	if result.Errors > 0 {
		t.Errorf("expected no errors for healthy requirements, got %d", result.Errors)
	}
}

func TestLintEARSFile_EP098(t *testing.T) {
	result := lintEARSFile("testdata/EP-098/ep-requirements.md")
	if result.Total != 4 {
		t.Errorf("total reqs = %d, want 4", result.Total)
	}
	if result.Errors == 0 {
		t.Error("expected errors for broken requirements")
	}
}
