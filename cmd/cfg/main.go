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
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/rawcg"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/srcnote"
)

var (
	theFlags = struct {
		mode       string
		in         string
		start      string
		ignoreFcns []string
		anFiles    []string
	}{}
)

func init() {
	flag.StringVar(&theFlags.mode, "mode", "", "rawcg | full")
	flag.StringVar(&theFlags.in, "in", "", "input file")
	flag.StringVar(&theFlags.start, "start", "", "Name of function to start call graph from")

	flag.Func("ignore", "Ignore function with this name. Can specify multiple times. Defaults to @default",
		func(fn string) error {
			theFlags.ignoreFcns = append(theFlags.ignoreFcns, fn)
			return nil
		})
	flag.Func("an", "Annotation file to read. See pkg/llvmp/srcnote for the file format.",
		func(fn string) error {
			theFlags.anFiles = append(theFlags.anFiles, fn)
			return nil
		})
}

func checkAndDefaultFlags() {
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
	if theFlags.ignoreFcns == nil {
		theFlags.ignoreFcns = []string{"@default"}
	}
}

func main() {
	flag.Parse()

	checkAndDefaultFlags()

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

		out, err := rawcg.Run(m, &rawcg.Params{
			Start:   theFlags.start,
			Ignored: ignored,
			SrcAn:   srcAn,
		})
		if err != nil {
			// TODO: error
			fmt.Printf("// ERROR: rawcg.Run() = %v\n", err)
		}
		fmt.Print(out)
	case "cg":
		// TODO
		fmt.Print(llvmp.Graphviz(m))
	}
}
