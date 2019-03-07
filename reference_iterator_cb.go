package gitdb

/*
#include <git2.h>
#include <git2/sys/refdb_backend.h>
extern void _go_git_populate_reference_iterator_cb(git_reference_iterator* it);
*/
import "C"
import (
	"unsafe"

	git "gopkg.in/libgit2/git2go.v27"
)

var refItHandles = NewHandleList()

type refIterWrapper struct {
	name *C.char
	iter ReferenceIterator
}

func (v *refIterWrapper) Next() (*C.char, git.ReferenceType, string, error) {
	name, typ, ref, err := v.iter.Next()
	if err != nil {
		return nil, 0, "", err
	}
	if v.name != nil {
		C.free(unsafe.Pointer(v.name))
	}
	v.name = C.CString(name)
	return v.name, typ, ref, nil
}

func (v *refIterWrapper) Free() {
	if v.name != nil {
		C.free(unsafe.Pointer(v.name))
		v.name = nil
	}
	v.iter.Free()
}

func getRefIt(handle *C.git_reference_iterator) *refIterWrapper {
	return refItHandles.Get(unsafe.Pointer(handle)).(*refIterWrapper)
}

//export cbRefdbBackendIterator
func cbRefdbBackendIterator(
	iter **C.git_reference_iterator,
	handle *C.git_refdb_backend,
	glob *C.char) C.int {

	backend := getRefdbBackend(handle)
	iterator, err := backend.Iterator(C.GoString(glob))
	if err != nil {
		return errToCode(err)
	}

	*iter = (*C.git_reference_iterator)(C.calloc(1, C.sizeof_struct_git_reference_iterator))
	C._go_git_populate_reference_iterator_cb(*iter)

	refItHandles.Track(unsafe.Pointer(*iter), &refIterWrapper{iter: iterator})

	return C.GIT_OK
}

//export cbReferenceIteratorFree
func cbReferenceIteratorFree(handle *C.git_reference_iterator) {
	ptr := unsafe.Pointer(handle)
	refItHandles.Untrack(ptr)
	C.free(ptr)
}

//export cbReferenceIteratorNext
func cbReferenceIteratorNext(ref **C.git_reference,
	handle *C.git_reference_iterator) C.int {
	iter := getRefIt(handle)
	name, typ, target, err := iter.Next()
	if err != nil {
		return errToCode(err)
	}
	*ref = allocReference(name, typ, target)
	return C.GIT_OK
}

//export cbReferenceIteratorNextName
func cbReferenceIteratorNextName(ref_name **C.char,
	handle *C.git_reference_iterator) C.int {
	iter := getRefIt(handle)
	name, _, _, err := iter.Next()
	if err != nil {
		return errToCode(err)
	}
	*ref_name = name
	return C.GIT_OK
}
