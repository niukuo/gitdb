package gitdb

/*
#include <string.h>
#include <git2.h>
#include <git2/sys/odb_backend.h>

extern void _go_git_populate_odb_backend_cb(git_odb_backend *backend);
extern int _go_git_odb_backend_foreach(git_odb_foreach_cb cb, const git_oid *id, void *payload);
*/
import "C"
import (
	"unsafe"

	git "gopkg.in/libgit2/git2go.v27"
)

var odbBackendHandles = NewHandleList()

func (v *OdbBackend) Track() *git.OdbBackend {
	handle := C.calloc(1, C.sizeof_struct_git_odb_backend)
	backend := (*C.git_odb_backend)(handle)
	C.git_odb_init_backend(backend, C.GIT_ODB_BACKEND_VERSION)
	C._go_git_populate_odb_backend_cb(backend)

	odbBackendHandles.Track(handle, v)
	return git.NewOdbBackendFromC(handle)
}

func getOdbBackend(handle *C.git_odb_backend) *OdbBackend {
	return odbBackendHandles.Get(unsafe.Pointer(handle)).(*OdbBackend)
}

//export cbOdbBackendFree
func cbOdbBackendFree(handle *C.git_odb_backend) {
	ptr := unsafe.Pointer(handle)
	backend := odbBackendHandles.Untrack(ptr).(*OdbBackend)
	C.free(ptr)
	if backend.Free != nil {
		backend.Free()
	}
}

//export cbOdbBackendRead
func cbOdbBackendRead(
	ptr *unsafe.Pointer,
	csize *C.size_t,
	ctyp *C.git_otype,
	handle *C.git_odb_backend,
	oid *C.git_oid,
) C.int {
	backend := getOdbBackend(handle)
	if data, typ, err := backend.Read(newOidFromC(oid)); err != nil {
		return errToCode(err)
	} else {
		*csize = C.size_t(len(data))
		*ptr = C.git_odb_backend_malloc(handle, *csize)
		C.memcpy(*ptr, unsafe.Pointer(&data[0]), *csize)
		*ctyp = C.git_otype(typ)
	}
	return C.int(git.ErrOk)
}

//export cbOdbBackendReadHeader
func cbOdbBackendReadHeader(
	csize *C.size_t,
	ctyp *C.git_otype,
	handle *C.git_odb_backend,
	oid *C.git_oid,
) C.int {
	backend := getOdbBackend(handle)
	if sz, typ, err := backend.ReadHeader(newOidFromC(oid)); err != nil {
		return errToCode(err)
	} else {
		*csize = C.size_t(sz)
		*ctyp = C.git_otype(typ)
	}
	return C.int(git.ErrOk)
}

//export cbOdbBackendWrite
func cbOdbBackendWrite(
	handle *C.git_odb_backend,
	oid *C.git_oid,
	data unsafe.Pointer,
	size C.size_t,
	typ C.git_otype,
) C.int {
	backend := getOdbBackend(handle)
	if err := backend.Write(
		newOidFromC(oid),
		C.GoBytes(data, C.int(size)),
		git.ObjectType(typ)); err != nil {
		return errToCode(err)
	}
	return C.int(git.ErrOk)
}

//export cbOdbBackendExists
func cbOdbBackendExists(
	handle *C.git_odb_backend,
	oid *C.git_oid,
) C.int {
	backend := getOdbBackend(handle)
	if exists := backend.Exists(newOidFromC(oid)); exists {
		return 1
	}
	return 0
}

//export cbOdbBackendForEach
func cbOdbBackendForEach(
	handle *C.git_odb_backend,
	cb C.git_odb_foreach_cb,
	payload unsafe.Pointer,
) C.int {
	backend := getOdbBackend(handle)
	if err := backend.ForEach(func(id *git.Oid) error {
		ret := C._go_git_odb_backend_foreach(cb, (*C.git_oid)(unsafe.Pointer(id)), payload)
		if ret != 0 {
			return git.MakeGitError2(int(ret))
		}
		return nil
	}); err != nil {
		return errToCode(err)
	}
	return 0
}
