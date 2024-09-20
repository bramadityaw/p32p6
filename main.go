package main

import (
	"bufio"
	"fmt"
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

func parseU8(r *bufio.Reader) (uint8, error) {
	n, err := parseU64(r)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}

func parseU64(r *bufio.Reader) (uint64, error) {
	b, err := r.ReadBytes(' ')
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(string(b[:len(b)-1]), 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func newImg(r *bufio.Reader) ppmImg {
	var p ppmImg

	w, err := parseU64(r)
	h, err := parseU64(r)

	if err != nil {
		fmt.Println("Unknown file size!")
		os.Exit(1)
	}

	p.w = w
	p.h = h

	return p
}

func main() {
	file, err := os.Open("p3.ppm")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	f_reader := bufio.NewReader(file)
	b, err := f_reader.ReadBytes(' ')

	ptype := string(b[:len(b)-1])

	if string(ptype) != "P3" {
		fmt.Printf("%s is an invalid ASCII image PPM file!\nWrong magic number %s.", file.Name(), string(ptype))
		os.Exit(1)
	}

	p := newImg(f_reader)
	fmt.Printf("image size:\n\tw: %v\n\th: %v\n", p.w, p.h)
	fmt.Printf("amount of pixels: %d\n", len(p.data))
}
