#include "_cgo_export.h"
#include <git2.h>
#include <git2/sys/refdb_backend.h>

void _go_git_populate_reference_iterator_cb(git_reference_iterator* it)
{
    it->next = &cbReferenceIteratorNext;
    it->next_name = (int (*)(const char **ref_name, git_reference_iterator *iter))&cbReferenceIteratorNextName;
    it->free = &cbReferenceIteratorFree;
}
