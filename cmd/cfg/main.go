// TODO:
//
// [x] generate different nodes for DSO as these are tail call targets.
// [x] generate start differently.
// [x] use subgraph to group the tail call together with the sub fns.
// [ ] some bug with tail_ipv6_ct_egress__CT_TAIL_CALL_BUFFER6
// [ ] policy call needs to be added

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bowei/cilium-bpf-hack/pkg/llvmp"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/ignore"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/srcnote"
)

var (
	theFlags = struct {
		mode       string
		in         string
		start      string
		ignoreFcns []string
		anFiles    []string
	}{
		ignoreFcns: []string{"@default"},
	}
)

func init() {
	flag.StringVar(&theFlags.mode, "mode", "", "rawcg | full")
	flag.StringVar(&theFlags.in, "in", "", "input file")
	flag.StringVar(&theFlags.start, "start", "", "name of function to start call graph from")

	flag.Func("ignore", "ignore function with this name. can specify multiple times",
		func(fn string) error {
			theFlags.ignoreFcns = append(theFlags.ignoreFcns, fn)
			return nil
		})
	flag.Func("an", "annotation file to read. See pkg/llvmp/srcnote for the file format.",
		func(fn string) error {
			theFlags.anFiles = append(theFlags.anFiles, fn)
			return nil
		})
}

func checkFlags() {
	switch theFlags.mode {
	case "rawcg":
		if theFlags.start == "" {
			fmt.Println("must specify -start", theFlags.mode)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid mode %q\n", theFlags.mode)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	checkFlags()

	m, err := llvmp.ParseLL(theFlags.in)
	if err != nil {
		panic(err)
	}

	switch theFlags.mode {
	case "rawcg":
		fmt.Printf("// Commandline: %+v\n", theFlags)
		ignored, err := ignore.Make(theFlags.ignoreFcns)
		if err != nil {
			panic(err)
		}
		srcAn, err := srcnote.Load(theFlags.anFiles...)
		if err != nil {
			panic(err)
		}

		out := llvmp.RawCG(m, &llvmp.RawCGParams{
			Start:   theFlags.start,
			Ignored: ignored,
			SrcAn:   srcAn,
		})
		fmt.Print(out)
	case "cg":
		// TODO
		fmt.Print(llvmp.Graphviz(m))
	}
}
