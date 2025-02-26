// Copyright 2023 The Wuffs Authors.
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

// ----------------

/*
This test program is typically run indirectly, by the "wuffs test" or "wuffs
bench" commands. These commands take an optional "-mimic" flag to check that
Wuffs' output mimics (i.e. exactly matches) other libraries' output, such as
giflib for GIF, libpng for PNG, etc.

To manually run this test:

for CC in clang gcc; do
  $CC -std=c99 -Wall -Werror netpbm.c && ./a.out
  rm -f a.out
done

Each edition should print "PASS", amongst other information, and exit(0).

Add the "wuffs mimic cflags" (everything after the colon below) to the C
compiler flags (after the .c file) to run the mimic tests.

To manually run the benchmarks, replace "-Wall -Werror" with "-O3" and replace
the first "./a.out" with "./a.out -bench". Combine these changes with the
"wuffs mimic cflags" to run the mimic benchmarks.
*/

// ¿ wuffs mimic cflags: -DWUFFS_MIMIC

// Wuffs ships as a "single file C library" or "header file library" as per
// https://github.com/nothings/stb/blob/master/docs/stb_howto.txt
//
// To use that single file as a "foo.c"-like implementation, instead of a
// "foo.h"-like header, #define WUFFS_IMPLEMENTATION before #include'ing or
// compiling it.
#define WUFFS_IMPLEMENTATION

// Defining the WUFFS_CONFIG__MODULE* macros are optional, but it lets users of
// release/c/etc.c choose which parts of Wuffs to build. That file contains the
// entire Wuffs standard library, implementing a variety of codecs and file
// formats. Without this macro definition, an optimizing compiler or linker may
// very well discard Wuffs code for unused codecs, but listing the Wuffs
// modules we use makes that process explicit. Preprocessing means that such
// code simply isn't compiled.
#define WUFFS_CONFIG__MODULES
#define WUFFS_CONFIG__MODULE__BASE
#define WUFFS_CONFIG__MODULE__NETPBM

// If building this program in an environment that doesn't easily accommodate
// relative includes, you can use the script/inline-c-relative-includes.go
// program to generate a stand-alone C file.
#include "../../../release/c/wuffs-unsupported-snapshot.c"
#include "../testlib/testlib.c"
#ifdef WUFFS_MIMIC
// No mimic library.
#endif

// ---------------- Netpbm Tests

const char*  //
test_wuffs_netpbm_decode_interface() {
  CHECK_FOCUS(__func__);
  wuffs_netpbm__decoder dec;
  CHECK_STATUS("initialize",
               wuffs_netpbm__decoder__initialize(
                   &dec, sizeof dec, WUFFS_VERSION,
                   WUFFS_INITIALIZE__LEAVE_INTERNAL_BUFFERS_UNINITIALIZED));
  return do_test__wuffs_base__image_decoder(
      wuffs_netpbm__decoder__upcast_as__wuffs_base__image_decoder(&dec),
      "test/data/hippopotamus.ppm", 0, SIZE_MAX, 36, 28, 0xFFF5F5F5);
}

const char*  //
test_wuffs_netpbm_decode_truncated_input() {
  CHECK_FOCUS(__func__);

  wuffs_base__io_buffer src =
      wuffs_base__ptr_u8__reader(g_src_array_u8, 0, false);
  wuffs_netpbm__decoder dec;
  CHECK_STATUS("initialize",
               wuffs_netpbm__decoder__initialize(
                   &dec, sizeof dec, WUFFS_VERSION,
                   WUFFS_INITIALIZE__LEAVE_INTERNAL_BUFFERS_UNINITIALIZED));

  wuffs_base__status status =
      wuffs_netpbm__decoder__decode_image_config(&dec, NULL, &src);
  if (status.repr != wuffs_base__suspension__short_read) {
    RETURN_FAIL("closed=false: have \"%s\", want \"%s\"", status.repr,
                wuffs_base__suspension__short_read);
  }

  src.meta.closed = true;
  status = wuffs_netpbm__decoder__decode_image_config(&dec, NULL, &src);
  if (status.repr != wuffs_netpbm__error__truncated_input) {
    RETURN_FAIL("closed=true: have \"%s\", want \"%s\"", status.repr,
                wuffs_netpbm__error__truncated_input);
  }
  return NULL;
}

