package engine

type key int

const (
	ContextOriginalPath key = iota
	ContextRequestStart
	ContextUserID
)
