#include "_cgo_export.h"
#include <git2.h>
#include <git2/sys/refdb_backend.h>

void _go_git_populate_refdb_backend_cb(git_refdb_backend *backend)
{
    backend->exists = (int (*)(int *exists, git_refdb_backend *backend, const char *ref_name))&cbRefdbBackendExists;
    backend->lookup = (int (*)(git_reference **out, git_refdb_backend *backend, const char *ref_name))&cbRefdbBackendLookup;
    backend->iterator = (int (*)(git_reference_iterator **iter, struct git_refdb_backend *backend, const char *glob))&cbRefdbBackendIterator;
    backend->write = (int (*)(git_refdb_backend *backend, const git_reference *ref, int force, const git_signature *who, const char *message, const git_oid *old, const char *old_target))&cbRefdbBackendWrite;
    backend->rename = (int (*)(
		git_reference **out, git_refdb_backend *backend,
		const char *old_name, const char *new_name, int force,
		const git_signature *who, const char *message))&cbRefdbBackendRename;
    backend->del = (int (*)(git_refdb_backend *backend, const char *ref_name, const git_oid *old_id, const char *old_target))&cbRefdbBackendDel;
    backend->free = &cbRefdbBackendFree;
}

