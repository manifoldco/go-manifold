package names

import (
	"encoding/binary"
	"strings"

	"github.com/manifoldco/go-base32"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/names/data"
)

const (
	entropy    = 16 * 8 // in bits
	base32Bits = 25     // Enough to fill 5 characters exactly.
)

var aShare, cShare, sShare int

func init() {
	aShare = bitsNeeded(len(data.Adjectives))
	cShare = bitsNeeded(len(data.Colors))
	sShare = bitsNeeded(len(data.Shapes))

	if aShare+cShare+sShare+base32Bits > entropy {
		panic("Word lists are now too big")
	}
}

// New returns a generated name based on the provided id, and its matching label.
// Names are title cased with spaces between words.
// Labels are lowercased with hyphens between words, and 5 trailing characters
// of base32 encoded data from the id.
func New(id manifold.ID) (string, string) {
	idBytes := id[2:]

	offset := 0
	adj, offset := fetchWord(idBytes, data.Adjectives, offset, aShare)
	color, offset := fetchWord(idBytes, data.Colors, offset, cShare)
	shape, offset := fetchWord(idBytes, data.Shapes, offset, sShare)

	num := fetchNum(idBytes, offset)

	name := strings.Title(adj + " " + color + " " + shape)
	label := strings.Replace(strings.ToLower(name), " ", "-", -1) + "-" + num
	return name, label
}

func fetchWord(idBytes []byte, wordList []string, offset, bitShare int) (string, int) {
	padding := make([]byte, (64-bitShare)/8)

	mask := byte(0xff)
	bytesNeeded := bitShare / 8
	if rem := bitShare % 8; rem > 0 {
		bytesNeeded++
		mask >>= uint(8 - rem)
	}

	vBytes := idBytes[offset : offset+bytesNeeded]
	vBytes[0] &= mask
	v := append(padding, vBytes...)
	idx := binary.BigEndian.Uint64(v)
	if idx >= uint64(len(wordList)) {
		idx -= uint64(len(wordList))
	}

	return wordList[idx], offset + bytesNeeded
}

func fetchNum(idBytes []byte, offset int) string {
	mask := byte(0xff)
	bytesNeeded := base32Bits / 8
	if rem := base32Bits % 8; rem > 0 {
		bytesNeeded++
		mask <<= uint(8 - rem)
	}

	vBytes := idBytes[offset : offset+bytesNeeded]
	vBytes[len(vBytes)-1] &= mask

	return base32.EncodeToString(vBytes)[:5]
}

func bitsNeeded(val int) int {
	var bits int
	for bits = 0; val > 0; val = val >> 1 {
		bits++
	}

	return bits
}
