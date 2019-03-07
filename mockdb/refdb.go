package mockdb

import (
	"sort"

	"github.com/niukuo/gitdb"

	git "gopkg.in/libgit2/git2go.v27"
)

type RefdbBackend = *refdbBackend
type refdbBackend struct {
	data map[string]reference
}

type reference struct {
	name   string
	typ    git.ReferenceType
	target string
}

type iterator struct {
	cur  int
	refs []reference
}

func (v *iterator) Next() (string, git.ReferenceType, string, error) {
	if v.cur >= len(v.refs) {
		return "", 0, "", git.MakeGitError2(int(git.ErrIterOver))
	}
	i := v.cur
	v.cur++
	return v.refs[i].name, git.ReferenceSymbolic, v.refs[i].target, nil
}

func (v *iterator) Free() {
}

func NewRefdbBackend() RefdbBackend {
	r := &refdbBackend{
		data: make(map[string]reference),
	}
	return r
}

func (r *refdbBackend) Track() *git.RefdbBackend {
	backend := gitdb.NewRefdbBackend()

	backend.Exists = r.Exists
	backend.Lookup = r.Lookup
	backend.Iterator = r.Iterator
	backend.Rename = r.Rename
	backend.Write = r.Write
	backend.Del = r.Del

	return backend.Track()
}

func (r *refdbBackend) Exists(ref_name string) (bool, error) {
	_, ok := r.data[ref_name]
	return ok, nil
}

func (r *refdbBackend) Lookup(ref_name string) (git.ReferenceType, string, error) {
	v, ok := r.data[ref_name]
	if !ok {
		return 0, "", git.MakeGitError2(int(git.ErrNotFound))
	}
	return v.typ, v.target, nil
}

func (r *refdbBackend) Iterator(glob string) (gitdb.ReferenceIterator, error) {
	refs := make([]reference, 0, len(r.data))
	for _, v := range r.data {
		refs = append(refs, v)
	}
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].name < refs[j].name
	})
	return &iterator{
		refs: refs,
	}, nil
}

func (r *refdbBackend) checkOld(name string, old_id *git.Oid, old_target string) error {
	switch {
	case old_id != nil:
		if v, ok := r.data[name]; !ok || v.typ != git.ReferenceOid || v.target != old_id.String() {
			return git.MakeGitError2(int(git.ErrModified))
		}
	case old_target != "":
		if v, ok := r.data[name]; !ok || v.typ != git.ReferenceSymbolic || v.target != old_target {
			return git.MakeGitError2(int(git.ErrModified))
		}
	}
	return nil
}

func (r *refdbBackend) Write(ref gitdb.Reference, force bool, who *git.Signature, message string, old_id *git.Oid, old_target string) error {
	if err := r.checkOld(ref.Name(), old_id, old_target); err != nil {
		return err
	}
	v := reference{
		name: ref.Name(),
		typ:  ref.Type(),
	}
	if v.typ == git.ReferenceOid {
		v.target = ref.Target().String()
	} else {
		v.target = ref.SymbolicTarget()
	}
	r.data[ref.Name()] = v
	return nil
}

func (r *refdbBackend) Rename(old_name, new_name string, force bool, who *git.Signature, message string) (git.ReferenceType, string, error) {
	v, ok := r.data[old_name]
	if !ok {
		return 0, "", git.MakeGitError2(int(git.ErrNotFound))
	}
	if !force {
		if _, ok := r.data[new_name]; ok {
			return 0, "", git.MakeGitError2(int(git.ErrExists))
		}
	}
	delete(r.data, old_name)
	r.data[new_name] = v
	return v.typ, v.target, nil
}

func (r *refdbBackend) Del(ref_name string, old_id *git.Oid, old_target string) error {
	if err := r.checkOld(ref_name, old_id, old_target); err != nil {
		return err
	}
	delete(r.data, ref_name)
	return nil
}
