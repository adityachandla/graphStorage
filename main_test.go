package main

import (
	"testing"
)

func getDummyGraph() *Graph {
	g := NewGraph()
	one := g.createNode()
	two := g.createNode()
	three := g.createNode()
	g.createEdge(one, three)
	g.createEdge(three, two)
	g.createEdge(two, one)
	g.createEdge(two, three)
	return g
}

func TestRepresentation(t *testing.T) {
	g := getDummyGraph()
	repr := g.convertToDiskFormat()
	if repr.startNodeId != 0 {
		t.Fail()
	}
	if repr.endNodeId != nodeId(len(g.adjacency)-1) {
		t.Fail()
	}
	if len(repr.nodeInfos) != len(g.adjacency) {
		t.Fail()
	}
	for i, info := range repr.nodeInfos {
		numOut := (info.incomingOffset - info.outgoingOffset) / SIZE_INT
		if numOut != len(g.adjacency[info.id]) {
			t.Fatalf("Wanted %d got %d\n", len(g.adjacency[info.id]), numOut)
			t.Fail()
		}
		if i != len(repr.nodeInfos)-1 {
			numIn := (repr.nodeInfos[i+1].outgoingOffset - info.incomingOffset) / SIZE_INT
			if numIn != len(g.reverseAdjacency[info.id]) {
				t.Fail()
			}
		}
	}
}

func TestConversions(t *testing.T) {
	arr := []int{2341234123, 2, 0, 11003}
	for _, v := range arr {
		res := byteArrayToInt(intToByteArray(v))
		if res != v {
			t.Fail()
		}
	}
}

func TestConvertToBytes(t *testing.T) {
	g := getDummyGraph()
	repr := g.convertToDiskFormat()
	bArray := repr.convertToBytes()

	//Try to access node 0's outgoing edges
	startIdx := repr.nodeInfos[0].outgoingOffset
	endIdx := repr.nodeInfos[0].incomingOffset
	arr := make([]int, 0)
	for i := startIdx; i < endIdx; i += 4 {
		v := byteArrayToInt(bArray[i : i+4])
		arr = append(arr, v)
	}
	if len(arr) != len(g.adjacency[0]) {
		t.Fatalf("Size mismatch, should be %d but is %d", len(g.adjacency[0]), len(arr))
	}
	for i := 0; i < len(arr); i++ {
		if nodeId(arr[i]) != g.adjacency[0][i] {
			t.Fatalf("Wanted %d got %d\n", g.adjacency[0][i], arr[i])
		}
	}

	//Try to access 3's incoming edges
	startIdx = repr.nodeInfos[2].incomingOffset
	arr = make([]int, 0)
	for i := startIdx; i < len(bArray); i += 4 {
		v := byteArrayToInt(bArray[i : i+4])
		arr = append(arr, v)
	}
	if len(arr) != len(g.reverseAdjacency[2]) {
		t.Fatalf("Size mismatch, should be %d but is %d", len(g.adjacency[2]), len(arr))
	}
	for i := 0; i < len(arr); i++ {
		if nodeId(arr[i]) != g.reverseAdjacency[2][i] {
			t.Fatalf("Wanted %d got %d\n", g.reverseAdjacency[2][i], arr[i])
		}
	}
}
