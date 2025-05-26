package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicode00/ics-proxy/internal/parser"
)

func main() {
	iPath := os.Args[1]

	iFile, err := os.Open(iPath)
	if err != nil {
		log.Fatal("failed to open " + iPath)
	}
	defer iFile.Close()

	ext := filepath.Ext(iPath)
	oPath := strings.TrimSuffix(iPath, ext) + "_output" + ext
	oFile, err := os.Create(oPath)
	if err != nil {
		log.Fatal("failed to create " + oPath)
	}
	defer oFile.Close()

	r := bufio.NewReader(iFile)
	w := bufio.NewWriter(oFile)

	parser.Trim(r, w)
}
