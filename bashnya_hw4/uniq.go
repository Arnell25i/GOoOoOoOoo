package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type options struct {
	count          bool 
	onlyDup        bool 
	onlyUniq       bool 
	ignoreCase     bool 
	skipFields     int
	skipChars      int
}

func (o options) validate() error {
	m := 0
	if o.count { m++ }
	if o.onlyDup { m++ }
	if o.onlyUniq { m++ }
	if m > 1 {
		return errors.New("flags -c, -d, -u are mutually exclusive; use only one")
	}
	if o.skipFields < 0 {
		return errors.New("-f must be >= 0")
	}
	if o.skipChars < 0 {
		return errors.New("-s must be >= 0")
	}
	return nil
}

func computeKey(line string, opt options) string {
	s := line
	i := 0
	n := len(s)
	for i < n && s[i] == ' ' {
		i++
	}

	for f := 0; f < opt.skipFields; f++ {
		for i < n && s[i] != ' ' {
			i++
		}
		for i < n && s[i] == ' ' {
			i++
		}
	}

	chars := opt.skipChars
	if chars > 0 {
		if i+chars > n {
			i = n
		} else {
			i += chars
		}
	}
	key := s[i:]
	if opt.ignoreCase {
		key = strings.ToLower(key)
	}
	return key
}

func process(r io.Reader, w io.Writer, opt options) error {
	if err := opt.validate(); err != nil {
		return err
	}
	in := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	in.Buffer(buf, 1024*1024)

	writer := bufio.NewWriter(w)
	defer writer.Flush()

	first := true
	var prevLine, prevKey string
	count := 0

	flush := func() error {
		if first {
			return nil
		}
		if opt.onlyDup {
			if count > 1 {
				if _, err := fmt.Fprintln(writer, prevLine); err != nil { return err }
			}
			return nil
		}
		if opt.onlyUniq {
			if count == 1 {
				if _, err := fmt.Fprintln(writer, prevLine); err != nil { return err }
			}
			return nil
		}
		if opt.count {
			if _, err := fmt.Fprintf(writer, "%d %s\n", count, prevLine); err != nil { return err }
			return nil
		}

		if _, err := fmt.Fprintln(writer, prevLine); err != nil { return err }
		return nil
	}

	for in.Scan() {
		line := in.Text()
		key := computeKey(line, opt)
		if first {
			prevLine = line
			prevKey = key
			count = 1
			first = false
			continue
		}
		if key == prevKey {
			count++
		} else {
			if err := flush(); err != nil { return err }
			prevLine = line
			prevKey = key
			count = 1
		}
	}
	if err := in.Err(); err != nil {
		return err
	}

	if err := flush(); err != nil { return err }
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]\n")
}

func main() {
	var opt options
	flag.BoolVar(&opt.count, "c", false, "prefix lines by the number of occurrences")
	flag.BoolVar(&opt.onlyDup, "d", false, "only print duplicate lines")
	flag.BoolVar(&opt.onlyUniq, "u", false, "only print unique lines")
	flag.BoolVar(&opt.ignoreCase, "i", false, "ignore differences in case when comparing")
	flag.IntVar(&opt.skipFields, "f", 0, "avoid comparing the first N fields")
	flag.IntVar(&opt.skipChars, "s", 0, "avoid comparing the first N characters (after fields)")
	flag.Usage = usage
	flag.Parse()

	if err := opt.validate(); err != nil {
		usage()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	args := flag.Args()
	var inFile, outFile string
	if len(args) >= 1 { inFile = args[0] }
	if len(args) >= 2 { outFile = args[1] }
	if len(args) > 2 {
		usage()
		fmt.Fprintln(os.Stderr, "Error: too many positional arguments")
		os.Exit(2)
	}

	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout
	var inF, outF *os.File
	var err error

	if inFile != "" {
		inF, err = os.Open(inFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening input:", err)
			os.Exit(1)
		}
		r = inF
		defer inF.Close()
	}
	if outFile != "" {
		outF, err = os.Create(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output:", err)
			if inF != nil { inF.Close() }
			os.Exit(1)
		}
		w = outF
		defer outF.Close()
	}

	if err := process(r, w, opt); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}