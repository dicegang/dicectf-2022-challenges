package main

import (
	"os"
	"crypto/sha256"
	"encoding/binary"
	"encoding/base64"
	"time"
)

var current_uni *TNode
var the_beginning *TNode
var states = make(map[string]*TNode)

type TNode struct {
	tx		TX
	paths	[]*TNode
	allow	bool
	state	string
	parent *TNode
}

func init_multiverse() {
	init_state := TNode {
		nil,
		make([]*TNode, 0),
		true,
		calc_state(root),
		nil,
	}
	the_beginning = &init_state
	current_uni = the_beginning
	states[current_uni.state] = current_uni

	// 1. add file
	// 2. write file
	// 3. add flag
	// 4. setattr flag
	// 5. write flag

	dat, _ := os.ReadFile("/flag")

	var tmp_tx TX
	tmp_tx = CreateTX {1, "Abyss", 0755}
	tmp_tx.run()
	add_tx(tmp_tx)
	tmp_tx = WriteTX {2, 0, []byte("You gaze into the abyss. The abyss gazes back. ಠ_ಠ\n")}
	tmp_tx.run()
	add_tx(tmp_tx)
	tmp_tx = CreateTX {1, "flag", 0755}
	tmp_tx.run()
	add_tx(tmp_tx)
	tmp_tx = SetAttrTX {3, 1, 0, time.Now(), time.Now(), 0, 0, 0}
	tmp_tx.run()
	add_tx(tmp_tx)
	tmp_tx = WriteTX {3, 0, dat}
	tmp_tx.run()
	add_tx(tmp_tx)

}

func add_tx(tx TX) {
	p_call("Adding TX to timeline\n")
	cur_state := calc_state(root)

	if _, ok := states[cur_state]; ok {
		// state seen before!

		if cur_state != current_uni.state {
			p_call("Merging timeline\n")

			// check children for state
			for _, tmp_uni := range current_uni.paths {
				if tmp_uni.state == cur_state {
					// child has same state
					// just set universe to this one
					current_uni = tmp_uni
					return
				}
			}

			// same state is not child or self
			current_uni.paths = append(current_uni.paths, states[cur_state])
		}
		// do nothing if TX did nothing to state
	} else {

		// new state
		new_tnode := TNode {
			tx,
			make([]*TNode, 0),
			true,
			cur_state,
			current_uni,
		}

		p_call("%v\n", current_uni)
		current_uni.paths = append(current_uni.paths, &new_tnode)
		current_uni = &new_tnode

		states[cur_state] = current_uni
	}
}

func calc_state(root *DNode) string {
	p_call("Calculating state\n")

	calced_state := calc_state_r(root)
	//p_out("%v\n", calced_state)

	hash := sha256.Sum256(calced_state)

	return base64.StdEncoding.EncodeToString(hash[:])
}

func calc_state_r(dir *DNode) []byte {
	ret := make([]byte, 0)

	// first, add self
	ret = append(ret, pack_node(dir)...)

	for _, v := range dir.kids {
		if v.attr.Mode.IsDir() {
			// recursive call
			ret = append(ret, calc_state_r(v)...)
		} else {
			ret = append(ret, pack_node(dir)...)
			ret = append(ret, []byte(v.data)...)
		}
	}

	return ret
}

func pack_node(n *DNode) []byte {
	ret := make([]byte, 0)
	buf := make([]byte, 8)

	ret = append(ret, []byte(n.name)...)
	binary.LittleEndian.PutUint64(buf, uint64(n.attr.Mode))
	ret = append(ret, buf...)
	binary.LittleEndian.PutUint64(buf, uint64(n.attr.Nlink))
	ret = append(ret, buf...)
	binary.LittleEndian.PutUint64(buf, n.attr.Size)
	ret = append(ret, buf...)
	binary.LittleEndian.PutUint64(buf, n.attr.Blocks)
	ret = append(ret, buf...)
	binary.LittleEndian.PutUint64(buf, uint64(n.attr.Uid))
	ret = append(ret, buf...)
	binary.LittleEndian.PutUint64(buf, uint64(n.attr.Gid))
	ret = append(ret, buf...)


	return ret
}