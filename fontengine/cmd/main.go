package main

//go:generate go run main.go

import (
	"bytes"
	"fmt"
	"os"

	"github.com/cdvelop/docpdf/fontengine"
	//"runtime/debug"
)

func main() {

	lenarg := len(os.Args)
	if lenarg < 5 {
		echoUsage()
		return
	}
	fmt.Println("Deprecated: No longer need to create font maps!!!!")
	i := 1
	encoding := os.Args[i+0]
	mappath := os.Args[i+1]
	fontpath := os.Args[i+2]
	outputpath := os.Args[i+3]

	fmk := fontengine.NewFontEngine(func(filename string, data []byte) error {
		os.MkdirAll(outputpath, os.ModePerm)
		f, err := os.Create(outputpath + "/" + filename)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(data)
		return err
	})
	err := fmk.MakeFont(fontpath, mappath, encoding, outputpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n\n", err.Error())
		echoUsage()
		return
	}

	//print result
	results := fmk.GetResults()
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Printf("Finish.\n")
}

func echoUsage() {
	var buff bytes.Buffer
	buff.WriteString("fontmaker is tool for making font file to use with docpdf.\n")
	buff.WriteString("\nUsage:\n")
	buff.WriteString("\tfontmaker encoding map_folder font_file output_folder\n")
	buff.WriteString("\nExample:\n")
	buff.WriteString("\tfontmaker cp874 /gopath/github.com/cdvelop/docpdf/fontmaker/map  ../ttf/Loma.ttf ./tmp\n")
	buff.WriteString("\n")
	fmt.Print(buff.String())
}
