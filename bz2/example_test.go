package bz2

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func ExampleUnpack() {
	file := "file.bz2"
	err := Unpack(file, "/home/bob/data", Options{})

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}

func ExampleRead() {
	file := "file.bz2"
	fd, err := os.Open(file)

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	err = Read(fd, "/home/bob/data", Options{MaxReadLimit: 15 * 1024 * 1024})

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}
