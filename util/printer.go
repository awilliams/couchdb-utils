package util

import (
	"fmt"
	"io"
	"os"
)

type PrettyPrinter interface {
	PP(printer Printer)
}

type Printer interface {
	Print(format string, args ...interface{})
}

func PrettyPrint(printers ...PrettyPrinter) {
	for _, printer := range printers {
		printer.PP(output)
	}
}

func PrintError(err error) {
	errorOutput.Print("Error: %s", err.Error())
}

type printer struct {
	writer *io.Writer
}

func (p printer) Print(format string, args ...interface{}) {
	fmt.Fprintf(*p.writer, format+"\n", args...)
}

var out io.Writer = io.Writer(os.Stdout)
var err io.Writer = io.Writer(os.Stderr)
var output printer = printer{&out}
var errorOutput printer = printer{&err}
