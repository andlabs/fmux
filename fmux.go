// fmux:  file multiplexer
package main

// TODO: really use flag? there's only one option

import (
	"fmt"
//	"io"
	"log"
	"flag"
	"os"
	"encoding/binary"
)

var infiles []*os.File
var outfile *os.File = os.Stdout

var outfname *string = flag.String("o", "", "output file")

// run:  the actual work
func run(fsize int64) {
//	var b byte
	var err error

	n := len(infiles)
	fullbuf := make([]byte, fsize*int64(n))
	filebuf := make([]byte, fsize)

	for i := 0; i < n; i++ {
		err = binary.Read(infiles[i], binary.BigEndian, &filebuf)
		if err != nil {
			log.Fatalf("%s: read failed: %v\n", flag.Arg(i), err)
		}
		j := int64(i)
		for k := int64(0); k < fsize; k++ {
			fullbuf[j] = filebuf[k]
			j += int64(n)
		}
	}

	err = binary.Write(outfile, binary.BigEndian, fullbuf)
	if err != nil {
		log.Fatalf("%s: write failed: %v\n", outfname, err)
	}

/*
	for j := int64(0); j < fsize; j++ {	// crash if I counted fsize down to zero?
		for i := 0; i < len(infiles); i++ {
			err = binary.Read(infiles[i], binary.BigEndian, &b)
			if err != nil {
				log.Fatalf("%s: read failed: %v\n", flag.Arg(i), err)
			}
			err = binary.Write(outfile, binary.BigEndian, b)
			if err != nil {
				log.Fatalf("%s: write failed: %v\n", *outfname, err)
			}
		}
	}
*/
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
	} else {
		*outfname = "<stdout>"	// for error reporting
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
	run(fsize)
}
