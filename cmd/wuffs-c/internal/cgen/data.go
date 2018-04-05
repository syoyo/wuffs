// Code generated by running "go generate". DO NOT EDIT.

// Copyright 2017 The Wuffs Authors.
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

const baseHeader = "" +
	"#ifndef WUFFS_BASE_HEADER_H\n#define WUFFS_BASE_HEADER_H\n\n// Copyright 2017 The Wuffs Authors.\n//\n// Licensed under the Apache License, Version 2.0 (the \"License\");\n// you may not use this file except in compliance with the License.\n// You may obtain a copy of the License at\n//\n//    https://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n// See the License for the specific language governing permissions and\n// limitations under the License.\n\n#include <stdbool.h>\n#include <stdint.h>\n#include <string.h>\n\n// Wuffs requires a word size of at least 32 bits because it assumes that\n// converting a u32 to usize will never overflow. For example, the size of a\n// decoded image is often represented, explicitly or implicitly in an image\n// file, as a u32, and it is convenient to compare that to a buffer size.\n//\n// Si" +
	"milarly, the word size is at most 64 bits because it assumes that\n// converting a usize to u64 will never overflow.\n#if __WORDSIZE < 32\n#error \"Wuffs requires a word size of at least 32 bits\"\n#elif __WORDSIZE > 64\n#error \"Wuffs requires a word size of at most 64 bits\"\n#endif\n\n// WUFFS_VERSION is the major.minor version number as a uint32. The major\n// number is the high 16 bits. The minor number is the low 16 bits.\n//\n// The intention is to bump the version number at least on every API / ABI\n// backwards incompatible change.\n//\n// For now, the API and ABI are simply unstable and can change at any time.\n//\n// TODO: don't hard code this in base-header.h.\n#define WUFFS_VERSION 0x00001\n\n// ---------------- Fundamentals\n\n// Flicks are a unit of time. One flick (frame-tick) is 1 / 705_600_000 of a\n// second. See https://github.com/OculusVR/Flicks\ntypedef uint64_t wuffs_base__flicks;\n\n#define WUFFS_BASE__FLICKS_PER_SECOND 705600000ULL\n\n// wuffs_base__rectangle is a rectangle on the integer grid. It contains all\n// p" +
	"oints (x, y) such that ((min_x <= x) && (x < max_x)) and likewise for y. It\n// is therefore empty if min_x >= max_x. There are multiple representations of\n// an empty rectangle.\n//\n// A value with all fields zero is a valid, empty rectangle.\n//\n// The X and Y axes increase right and down.\ntypedef struct {\n  uint32_t min_x;\n  uint32_t min_y;\n  uint32_t max_x;\n  uint32_t max_y;\n} wuffs_base__rectangle;\n\nstatic inline uint32_t wuffs_base__rectangle__width(wuffs_base__rectangle r) {\n  return r.max_x > r.min_x ? r.max_x - r.min_x : 0;\n}\n\nstatic inline uint32_t wuffs_base__rectangle__height(wuffs_base__rectangle r) {\n  return r.max_y > r.min_y ? r.max_y - r.min_y : 0;\n}\n\n// wuffs_base__slice_u8 is a 1-dimensional buffer.\n//\n// A value with all fields NULL or zero is a valid, empty slice.\ntypedef struct {\n  uint8_t* ptr;\n  size_t len;\n} wuffs_base__slice_u8;\n\n// wuffs_base__table_u8 is a 2-dimensional buffer.\n//\n// A value with all fields NULL or zero is a valid, empty table.\ntypedef struct {\n  uint8_t* ptr;\n  size_" +
	"t width;\n  size_t height;\n  size_t stride;\n} wuffs_base__table_u8;\n\n// ---------------- I/O\n\n// TODO: rename buf1 to io_buffer, writer1 to io_writer, etc.\n\n// wuffs_base__buf1 is a 1-dimensional buffer (a pointer and length), plus\n// additional indexes into that buffer, plus an opened / closed flag.\n//\n// A value with all fields NULL or zero is a valid, empty buffer.\ntypedef struct {\n  uint8_t* ptr;  // Pointer.\n  size_t len;    // Length.\n  size_t wi;     // Write index. Invariant: wi <= len.\n  size_t ri;     // Read  index. Invariant: ri <= wi.\n  bool closed;   // No further writes are expected.\n} wuffs_base__buf1;\n\n// wuffs_base__limit1 provides a limited view of a 1-dimensional byte stream:\n// its first N bytes. That N can be greater than a buffer's current read or\n// write capacity. N decreases naturally over time as bytes are read from or\n// written to the stream.\n//\n// A value with all fields NULL or zero is a valid, unlimited view.\ntypedef struct wuffs_base__limit1 {\n  uint64_t* ptr_to_len;           " +
	"  // Pointer to N.\n  struct wuffs_base__limit1* next;  // Linked list of limits.\n} wuffs_base__limit1;\n\ntypedef struct {\n  // TODO: move buf into private_impl? As it is, it looks like users can modify\n  // the buf field to point to a different buffer, which can turn the limit and\n  // mark fields into dangling pointers.\n  wuffs_base__buf1* buf;\n  // Do not access the private_impl's fields directly. There is no API/ABI\n  // compatibility or safety guarantee if you do so.\n  struct {\n    wuffs_base__limit1 limit;\n    uint8_t* mark;\n  } private_impl;\n} wuffs_base__reader1;\n\ntypedef struct {\n  // TODO: move buf into private_impl? As it is, it looks like users can modify\n  // the buf field to point to a different buffer, which can turn the limit and\n  // mark fields into dangling pointers.\n  wuffs_base__buf1* buf;\n  // Do not access the private_impl's fields directly. There is no API/ABI\n  // compatibility or safety guarantee if you do so.\n  struct {\n    wuffs_base__limit1 limit;\n    uint8_t* mark;\n  } private_impl" +
	";\n} wuffs_base__writer1;\n\n// ---------------- Images\n\n// wuffs_base__pixel_format encodes the format of the bytes that constitute an\n// image frame's pixel data. Its bits:\n//  - bit        31  is reserved.\n//  - bits 30 .. 28 encodes color (and channel order, in terms of memory).\n//  - bits 27 .. 26 are reserved.\n//  - bits 25 .. 24 encodes transparency.\n//  - bit        23 indicates big-endian/MSB-first (as opposed to little/LSB).\n//  - bit        22 indicates floating point (as opposed to integer).\n//  - bits 21 .. 20 are the number of planes, minus 1. Zero means packed.\n//  - bits 19 .. 16 encodes the number of bits (depth) in an index value.\n//                  Zero means direct, not palette-indexed.\n//  - bits 15 .. 12 encodes the number of bits (depth) in the 3rd channel.\n//  - bits 11 ..  8 encodes the number of bits (depth) in the 2nd channel.\n//  - bits  7 ..  4 encodes the number of bits (depth) in the 1st channel.\n//  - bits  3 ..  0 encodes the number of bits (depth) in the 0th channel.\n//\n// The " +
	"bit fields of a wuffs_base__pixel_format are not independent. For\n// example, the number of planes should not be greater than the number of\n// channels. Similarly, bits 15..4 are unused (and should be zero) if bits\n// 31..24 (color and transparency) together imply only 1 channel (gray, no\n// alpha) and floating point samples should mean a bit depth of 16, 32 or 64.\n//\n// Formats hold between 1 and 4 channels. For example: Y (1 channel: gray), YA\n// (2 channels: gray and alpha), BGR (3 channels: blue, green, red) or CMYK (4\n// channels: cyan, magenta, yellow, black).\n//\n// For direct formats with N > 1 channels, those channels can be laid out in\n// either 1 (packed) or N (planar) planes. For example, RGBA data is usually\n// packed, but YUV data is usually planar, due to chroma subsampling (for\n// details, see the wuffs_base__pixel_subsampling type). For indexed formats,\n// the palette (always 256 × 4 bytes) holds up to 4 packed bytes of color data\n// per index value, and there is only 1 plane (for the index)." +
	" The distance\n// between successive palette elements is always 4 bytes.\n//\n// The color field is encoded in 3 bits:\n//  - 0 means                 A (Alpha).\n//  - 1 means   Y       or   YA (Gray, Alpha).\n//  - 2 means BGR, BGRX or BGRA (Blue, Green, Red, X-padding or Alpha).\n//  - 3 means RGB, RGBX or RGBA (Red, Green, Blue, X-padding or Alpha).\n//  - 4 means YUV       or YUVA (Luma, Chroma-blue, Chroma-red, Alpha).\n//  - 5 means CMY       or CMYK (Cyan, Magenta, Yellow, Black).\n//  - all other values are reserved.\n//\n// In Wuffs, channels are given in memory order, regardless of endianness,\n// since the C type for the pixel data is an array of bytes, not an array of\n// uint32_t. For example, packed BGRA with 8 bits per channel means that the\n// bytes in memory are always Blue, Green, Red then Alpha. On big-endian\n// systems, that is the uint32_t 0xBBGGRRAA. On little-endian, 0xAARRGGBB.\n//\n// When the color field (3 bits) encodes multiple options, the transparency\n// field (2 bits) distinguishes them:\n//  - " +
	"0 means fully opaque, no extra channels\n//  - 1 means fully opaque, one extra channel (X or K, padding or black).\n//  - 2 means one extra alpha channel, other channels are non-premultiplied.\n//  - 3 means one extra alpha channel, other channels are     premultiplied.\n//\n// The zero wuffs_base__pixel_format value is an invalid pixel format, as it is\n// invalid to combine the zero color (alpha only) with the zero transparency.\n//\n// Bit depth is encoded in 4 bits:\n//  -  0 means the channel or index is unused.\n//  -  x means a bit depth of  x, for x in the range 1..8.\n//  -  9 means a bit depth of 10.\n//  - 10 means a bit depth of 12.\n//  - 11 means a bit depth of 16.\n//  - 12 means a bit depth of 24.\n//  - 13 means a bit depth of 32.\n//  - 14 means a bit depth of 48.\n//  - 15 means a bit depth of 64.\n//\n// For example, wuffs_base__pixel_format 0x3280BBBB is a natural format for\n// decoding a PNG image - network byte order (also known as big-endian),\n// packed, non-premultiplied alpha - that happens to be 16-bi" +
	"t-depth truecolor\n// with alpha (RGBA). In memory order:\n//\n//  ptr+0  ptr+1  ptr+2  ptr+3  ptr+4  ptr+5  ptr+6  ptr+7\n//  Rhi    Rlo    Ghi    Glo    Bhi    Blo    Ahi    Alo\n//\n// For example, the value wuffs_base__pixel_format 0x20000565 means BGR with no\n// alpha or padding, 5/6/5 bits for blue/green/red, packed 2 bytes per pixel,\n// laid out LSB-first in memory order:\n//\n//  ptr+0...........  ptr+1...........\n//  MSB          LSB  MSB          LSB\n//  G₂G₁G₀B₄B₃B₂B₁B₀  R₄R₃R₂R₁R₀G₅G₄G₃\n//\n// On little-endian systems (but not big-endian), this Wuffs pixel format value\n// (0x20000565) corresponds to the Cairo library's CAIRO_FORMAT_RGB16_565, the\n// SDL2 (Simple DirectMedia Layer 2) library's SDL_PIXELFORMAT_RGB565 and the\n// Skia library's kRGB_565_SkColorType. Note BGR in Wuffs versus RGB in the\n// other libraries.\n//\n// Regardless of endianness, this Wuffs pixel format value (0x20000565)\n// corresponds to the V4L2 (Video For Linux 2) library's V4L2_PIX_FMT_RGB565\n// and t" +
	"he Wayland-DRM library's WL_DRM_FORMAT_RGB565.\n//\n// Different software libraries name their pixel formats (and especially their\n// channel order) either according to memory layout or as bits of a native\n// integer type like uint32_t. The two conventions differ because of a system's\n// endianness. As mentioned earlier, Wuffs pixel formats are always in memory\n// order. More detail of other software libraries' naming conventions is in the\n// Pixel Format Guide at https://afrantzis.github.io/pixel-format-guide/\n//\n// Do not manipulate these bits directly; they are private implementation\n// details. Use methods such as wuffs_base__pixel_format__num_planes instead.\ntypedef uint32_t wuffs_base__pixel_format;\n\n// Common 8-bit-depth pixel formats. This list is not exhaustive; not all valid\n// wuffs_base__pixel_format values are present.\n\n#define WUFFS_BASE__PIXEL_FORMAT__INVALID ((wuffs_base__pixel_format)0x00000000)\n\n#define WUFFS_BASE__PIXEL_FORMAT__A ((wuffs_base__pixel_format)0x02000008)\n\n#define WUFFS_BASE__PIX" +
	"EL_FORMAT__Y ((wuffs_base__pixel_format)0x10000008)\n#define WUFFS_BASE__PIXEL_FORMAT__YA_NONPREMUL \\\n  ((wuffs_base__pixel_format)0x12000008)\n#define WUFFS_BASE__PIXEL_FORMAT__YA_PREMUL \\\n  ((wuffs_base__pixel_format)0x13000008)\n\n#define WUFFS_BASE__PIXEL_FORMAT__BGR ((wuffs_base__pixel_format)0x20000888)\n#define WUFFS_BASE__PIXEL_FORMAT__BGRX ((wuffs_base__pixel_format)0x21008888)\n#define WUFFS_BASE__PIXEL_FORMAT__BGRX_INDEXED \\\n  ((wuffs_base__pixel_format)0x21088888)\n#define WUFFS_BASE__PIXEL_FORMAT__BGRA_NONPREMUL \\\n  ((wuffs_base__pixel_format)0x22008888)\n#define WUFFS_BASE__PIXEL_FORMAT__BGRA_NONPREMUL_INDEXED \\\n  ((wuffs_base__pixel_format)0x22088888)\n#define WUFFS_BASE__PIXEL_FORMAT__BGRA_PREMUL \\\n  ((wuffs_base__pixel_format)0x23008888)\n\n#define WUFFS_BASE__PIXEL_FORMAT__RGB ((wuffs_base__pixel_format)0x30000888)\n#define WUFFS_BASE__PIXEL_FORMAT__RGBX ((wuffs_base__pixel_format)0x31008888)\n#define WUFFS_BASE__PIXEL_FORMAT__RGBX_INDEXED \\\n  ((wuffs_base__pixel_format)0x31088888)\n#define WUFFS_BASE__PI" +
	"XEL_FORMAT__RGBA_NONPREMUL \\\n  ((wuffs_base__pixel_format)0x32008888)\n#define WUFFS_BASE__PIXEL_FORMAT__RGBA_NONPREMUL_INDEXED \\\n  ((wuffs_base__pixel_format)0x32088888)\n#define WUFFS_BASE__PIXEL_FORMAT__RGBA_PREMUL \\\n  ((wuffs_base__pixel_format)0x33008888)\n\n#define WUFFS_BASE__PIXEL_FORMAT__YUV ((wuffs_base__pixel_format)0x40200888)\n#define WUFFS_BASE__PIXEL_FORMAT__YUVK ((wuffs_base__pixel_format)0x41308888)\n#define WUFFS_BASE__PIXEL_FORMAT__YUVA_NONPREMUL \\\n  ((wuffs_base__pixel_format)0x42308888)\n\n#define WUFFS_BASE__PIXEL_FORMAT__CMY ((wuffs_base__pixel_format)0x50200888)\n#define WUFFS_BASE__PIXEL_FORMAT__CMYK ((wuffs_base__pixel_format)0x51308888)\n\nstatic inline bool wuffs_base__pixel_format__is_valid(\n    wuffs_base__pixel_format f) {\n  return f != 0;\n}\n\nstatic inline bool wuffs_base__pixel_format__is_indexed(\n    wuffs_base__pixel_format f) {\n  return ((f >> 16) & 0x0F) != 0;\n}\n\n#define WUFFS_BASE__PIXEL_FORMAT__NUM_PLANES_MAX 4\n\nstatic inline uint32_t wuffs_base__pixel_format__num_planes(\n    wuffs_" +
	"base__pixel_format f) {\n  return f ? (((f >> 20) & 0x03) + 1) : 0;\n}\n\ntypedef struct {\n  wuffs_base__table_u8 planes[WUFFS_BASE__PIXEL_FORMAT__NUM_PLANES_MAX];\n} wuffs_base__pixel_buffer;\n\n// wuffs_base__pixel_subsampling encodes the mapping of pixel space coordinates\n// (x, y) to pixel buffer indices (i, j). That mapping can differ for each\n// plane p. For a depth of 8 bits (1 byte), the p'th plane's sample starts at\n// (planes[p].ptr + (j * planes[p].stride) + i).\n//\n// For packed pixel formats, the mapping is trivial: i = x and j = y. For\n// planar pixel formats, the mapping can differ due to chroma subsampling. For\n// example, consider a three plane YUV pixel format with 4:2:2 subsampling. For\n// the luma (Y) channel, there is one sample for every pixel, but for the\n// chroma (U, V) channels, there is one sample for every two pixels: pairs of\n// horizontally adjacent pixels form one macropixel, i = x / 2 and j == y. In\n// general, for a given p:\n//  - i = (x + bias_x) >> shift_x.\n//  - j = (y + bias_y) >>" +
	" shift_y.\n// where biases and shifts are in the range 0..3 and 0..2 respectively.\n//\n// In general, the biases will be zero after decoding an image. However, making\n// a sub-image may change the bias, since the (x, y) coordinates are relative\n// to the sub-image's top-left origin, but the backing pixel buffers were\n// created relative to the original image's origin.\n//\n// For each plane p, each of those four numbers (biases and shifts) are encoded\n// in two bits, which combine to form an 8 bit unsigned integer:\n//\n//  e_p = (bias_x << 6) | (shift_x << 4) | (bias_y << 2) | (shift_y << 0)\n//\n// Those e_p values (e_0 for the first plane, e_1 for the second plane, etc)\n// combine to form a wuffs_base__pixel_subsampling value:\n//\n//  pixsub = (e_3 << 24) | (e_2 << 16) | (e_1 << 8) | (e_0 << 0)\n//\n// Do not manipulate these bits directly; they are private implementation\n// details. Use methods such as wuffs_base__pixel_subsampling__bias_x instead.\ntypedef uint32_t wuffs_base__pixel_subsampling;\n\n#define WUFFS_BASE_" +
	"_PIXEL_SUBSAMPLING__NONE ((wuffs_base__pixel_subsampling)0);\n\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__444 \\\n  ((wuffs_base__pixel_subsampling)0x000000);\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__440 \\\n  ((wuffs_base__pixel_subsampling)0x010100);\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__422 \\\n  ((wuffs_base__pixel_subsampling)0x101000);\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__420 \\\n  ((wuffs_base__pixel_subsampling)0x111100);\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__411 \\\n  ((wuffs_base__pixel_subsampling)0x202000);\n#define WUFFS_BASE__PIXEL_SUBSAMPLING__410 \\\n  ((wuffs_base__pixel_subsampling)0x212100);\n\nstatic inline uint32_t wuffs_base__pixel_subsampling__bias_x(\n    wuffs_base__pixel_subsampling s,\n    uint32_t plane) {\n  uint32_t shift = ((plane & 0x03) * 8) + 6;\n  return (s >> shift) & 0x03;\n}\n\nstatic inline uint32_t wuffs_base__pixel_subsampling__shift_x(\n    wuffs_base__pixel_subsampling s,\n    uint32_t plane) {\n  uint32_t shift = ((plane & 0x03) * 8) + 4;\n  return (s >> shift) & 0x03;\n}\n\nstatic inline uint32_t" +
	" wuffs_base__pixel_subsampling__bias_y(\n    wuffs_base__pixel_subsampling s,\n    uint32_t plane) {\n  uint32_t shift = ((plane & 0x03) * 8) + 2;\n  return (s >> shift) & 0x03;\n}\n\nstatic inline uint32_t wuffs_base__pixel_subsampling__shift_y(\n    wuffs_base__pixel_subsampling s,\n    uint32_t plane) {\n  uint32_t shift = ((plane & 0x03) * 8) + 0;\n  return (s >> shift) & 0x03;\n}\n\ntypedef struct {\n  // Do not access the private_impl's fields directly. There is no API/ABI\n  // compatibility or safety guarantee if you do so.\n  struct {\n    wuffs_base__pixel_format pixfmt;\n    wuffs_base__pixel_subsampling pixsub;\n    uint32_t width;\n    uint32_t height;\n    uint32_t num_loops;\n  } private_impl;\n} wuffs_base__image_config;\n\nstatic inline void wuffs_base__image_config__initialize(\n    wuffs_base__image_config* c,\n    wuffs_base__pixel_format pixfmt,\n    wuffs_base__pixel_subsampling pixsub,\n    uint32_t width,\n    uint32_t height,\n    uint32_t num_loops) {\n  if (!c) {\n    return;\n  }\n  // TODO: move the check from wuffs" +
	"_base__image_config__is_valid here. Should\n  // this function return bool? An error type?\n  c->private_impl.pixfmt = pixfmt;\n  c->private_impl.pixsub = pixsub;\n  c->private_impl.width = width;\n  c->private_impl.height = height;\n  c->private_impl.num_loops = num_loops;\n}\n\nstatic inline void wuffs_base__image_config__invalidate(\n    wuffs_base__image_config* c) {\n  if (c) {\n    *c = ((wuffs_base__image_config){});\n  }\n}\n\nstatic inline bool wuffs_base__image_config__is_valid(\n    wuffs_base__image_config* c) {\n  if (!c || !c->private_impl.pixfmt) {\n    return false;\n  }\n  uint64_t wh =\n      ((uint64_t)c->private_impl.width) * ((uint64_t)c->private_impl.height);\n  // TODO: handle things other than 1 byte per pixel.\n  //\n  // TODO: move the check to wuffs_base__image_config__initialize.\n  return wh <= ((uint64_t)SIZE_MAX);\n}\n\nstatic inline wuffs_base__pixel_format wuffs_base__image_config__pixel_format(\n    wuffs_base__image_config* c) {\n  return c ? c->private_impl.pixfmt : 0;\n}\n\nstatic inline wuffs_base__pixel_" +
	"subsampling\nwuffs_base__image_config__pixel_subsampling(wuffs_base__image_config* c) {\n  return wuffs_base__image_config__is_valid(c) ? c->private_impl.pixsub : 0;\n}\n\nstatic inline uint32_t wuffs_base__image_config__width(\n    wuffs_base__image_config* c) {\n  return wuffs_base__image_config__is_valid(c) ? c->private_impl.width : 0;\n}\n\nstatic inline uint32_t wuffs_base__image_config__height(\n    wuffs_base__image_config* c) {\n  return wuffs_base__image_config__is_valid(c) ? c->private_impl.height : 0;\n}\n\nstatic inline uint32_t wuffs_base__image_config__num_loops(\n    wuffs_base__image_config* c) {\n  return wuffs_base__image_config__is_valid(c) ? c->private_impl.num_loops : 0;\n}\n\n// TODO: this is the right API for planar (not packed) pixbufs? Should it allow\n// decoding into a color model different from the format's intrinsic one? For\n// example, decoding a JPEG image straight to RGBA instead of to YCbCr?\nstatic inline size_t wuffs_base__image_config__pixbuf_size(\n    wuffs_base__image_config* c) {\n  if (wuffs_" +
	"base__image_config__is_valid(c)) {\n    uint64_t wh =\n        ((uint64_t)c->private_impl.width) * ((uint64_t)c->private_impl.height);\n    // TODO: handle things other than 1 byte per pixel.\n    return (size_t)wh;\n  }\n  return 0;\n}\n\ntypedef struct {\n  // Do not access the private_impl's fields directly. There is no API/ABI\n  // compatibility or safety guarantee if you do so.\n  struct {\n    wuffs_base__image_config config;\n    uint32_t loop_count;  // 0-based count of the current loop.\n    wuffs_base__pixel_buffer pixbuf;\n    // TODO: color spaces.\n    wuffs_base__rectangle dirty_rect;\n    wuffs_base__flicks duration;\n    uint8_t palette[1024];\n  } private_impl;\n} wuffs_base__image_buffer;\n\nstatic inline void wuffs_base__image_buffer__initialize(\n    wuffs_base__image_buffer* f,\n    wuffs_base__image_config config,\n    wuffs_base__pixel_buffer pixbuf) {\n  if (!f) {\n    return;\n  }\n  *f = ((wuffs_base__image_buffer){});\n  f->private_impl.config = config;\n  f->private_impl.pixbuf = pixbuf;\n}\n\nstatic inline void wu" +
	"ffs_base__image_buffer__update(\n    wuffs_base__image_buffer* f,\n    wuffs_base__rectangle dirty_rect,\n    wuffs_base__flicks duration,\n    uint8_t* palette_ptr,\n    size_t palette_len) {\n  if (!f) {\n    return;\n  }\n  f->private_impl.dirty_rect = dirty_rect;\n  f->private_impl.duration = duration;\n  if (palette_ptr) {\n    memmove(f->private_impl.palette, palette_ptr,\n            palette_len <= 1024 ? palette_len : 1024);\n  }\n}\n\n// wuffs_base__image_buffer__loop returns whether the image decoder should loop\n// back to the beginning of the animation, assuming that we've reached the end\n// of the encoded stream. If so, it increments f's count of the animation loops\n// played so far.\nstatic inline bool wuffs_base__image_buffer__loop(wuffs_base__image_buffer* f) {\n  if (!f) {\n    return false;\n  }\n  uint32_t n = f->private_impl.config.private_impl.num_loops;\n  if (n == 0) {\n    return true;\n  }\n  if (f->private_impl.loop_count < n - 1) {\n    f->private_impl.loop_count++;\n    return true;\n  }\n  return false;\n}\n\n// w" +
	"uffs_base__image_buffer__dirty_rect returns an upper bound for what part of\n// this frame's pixels differs from the previous frame.\nstatic inline wuffs_base__rectangle wuffs_base__image_buffer__dirty_rect(\n    wuffs_base__image_buffer* f) {\n  return f ? f->private_impl.dirty_rect : ((wuffs_base__rectangle){0});\n}\n\n// wuffs_base__image_buffer__duration returns the amount of time to display\n// this frame. Zero means to display forever - a still (non-animated) image.\nstatic inline wuffs_base__flicks wuffs_base__image_buffer__duration(\n    wuffs_base__image_buffer* f) {\n  return f ? f->private_impl.duration : 0;\n}\n\n// wuffs_base__image_buffer__palette returns the palette that the pixel data\n// can index. The backing array is inside f and has length 1024.\nstatic inline uint8_t* wuffs_base__image_buffer__palette(\n    wuffs_base__image_buffer* f) {\n  return f ? f->private_impl.palette : NULL;\n}\n\n#endif  // WUFFS_BASE_HEADER_H\n" +
	""

