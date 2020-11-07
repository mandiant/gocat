// +build ignore

package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	rxpHashName    = regexp.MustCompile(`static const char \*HASH_NAME\s+=\s+\"(.*?)\";`)
	rxpKernelType  = regexp.MustCompile(`static const u64\s+KERN_TYPE\s+=\s+(.*?);`)
	rxpExampleHash = regexp.MustCompile(`static const char\s+\*ST_HASH\s+=\s+"(.*?)\";`)
)

func locateHashName(buff []byte) string {
	matches := rxpHashName.FindSubmatch(buff)
	if len(matches) != 2 {
		return ""
	}

	return string(matches[1])
}

func locateKernelType(buff []byte) (int, error) {
	matches := rxpKernelType.FindSubmatch(buff)
	if len(matches) != 2 {
		return 0, errors.New("could not locate kernel type")
	}

	return strconv.Atoi(string(matches[1]))
}

func locateExample(buff []byte) string {
	matches := rxpExampleHash.FindSubmatch(buff)
	if len(matches) != 2 {
		return ""
	}

	return string(matches[1])
}

type customHash struct {
	Name    string
	Type    int
	Example string
}

var knownDynamic = map[int][]customHash{
	16511: []customHash{
		{
			Name: "JWT (JSON Web Token) HS256",
			Type: 16511,
		},
		{
			Name: "JWT (JSON Web Token) HS384",
			Type: 16512,
		},
		{
			Name: "JWT (JSON Web Token) HS512",
			Type: 16513,
		},
	},
}

func main() {
	srcPath := os.Getenv("HASHCAT_SRC_PATH")
	if srcPath != "" {
		log.Fatal("HASHCAT_SRC_PATH must be set to hashcat's src/modules directory to generate code")
	}

	b := new(bytes.Buffer)
	b.WriteString("// Code automatically generated; DO NOT EDIT.\n")
	b.WriteString("\n")
	b.WriteString("package types")
	b.WriteString("\n")
	b.WriteString("// Hash describes information about supported file hashes\n")
	b.WriteString("type Hash struct {\n")
	b.WriteString("\tName string\n")
	b.WriteString("\tExample string\n")
	b.WriteString("\tType int\n")
	b.WriteString("}\n")
	b.WriteString("var hashes = []Hash{\n")

	filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		bytez, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf(err.Error())
		}

		hashName := locateHashName(bytez)
		if hashName == "" {
			log.Fatalf("Could not locate hash name in %s", info.Name())
		}

		kernelType, err := locateKernelType(bytez)
		if err != nil {
			log.Fatalf("Could not locate kernel type in %s", info.Name())
		}

		example := locateExample(bytez)

		if dyn, ok := knownDynamic[kernelType]; ok {
			for _, hashType := range dyn {
				b.WriteString("\t{\n")
				b.WriteString(fmt.Sprintf("\t Name: \"%s\",\n", hashType.Name))
				b.WriteString(fmt.Sprintf("\t Type: %d,\n", hashType.Type))

				if hashType.Example != "" {
					b.WriteString(fmt.Sprintf("\t Example: \"%s\",\n", hashType.Example))
				} else if hashType.Example == "" && example != "" {
					b.WriteString(fmt.Sprintf("\t Example: \"%s\",\n", example))
				}

				b.WriteString("\t},\n")
			}

			return nil
		}

		b.WriteString("\t{\n")
		b.WriteString(fmt.Sprintf("\t Name: \"%s\",\n", hashName))
		b.WriteString(fmt.Sprintf("\t Type: %d,\n", kernelType))
		if example != "" {
			b.WriteString(fmt.Sprintf("\t Example: \"%s\",\n", example))
		}
		b.WriteString("\t},\n")

		return nil
	})

	b.WriteString("\t}\n\n")

	b.WriteString("// SupportedHashes returns a list of available hashes supported by Hashcat\n")
	b.WriteString("func SupportedHashes() []Hash {\n")
	b.WriteString("\t return hashes\n")
	b.WriteString("}\n")

	formattedSource, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalf("Could not format generated source: %s", err)
	}

	fd, err := os.Create("./hash_types.go")
	if err != nil {
		log.Fatalf("Could not create destination file: %s", err)
	}
	defer fd.Close()

	if _, err := fd.Write(formattedSource); err != nil {
		log.Fatalf("Could not write destination file: %s", err)
	}
}
