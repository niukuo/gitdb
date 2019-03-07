package gitdb

/*
#include <git2.h>
*/
import "C"
import git "gopkg.in/libgit2/git2go.v27"

func errToCode(err error) C.int {
	if e, ok := err.(*git.GitError); ok {
		return C.int(e.Code)
	}
	return C.int(git.ErrGeneric)
}
