package main

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestAnalyzeFile_DetectsInitLangCall(t *testing.T) {
	src := `package ui

import "fyne.io/fyne/v2/lang"

func init() {
	_ = lang.X("metric.help", "fallback")
}
`

	violations, warnings := analyzeFromSource(t, src)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].message, "must not be called during package init") {
		t.Fatalf("unexpected violation message: %q", violations[0].message)
	}
}

func TestAnalyzeFile_DetectsPackageVarLangCall(t *testing.T) {
	src := `package ui

import "fyne.io/fyne/v2/lang"

var defs = []string{
	lang.X("metric.help", "fallback"),
}
`

	violations, warnings := analyzeFromSource(t, src)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestAnalyzeFile_AllowsRuntimeLangCall(t *testing.T) {
	src := `package ui

import "fyne.io/fyne/v2/lang"

func labelText() string {
	return lang.X("metric.help", "fallback")
}
`

	violations, warnings := analyzeFromSource(t, src)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestAnalyzeFile_IgnoreWithoutReasonWarns(t *testing.T) {
	src := `package ui

import "fyne.io/fyne/v2/lang"

func init() {
	//i18n:ignore
	_ = lang.X("metric.help", "fallback")
}
`

	violations, warnings := analyzeFromSource(t, src)
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning line, got %v", warnings)
	}
}

func analyzeFromSource(t *testing.T, src string) ([]violation, []int) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "internal/ui/sample.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	commentsByLine := collectIgnoreTags(fset, file)
	return analyzeFile(fset, file, "internal/ui/sample.go", commentsByLine)
}
