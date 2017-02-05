package main

import "github.com/posteo/go-agentx/value"

type OIDSorter []value.OID

func (os OIDSorter) Len() int {
	return len(os)
}
func (os OIDSorter) Swap(i, j int) {
	os[i], os[j] = os[j], os[i]
}

func (os OIDSorter) Less(i, j int) bool {
	return OIDLessThan(os[i], os[j])
}

func OIDLessThan(oid value.OID, other value.OID) bool {
	len_left := len(oid)
	len_right := len(other)
	to := len_left
	if len_right < len_left {
		to = len_right
	}
	for k := 0; k < to; k++ {
		if oid[k] < other[k] {
			return true
		} else if oid[k] > other[k] {
			return false
		}
	}
	return len_left < len_right
}

func OIDGreaterThan(oid value.OID, other value.OID) bool {
	len_left := len(oid)
	len_right := len(other)
	to := len_left
	if len_right < len_left {
		to = len_right
	}
	for k := 0; k < to; k++ {
		if oid[k] > other[k] {
			return true
		} else if oid[k] < other[k] {
			return false
		}
	}
	return len_left > len_right
}