const char*  //
test_wuffs_netpbm_decode_frame_config() {
  CHECK_FOCUS(__func__);
  wuffs_netpbm__decoder dec;
  CHECK_STATUS("initialize",
               wuffs_netpbm__decoder__initialize(
                   &dec, sizeof dec, WUFFS_VERSION,
                   WUFFS_INITIALIZE__LEAVE_INTERNAL_BUFFERS_UNINITIALIZED));

  wuffs_base__frame_config fc = ((wuffs_base__frame_config){});
  wuffs_base__io_buffer src = ((wuffs_base__io_buffer){
      .data = g_src_slice_u8,
  });
  CHECK_STRING(read_file(&src, "test/data/hippopotamus.pgm"));
  CHECK_STATUS("decode_frame_config #0",
               wuffs_netpbm__decoder__decode_frame_config(&dec, &fc, &src));

  wuffs_base__status status =
      wuffs_netpbm__decoder__decode_frame_config(&dec, &fc, &src);
  if (status.repr != wuffs_base__note__end_of_data) {
    RETURN_FAIL("decode_frame_config #1: have \"%s\", want \"%s\"", status.repr,
                wuffs_base__note__end_of_data);
  }
  return NULL;
}

const char*  //
test_wuffs_netpbm_decode_image_config() {
  CHECK_FOCUS(__func__);
  wuffs_netpbm__decoder dec;
  CHECK_STATUS("initialize",
               wuffs_netpbm__decoder__initialize(
                   &dec, sizeof dec, WUFFS_VERSION,
                   WUFFS_INITIALIZE__LEAVE_INTERNAL_BUFFERS_UNINITIALIZED));

  wuffs_base__image_config ic = ((wuffs_base__image_config){});
  wuffs_base__io_buffer src = ((wuffs_base__io_buffer){
      .data = g_src_slice_u8,
  });
  CHECK_STRING(read_file(&src, "test/data/hippopotamus.pgm"));
  CHECK_STATUS("decode_image_config",
               wuffs_netpbm__decoder__decode_image_config(&dec, &ic, &src));

  uint32_t have_width = wuffs_base__pixel_config__width(&ic.pixcfg);
  uint32_t want_width = 36;
  if (have_width != want_width) {
    RETURN_FAIL("width: have %" PRIu32 ", want %" PRIu32, have_width,
                want_width);
  }
  uint32_t have_height = wuffs_base__pixel_config__height(&ic.pixcfg);
  uint32_t want_height = 28;
  if (have_height != want_height) {
    RETURN_FAIL("height: have %" PRIu32 ", want %" PRIu32, have_height,
                want_height);
  }
  return NULL;
}

// ---------------- Mimic Tests

#ifdef WUFFS_MIMIC

// No mimic tests.

#endif  // WUFFS_MIMIC

// ---------------- Netpbm Benches

// No Netpbm benches.

// ---------------- Mimic Benches

#ifdef WUFFS_MIMIC

// No mimic benches.

#endif  // WUFFS_MIMIC

// ---------------- Manifest

proc g_tests[] = {

    test_wuffs_netpbm_decode_frame_config,
    test_wuffs_netpbm_decode_image_config,
    test_wuffs_netpbm_decode_interface,
    test_wuffs_netpbm_decode_truncated_input,

#ifdef WUFFS_MIMIC

// No mimic tests.

#endif  // WUFFS_MIMIC

    NULL,
};

proc g_benches[] = {

// No Netpbm benches.

#ifdef WUFFS_MIMIC

// No mimic benches.

#endif  // WUFFS_MIMIC

    NULL,
};

int  //
main(int argc, char** argv) {
  g_proc_package_name = "std/netpbm";
  return test_main(argc, argv, g_tests, g_benches);
}
