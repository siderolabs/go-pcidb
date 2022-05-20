# go-pcidb

`go-pcidb` provides a Go library for converting PCI IDs to descriptive names using
[PCI ID database](https://pci-ids.ucw.cz/).

The library embeds the database as compact static Go `map` structures (using `go generate`) to avoid
parsing costs and optimize memory usage.

An alternative Go library that does parsing on load from `pci.ids` is [github.com/jaypipes/pcidb](https://github.com/jaypipes/pcidb).

## Updating PCI ID database

Download new `pci.ids` from [https://pci-ids.ucw.cz/](https://pci-ids.ucw.cz/) and run `make generate` to update the database.
