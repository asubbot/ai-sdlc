package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type testASTFunc struct {
	Name           string
	DocStartLine   int
	DocEndLine     int
	BodyStartLine  int
	BodyEndLine    int
	HasDirectTSkip bool
}

type testASTIndex struct {
	filename     string
	funcs        []testASTFunc
	byName       map[string]testASTFunc
	commentLines map[int]struct{}
}

func parseTestASTIndex(src []byte, filename string) (*testASTIndex, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	index := &testASTIndex{
		filename:     filename,
		byName:       make(map[string]testASTFunc),
		commentLines: make(map[int]struct{}),
	}

	for _, group := range file.Comments {
		for _, comment := range group.List {
			if len(comment.Text) >= 2 && comment.Text[:2] == "//" {
				index.commentLines[fset.Position(comment.Slash).Line] = struct{}{}
			}
		}
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name == nil || fn.Body == nil {
			continue
		}

		name := fn.Name.Name
		if len(name) < 5 || name[:4] != "Test" {
			continue
		}

		info := testASTFunc{
			Name:           name,
			BodyStartLine:  fset.Position(fn.Body.Lbrace).Line,
			BodyEndLine:    fset.Position(fn.Body.Rbrace).Line,
			HasDirectTSkip: funcBodyHasTSkip(fn.Body),
		}
		if fn.Doc != nil {
			info.DocStartLine = fset.Position(fn.Doc.Pos()).Line
			info.DocEndLine = fset.Position(fn.Doc.End()).Line
		}

		index.funcs = append(index.funcs, info)
		index.byName[name] = info
	}

	return index, nil
}

func (idx *testASTIndex) bindTraceLine(line int) (string, error) {
	if !idx.isActualCommentLine(line) {
		return "", fmt.Errorf("%s:%d: AC trace comment is not attached to a top-level Test* doc comment or body", idx.filename, line)
	}
	for _, fn := range idx.funcs {
		if fn.DocStartLine > 0 && line >= fn.DocStartLine && line <= fn.DocEndLine {
			return fn.Name, nil
		}
		if line >= fn.BodyStartLine && line <= fn.BodyEndLine {
			return fn.Name, nil
		}
	}

	return "", fmt.Errorf("%s:%d: AC trace comment is not attached to a top-level Test* doc comment or body", idx.filename, line)
}

func (idx *testASTIndex) isActualCommentLine(line int) bool {
	_, ok := idx.commentLines[line]
	return ok
}

func (idx *testASTIndex) hasDirectTSkip(name string) bool {
	fn, ok := idx.byName[name]
	return ok && fn.HasDirectTSkip
}

func (idx *testASTIndex) countTestsWithSkip() int {
	count := 0
	for _, fn := range idx.funcs {
		if fn.HasDirectTSkip {
			count++
		}
	}
	return count
}

func (idx *testASTIndex) topLevelTestNames() []string {
	names := make([]string, 0, len(idx.funcs))
	for _, fn := range idx.funcs {
		names = append(names, fn.Name)
	}
	return names
}

// parseTestFuncsWithTSkip keeps existing tests focused on skip detection while
// delegating to the shared AST index.
func parseTestFuncsWithTSkip(src []byte, filename string) (map[string]bool, error) {
	idx, err := parseTestASTIndex(src, filename)
	if err != nil {
		return nil, err
	}

	out := make(map[string]bool, len(idx.funcs))
	for _, fn := range idx.funcs {
		out[fn.Name] = fn.HasDirectTSkip
	}
	return out, nil
}

func funcBodyHasTSkip(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}

	var found bool
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if !ok || id.Name != "t" || sel.Sel == nil {
			return true
		}
		if sel.Sel.Name == "Skip" {
			found = true
			return false
		}
		return true
	})

	return found
}
