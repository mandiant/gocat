# gocat

gocat is a cgo library for interacting with libhashcat. gocat enables you to create purpose-built password cracking tools that leverage the capabilities of [hashcat](https://hashcat.net/hashcat/).

Below is a matrix that details which versions of hashcat we support:

| Branch        | Hashcat Version | Go Import                     |
| ------------- | --------------- | ----------------------------- |
| `master `     | `v6.1.1`        | `github.com/mandiant/gocat/v6` |
| `v5`          | `v5.X`          | `github.com/fireeye/gocat`    |


## Installation (Please Read)

gocat requires hashcat [v6.X](https://github.com/hashcat/hashcat/releases) or higher to be compiled as a shared library. This can be accomplished by modifying hashcat's `src/Makefile` and setting `SHARED` to `1` . At this time, we also recommend disabling the brain functionality by setting `ENABLE_BRAIN` to `0`

    git clone https://github.com/hashcat/hashcat.git
    git checkout v6.1.1
    make install SHARED=1 ENABLE_BRAIN=0
    cp deps/LZMA-SDK/C/LzmaDec.h /usr/local/include/hashcat/
    cp deps/LZMA-SDK/C/7zTypes.h /usr/local/include/hashcat/
    cp deps/LZMA-SDK/C/Lzma2Dec.h /usr/local/include/hashcat/
    cp -r ./OpenCL/inc_types.h /usr/local/include/hashcat/
    cp -r ./deps/zlib/contrib /usr/local/include/hashcat
    ln -s /usr/local/lib/libhashcat.so.6.1.1 /usr/local/lib/libhashcat.so

At this time, you will also need to set the following environment variables when compiling code that uses this library:

    $ export HASHCAT_SRC_PATH=<Place path here>
    $ export CGO_CFLAGS="-I$HASHCAT_SRC_PATH/OpenCL -I$HASHCAT_SRC_PATH/deps/LZMA-SDK/C -I$HASHCAT_SRC_PATH/deps/zlib -I$HASHCAT_SRC_PATH/deps/zlib/contrib -I$HASHCAT_SRC_PATH/deps/OpenCL-Headers $CGO_CFLAGS"

## Testing

gocat tests need to be run from the `/usr/local/share/hashcat/` directory to access all of the shared files needed for testing.

    go test -c
    cp gocat.test /usr/local/share/hashcat
    cp -r testdata /usr/local/share/hashcat
    /usr/local/share/hashcat/gocat.test

## Known Issues

* Lack of Windows Support: This won't work on windows as I haven't figured out how to build hashcat on windows
* Memory Leaks: hashcat has several (small) memory leaks that could cause increase of process memory over time

## Contributing

Contributions are welcome via pull requests provided they meet the following criteria:

1. One feature or bug fix per PR
1. Code should be properly formatted (using go fmt)
1. Tests coverage should rarely decrease. All new features should have proper coverage
