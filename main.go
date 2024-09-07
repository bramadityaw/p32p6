package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

type pixel struct {
	r, g, b uint8
}

type ppmImg struct {
	w, h uint64
	data []*pixel
}

func newImg(f *os.File) ppmImg {
	var p ppmImg

	size := make([]byte, 3)

	//TODO: Find a way to read until a certain rune; which in this case is ' '.
	f.Read(size)
	w, err := strconv.ParseUint(string(size)[1:],
		10, 64)
	f.Read(size)
	h, err := strconv.ParseUint(string(size)[1:], 10, 64)

	if err != nil {
		log.Fatal("Unknown file size!")
	}

	p.w = w
	p.h = h

	return p
}

func main() {
	file, err := os.Open("p3.ppm")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ptype := make([]byte, 2)
	file.Read(ptype)

	if string(ptype) != "P3" {
		fmt.Printf("%s is an invalid ASCII image PPM file!\nWrong magic number %s.", file.Name(), string(ptype))
		os.Exit(1)
	}

	p := newImg(file)
	fmt.Printf("image size:\n\tw: %v\n\th: %v\n", p.w, p.h)

	buf := bytes.Buffer{}
	buf.WriteString("P6 ")
}
