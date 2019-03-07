package mockdb

import (
	"sort"

	"github.com/niukuo/gitdb"

	git "gopkg.in/libgit2/git2go.v27"
)

type OdbBackend = *odbBackend

type odbBackend struct {
	data map[string]object
}

type object struct {
	otype git.ObjectType
	data  []byte
}

func NewOdbBackend() OdbBackend {
	o := &odbBackend{
		data: make(map[string]object),
	}
	return o
}

func (o *odbBackend) Track() *git.OdbBackend {
	backend := gitdb.NewOdbBackend()

	backend.Read = o.Read
	backend.ReadHeader = o.ReadHeader
	backend.Exists = o.Exists
	backend.ForEach = o.ForEach
	backend.Write = o.Write

	return backend.Track()
}

func (o *odbBackend) Read(oid *git.Oid) ([]byte, git.ObjectType, error) {
	v, ok := o.data[oid.String()]
	if !ok {
		return nil, 0, git.MakeGitError2(int(git.ErrNotFound))
	}
	return v.data, v.otype, nil
}

func (o *odbBackend) ReadHeader(oid *git.Oid) (uint64, git.ObjectType, error) {
	v, ok := o.data[oid.String()]
	if !ok {
		return 0, 0, git.MakeGitError2(int(git.ErrNotFound))
	}
	return uint64(len(v.data)), v.otype, nil
}

func (o *odbBackend) Write(oid *git.Oid, data []byte, otype git.ObjectType) error {
	o.data[oid.String()] = object{
		otype: otype,
		data:  append(make([]byte, 0, len(data)), data...),
	}
	return nil
}

func (o *odbBackend) Exists(oid *git.Oid) bool {
	_, ok := o.data[oid.String()]
	return ok
}

func (o *odbBackend) ForEach(callback git.OdbForEachCallback) error {
	keys := make([]string, 0, len(o.data))
	for key := range o.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		oid, err := git.NewOid(key)
		if err != nil {
			return err
		}
		if err := callback(oid); err != nil {
			return err
		}
	}
	return nil
}
