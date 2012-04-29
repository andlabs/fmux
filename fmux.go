// fmux:  file multiplexer
package main

// TODO: really use flag? there's only one option

import (
	"fmt"
//	"io"
	"log"
	"flag"
	"os"
)

var infiles []*os.File
var outfile *os.File = os.Stdout

var outfname *string = flag.String("o", "", "output file")

// run:  the actual work
func run() {
	// TODO
}

// getsize:  get file size
func getsize(name string) int64 {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatalf("%s: stat (to get file size) failed: %v", name, err)
	}
	return info.Size()
}

// setup:  handle command line and open files
func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [-o outfile] files\n", os.Args[0])
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	if *outfname != "" {
		outfile, err = os.Create(*outfname)
		if err != nil {
			log.Fatalf("%s: create failed: %v\n", outfname, err)
		}
		defer outfile.Close()
	}

	first := flag.Arg(0)
	fsize := getsize(first)

	infiles = make([]*os.File, flag.NArg())
	for i := 0; i < flag.NArg(); i++ {
		fn := flag.Arg(i)
		if l := getsize(fn); l != fsize {
			log.Fatalf("%s: size (%d) does not match that of %s (%d)",
				fn, l, first, fsize)
		}
		infiles[i], err = os.Open(fn)
		if err != nil {
			log.Fatalf("%s: open failed: %v\n", fn, err)
		}
		defer infiles[i].Close()
	}

	// now do the work
	run()
}
