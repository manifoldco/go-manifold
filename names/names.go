package names

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/names/data"
)

const (
	entropy = 16 * 8 // in bits
)

var aShare, cShare, sShare int

func init() {
	aShare = bitsNeeded(len(data.Adjectives))
	cShare = bitsNeeded(len(data.Colors))
	sShare = bitsNeeded(len(data.Shapes))

	if aShare+cShare+sShare > entropy {
		panic("Word lists are now too big")
	}
}

// New returns a generated name based on the provided id, and its matching label.
// Names are title cased with spaces between words.
// Labels are lowercased with hyphens between words
// Deprecated: for resource names use `ForResource` instead
func New(id manifold.ID) (string, string) {
	idBytes := id[2:]

	offset := 0
	adj, offset := fetchWord(idBytes, data.Adjectives, offset, aShare)
	color, offset := fetchWord(idBytes, data.Colors, offset, cShare)
	shape, _ := fetchWord(idBytes, data.Shapes, offset, sShare)

	name := strings.Title(adj + " " + color + " " + shape)
	label := strings.Replace(strings.ToLower(name), " ", "-", -1)
	return name, label
}

// ForResource returns a new label combining a product label with a random label
// for the resource based on its id.
func ForResource(product manifold.Label, id manifold.ID) manifold.Label {
	idBytes := id[2:]

	offset := 0
	adj, offset := fetchWord(idBytes, data.Adjectives, offset, aShare)
	color, offset := fetchWord(idBytes, data.Colors, offset, cShare)
	shape, _ := fetchWord(idBytes, data.Shapes, offset, sShare)

	label := fmt.Sprintf("%s-%s-%s-%s", product, adj, color, shape)
	label = strings.ToLower(strings.Replace(label, " ", "-", -1))

	return manifold.Label(label)
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

func bitsNeeded(val int) int {
	var bits int
	for bits = 0; val > 0; val = val >> 1 {
		bits++
	}

	return bits
}
