package data_struct

import (
	"bufio"
	"crypto/sha256"
	"os"
	"strconv"
)

type Node struct {
	Index    int
	Value    string
	RootTree *MHTree
}

// the struct for MHTree
type MHTree struct {
	Length   int
	Nodes    []Node
	rootHash string
}

func (t *MHTree) GetRootHash() string {
	// yes or not storage,save in hash
	t.rootHash = t.Nodes[1].getNodeHash()
	return t.rootHash
}

func (n *Node) getNodeHash() string {
	// the node for leaf hash
	if n.Value == "" {
		return calDataHash(n.Value)
	}
	//if is not the node for leaf,then cal the while hash
	return calDataHash(n.RootTree.Nodes[n.Index*2].getNodeHash() + n.RootTree.Nodes[n.Index*2+1].getNodeHash())
}

// cal the data for hash
func calDataHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return string(hash.Sum(nil))
}

//create the merkel tree from the struct file

func CreateMHTree(filename string) MHTree {
	var tree MHTree
	//open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// get the file reader
	buf := bufio.NewReader(file)
	//read the first line, get the num for the leaf(require the count for leaf is 2*pow())
	dataCountStr, _, _ := buf.ReadLine()
	dataCount, _ := strconv.Atoi(string(dataCountStr))

	// judge the 2*pow()
	level := 0
	for i := 1; ; i++ {
		if 2<<i == dataCount {
			level = i
			break
		}
	}
	//create the merkel tree

	//set the value for the leaf
	for i := 2 << level; i < tree.Length; i++ {
		str, _, _ := buf.ReadLine()
		tree.Nodes[i].Index = i
		tree.Nodes[i].RootTree = &tree
		tree.Nodes[i].Value = string(str)
	}
	return tree
}
