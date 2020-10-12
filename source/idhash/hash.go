// Package idhash contains functions for deriving MD5 hash IDs
// for Individual objects from Contribution and Disbursement ojects.
package idhash

import (
	"crypto/md5"
	"encoding/hex"
)

// FormatIndvInput returns an input string for NewHash derived from the Name, Employer, Occupation, & Zip fields.
func FormatIndvInput(name, employer, occupation, zip string) string {
	return name + " - " + employer + " - " + occupation + " - " + zip
}

// FormatOrgInput returns an input string for NewHash derived from the Name and Zip fields.
func FormatOrgInput(name, zip string) string {
	return name + " - " + zip
}

// NewHash creates a new MD5 hash and returns the hash encoded as a string.
func NewHash(input string) string {
	sum := md5.Sum([]byte(input))
	pass := hex.EncodeToString(sum[:])
	return pass
}
