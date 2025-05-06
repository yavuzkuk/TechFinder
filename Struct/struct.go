package Struct

type HeaderInfo struct {
	Type string `json:"type"`
}
type HeaderInfoWithVersion struct {
	HeaderInfo
	Version string `json:"version"`
}

type ComponentDetail struct {
	Component string
	Version   string
}
