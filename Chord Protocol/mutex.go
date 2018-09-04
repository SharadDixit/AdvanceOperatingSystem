package main

import "sync"
import (
	"./main_struct"
	"sort"
)

var write_lock sync.Mutex
type NodeDictionary struct{
	sync.RWMutex
	node_dictionary map[int64](main_struct.Global_nodeinfo)
}
func node_dictionary_Map() *NodeDictionary {
	return &NodeDictionary{
		node_dictionary : make(map[int64](main_struct.Global_nodeinfo)),
	}
}
//Map consist of all the nodes presently in chord ring,map does not replace routing through choord,
// therefore no functions in this program depend on this map for routing. Routing is done through a combination of predecessor, successor pointer traversal and finger table lookups.
var node_dictionary = node_dictionary_Map()

//Loading Key From Map    //Taking nodeID:key value and obtaining the node from the map
func (NodeDicStruct *NodeDictionary) Loading(map_key int64)(map_value main_struct.Global_nodeinfo,yes bool){
	NodeDicStruct.RLock()
	node_obtain,yes:=NodeDicStruct.node_dictionary[map_key]
	NodeDicStruct.RUnlock()
	return node_obtain,yes
}
//Loading all keys from the map in an slice

func (NodeDicStruct *NodeDictionary) LoadAllKeys()([]int64){
	NodeDicStruct.RLock()
	nodes_ring:=make([]int64,0,len(node_dictionary.node_dictionary))
	for k:=range node_dictionary.node_dictionary{
		nodes_ring=append(nodes_ring,k)
	}
	sort.Sort(node_list_order(nodes_ring))
	NodeDicStruct.RUnlock()
	return nodes_ring
}

//Storing the node in the map, that is joining the ring  //Taking both nodeID and node and storing them in the map
func(NodeDicStruct *NodeDictionary) Storing(map_key int64, map_value main_struct.Global_nodeinfo){
	NodeDicStruct.Lock()
	NodeDicStruct.node_dictionary[map_key] = map_value
	NodeDicStruct.Unlock()
}

//Deleting the node in the map, leaving the chord ring  //Taking map_key:nodeID value and remove from map
func(NodeDicStruct *NodeDictionary) Deleting(map_key int64){
	NodeDicStruct.Lock()
	delete(NodeDicStruct.node_dictionary,map_key)
	NodeDicStruct.Unlock()
}
