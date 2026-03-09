package tlz4

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"

	"github.com/essentialkaos/npck/tar"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func ExampleUnpack() {
	file := "file.tlz4"
	err := Unpack(file, "/home/bob/data", tar.DefaultOptions)

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}

func ExampleRead() {
	file := "file.tlz4"
	fd, err := os.Open(file)

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	err = Read(fd, "/home/bob/data", tar.DefaultOptions)

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}
