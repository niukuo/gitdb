package gitdb

/*
#include <git2.h>
#include <git2/sys/refs.h>
*/
import "C"
import (
	"runtime"
	"time"
	"unsafe"

	git "gopkg.in/libgit2/git2go.v27"
)

func newOidFromC(coid *C.git_oid) *git.Oid {
	if coid == nil {
		return nil
	}

	return git.NewOidFromBytes(C.GoBytes(unsafe.Pointer(coid), 20))
}

type reference struct {
	ptr *C.git_reference
}

// weak reference, only for view
func newReferenceFromC(ptr *C.git_reference) *reference {
	return &reference{
		ptr: ptr,
	}
}

func (v *reference) Name() string {
	ret := C.GoString(C.git_reference_name(v.ptr))
	runtime.KeepAlive(v)
	return ret
}

func (v *reference) Type() git.ReferenceType {
	ret := git.ReferenceType(C.git_reference_type(v.ptr))
	runtime.KeepAlive(v)
	return ret
}

func (v *reference) Target() *git.Oid {
	ret := newOidFromC(C.git_reference_target(v.ptr))
	runtime.KeepAlive(v)
	return ret
}

func (v *reference) SymbolicTarget() string {
	var ret string
	cstr := C.git_reference_symbolic_target(v.ptr)

	if cstr != nil {
		return C.GoString(cstr)
	}

	runtime.KeepAlive(v)
	return ret
}

func newSignatureFromC(sig *C.git_signature) *git.Signature {
	// git stores minutes, go wants seconds
	loc := time.FixedZone("", int(sig.when.offset)*60)
	return &git.Signature{
		C.GoString(sig.name),
		C.GoString(sig.email),
		time.Unix(int64(sig.when.time), 0).In(loc),
	}
}
