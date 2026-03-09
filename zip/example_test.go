package zip

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
	file := "file.zip"
	err := Unpack(file, "/home/bob/data")

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}

func ExampleRead() {
	file := "file.zip"

	fi, err := os.Stat(file)

	if err != nil {
		fmt.Printf("Error: Can't check file %s stat: %v\n", file, err)
		return
	}

	fd, err := os.Open(file)

	if err != nil {
		fmt.Printf("Error: Can't open file %s: %v\n", file, err)
		return
	}

	err = Read(fd, fi.Size(), "/home/bob/data")

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}
