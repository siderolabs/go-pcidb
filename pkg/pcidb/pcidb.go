// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package pcidb provides generated Go structs describing PCI ID database optimized for performance and low GC pressure.
package pcidb

//go:generate go run gen.go

// Class is a PCI device class.
type Class = uint8

// Subclass is a PCI device subclass.
type Subclass = uint8

// ProgrammingInterface is a PCI programming interface.
type ProgrammingInterface = uint8

// Vendor is a PCI vendor.
type Vendor = uint16

// Product is a PCI product.
type Product = uint16

// Subsystem is a PCI subsystem.
type Subsystem = uint16

// SubsystemInfo describes a PCI subsystem.
type SubsystemInfo struct {
	Name   string
	Vendor Vendor
}

// ClassSubclass is a PCI device class + subclass combined into 16 bits.
type ClassSubclass = uint16

// ClassSubclassProgrammingInterface is a class + subclass + programming interface combined into 24 bits.
type ClassSubclassProgrammingInterface = uint32

// VendorProduct is a PCI vendor and product combined into 32 bits.
type VendorProduct = uint32

// VendorProductSubsystem is a PCI vendor, product and subsystem combined into 48 bits.
type VendorProductSubsystem = uint64

// LookupClass by ID.
func LookupClass(classID Class) (string, bool) {
	return lookupClass(classID)
}

// LookupSubclass by ID.
func LookupSubclass(classID Class, subclassID Subclass) (string, bool) {
	key := uint16(classID)<<8 | uint16(subclassID)

	return lookupSubclass(key)
}

// LookupVendor by ID.
func LookupVendor(vendorID Vendor) (string, bool) {
	return lookupVendor(vendorID)
}

// LookupProduct by vendor and product IDs.
func LookupProduct(vendorID Vendor, productID Product) (string, bool) {
	key := uint32(vendorID)<<16 | uint32(productID)

	return lookupProduct(key)
}
