package main

import (
	"bazil.org/fuse"
	"time"
	"os"
)

var _ TX = (*WriteTX)(nil)
var _ TX = (*SetAttrTX)(nil)
var _ TX = (*MkdirTX)(nil)
var _ TX = (*RemoveTX)(nil)
var _ TX = (*RenameTX)(nil)
var _ TX = (*CreateTX)(nil)

// transaction interface
type TX interface {
	run()	(uint64, error)
}

type CreateTX struct {
	nid		uint64
	Name	string
	Mode	os.FileMode
}

func (tx CreateTX) run() (uint64, error) {
	p_call("Create TX\n")

	p := nodeMap[tx.nid]

	if p.kids[tx.Name] != nil {
		// already exists...
		return 0, fuse.EEXIST
	}

	// make the new file
	var n = new(DNode)
	n.init_node(tx.Name, tx.Mode)
	p.kids[tx.Name] = n

	return n.nid, nil
}

type WriteTX struct {
	nid		uint64
	Offset	int64
	Data	[]byte
}

func (tx WriteTX) run() (uint64, error) {
	p_call("Write TX\n")

	n := nodeMap[tx.nid]
	p_out("%v %v %v %v\n", n.attr.Size, tx.Offset, len(tx.Data), tx.Data)
	newlen := n.attr.Size

	if uint64(tx.Offset+int64(len(tx.Data))) > newlen {
		newlen = uint64(tx.Offset + int64(len(tx.Data)))
	}

	// maybe not the best in terms of efficiency but it makes the code really easy
	// nicely takes care of the case where offset > file len
	// basically just makes a new slice and copies the write data in
	new_data := n.data
	if newlen > n.attr.Size {
		new_data = append(new_data, make([]uint8, newlen-n.attr.Size)...)
	}
	copy(new_data[tx.Offset:], tx.Data)
	n.attr.Size = uint64(len(new_data))
	n.attr.Blocks = n.attr.Size/512 + 1
	n.data = new_data
	n.attr.Mtime = time.Now() // update mtime

	return uint64(len(tx.Data)), nil
}

type SetAttrTX struct {
	nid		uint64
	Valid	fuse.SetattrValid
	Size	uint64
	Atime	time.Time
	Mtime time.Time
	Mode	os.FileMode
	Uid		uint32
	Gid		uint32
}

func (tx SetAttrTX) run() (uint64, error) {
	p_call("SetAttr TX\n")

	n := nodeMap[tx.nid]

	// this code is really stupid, is there a better way to do this?
	if tx.Valid.Size() {
		if tx.Size < n.attr.Size { // update size if we are smaller
			n.data = n.data[:tx.Size]
		}

		n.attr.Size = tx.Size
		n.attr.Blocks = tx.Size/512 + 1
	}
	if tx.Valid.Atime() {
		n.attr.Atime = tx.Atime
	}
	if tx.Valid.Mtime() {
		n.attr.Mtime = tx.Mtime
	}
	if tx.Valid.Mode() {
		n.attr.Mode = tx.Mode
	}
	if tx.Valid.Uid() {
		n.attr.Uid = tx.Uid
	}
	if tx.Valid.Gid() {
		n.attr.Gid = tx.Gid
	}

	return 0, nil
}

type MkdirTX struct {
	nid		uint64
	Name	string
	Mode	os.FileMode
}

func (tx MkdirTX) run() (uint64, error) {
	p_call("Mkdir TX\n")
	
	p := nodeMap[tx.nid]

	if p.kids[tx.Name] != nil {
		// already exists...
		return 0, fuse.EEXIST
	}

	// make the new dir
	var n = new(DNode)
	n.init_node(tx.Name, os.ModeDir|tx.Mode)
	p.kids[tx.Name] = n

	return n.nid, nil
}

type RemoveTX struct {
	nid		uint64
	Name	string
	Dir		bool
}

func (tx RemoveTX) run() (uint64, error) {
	p_call("Remove TX\n")

	n := nodeMap[tx.nid]

	if n.kids[tx.Name] == nil {
		// if file doesn't exist, error
		return 0, fuse.ENOENT
	}

	if n.kids[tx.Name].attr.Mode.IsDir() != tx.Dir {
		// make sure IsDir and Dir match
		return 0, fuse.EPERM
	}

	// update nlinks
	n.kids[tx.Name].attr.Nlink -= 1

	// delete the file by removing from kids
	delete(n.kids, tx.Name)

	return 0, nil
}

type RenameTX struct {
	nid		uint64
	OldName	string
	NewName	string
	ndid  uint64 // nid of dir to move to
}

func (tx RenameTX) run() (uint64, error) {
	p_call("Rename TX\n")

	n := nodeMap[tx.nid]
	dir := nodeMap[tx.ndid]

	if n.kids[tx.OldName] == nil {
		// if file doesnt exist, error
		return 0, fuse.ENOENT
	}

	// move from olddir/oldname to newdir/newname
	dir.kids[tx.NewName] = n.kids[tx.OldName]
	dir.kids[tx.NewName].name = tx.NewName
	delete(n.kids, tx.OldName)

	return 0, nil
}