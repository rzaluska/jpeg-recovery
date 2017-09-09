package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"strconv"
)

type CounterReader struct {
	r       io.Reader
	Counter int
}

func NewCounterReader(r io.Reader) *CounterReader {
	return &CounterReader{
		r:       r,
		Counter: 0,
	}
}

func (cr *CounterReader) Read(b []byte) (n int, err error) {
	n, err = cr.r.Read(b)
	cr.Counter += n
	return
}

func main() {
	inputFileName := flag.String("f", "", "Input file name")
	blockSize := flag.Int("b", 512, "Size of block (cluster) of allocation")
	flag.Parse()

	if *inputFileName == "" {
		fmt.Printf("Usage:\n%s -f inputfile [-b blockSize]\n", os.Args[0])
		return
	}

	f, err := os.Open(*inputFileName)
	if err != nil {
		panic("Can't open input file")
	}
	defer f.Close()

	address := 0

	cluserSize := *blockSize

	for {
		cluserReader := NewCounterReader(f)
		_, err = jpeg.Decode(cluserReader)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			address += cluserSize
			f.Seek(int64(address), os.SEEK_SET)
		} else {
			numWholeClusters := cluserReader.Counter / cluserSize

			partOfCluster := cluserReader.Counter - numWholeClusters*cluserSize

			if partOfCluster > 0 {
				numWholeClusters++
			}

			outFile, err := os.Create(strconv.Itoa(address) + ".jpg")
			if err != nil {
				panic(err)
			}

			f.Seek(int64(address), os.SEEK_SET)

			_, err = io.CopyN(outFile, f, int64(cluserReader.Counter))

			if err != nil {
				panic(err)
			}

			outFile.Close()

			address += numWholeClusters * cluserSize

			f.Seek(int64(address), os.SEEK_SET)
			fmt.Printf("JPEG found at 0x%x. Size: %d bytes; Num clusters: %d\n", address, cluserReader.Counter, numWholeClusters)
		}
	}
}
