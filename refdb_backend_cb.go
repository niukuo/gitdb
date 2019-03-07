package gitdb

/*
#include <stdlib.h>
#include <git2.h>
#include <git2/sys/refs.h>
#include <git2/sys/refdb_backend.h>

extern void _go_git_populate_refdb_backend_cb(git_refdb_backend *backend);
*/
import "C"
import (
	"unsafe"

	git "gopkg.in/libgit2/git2go.v27"
)

var refdbBackendHandles = NewHandleList()

func (v *RefdbBackend) Track() *git.RefdbBackend {
	handle := C.calloc(1, C.sizeof_struct_git_refdb_backend)
	backend := (*C.git_refdb_backend)(handle)
	C.git_refdb_init_backend(backend, C.GIT_REFDB_BACKEND_VERSION)
	C._go_git_populate_refdb_backend_cb(backend)

	refdbBackendHandles.Track(handle, v)
	return git.NewRefdbBackendFromC(handle)
}

func getRefdbBackend(handle *C.git_refdb_backend) *RefdbBackend {
	return refdbBackendHandles.Get(unsafe.Pointer(handle)).(*RefdbBackend)
}

//export cbRefdbBackendFree
func cbRefdbBackendFree(handle *C.git_refdb_backend) {
	ptr := unsafe.Pointer(handle)
	backend := refdbBackendHandles.Untrack(ptr).(*RefdbBackend)
	C.free(ptr)
	if backend.Free != nil {
		backend.Free()
	}
}

//export cbRefdbBackendExists
func cbRefdbBackendExists(exists *C.int, handle *C.git_refdb_backend, ref_name *C.char) C.int {
	backend := getRefdbBackend(handle)
	if ret := backend.Exists(C.GoString(ref_name)); ret {
		*exists = 1
	} else {
		*exists = 0
	}
	return C.GIT_OK
}

func allocReference(ref_name *C.char, typ git.ReferenceType, ref string) *C.git_reference {
	cref := C.CString(ref)
	defer C.free(unsafe.Pointer(cref))

	if typ == git.ReferenceSymbolic {
		return C.git_reference__alloc_symbolic(ref_name, cref)
	} else {
		var oid C.git_oid
		C.git_oid_fromstr(&oid, cref)
		return C.git_reference__alloc(ref_name, &oid, nil)
	}
}

//export cbRefdbBackendLookup
func cbRefdbBackendLookup(
	out **C.git_reference,
	handle *C.git_refdb_backend,
	ref_name *C.char,
) C.int {
	backend := getRefdbBackend(handle)
	if typ, ref, err := backend.Lookup(C.GoString(ref_name)); err != nil {
		return errToCode(err)
	} else {
		*out = allocReference(ref_name, typ, ref)
	}
	return C.GIT_OK
}

//export cbRefdbBackendRename
func cbRefdbBackendRename(
	out **C.git_reference,
	handle *C.git_refdb_backend,
	old_name *C.char,
	new_name *C.char,
	force C.int,
	who *C.git_signature,
	message *C.char,
) C.int {
	backend := getRefdbBackend(handle)
	if typ, ref, err := backend.Rename(
		C.GoString(old_name),
		C.GoString(new_name),
		force == 1,
		newSignatureFromC(who),
		C.GoString(message),
	); err != nil {
		return errToCode(err)
	} else {
		*out = allocReference(new_name, typ, ref)
	}
	return C.GIT_OK
}

//export cbRefdbBackendWrite
func cbRefdbBackendWrite(
	handle *C.git_refdb_backend,
	ref *C.git_reference,
	force C.int,
	who *C.git_signature,
	message *C.char,
	old_id *C.git_oid,
	old_target *C.char) C.int {

	backend := getRefdbBackend(handle)
	if err := backend.Write(
		newReferenceFromC(ref),
		force == 1,
		newSignatureFromC(who),
		C.GoString(message),
		newOidFromC(old_id),
		C.GoString(old_target)); err != nil {
		return errToCode(err)
	}
	return C.GIT_OK
}

//export cbRefdbBackendDel
func cbRefdbBackendDel(
	handle *C.git_refdb_backend,
	ref_name *C.char,
	old_id *C.git_oid,
	old_target *C.char) C.int {
	backend := getRefdbBackend(handle)
	if err := backend.Del(
		C.GoString(ref_name),
		newOidFromC(old_id),
		C.GoString(old_target)); err != nil {
		return errToCode(err)
	}
	return C.GIT_OK
}
