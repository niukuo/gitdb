package gitdb

import git "gopkg.in/libgit2/git2go.v27"

type OdbBackend struct {
	Read       func(oid *git.Oid) ([]byte, git.ObjectType, error)
	ReadHeader func(oid *git.Oid) (uint64, git.ObjectType, error)
	Exists     func(oid *git.Oid) bool
	ForEach    func(callback git.OdbForEachCallback) error
	Write      func(oid *git.Oid, data []byte, otype git.ObjectType) error
	Free       func()
}

func NewOdbBackend() *OdbBackend {
	b := &OdbBackend{}
	return b
}
