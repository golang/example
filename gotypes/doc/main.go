// The doc command prints the doc comment of a package-level object.
package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: doc <package> <object>")
	}
	//!+part1
	pkgpath, name := os.Args[1], os.Args[2]

	// Load complete type information for the specified packages,
	// along with type-annotated syntax.
	// Types for dependencies are loaded from export data.
	conf := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(conf, pkgpath)
	if err != nil {
		log.Fatal(err) // failed to load anything
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1) // some packages contained errors
	}

	// Find the package and package-level object.
	pkg := pkgs[0]
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		log.Fatalf("%s.%s not found", pkg.Types.Path(), name)
	}
	//!-part1
	//!+part2

	// Print the object and its methods (incl. location of definition).
	fmt.Println(obj)
	for _, sel := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
		fmt.Printf("%s: %s\n", pkg.Fset.Position(sel.Obj().Pos()), sel)
	}

	// Find the path from the root of the AST to the object's position.
	// Walk up to the enclosing ast.Decl for the doc comment.
	for _, file := range pkg.Syntax {
		pos := obj.Pos()
		if !(file.FileStart <= pos && pos < file.FileEnd) {
			continue // not in this file
		}
		path, _ := astutil.PathEnclosingInterval(file, pos, pos)
		for _, n := range path {
			switch n := n.(type) {
			case *ast.GenDecl:
				fmt.Println("\n", n.Doc.Text())
				return
			case *ast.FuncDecl:
				fmt.Println("\n", n.Doc.Text())
				return
			}
		}
	}
	//!-part2
}

// (The $GOROOT below is the actual string that appears in file names
// loaded from export data for packages in the standard library.)

/*
//!+output
$ ./doc net/http File
type net/http.File interface{Readdir(count int) ([]os.FileInfo, error); Seek(offset int64, whence int) (int64, error); Stat() (os.FileInfo, error); io.Closer; io.Reader}
$GOROOT/src/io/io.go:92:2: method (net/http.File) Close() error
$GOROOT/src/io/io.go:71:2: method (net/http.File) Read(p []byte) (n int, err error)
/go/src/net/http/fs.go:65:2: method (net/http.File) Readdir(count int) ([]os.FileInfo, error)
$GOROOT/src/net/http/fs.go:66:2: method (net/http.File) Seek(offset int64, whence int) (int64, error)
/go/src/net/http/fs.go:67:2: method (net/http.File) Stat() (os.FileInfo, error)

 A File is returned by a FileSystem's Open method and can be
served by the FileServer implementation.

The methods should behave the same as those on an *os.File.
//!-output
*/
