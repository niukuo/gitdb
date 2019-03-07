package gitdb

import git "gopkg.in/libgit2/git2go.v27"

type Reference interface {
	Name() string
	Type() git.ReferenceType
	Target() *git.Oid
	SymbolicTarget() string
}

type RefdbBackend struct {
	Exists   func(ref_name string) bool
	Lookup   func(ref_name string) (git.ReferenceType, string, error)
	Iterator func(glob string) (ReferenceIterator, error)
	Write    func(ref Reference, force bool, who *git.Signature, message string, old_id *git.Oid, old_target string) error
	Rename   func(old_name string, new_name string, force bool, who *git.Signature, message string) (git.ReferenceType, string, error)
	Del      func(ref_name string, old_id *git.Oid, old_target string) error
	Free     func()
}

func NewRefdbBackend() *RefdbBackend {
	b := &RefdbBackend{}
	return b
}
