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
	Counter int64
}

func NewCounterReader(r io.Reader) *CounterReader {
	return &CounterReader{
		r:       r,
		Counter: 0,
	}
}

func (cr *CounterReader) Read(b []byte) (n int, err error) {
	n, err = cr.r.Read(b)
	cr.Counter += int64(n)
	return
}

func main() {
	inputFileName := flag.String("f", "", "Input file name")
	blockSize := flag.Int64("b", 512, "Size of disc block (cluster)")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	if *inputFileName == "" {
		fmt.Printf("Usage:\n%s -f inputfile [-v -b blockSize]\n", os.Args[0])
		return
	}

	inputFile, err := os.Open(*inputFileName)
	if err != nil {
		fmt.Printf("Can't open input file: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		err := inputFile.Close()
		if err != nil {
			fmt.Printf("Can't close input file: %s\n", err)
			os.Exit(1)
		}
	}()

	address := int64(0)

	for {
		counterReader := NewCounterReader(inputFile)
		_, err = jpeg.Decode(counterReader)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			address += *blockSize
			_, err := inputFile.Seek(address, io.SeekStart)
			if err != nil {
				fmt.Printf("Can't seek to next cluster: %s\n", err)
				os.Exit(1)
			}
		} else {
			numWholeClusters := counterReader.Counter / (*blockSize)

			partOfCluster := counterReader.Counter - numWholeClusters*(*blockSize)

			if partOfCluster > 0 {
				numWholeClusters++
			}

			outFile, err := os.Create(strconv.FormatInt(address, 10) + ".jpg")
			if err != nil {
				fmt.Printf("Can't create JPEG file: %s\n", err)
				os.Exit(1)
			}

			_, err = inputFile.Seek(address, io.SeekStart)
			if err != nil {
				fmt.Printf("Can't seek back to file start cluster: %s\n", err)
				os.Exit(1)
			}

			_, err = io.CopyN(outFile, inputFile, int64(counterReader.Counter))

			if err != nil {
				fmt.Printf("Can't save JPEG data to file: %s\n", err)
				os.Exit(1)
			}

			err = outFile.Close()
			if err != nil {
				fmt.Printf("Can't close JPEG file: %s\n", err)
				os.Exit(1)
			}

			if *verbose {
				fmt.Printf("JPEG file found at address 0x%x. Size: %d bytes; Num clusters: %d\n", address, counterReader.Counter, numWholeClusters)
				os.Exit(1)
			}

			address += numWholeClusters * (*blockSize)

			_, err = inputFile.Seek(address, io.SeekStart)
			if err != nil {
				fmt.Printf("Can't seek to next cluster: %s\n", err)
				os.Exit(1)
			}
		}
	}
}
