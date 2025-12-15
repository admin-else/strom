package main

type EntryHolderSet struct {
	BaseName  string
	Otherwise struct {
		Name string
		Type any
	}
}
