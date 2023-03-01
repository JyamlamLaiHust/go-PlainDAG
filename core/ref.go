package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

func (r *Ref) Hash() int {
	// hash function that return integer
	h := sha256.Sum256([]byte(fmt.Sprintf("%v", r)))
	return int(binary.BigEndian.Uint32(h[:4]))
}

func (r *Ref) Equals(other *Ref) bool {
	return bytes.Equal(r.H, other.H) && r.Index == other.Index
}
