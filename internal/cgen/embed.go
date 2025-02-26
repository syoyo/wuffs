// Copyright 2021 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cgen

import (
	_ "embed"
	"strings"
)

// EmbeddedString holds hand-written C/C++ code, exposed by this Go package by
// a "go:embed" directive.
type EmbeddedString string

// Trim removes the leading "Copyright etc" boilerplate.
func (e EmbeddedString) Trim() string {
	s := string(e)
	if strings.HasPrefix(s, "// Copyright ") {
		if i := strings.Index(s, "\n\n"); i >= 0 {
			s = s[i+2:]
		}
	}
	return s
}

// ----

//go:embed base/all-impl.c
var embedBaseAllImplC EmbeddedString

// ----

//go:embed base/fundamental-private.h
var embedBaseFundamentalPrivateH EmbeddedString

//go:embed base/fundamental-public.h
var embedBaseFundamentalPublicH EmbeddedString

//go:embed base/memory-private.h
var embedBaseMemoryPrivateH EmbeddedString

//go:embed base/memory-public.h
var embedBaseMemoryPublicH EmbeddedString

//go:embed base/image-private.h
var embedBaseImagePrivateH EmbeddedString

//go:embed base/image-public.h
var embedBaseImagePublicH EmbeddedString

//go:embed base/io-private.h
var embedBaseIOPrivateH EmbeddedString

//go:embed base/io-public.h
var embedBaseIOPublicH EmbeddedString

//go:embed base/range-private.h
var embedBaseRangePrivateH EmbeddedString

//go:embed base/range-public.h
var embedBaseRangePublicH EmbeddedString

//go:embed base/strconv-private.h
var embedBaseStrConvPrivateH EmbeddedString

//go:embed base/strconv-public.h
var embedBaseStrConvPublicH EmbeddedString

//go:embed base/token-private.h
var embedBaseTokenPrivateH EmbeddedString

//go:embed base/token-public.h
var embedBaseTokenPublicH EmbeddedString

// ----

//go:embed base/floatconv-submodule-code.c
var embedBaseFloatConvSubmoduleCodeC EmbeddedString

//go:embed base/floatconv-submodule-data.c
var embedBaseFloatConvSubmoduleDataC EmbeddedString

//go:embed base/intconv-submodule.c
var embedBaseIntConvSubmoduleC EmbeddedString

//go:embed base/magic-submodule.c
var embedBaseMagicSubmoduleC EmbeddedString

//go:embed base/pixconv-submodule-regular.c
var embedBasePixConvSubmoduleRegularC EmbeddedString

//go:embed base/pixconv-submodule-ycck.c
var embedBasePixConvSubmoduleYcckC EmbeddedString

//go:embed base/utf8-submodule.c
var embedBaseUTF8SubmoduleC EmbeddedString

// ----

//go:embed auxiliary/base.cc
var EmbeddedString_AuxBaseCc EmbeddedString

//go:embed auxiliary/base.hh
var EmbeddedString_AuxBaseHh EmbeddedString

//go:embed auxiliary/cbor.cc
var embedAuxCborCc EmbeddedString

//go:embed auxiliary/cbor.hh
var embedAuxCborHh EmbeddedString

//go:embed auxiliary/image.cc
var embedAuxImageCc EmbeddedString

//go:embed auxiliary/image.hh
var embedAuxImageHh EmbeddedString

//go:embed auxiliary/json.cc
var embedAuxJsonCc EmbeddedString

//go:embed auxiliary/json.hh
var embedAuxJsonHh EmbeddedString

var EmbeddedStrings_AuxNonBaseCcFiles = []EmbeddedString{
	embedAuxCborCc,
	embedAuxImageCc,
	embedAuxJsonCc,
}

var EmbeddedStrings_AuxNonBaseHhFiles = []EmbeddedString{
	embedAuxCborHh,
	embedAuxImageHh,
	embedAuxJsonHh,
}
