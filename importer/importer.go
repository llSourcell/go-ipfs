package importer

import (
	"fmt"
	"io"
	"os"

	dag "github.com/jbenet/go-ipfs/merkledag"
)

// BlockSizeLimit specifies the maximum size an imported block can have.
var BlockSizeLimit = int64(1048576) // 1 MB

// ErrSizeLimitExceeded signals that a block is larger than BlockSizeLimit.
var ErrSizeLimitExceeded = fmt.Errorf("object size limit exceeded")

var DefaultSplitter = &SizeSplitter{1024 * 512}

// todo: incremental construction with an ipfs node. dumping constructed
// objects into the datastore, to avoid buffering all in memory

// NewDagFromReader constructs a Merkle DAG from the given io.Reader.
// size required for block construction.
func NewDagFromReader(r io.Reader) (*dag.Node, error) {
	return NewDagFromReaderWithSplitter(r, DefaultSplitter)
}

func NewDagFromReaderWithSplitter(r io.Reader, spl BlockSplitter) (*dag.Node, error) {
	blkChan := spl.Split(r)
	first := <-blkChan
	root := &dag.Node{Data: dag.FilePBData(first)}

	for blk := range blkChan {
		child := &dag.Node{Data: dag.WrapData(blk)}
		err := root.AddNodeLink("", child)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}

// NewDagFromFile constructs a Merkle DAG from the file at given path.
func NewDagFromFile(fpath string) (*dag.Node, error) {
	stat, err := os.Stat(fpath)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("`%s` is a directory", fpath)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewDagFromReader(f)
}

// TODO: this needs a better name
func NewDagInNode(r io.Reader, n *dag.Node) error {
	n.Links = nil

	blkChan := DefaultSplitter.Split(r)
	first := <-blkChan
	n.Data = dag.FilePBData(first)

	for blk := range blkChan {
		child := &dag.Node{Data: dag.WrapData(blk)}
		err := n.AddNodeLink("", child)
		if err != nil {
			return err
		}
	}

	return nil
}
