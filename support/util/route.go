package util

import (
	"fmt"
	"hash/fnv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

// ShortenRouteHostnameIfNeeded will return a shortened hostname if the route hostname will exceed
// the allowed DNS name size. If the hostname is not too long, an empty string is returned so that
// the default can be used.
func ShortenRouteHostnameIfNeeded(name, namespace string, baseDomain string) string {
	if baseDomain == "" {
		return ""
	}
	if len(name)+len(namespace)+1 < validation.DNS1123LabelMaxLength {
		return ""
	} else {
		return fmt.Sprintf("%s.%s", strings.TrimSuffix(shortenName(name+"-"+namespace, "", validation.DNS1123LabelMaxLength), "-"), baseDomain)
	}
}

// shortenName returns a name given a base ("deployment-5") and a suffix ("deploy")
// It will first attempt to join them with a dash. If the resulting name is longer
// than maxLength: if the suffix is too long, it will truncate the base name and add
// an 8-character hash of the [base]-[suffix] string.  If the suffix is not too long,
// it will truncate the base, add the hash of the base and return [base]-[hash]-[suffix]
// Source: openshift/origin v3.9.0 pkg/api/apihelpers/namer.go
func shortenName(base, suffix string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	name := fmt.Sprintf("%s-%s", base, suffix)
	if len(name) <= maxLength {
		return name
	}

	baseLength := maxLength - 10 /*length of -hash-*/ - len(suffix)

	// if the suffix is too long, ignore it
	if baseLength < 0 {
		prefix := base[0:min(len(base), max(0, maxLength-9))]
		// Calculate hash on initial base-suffix string
		shortName := fmt.Sprintf("%s-%s", prefix, hash(name))
		return shortName[:min(maxLength, len(shortName))]
	}

	prefix := base[0:baseLength]
	// Calculate hash on initial base-suffix string
	return fmt.Sprintf("%s-%s-%s", prefix, hash(base), suffix)
}

// max returns the greater of its 2 inputs
func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// min returns the lesser of its 2 inputs
func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// hash calculates the hexadecimal representation (8-chars)
// of the hash of the passed in string using the FNV-a algorithm
func hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	intHash := hash.Sum32()
	result := fmt.Sprintf("%08x", intHash)
	return result
}
