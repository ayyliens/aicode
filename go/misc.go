package main

type CmdOaiCommon struct {
	Watch bool `flag:"--watch" desc:"watch and rerun"`
	Init  bool `flag:"--init" desc:"perform initial run before watching"`
}
