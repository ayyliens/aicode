package main

type CmdOaiCommon struct {
	Path  string `flag:"--path" desc:"target path (required)"`
	Watch bool   `flag:"--watch" desc:"watch and rerun"`
	Init  bool   `flag:"--init" desc:"perform initial run before watching"`
}
