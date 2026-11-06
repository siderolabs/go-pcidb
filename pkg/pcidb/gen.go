// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build exclude

package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Source data: pci.ids
// Source can be updated by downloading new version from https://pci-ids.ucw.cz/.

// Some code in this file is based on https://github.com/jaypipes/pcidb.

// Output: names.bin (deduplicated, length-prefixed names) and one file per lookup table
// with sorted fixed-size records: big-endian key followed by a big-endian offset into names.bin.
//
// See table.go for the lookup side.

const (
	offsetSize    = 3
	maxNameLength = 255
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pcidb-gen: ")

	if err := parseDBFile(); err != nil {
		log.Fatalf("error: %s", err)
	}

	var names namesBuilder

	tables := []struct {
		file    string
		keySize int
		entries map[uint64]string
	}{
		{"classes.bin", 1, toEntries(Classes)},
		{"subclasses.bin", 2, toEntries(Subclasses)},
		{"vendors.bin", 2, toEntries(Vendors)},
		{"products.bin", 4, toEntries(Products)},
	}

	for _, table := range tables {
		keys := make([]uint64, 0, len(table.entries))

		for key := range table.entries {
			keys = append(keys, key)
		}

		slices.Sort(keys)

		var buf bytes.Buffer

		for _, key := range keys {
			writeUint(&buf, key, table.keySize)
			writeUint(&buf, names.add(table.entries[key]), offsetSize)
		}

		writeFile(table.file, buf.Bytes())
	}

	if names.buf.Len() >= 1<<(8*offsetSize) {
		log.Fatalf("names blob is too large: %d bytes", names.buf.Len())
	}

	writeFile("names.bin", names.buf.Bytes())
}

func toEntries[K uint8 | uint16 | uint32](m map[K]string) map[uint64]string {
	entries := make(map[uint64]string, len(m))

	for k, v := range m {
		entries[uint64(k)] = v
	}

	return entries
}

type namesBuilder struct {
	buf     bytes.Buffer
	offsets map[string]uint64
}

// add appends the name to the blob (unless it's already there) and returns its offset.
func (b *namesBuilder) add(name string) uint64 {
	if offset, ok := b.offsets[name]; ok {
		return offset
	}

	if len(name) > maxNameLength {
		log.Fatalf("name is too long (%d bytes): %q", len(name), name)
	}

	if b.offsets == nil {
		b.offsets = map[string]uint64{}
	}

	offset := uint64(b.buf.Len())

	b.buf.WriteByte(byte(len(name)))
	b.buf.WriteString(name)

	b.offsets[name] = offset

	return offset
}

func writeUint(buf *bytes.Buffer, v uint64, size int) {
	for i := size - 1; i >= 0; i-- {
		buf.WriteByte(byte(v >> (8 * i)))
	}
}

func writeFile(name string, data []byte) {
	if err := os.WriteFile(name, data, 0o644); err != nil {
		log.Fatalf("writing %s: %s", name, err)
	}
}

type (
	Class         = uint8
	ClassSubclass = uint16 // Class + Subclass
	Vendor        = uint16
	VendorProduct = uint32 // Vendor + Product
)

var (
	Classes    = make(map[Class]string)
	Vendors    = make(map[Vendor]string)
	Subclasses = make(map[ClassSubclass]string)
	Products   = make(map[VendorProduct]string)
)

func mustParse8(in []rune) uint8 {
	v, err := strconv.ParseUint(string(in), 16, 8)
	if err != nil {
		panic(err)
	}

	return uint8(v)
}

func mustParse16(in []rune) uint16 {
	v, err := strconv.ParseUint(string(in), 16, 16)
	if err != nil {
		panic(err)
	}

	return uint16(v)
}

func parseDBFile() error {
	in, err := os.Open("pci.ids")
	if err != nil {
		return err
	}

	defer in.Close()

	scanner := bufio.NewScanner(in)

	inClassBlock := false

	var (
		curClass  Class
		curVendor Vendor
	)

	for scanner.Scan() {
		line := scanner.Text()
		// skip comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lineBytes := []rune(line)

		// Lines starting with an uppercase "C" indicate a PCI top-level class
		// information block. These lines look like this:
		//
		// C 02  Network controller
		if lineBytes[0] == 'C' {
			inClassBlock = true

			curClass = mustParse8(lineBytes[2:4])
			className := string(lineBytes[6:])

			Classes[curClass] = className

			continue
		}

		// Lines not beginning with an uppercase "C" or a TAB character
		// indicate a top-level vendor information block. These lines look like
		// this:
		//
		// 0a89  BREA Technologies Inc
		if lineBytes[0] != '\t' {
			inClassBlock = false

			curVendor = mustParse16(lineBytes[0:4])
			vendorName := string(lineBytes[6:])

			Vendors[curVendor] = vendorName

			continue
		}

		// Lines beginning with only a single TAB character are *either* a
		// subclass OR are a device information block. If we're in a class
		// block (i.e. the last parsed block header was for a PCI class), then
		// we parse a subclass block. Otherwise, we parse a device information
		// block.
		//
		// A subclass information block looks like this:
		//
		// \t00  Non-VGA unclassified device
		//
		// A device information block looks like this:
		//
		// \t0002  PCI to MCA Bridge
		//
		// Lines beginning with two TAB characters are programming interfaces
		// (in a class block) or subsystems (in a vendor block); they are not
		// exposed by the package, so they are skipped.
		if len(lineBytes) > 1 && lineBytes[1] != '\t' {
			if inClassBlock {
				subclassID := mustParse8(lineBytes[1:3])
				subclassName := string(lineBytes[5:])

				Subclasses[ClassSubclass(uint16(curClass)<<8|uint16(subclassID))] = subclassName
			} else {
				productID := mustParse16(lineBytes[1:5])
				productName := string(lineBytes[7:])

				Products[VendorProduct(uint32(curVendor)<<16|uint32(productID))] = productName
			}
		}
	}

	return scanner.Err()
}
