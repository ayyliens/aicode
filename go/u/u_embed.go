package u

import (
	"path/filepath"

	"github.com/mitranim/gg"
)

type Pathed struct {
	Path string `flag:"--path" desc:"path to use/run/watch" json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
}

func (self Pathed) PathJoin(path string) string {
	return filepath.Join(self.Path, path)
}

func (self Pathed) HasFile(name string) bool {
	return gg.FileExists(self.PathJoin(name))
}

func (self Pathed) ReadFile(name string) []byte {
	return gg.ReadFile[[]byte](self.PathJoin(name))
}

func (self Pathed) WriteFile(name string, body []byte) {
	gg.WriteFile(self.PathJoin(name), body)
}

func (self Pathed) DeleteFile(name string) {
	RemoveFileOrDir(self.PathJoin(name))
}

func (self Pathed) DeleteFileOrSkip(name string) {
	RemoveFileOrDirOrSkip(self.PathJoin(name))
}

func (self Pathed) TouchFile(name string) { self.TouchedFile(name) }

func (self Pathed) TouchedFile(name string) bool {
	return TouchedFile(self.PathJoin(name))
}

type Verbose struct {
	Verb bool `flag:"--verb" desc:"enable verbose logging" json:"verb,omitempty" yaml:"verb,omitempty" toml:"verb,omitempty"`
}

type Inited struct {
	Init bool `flag:"--init" desc:"run once before watching" json:"init,omitempty" yaml:"init,omitempty" toml:"init,omitempty"`
}

type Ignored struct {
	Ignore []string `flag:"--ignore" desc:"paths to ignore when watching" json:"ignore,omitempty" yaml:"ignore,omitempty" toml:"ignore,omitempty"`
}

type Watched struct {
	Watch bool `flag:"--watch" desc:"watch and rerun" json:"watch,omitempty" yaml:"watch,omitempty" toml:"watch,omitempty"`
}

type Named struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
}

type TermClearer struct {
	Clear bool `flag:"--clear" desc:"clear terminal on restart" json:"clear,omitempty" yaml:"clear,omitempty" toml:"clear,omitempty"`
}
