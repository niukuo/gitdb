#include "_cgo_export.h"
#include <git2.h>
#include <git2/sys/odb_backend.h>

void _go_git_populate_odb_backend_cb(git_odb_backend *backend)
{
    backend->read = (int (*)(void **, size_t *, git_otype *, git_odb_backend *, const git_oid *))&cbOdbBackendRead;
    backend->read_header = (int (*)(size_t *, git_otype *, git_odb_backend *, const git_oid *))&cbOdbBackendReadHeader;
    backend->write = (int (*)(git_odb_backend *, const git_oid *, const void *, size_t, git_otype))&cbOdbBackendWrite;
    backend->exists = (int (*)(git_odb_backend *, const git_oid *))&cbOdbBackendExists;
    backend->foreach = &cbOdbBackendForEach;
    backend->free = &cbOdbBackendFree;
}

int _go_git_odb_backend_foreach(git_odb_foreach_cb cb, const git_oid *id, void *payload)
{
    return cb(id, payload);
}
