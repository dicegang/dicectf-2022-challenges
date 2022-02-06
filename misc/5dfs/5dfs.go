package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
	"encoding/binary"
	"golang.org/x/sys/unix"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	. "github.com/mattn/go-getopt"
	"golang.org/x/net/context"
	"sync"
)

//=============================================================================

type DNode struct {
	nid   uint64
	name  string
	attr  fuse.Attr
	valid bool
	kids  map[string]*DNode
	data  []uint8
}

type Dfs struct{}

var root *DNode

var _ fs.Node = (*DNode)(nil)
var _ fs.FS = (*Dfs)(nil)

var debug = false
var mountPoint = "5dfs"
var nodeMap = make(map[uint64]*DNode)
var nextNid uint64 = 1
var conn *fuse.Conn
var lock sync.Mutex

func p_out(s string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(s, args...)
}

func p_err(s string, args ...interface{}) {
	fmt.Printf(s, args...)
}

func p_call(s string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(s, args...)
}

// Implement:
func (Dfs) Root() (n fs.Node, err error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("Root\n")
	// return root
	return root, nil
}

func (n *DNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Attr\n")

	if !n.valid {
		return fuse.ENOENT
	}
	// same as getattr, not really sure why we have both...
	*attr = n.attr
	return nil
}

func (n *DNode) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("Open\n")
	if !n.valid {
		return nil, fuse.ENOENT
	}

	resp.Flags |= fuse.OpenDirectIO
	return n, nil
}


func (n *DNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("Lookup\n")

	if !n.valid {
		return nil, fuse.ENOENT
	}

	// return kid if exists
	if n.kids[name] != nil {
		return n.kids[name], nil
	}
	return nil, fuse.ENOENT
}

func (n *DNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("ReadDirAll\n")

	if !n.valid {
		return nil, fuse.ENOENT
	}

	if n.attr.Mode|os.ModeDir == 0 {
		return nil, fuse.EPERM
	}

	ret := make([]fuse.Dirent, 0)

	// iterate over kids and return all of them as Dirents
	for k, v := range n.kids {
		dirent := new(fuse.Dirent)
		dirent.Inode = v.attr.Inode
		if v.attr.Mode.IsDir() {
			dirent.Type = fuse.DT_Dir
		} else {
			dirent.Type = fuse.DT_File
		}
		dirent.Name = k
		ret = append(ret, *dirent)
	}

	return ret, nil
}

func (n *DNode) Getattr(ctx context.Context, req *fuse.GetattrRequest, resp *fuse.GetattrResponse) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Getattr\n")

	if !n.valid {
		return fuse.ENOENT
	}

	// just return attr
	resp.Attr = n.attr
	return nil
}

func (n *DNode) Getxattr(ctx context.Context, req *fuse.GetxattrRequest, resp *fuse.GetxattrResponse) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Getxattr\n")

	if n.nid != root.nid || req.Name != "timelines" {
		p_out("%v %v %v\n", n.nid, root.nid, req.Name)
		return fuse.ENOSYS
	}

	// return # of child universes
	// check buf is big enough
	if req.Size < 8 {
		return unix.E2BIG
	}

	resp.Xattr = make([]byte, 8)
	binary.LittleEndian.PutUint64(resp.Xattr, uint64(len(current_uni.paths)))

	return nil
}

// must be defined or editing w/ vi or emacs fails. Doesn't have to do anything
func (n *DNode) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Fsync\n")

	if !n.valid {
		return fuse.ENOENT
	}

	// do nothing...
	return nil
}

func (n *DNode) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Setattr\n")

	if !n.valid {
		return fuse.ENOENT
	}

	if req.Valid.Mode() && req.Mode != 000 {
		return fuse.EPERM // no setting flag to readable
	}

	setattrTX := SetAttrTX {
		n.nid,
		req.Valid,
		req.Size,
		req.Atime,
		req.Mtime,
		req.Mode,
		req.Uid,
		req.Gid,
	}

	setattrTX.run() // no error to check

	resp.Attr = n.attr

	// add TX to timeline
	add_tx(setattrTX)

	return nil
}

func (n *DNode) Setxattr(ctx context.Context, req *fuse.SetxattrRequest) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Setxattr\n")

	if len(req.Xattr) != 8 {
		return unix.EINVAL
	}

	arg := binary.LittleEndian.Uint64(req.Xattr)

	if req.Name == "timeleap" {
		// jump to past (parent universe)
		// Xattr is packed int denoting how much to jump up
		if arg <= 0 {
			return unix.EINVAL
		}

		new_uni := current_uni
		for i := uint64(0); i < arg; i++ {
			if new_uni.parent == nil {
				return unix.EINVAL
			}
			new_uni = new_uni.parent
		}

		// must reset FS state and rerun TXs

		// invalidate all nodes
		for _, n := range nodeMap {
			if n.nid != root.nid {
				n.valid = false
			}
		}

		nodeMap = make(map[uint64]*DNode)
		nextNid = 1
		root.init_node("", os.ModeDir|0755)

		// run txs, rn it just resets fs entirely
		txs := make([]TX, 0)
		for tmp_uni := new_uni; tmp_uni.parent != nil; tmp_uni = tmp_uni.parent {
			txs = append(txs, tmp_uni.tx)
		}

		for i := len(txs)-1; i >= 0; i-- {
			txs[i].run()
		}
		
		current_uni = new_uni

	}	else if req.Name == "backtothefuture" {
		// jump to specific child universe
		// Xattr is packed int denoting which future universe to go to
		if arg >= uint64(len(current_uni.paths)) {
			return unix.EINVAL
		}

		// when travelling forward, we can simply run the TX
		new_uni := current_uni.paths[arg]

		if !new_uni.allow {
			return fuse.EPERM
		}

		new_uni.tx.run()
		current_uni = new_uni
	}

	// clears kernel cache
	exec.Command("sync").Run()
	os.WriteFile("/proc/sys/vm/drop_caches", []byte("3\n"), 0755)

	return nil
}

