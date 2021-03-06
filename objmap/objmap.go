// objmap package stores objects as a mapping of their distance from
// the node. Each object is stored in a vector contained by ObjNode at
// its correct distance bucket. The `key' in the <Key, Value> pair is
// the hash of the object being referred to while the `Value' is the
// address of the peer where it can be found.
package objmap

import (
	"fmt"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// ObjNode is a single bucket that stores the objects in a vector.
// Each object has its index mapped using its hash within the same node.
type ObjNode struct {
	Nmap map[[utils.HASHSZ]byte]int
	Nvec []*pkt.ObjAddr
}

// ObjMap is the high level mapping of distances from the node to seprate buckets.
type ObjMap struct {
	omap map[int]*ObjNode
}

// Init initialized the ObjMap.
func (omap_p *ObjMap) Init() {
	omap_p.omap = make(map[int]*ObjNode)
}

// Insert inserts the object accoring to its distance from the node.
func (omap_p *ObjMap) Insert(srchash [utils.HASHSZ]byte, obj pkt.ObjAddr) {
	var indx int = utils.GetDist(&srchash, &obj.Hash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if _, ok := node_p.Nmap[obj.Hash]; !ok {
			node_p.Nmap[obj.Hash] = len(node_p.Nvec)
			node_p.Nvec = append(node_p.Nvec, &pkt.ObjAddr{Hash: obj.Hash, Addr: obj.Addr})
		}
	} else {
		var node_p *ObjNode = &ObjNode{Nmap: make(map[[utils.HASHSZ]byte]int),
			Nvec: []*pkt.ObjAddr{&pkt.ObjAddr{Hash: obj.Hash, Addr: obj.Addr}}}
		node_p.Nmap[obj.Hash] = 0
		omap_p.omap[indx] = node_p
	}
}

// Get fetches the object if it exists.
func (omap_p *ObjMap) Get(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if k, ok := node_p.Nmap[*dsthash]; ok {
			return node_p.Nvec[k], nil
		}
	}

	return nil, fmt.Errorf("Not Found!")
}

// GetAll returns all the objects that are closer to the given hash from the src hash as a list of
// pointers to a slice of objects.
func (omap_p *ObjMap) GetAll(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) []*[]*pkt.ObjAddr {
	var indx int = utils.GetDist(srchash, dsthash)
	var ret []*[]*pkt.ObjAddr
	for i := 1; i <= indx; i++ {
		if node_p, ok := omap_p.omap[i]; ok {
			ret = append(ret, &node_p.Nvec)
		}
	}

	return ret
}
