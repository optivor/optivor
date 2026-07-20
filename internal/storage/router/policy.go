package router

import "strings"

type AccessPolicy int

const (
	PolicyPublic AccessPolicy = iota
	PolicySigned
	PolicyPrivate
)

func (p AccessPolicy) String() string {
	switch p {
	case PolicyPublic:
		return "public"
	case PolicySigned:
		return "signed"
	case PolicyPrivate:
		return "private"
	default:
		return "unknown"
	}
}

func ParseAccessPolicy(s string) AccessPolicy {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "signed":
		return PolicySigned
	case "private":
		return PolicyPrivate
	case "public":
		fallthrough
	default:
		return PolicyPublic
	}
}
