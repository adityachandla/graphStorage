package main

import (
	"fmt"
)

type nodeId int

const SIZE_INT = 4

type Graph struct {
	adjacency        [][]nodeId
	reverseAdjacency [][]nodeId
}

func NewGraph() *Graph {
	return &Graph{}
}

// In this model, how do you remove vertices?
// Well just mark them deleted and then move on. We can later
// reclaim the space if needed.
func (g *Graph) createNode() nodeId {
	curr := nodeId(len(g.adjacency)) //Length is one more than the highest index
	g.adjacency = append(g.adjacency, []nodeId{})
	g.reverseAdjacency = append(g.reverseAdjacency, []nodeId{})
	return curr
}

func (g *Graph) createEdge(src, dest nodeId) {
	g.adjacency[src] = append(g.adjacency[src], dest)
	g.reverseAdjacency[dest] = append(g.reverseAdjacency[dest], src)
}

func (g *Graph) convertToDiskFormat() *graphRepr {
	result := &graphRepr{}
	//The following two may be different for segmented graphs
	result.startNodeId = 0
	result.endNodeId = nodeId(len(g.adjacency)) - 1

	result.numNodes = int(result.endNodeId) - int(result.startNodeId) + 1
	result.nodeInfos = make([]nodeInfo, result.numNodes)
	//There are three header fields for graphRepr and 3 fields for nodeInfo
	startOffset := SIZE_INT * (GRAPH_REPR_NUM_INTS + NODE_INFO_NUM_INTS*result.numNodes)
	totalEdges := 0
	for i := 0; i < result.numNodes; i++ { //TODO This works but doesn't look right
		result.nodeInfos[i].id = nodeId(i) //Not the case for segmented graphs
		result.nodeInfos[i].outgoingOffset = startOffset
		startOffset += SIZE_INT * (len(g.adjacency[i]))
		totalEdges += len(g.adjacency[i])
		result.nodeInfos[i].incomingOffset = startOffset
		startOffset += SIZE_INT * (len(g.reverseAdjacency[i]))
	}
	result.nodes = make([]nodeId, 2*totalEdges)
	nodesIdx := 0
	for i := 0; i < result.numNodes; i++ {
		nodesIdx += copy(result.nodes[nodesIdx:], g.adjacency[i])
		nodesIdx += copy(result.nodes[nodesIdx:], g.reverseAdjacency[i])
	}
	return result
}

// There are three header fields
const GRAPH_REPR_NUM_INTS = 3

// This contains the on-disk representation of the forward and backward
// adjacency lists.
type graphRepr struct {
	startNodeId nodeId     //L1 Header: Min node Id
	endNodeId   nodeId     //L1 Header: Max node Id
	numNodes    int        //L1 Header: Number of nodes to read
	nodeInfos   []nodeInfo //L2 Header: Node information
	nodes       []nodeId   //Body: CSR
}

// Number of ints in nodeInfo field
const NODE_INFO_NUM_INTS = 3

// This stores the node id, byte offset of outgoing nodes, byte
// offset of incoming nodes and the end offset.
type nodeInfo struct {
	id             nodeId //Needed in case of deleted nodes
	outgoingOffset int
	incomingOffset int
}

func (ni *nodeInfo) convertToBytes() []byte {
	res := make([]byte, 0, NODE_INFO_NUM_INTS*SIZE_INT)
	res = append(res, intToByteArray(ni.id)...)
	res = append(res, intToByteArray(ni.outgoingOffset)...)
	res = append(res, intToByteArray(ni.incomingOffset)...)
	return res
}

func (gr *graphRepr) computeSize() int {
	size := GRAPH_REPR_NUM_INTS * SIZE_INT
	size += NODE_INFO_NUM_INTS * SIZE_INT * len(gr.nodeInfos) * SIZE_INT
	size += len(gr.nodes) * SIZE_INT
	return size
}

func (gr *graphRepr) convertToBytes() []byte {
	res := make([]byte, 0, gr.computeSize())
	res = append(res, intToByteArray(gr.startNodeId)...)
	res = append(res, intToByteArray(gr.endNodeId)...)
	res = append(res, intToByteArray(gr.endNodeId)...)
	for _, nInfo := range gr.nodeInfos {
		res = append(res, nInfo.convertToBytes()...)
	}
	for _, v := range gr.nodes {
		res = append(res, intToByteArray(v)...)
	}
	return res
}

// The function is generic because it should work for ints
// as well as types whose underlying type is int like nodeId
func intToByteArray[T ~int](i T) []byte {
	v := T((1 << 9) - 1)
	res := make([]byte, 4)
	res[0] = byte(i & v)
	v <<= 8
	res[1] = byte((i & v) >> 8)
	v <<= 8
	res[2] = byte((i & v) >> 16)
	v <<= 8
	res[3] = byte((i & v) >> 24)
	return res
}

func byteArrayToInt(arr []byte) int {
	res := int(0)
	res |= int(arr[0])
	res |= int(arr[1]) << 8
	res |= int(arr[2]) << 16
	res |= int(arr[3]) << 24
	return res
}

func main() {
	fmt.Println("vim-go")
}
