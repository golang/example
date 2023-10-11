// The hugeparam command identifies by-value parameters that are larger than n bytes.
//
// Example:
//
//	$ ./hugeparams encoding/xml
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

// !+
var bytesFlag = flag.Int("bytes", 48, "maximum parameter size in bytes")

func PrintHugeParams(fset *token.FileSet, info *types.Info, sizes types.Sizes, files []*ast.File) {
	checkTuple := func(descr string, tuple *types.Tuple) {
		for i := 0; i < tuple.Len(); i++ {
			v := tuple.At(i)
			if sz := sizes.Sizeof(v.Type()); sz > int64(*bytesFlag) {
				fmt.Printf("%s: %q %s: %s = %d bytes\n",
					fset.Position(v.Pos()),
					v.Name(), descr, v.Type(), sz)
			}
		}
	}
	checkSig := func(sig *types.Signature) {
		checkTuple("parameter", sig.Params())
		checkTuple("result", sig.Results())
	}
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncDecl:
				checkSig(info.Defs[n.Name].Type().(*types.Signature))
			case *ast.FuncLit:
				checkSig(info.Types[n.Type].Type.(*types.Signature))
			}
			return true
		})
	}
}

//!-

func main() {
	flag.Parse()

	// Load complete type information for the specified packages,
	// along with type-annotated syntax and the "sizeof" function.
	// Types for dependencies are loaded from export data.
	conf := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(conf, flag.Args()...)
	if err != nil {
		log.Fatal(err) // failed to load anything
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1) // some packages contained errors
	}

	for _, pkg := range pkgs {
		PrintHugeParams(pkg.Fset, pkg.TypesInfo, pkg.TypesSizes, pkg.Syntax)
	}
}

/*
//!+output
% ./hugeparam encoding/xml
/go/src/encoding/xml/marshal.go:167:50: "start" parameter: encoding/xml.StartElement = 56 bytes
/go/src/encoding/xml/marshal.go:734:97: "" result: encoding/xml.StartElement = 56 bytes
/go/src/encoding/xml/marshal.go:761:51: "start" parameter: encoding/xml.StartElement = 56 bytes
/go/src/encoding/xml/marshal.go:781:68: "start" parameter: encoding/xml.StartElement = 56 bytes
/go/src/encoding/xml/xml.go:72:30: "" result: encoding/xml.StartElement = 56 bytes
//!-output
*/
