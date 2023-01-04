package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"archive/tar"
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return UnpackReader(bufio.NewReader(fd), dir)
}

// UnpackReader reads packed data using given reader and unpacks it to
// the given directory
func UnpackReader(r io.Reader, dir string) error {
	_, err := os.Stat(dir)

	if err != nil {
		return err
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dir, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())

			if err != nil {
				return err
			}

			continue
		}

		fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		if err != nil {
			return err
		}

		_, err = io.Copy(fd, tr)

		fd.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