const baseImpl = "" +
	"// Copyright 2017 The Wuffs Authors.\n//\n// Licensed under the Apache License, Version 2.0 (the \"License\");\n// you may not use this file except in compliance with the License.\n// You may obtain a copy of the License at\n//\n//    https://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n// See the License for the specific language governing permissions and\n// limitations under the License.\n\n// wuffs_base__empty_struct is used when a Wuffs function returns an empty\n// struct. In C, if a function f returns void, you can't say \"x = f()\", but in\n// Wuffs, if a function g returns empty, you can say \"y = g()\".\ntypedef struct {\n} wuffs_base__empty_struct;\n\n#define WUFFS_BASE__IGNORE_POTENTIALLY_UNUSED_VARIABLE(x) (void)(x)\n\n// WUFFS_BASE__MAGIC is a magic number to check that initializers are called.\n// It's not foolp" +
	"roof, given C doesn't automatically zero memory before use,\n// but it should catch 99.99% of cases.\n//\n// Its (non-zero) value is arbitrary, based on md5sum(\"wuffs\").\n#define WUFFS_BASE__MAGIC 0x3CCB6C71U\n\n// WUFFS_BASE__ALREADY_ZEROED is passed from a container struct's initializer\n// to a containee struct's initializer when the container has already zeroed\n// the containee's memory.\n//\n// Its (non-zero) value is arbitrary, based on md5sum(\"zeroed\").\n#define WUFFS_BASE__ALREADY_ZEROED 0x68602EF1U\n\n// Denote intentional fallthroughs for -Wimplicit-fallthrough.\n//\n// The order matters here. Clang also defines \"__GNUC__\".\n#if defined(__clang__) && __cplusplus >= 201103L\n#define WUFFS_BASE__FALLTHROUGH [[clang::fallthrough]]\n#elif !defined(__clang__) && defined(__GNUC__) && (__GNUC__ >= 7)\n#define WUFFS_BASE__FALLTHROUGH __attribute__((fallthrough))\n#else\n#define WUFFS_BASE__FALLTHROUGH\n#endif\n\n// Use switch cases for coroutine suspension points, similar to the technique\n// in https://www.chiark.greenend.org.uk/" +
	"~sgtatham/coroutines.html\n//\n// We use trivial macros instead of an explicit assignment and case statement\n// so that clang-format doesn't get confused by the unusual \"case\"s.\n#define WUFFS_BASE__COROUTINE_SUSPENSION_POINT_0 case 0:;\n#define WUFFS_BASE__COROUTINE_SUSPENSION_POINT(n) \\\n  coro_susp_point = n;                            \\\n  WUFFS_BASE__FALLTHROUGH;                        \\\n  case n:;\n\n#define WUFFS_BASE__COROUTINE_SUSPENSION_POINT_MAYBE_SUSPEND(n) \\\n  if (status < 0) {                                             \\\n    goto exit;                                                  \\\n  } else if (status == 0) {                                     \\\n    goto ok;                                                    \\\n  }                                                             \\\n  coro_susp_point = n;                                          \\\n  goto suspend;                                                 \\\n  case n:;\n\n// Clang also defines \"__GNUC__\".\n#if defined(__GNUC__)\n#define WUFFS_BASE__LIKELY" +
	"(expr) (__builtin_expect(!!(expr), 1))\n#define WUFFS_BASE__UNLIKELY(expr) (__builtin_expect(!!(expr), 0))\n#else\n#define WUFFS_BASE__LIKELY(expr) (expr)\n#define WUFFS_BASE__UNLIKELY(expr) (expr)\n#endif\n\n// Uncomment this #include for printf-debugging.\n// #include <stdio.h>\n\n// ---------------- Static Inline Functions\n//\n// The helpers below are functions, instead of macros, because their arguments\n// can be an expression that we shouldn't evaluate more than once.\n//\n// They are in base-impl.h and hence copy/pasted into every generated C file,\n// instead of being in some \"base.c\" file, since a design goal is that users of\n// the generated C code can often just #include a single .c file, such as\n// \"gif.c\", without having to additionally include or otherwise build and link\n// a \"base.c\" file.\n//\n// They are static, so that linking multiple wuffs .o files won't complain about\n// duplicate function definitions.\n//\n// They are explicitly marked inline, even if modern compilers don't use the\n// inline attribute to g" +
	"uide optimizations such as inlining, to avoid the\n// -Wunused-function warning, and we like to compile with -Wall -Werror.\n\nstatic inline uint16_t wuffs_base__load_u16be(uint8_t* p) {\n  return ((uint16_t)(p[0]) << 8) | ((uint16_t)(p[1]) << 0);\n}\n\nstatic inline uint16_t wuffs_base__load_u16le(uint8_t* p) {\n  return ((uint16_t)(p[0]) << 0) | ((uint16_t)(p[1]) << 8);\n}\n\nstatic inline uint32_t wuffs_base__load_u32be(uint8_t* p) {\n  return ((uint32_t)(p[0]) << 24) | ((uint32_t)(p[1]) << 16) |\n         ((uint32_t)(p[2]) << 8) | ((uint32_t)(p[3]) << 0);\n}\n\nstatic inline uint32_t wuffs_base__load_u32le(uint8_t* p) {\n  return ((uint32_t)(p[0]) << 0) | ((uint32_t)(p[1]) << 8) |\n         ((uint32_t)(p[2]) << 16) | ((uint32_t)(p[3]) << 24);\n}\n\nstatic inline wuffs_base__slice_u8 wuffs_base__slice_u8__subslice_i(\n    wuffs_base__slice_u8 s,\n    uint64_t i) {\n  if ((i <= SIZE_MAX) && (i <= s.len)) {\n    return ((wuffs_base__slice_u8){\n        .ptr = s.ptr + i,\n        .len = s.len - i,\n    });\n  }\n  return ((wuffs_base__sli" +
	"ce_u8){});\n}\n\nstatic inline wuffs_base__slice_u8 wuffs_base__slice_u8__subslice_j(\n    wuffs_base__slice_u8 s,\n    uint64_t j) {\n  if ((j <= SIZE_MAX) && (j <= s.len)) {\n    return ((wuffs_base__slice_u8){.ptr = s.ptr, .len = j});\n  }\n  return ((wuffs_base__slice_u8){});\n}\n\nstatic inline wuffs_base__slice_u8 wuffs_base__slice_u8__subslice_ij(\n    wuffs_base__slice_u8 s,\n    uint64_t i,\n    uint64_t j) {\n  if ((i <= j) && (j <= SIZE_MAX) && (j <= s.len)) {\n    return ((wuffs_base__slice_u8){\n        .ptr = s.ptr + i,\n        .len = j - i,\n    });\n  }\n  return ((wuffs_base__slice_u8){});\n}\n\n// wuffs_base__slice_u8__prefix returns up to the first up_to bytes of s.\nstatic inline wuffs_base__slice_u8 wuffs_base__slice_u8__prefix(\n    wuffs_base__slice_u8 s,\n    uint64_t up_to) {\n  if ((uint64_t)(s.len) > up_to) {\n    s.len = up_to;\n  }\n  return s;\n}\n\n// wuffs_base__slice_u8__suffix returns up to the last up_to bytes of s.\nstatic inline wuffs_base__slice_u8 wuffs_base__slice_u8_suffix(\n    wuffs_base__slice_u8 s,\n " +
	"   uint64_t up_to) {\n  if ((uint64_t)(s.len) > up_to) {\n    s.ptr += (uint64_t)(s.len) - up_to;\n    s.len = up_to;\n  }\n  return s;\n}\n\n// wuffs_base__slice_u8__copy_from_slice calls memmove(dst.ptr, src.ptr,\n// length) where length is the minimum of dst.len and src.len.\n//\n// Passing a wuffs_base__slice_u8 with all fields NULL or zero (a valid, empty\n// slice) is valid and results in a no-op.\nstatic inline uint64_t wuffs_base__slice_u8__copy_from_slice(\n    wuffs_base__slice_u8 dst,\n    wuffs_base__slice_u8 src) {\n  size_t length = dst.len < src.len ? dst.len : src.len;\n  if (length > 0) {\n    memmove(dst.ptr, src.ptr, length);\n  }\n  return length;\n}\n\nstatic inline uint32_t wuffs_base__writer1__copy_from_history32(\n    uint8_t** ptr_ptr,\n    uint8_t* start,  // May be NULL, meaning an unmarked writer1.\n    uint8_t* end,\n    uint32_t distance,\n    uint32_t length) {\n  if (!start || !distance) {\n    return 0;\n  }\n  uint8_t* ptr = *ptr_ptr;\n  if ((size_t)(ptr - start) < (size_t)(distance)) {\n    return 0;\n  }\n  s" +
	"tart = ptr - distance;\n  size_t n = end - ptr;\n  if ((size_t)(length) > n) {\n    length = n;\n  } else {\n    n = length;\n  }\n  // TODO: unrolling by 3 seems best for the std/deflate benchmarks, but that\n  // is mostly because 3 is the minimum length for the deflate format. This\n  // function implementation shouldn't overfit to that one format. Perhaps the\n  // copy_from_history32 Wuffs method should also take an unroll hint argument,\n  // and the cgen can look if that argument is the constant expression '3'.\n  //\n  // See also wuffs_base__writer1__copy_from_history32__bco below.\n  //\n  // Alternatively, or additionally, have a sloppy_copy_from_history32 method\n  // that copies 8 bytes at a time, possibly writing more than length bytes?\n  for (; n >= 3; n -= 3) {\n    *ptr++ = *start++;\n    *ptr++ = *start++;\n    *ptr++ = *start++;\n  }\n  for (; n; n--) {\n    *ptr++ = *start++;\n  }\n  *ptr_ptr = ptr;\n  return length;\n}\n\n// wuffs_base__writer1__copy_from_history32__bco is a Bounds Check Optimized\n// version of the " +
	"wuffs_base__writer1__copy_from_history32 function above. The\n// caller needs to prove that:\n//  - start    != NULL\n//  - distance >  0\n//  - distance <= (*ptr_ptr - start)\n//  - length   <= (end      - *ptr_ptr)\nstatic inline uint32_t wuffs_base__writer1__copy_from_history32__bco(\n    uint8_t** ptr_ptr,\n    uint8_t* start,\n    uint8_t* end,\n    uint32_t distance,\n    uint32_t length) {\n  uint8_t* ptr = *ptr_ptr;\n  start = ptr - distance;\n  uint32_t n = length;\n  for (; n >= 3; n -= 3) {\n    *ptr++ = *start++;\n    *ptr++ = *start++;\n    *ptr++ = *start++;\n  }\n  for (; n; n--) {\n    *ptr++ = *start++;\n  }\n  *ptr_ptr = ptr;\n  return length;\n}\n\nstatic inline uint32_t wuffs_base__writer1__copy_from_reader32(\n    uint8_t** ptr_wptr,\n    uint8_t* wend,\n    uint8_t** ptr_rptr,\n    uint8_t* rend,\n    uint32_t length) {\n  uint8_t* wptr = *ptr_wptr;\n  size_t n = length;\n  if (n > wend - wptr) {\n    n = wend - wptr;\n  }\n  uint8_t* rptr = *ptr_rptr;\n  if (n > rend - rptr) {\n    n = rend - rptr;\n  }\n  if (n > 0) {\n    memm" +
	"ove(wptr, rptr, n);\n    *ptr_wptr += n;\n    *ptr_rptr += n;\n  }\n  return n;\n}\n\nstatic inline uint64_t wuffs_base__writer1__copy_from_slice(\n    uint8_t** ptr_wptr,\n    uint8_t* wend,\n    wuffs_base__slice_u8 src) {\n  uint8_t* wptr = *ptr_wptr;\n  size_t n = src.len;\n  if (n > wend - wptr) {\n    n = wend - wptr;\n  }\n  if (n > 0) {\n    memmove(wptr, src.ptr, n);\n    *ptr_wptr += n;\n  }\n  return n;\n}\n\nstatic inline uint32_t wuffs_base__writer1__copy_from_slice32(\n    uint8_t** ptr_wptr,\n    uint8_t* wend,\n    wuffs_base__slice_u8 src,\n    uint32_t length) {\n  uint8_t* wptr = *ptr_wptr;\n  size_t n = src.len;\n  if (n > length) {\n    n = length;\n  }\n  if (n > wend - wptr) {\n    n = wend - wptr;\n  }\n  if (n > 0) {\n    memmove(wptr, src.ptr, n);\n    *ptr_wptr += n;\n  }\n  return n;\n}\n\n// Note that the *__limit and *__mark methods are private (in base-impl.h) not\n// public (in base-header.h). We assume that, at the boundary between user code\n// and Wuffs code, the reader1 and writer1's private_impl fields (including\n// " +
	"limit and mark) are NULL. Otherwise, some internal assumptions break down.\n// For example, limits could be represented as pointers, even though\n// conceptually they are counts, but that pointer-to-count correspondence\n// becomes invalid if a buffer is re-used (e.g. on resuming a coroutine).\n//\n// Admittedly, some of the Wuffs test code calls these methods, but that test\n// code is still Wuffs code, not user code. Other Wuffs test code modifies\n// private_impl fields directly.\n\nstatic inline wuffs_base__reader1 wuffs_base__reader1__limit(\n    wuffs_base__reader1* o,\n    uint64_t* ptr_to_len) {\n  wuffs_base__reader1 ret = *o;\n  ret.private_impl.limit.ptr_to_len = ptr_to_len;\n  ret.private_impl.limit.next = &o->private_impl.limit;\n  return ret;\n}\n\nstatic inline wuffs_base__empty_struct wuffs_base__reader1__mark(\n    wuffs_base__reader1* o,\n    uint8_t* mark) {\n  o->private_impl.mark = mark;\n  return ((wuffs_base__empty_struct){});\n}\n\n// TODO: static inline wuffs_base__writer1 wuffs_base__writer1__limit()\n\nstatic" +
	" inline wuffs_base__empty_struct wuffs_base__writer1__mark(\n    wuffs_base__writer1* o,\n    uint8_t* mark) {\n  o->private_impl.mark = mark;\n  return ((wuffs_base__empty_struct){});\n}\n" +
	""

type template_args_short_read struct {
	PKGPREFIX string
	name      string
}

func template_short_read(b *buffer, args template_args_short_read) error {
	b.printf("short_read_%s:\nif (a_%s.buf && a_%s.buf->closed &&\n!a_%s.private_impl.limit.ptr_to_len) {\nstatus = %sERROR_UNEXPECTED_EOF;\ngoto exit;\n}\nstatus = %sSUSPENSION_SHORT_READ;\ngoto suspend;\n",
		args.name,
		args.name,
		args.name,
		args.name,
		args.PKGPREFIX,
		args.PKGPREFIX,
	)
	return nil
}
