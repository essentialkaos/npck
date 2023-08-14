package gz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func ExampleUnpack() {
	file := "file.gz"
	err := Unpack(file, "/home/bob/data")

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}

func ExampleRead() {
	file := "file.gz"
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	err = Read(fd, "/home/bob/data")

	if err != nil {
		fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
		return
	}

	fmt.Printf("File %s successfully unpacked!\n", file)
}
