package integrity

/*
#include "postgres.h"
#include "catalog/pg_control.h"
#include "storage/bufpage.h"
#include "storage/checksum_impl.h"

static int block_size() {
	return BLCKSZ;
}

static uint16 page_checksum(char *page, int number) {
	return pg_checksum_page(page, number);
}

static uint16 header_checksum(char *page) {
	PageHeader header = (PageHeader) page;
	return header->pd_checksum;
}
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

// ErrChecksum is returned if a page checksum didn't match.
var ErrChecksum = errors.New("checksum mismatch")

func blockSize() uint {
	return uint(C.block_size())
}

func pageChecksum(p []byte, n int) uint16 {
	page := (*C.char)(C.CBytes(p))
	defer C.free(unsafe.Pointer(page))
	return uint16(C.page_checksum(page, C.int(n)))
}

func headerChecksum(p []byte) uint16 {
	page := (*C.char)(C.CBytes(p))
	defer C.free(unsafe.Pointer(page))
	return uint16(C.header_checksum(page))
}

// Verify block checksums
func Verify(r io.Reader) error {
	blockNumber := 0
	blockSize := blockSize()
	block := make([]byte, blockSize)
	for {
		_, err := io.ReadFull(r, block)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		checksum := pageChecksum(block, blockNumber)
		headerChecksum := headerChecksum(block)
		if headerChecksum == 0 {
			// Checksum not enabled
			return nil
		}
		if checksum != headerChecksum {
			// Checksum mistmatch
			return ErrChecksum
		}
		blockNumber++
	}
	return nil
}
