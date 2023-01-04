package npck

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/npck/tar"
	"github.com/essentialkaos/npck/tbz"
	"github.com/essentialkaos/npck/tgz"
	"github.com/essentialkaos/npck/txz"
	"github.com/essentialkaos/npck/tzst"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func stub() {
	tar.Unpack("", "")
	tbz.Unpack("", "")
	tgz.Unpack("", "")
	txz.Unpack("", "")
	tzst.Unpack("", "")
}
