package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type pixel struct {
	r, g, b uint8
}

type ppmImg struct {
	w, h uint64
	max  uint64
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

func newImg(rd *bufio.Reader) (ppmImg, error) {
	var p ppmImg

	w, err := parseU64(rd)
	h, err := parseU64(rd)

	if err != nil {
		fmt.Println("Unknown file size!")
		os.Exit(1)
	}

	p.w = w
	p.h = h

	m, err := parseU64(rd)
	if err != nil {
		fmt.Println("Maximum value cannot be determined!")
		os.Exit(1)
	}
	p.max = m

	for {
		px := new(pixel)

		r, _ := parseU8(rd)
		g, _ := parseU8(rd)
		b, err := parseU8(rd)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if (r+g+b)/3 > uint8(p.max) {
			fmt.Println("Pixel too large!")
			continue
		}

		px.r = r
		px.g = g
		px.b = b

		p.data = append(p.data, px)
	}

	if len(p.data) != int(p.w*p.h) {
		return p, fmt.Errorf("Size of image not as promised!")
	}

	return p, nil
}

func writeImg(p ppmImg, filename string) error {
	var data []byte

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("P6 %d %d %d ", p.w, p.h, p.max)

	file.WriteString(header)

	for _, e := range p.data {
		data = append(data, e.r)
		data = append(data, e.g)
		data = append(data, e.b)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
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

	p, err := newImg(f_reader)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("image size:\n\tw: %v\n\th: %v\n", p.w, p.h)
	fmt.Printf("maximum value for pixel: %d\n", p.max)
	fmt.Printf("amount of pixels: %d\n", len(p.data))

	err = writeImg(p, "p6.ppm")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
