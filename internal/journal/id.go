package journal

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultHashLength     = 6
	ExpandedHashLength    = 8
	maxIDCollisionRetries = 32
)

type IDExistsFunc func(string) (bool, error)

var hashSuffixGenerator = generateHashSuffix

func GenerateID(prefix string, exists IDExistsFunc) (string, error) {
	return GenerateIDWithLength(prefix, DefaultHashLength, exists)
}

func GenerateIDWithLength(prefix string, hashLength int, exists IDExistsFunc) (string, error) {
	if strings.TrimSpace(prefix) == "" {
		return "", fmt.Errorf("generate id: prefix is required")
	}
	if hashLength <= 0 {
		return "", fmt.Errorf("generate id: hash length must be positive")
	}
	if exists == nil {
		exists = func(string) (bool, error) {
			return false, nil
		}
	}

	for attempt := 0; attempt < maxIDCollisionRetries; attempt++ {
		suffix, err := hashSuffixGenerator(hashLength)
		if err != nil {
			return "", fmt.Errorf("generate id suffix: %w", err)
		}

		candidate := prefix + "-" + suffix
		taken, err := exists(candidate)
		if err != nil {
			return "", fmt.Errorf("check id collision for %q: %w", candidate, err)
		}
		if !taken {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("generate id with prefix %q: exhausted collision retries", prefix)
}

func FormatEventID(occurredAt time.Time, sequence int) (string, error) {
	if sequence <= 0 {
		return "", fmt.Errorf("format event id: sequence must be positive")
	}

	date := occurredAt.UTC().Format("20060102")
	sequenceText := strconv.Itoa(sequence)
	width := max(6, len(sequenceText))
	return fmt.Sprintf("EVT-%s-%0*d", date, width, sequence), nil
}

func generateHashSuffix(length int) (string, error) {
	uuid4, err := newUUID4String()
	if err != nil {
		return "", err
	}

	sum := sha1.Sum([]byte(uuid4))
	full := hex.EncodeToString(sum[:])
	if length > len(full) {
		return "", fmt.Errorf("hash length %d exceeds available hex characters", length)
	}

	return full[:length], nil
}

func newUUID4String() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	), nil
}

func max(first, second int) int {
	if first > second {
		return first
	}
	return second
}