func (n *DNode) init_node(name string, mode os.FileMode) {
	p_call("init_node\n")

	n.nid = uint64(nextNid)
	nextNid++
	n.name = name
	n.valid = true
	n.data = make([]uint8, 0)
	n.kids = make(map[string]*DNode)

	n.attr.Inode = n.nid
	n.attr.Mode = mode
	time := time.Now()
	n.attr.Atime = time
	n.attr.Ctime = time
	n.attr.Mtime = time
	n.attr.Nlink = 1
	n.attr.Size = 0
	n.attr.Blocks = 0

	n.attr.Uid = uint32(os.Geteuid())
	n.attr.Gid = uint32(os.Geteuid())

	nodeMap[n.nid] = n
}

func (p *DNode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("Mkdir\n")

	if !p.valid {
		return nil, fuse.ENOENT
	}

	mkdirTX := MkdirTX {
		p.nid,
		req.Name,
		req.Mode,
	}

	newnid, err := mkdirTX.run()

	if err != nil {
		return nil, err
	}

	// add TX to timeline
	add_tx(mkdirTX)

	return nodeMap[newnid], nil
}

func (p *DNode) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("Create\n")

	if !p.valid {
		return nil, nil, fuse.ENOENT
	}

	createTX := CreateTX {
		p.nid,
		req.Name,
		req.Mode,
	}

	newnid, err := createTX.run()

	if err != nil {
		return nil, nil, err
	}

	// add TX to timeline
	add_tx(createTX)

	return nodeMap[newnid], nodeMap[newnid], nil
}

func (n *DNode) ReadAll(ctx context.Context) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()
	p_call("ReadAll\n")

	if !n.valid {
		return nil, fuse.ENOENT
	}

	if n.attr.Mode.Perm() == 0 {
		return nil, fuse.EPERM
	}

	n.attr.Atime = time.Now() // update atime
	// just return our entire data
	return n.data, nil
}

func (n *DNode) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Write %v %v %v\n", n.attr.Size, req.Offset, len(req.Data))

	if !n.valid {
		return fuse.ENOENT
	}

	writeTX := WriteTX {
		n.nid,
		req.Offset,
		append(make([]byte, 0, len(req.Data)), req.Data...),
	}

	// no error to check...
	sz, _ := writeTX.run()

	resp.Size = int(sz)

	// add TX to timeline
	add_tx(writeTX)

	return nil
}

func (n *DNode) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Flush\n")

	if !n.valid {
		return fuse.ENOENT
	}
	// no need to do anything here, as our writes are not buffered...
	return nil
}

func (n *DNode) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Remove\n")

	if !n.valid {
		return fuse.ENOENT
	}

	removeTX := RemoveTX {
		n.nid,
		req.Name,
		req.Dir,
	}

	_, err := removeTX.run()

	if err != nil {
		return err
	}

	// add TX to timeline
	add_tx(removeTX)

	return nil
}

func (n *DNode) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	lock.Lock()
	defer lock.Unlock()
	p_call("Rename\n")

	if !n.valid {
		return fuse.ENOENT
	}

	renameTX := RenameTX {
		n.nid,
		req.OldName,
		req.NewName,
		newDir.(*DNode).nid,
	}
	
	_, err := renameTX.run()

	if err != nil {
		return err
	}

	// add TX to timeline
	add_tx(renameTX)

	return nil
}

func main() {
	var flag int

	for {
		if flag = Getopt("dm:"); flag == EOF {
			break
		}

		switch flag {
		case 'd':
			debug = !debug
		case 'm':
			mountPoint = OptArg
		default:
			println("usage: main.go [-d | -m <mountpt>]", flag)
			os.Exit(1)
		}
	}
	p_out("mounting on %q, debug %v\n", mountPoint, debug)

	// create root node
	root = new(DNode)
	root.init_node("", os.ModeDir|0755)
	init_multiverse()

	p_out("root inode %d", int(root.attr.Inode))

	if _, err := os.Stat(mountPoint); err != nil {
		os.Mkdir(mountPoint, 0755)
	}
	fuse.Unmount(mountPoint)
	c, err := fuse.Mount(mountPoint, fuse.FSName("5dfs"), fuse.Subtype("5dfs"),
		fuse.LocalVolume(), fuse.VolumeName("5dfs"), fuse.AllowOther())
	if err != nil {
		log.Fatal(err)
	}

	conn = c

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func() {
		<-ch
		defer conn.Close()
		fuse.Unmount(mountPoint)
		os.Exit(1)
	}()

	err = fs.Serve(conn, Dfs{})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-conn.Ready
	if err := conn.MountError; err != nil {
		log.Fatal(err)
	}
}
