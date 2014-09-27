package merkledag

import (
	"fmt"

	"github.com/jbenet/go-ipfs/Godeps/_workspace/src/code.google.com/p/goprotobuf/proto"
	"github.com/op/go-logging"

	mh "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-multihash"
	blocks "github.com/jbenet/go-ipfs/blocks"
	bserv "github.com/jbenet/go-ipfs/blockservice"
	u "github.com/jbenet/go-ipfs/util"
)

var log = logging.MustGetLogger("merkledag")

// NodeMap maps u.Keys to Nodes.
// We cannot use []byte/Multihash for keys :(
// so have to convert Multihash bytes to string (u.Key)
type NodeMap map[u.Key]*Node

// Node represents a node in the IPFS Merkle DAG.
// nodes have opaque data and a set of navigable links.
type Node struct {
	Links []*Link
	Data  []byte

	// cache encoded/marshaled value
	encoded []byte
}

// Link represents an IPFS Merkle DAG Link between Nodes.
type Link struct {
	// utf string name. should be unique per object
	Name string // utf8

	// cumulative size of target object
	Size uint64

	// multihash of the target object
	Hash mh.Multihash

	// a ptr to the actual node for graph manipulation
	Node *Node
}

// AddNodeLink adds a link to another node.
func (n *Node) AddNodeLink(name string, that *Node) error {
	s, err := that.Size()
	if err != nil {
		return err
	}

	h, err := that.Multihash()
	if err != nil {
		return err
	}

	n.Links = append(n.Links, &Link{
		Name: name,
		Size: s,
		Hash: h,
		Node: that,
	})
	return nil
}

func (n *Node) RemoveNodeLink(name string) error {
	for i, l := range n.Links {
		if l.Name == name {
			n.Links = append(n.Links[:i], n.Links[i+1:]...)
			return nil
		}
	}
	return u.ErrNotFound
}

// Size returns the total size of the data addressed by node,
// including the total sizes of references.
func (n *Node) Size() (uint64, error) {
	b, err := n.Encoded(false)
	if err != nil {
		return 0, err
	}

	s := uint64(len(b))
	for _, l := range n.Links {
		s += l.Size
	}
	return s, nil
}

// Multihash hashes the encoded data of this node.
func (n *Node) Multihash() (mh.Multihash, error) {
	b, err := n.Encoded(false)
	if err != nil {
		return nil, err
	}

	return u.Hash(b)
}

// Key returns the Multihash as a key, for maps.
func (n *Node) Key() (u.Key, error) {
	h, err := n.Multihash()
	return u.Key(h), err
}

// Recursively update all hash links and size values in the tree
func (n *Node) Update() error {
	log.Debug("node update")
	for _, l := range n.Links {
		if l.Node != nil {
			err := l.Node.Update()
			if err != nil {
				return err
			}
			nhash, err := l.Node.Multihash()
			if err != nil {
				return err
			}
			l.Hash = nhash
			size, err := l.Node.Size()
			if err != nil {
				return err
			}
			l.Size = size
		}
	}
	_, err := n.Encoded(true)
	return err
}

// DAGService is an IPFS Merkle DAG service.
// - the root is virtual (like a forest)
// - stores nodes' data in a BlockService
// TODO: should cache Nodes that are in memory, and be
//       able to free some of them when vm pressure is high
type DAGService struct {
	Blocks *bserv.BlockService
}

// Add adds a node to the DAGService, storing the block in the BlockService
func (n *DAGService) Add(nd *Node) (u.Key, error) {
	k, _ := nd.Key()
	u.DOut("DagService Add [%s]\n", k.Pretty())
	if n == nil {
		return "", fmt.Errorf("DAGService is nil")
	}

	d, err := nd.Encoded(false)
	if err != nil {
		return "", err
	}

	b, err := blocks.NewBlock(d)
	if err != nil {
		return "", err
	}

	return n.Blocks.AddBlock(b)
}

func (n *DAGService) AddRecursive(nd *Node) error {
	_, err := n.Add(nd)
	if err != nil {
		fmt.Printf("AddRecursive Error: %s\n", err)
		return err
	}

	for _, link := range nd.Links {
		if link.Node != nil {
			err := n.AddRecursive(link.Node)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Get retrieves a node from the DAGService, fetching the block in the BlockService
func (n *DAGService) Get(k u.Key) (*Node, error) {
	if n == nil {
		return nil, fmt.Errorf("DAGService is nil")
	}

	b, err := n.Blocks.GetBlock(k)
	if err != nil {
		return nil, err
	}

	return Decoded(b.Data)
}

func FilePBData(data []byte) []byte {
	pbfile := new(PBData)
	typ := PBData_File
	pbfile.Type = &typ
	pbfile.Data = data

	data, err := proto.Marshal(pbfile)
	if err != nil {
		//this really shouldnt happen, i promise
		panic(err)
	}
	return data
}

func FolderPBData() []byte {
	pbfile := new(PBData)
	typ := PBData_Directory
	pbfile.Type = &typ

	data, err := proto.Marshal(pbfile)
	if err != nil {
		//this really shouldnt happen, i promise
		panic(err)
	}
	return data
}

func WrapData(b []byte) []byte {
	pbdata := new(PBData)
	typ := PBData_Raw
	pbdata.Data = b
	pbdata.Type = &typ

	out, err := proto.Marshal(pbdata)
	if err != nil {
		// This shouldnt happen. seriously.
		panic(err)
	}

	return out
}
