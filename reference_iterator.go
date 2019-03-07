package gitdb

import git "gopkg.in/libgit2/git2go.v27"

type ReferenceIterator interface {
	Next() (string, git.ReferenceType, string, error)
	Free()
}
