package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Rule struct {
	Pattern  string `json:"pattern"`
	Function string `json:"function"`
}

type SymbolSpec struct {
	Name       string   `json:"-"`                    // "Plus" - the symbol name
	Attributes []string `json:"attributes,omitempty"` // ["Flat", "Orderless"] - symbol attributes
	Functions  []Rule   `json:"functions,omitempty"`  // "(x__Integer)" , stdlib.PlusIntegers
}

type nodeWithPos struct {
	node ast.Node
	pos  token.Pos
}

func ParseSymbolSpecs(source []byte) ([]SymbolSpec, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	packageName := file.Name.Name

	// Collect all comments and declarations with their positions
	var nodes []nodeWithPos

	for _, commentGroup := range file.Comments {
		nodes = append(nodes, nodeWithPos{commentGroup, commentGroup.Pos()})
	}

	for _, decl := range file.Decls {
		nodes = append(nodes, nodeWithPos{decl, decl.Pos()})
	}

	// Sort by position
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].pos < nodes[j].pos
	})

	var symbols []SymbolSpec
	var currentSymbol *SymbolSpec

	// Process nodes in source order
	for _, nodeWrapper := range nodes {
		switch node := nodeWrapper.node.(type) {
		case *ast.CommentGroup:
			symbolName, attributes := parseTopLevelComments(node)
			if symbolName != "" {
				// Finalize previous symbol
				if currentSymbol != nil {
					symbols = append(symbols, *currentSymbol)
				}
				// Start new symbol
				currentSymbol = &SymbolSpec{
					Name:       symbolName,
					Attributes: attributes,
					Functions:  []Rule{},
				}
			} else if currentSymbol != nil && len(attributes) > 0 {
				// Add attributes to current symbol
				currentSymbol.Attributes = append(currentSymbol.Attributes, attributes...)
			}

		case *ast.FuncDecl:
			pattern := parseExprPattern(node)
			if pattern != "" && currentSymbol != nil {
				functionName := packageName + "." + node.Name.Name
				rule := Rule{
					Pattern:  pattern,
					Function: functionName,
				}
				currentSymbol.Functions = append(currentSymbol.Functions, rule)
			}
		}
	}

	// Add final symbol
	if currentSymbol != nil {
		symbols = append(symbols, *currentSymbol)
	}

	return symbols, nil
}

func parseTopLevelComments(commentGroup *ast.CommentGroup) (symbolName string, attributes []string) {
	for _, comment := range commentGroup.List {
		text := strings.TrimSpace(comment.Text)
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "@ExprSymbol") {
			parts := strings.Fields(text)
			if len(parts) > 1 {
				symbolName = parts[1]
			}
		} else if strings.HasPrefix(text, "@ExprAttributes") {
			parts := strings.Fields(text)
			if len(parts) > 1 {
				attributes = append(attributes, parts[1:]...)
			}
		}
	}
	return symbolName, attributes
}

func parseExprPattern(funcDecl *ast.FuncDecl) string {
	if funcDecl.Doc == nil {
		return ""
	}

	for _, comment := range funcDecl.Doc.List {
		text := strings.TrimSpace(comment.Text)
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "@ExprPattern") {
			parts := strings.Fields(text)
			if len(parts) > 1 {
				return strings.Join(parts[1:], " ")
			}
		}
	}
	return ""
}

func ParseSymbolSpecsFromDirectory(dirPath string) ([]SymbolSpec, error) {
	var allSymbols []SymbolSpec

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process .go files, but skip test files
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse the file
		symbols, err := ParseSymbolSpecs(content)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		// Add to combined results
		allSymbols = append(allSymbols, symbols...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return allSymbols, nil
}
