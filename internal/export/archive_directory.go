package export

import (
	"archive/zip"
	"encoding/binary"
	"io"
)

// inspectArchiveEntryCount walks actual central-directory records before the
// standard ZIP reader allocates a File for every entry. Counts from the end
// record alone are not authoritative (ZIP permits truncated 16-bit counts).
func inspectArchiveEntryCount(input io.ReaderAt, size int64) error {
	start, end, err := archiveDirectoryBounds(input, size)
	if err != nil {
		return err
	}
	var header [46]byte
	count := 0
	for offset := start; offset < end; {
		if end-offset < int64(len(header)) {
			return zip.ErrFormat
		}
		if _, err := input.ReadAt(header[:], offset); err != nil {
			return err
		}
		if binary.LittleEndian.Uint32(header[:4]) != 0x02014b50 {
			return zip.ErrFormat
		}
		count++
		if count > maxArchiveEntries {
			return errArchiveResourceLimit
		}
		length := int64(len(header)) + int64(binary.LittleEndian.Uint16(header[28:30])) + int64(binary.LittleEndian.Uint16(header[30:32])) + int64(binary.LittleEndian.Uint16(header[32:34]))
		if length > end-offset {
			return zip.ErrFormat
		}
		offset += length
	}
	return nil
}

func archiveDirectoryBounds(input io.ReaderAt, size int64) (int64, int64, error) {
	length := int64(65 * 1024)
	if size < length {
		length = size
	}
	if length < 22 {
		return 0, 0, zip.ErrFormat
	}
	buffer := make([]byte, length)
	if _, err := input.ReadAt(buffer, size-length); err != nil {
		return 0, 0, err
	}
	for i := len(buffer) - 22; i >= 0; i-- {
		if binary.LittleEndian.Uint32(buffer[i:i+4]) != 0x06054b50 {
			continue
		}
		if int(binary.LittleEndian.Uint16(buffer[i+20:i+22])) > len(buffer)-i-22 {
			continue
		}
		end := size - length + int64(i)
		directorySize := uint64(binary.LittleEndian.Uint32(buffer[i+12 : i+16]))
		directoryOffset := uint64(binary.LittleEndian.Uint32(buffer[i+16 : i+20]))
		if end >= 20 {
			var locator [20]byte
			if _, err := input.ReadAt(locator[:], end-20); err != nil {
				return 0, 0, err
			}
			if binary.LittleEndian.Uint32(locator[:4]) == 0x07064b50 {
				position := binary.LittleEndian.Uint64(locator[8:16])
				if position > uint64(size-56) {
					return 0, 0, zip.ErrFormat
				}
				var record [56]byte
				if _, err := input.ReadAt(record[:], int64(position)); err != nil {
					return 0, 0, err
				}
				if binary.LittleEndian.Uint32(record[:4]) != 0x06064b50 {
					return 0, 0, zip.ErrFormat
				}
				end = int64(position)
				directorySize = binary.LittleEndian.Uint64(record[40:48])
				directoryOffset = binary.LittleEndian.Uint64(record[48:56])
			}
		}
		if directorySize > uint64(end) || directoryOffset > uint64(size) {
			return 0, 0, zip.ErrFormat
		}
		start := end - int64(directorySize)
		// Match archive/zip's accepted prefix-offset interpretation.
		if start > int64(directoryOffset) {
			var signature [4]byte
			if _, err := input.ReadAt(signature[:], int64(directoryOffset)); err == nil && binary.LittleEndian.Uint32(signature[:]) == 0x02014b50 {
				start = int64(directoryOffset)
				end = start + int64(directorySize)
			}
		}
		if start < 0 || end > size {
			return 0, 0, zip.ErrFormat
		}
		return start, end, nil
	}
	return 0, 0, zip.ErrFormat
}
