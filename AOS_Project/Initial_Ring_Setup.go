package main

import "./main_struct"

func intialRingSetup(nodeID int64,node_info main_struct.Global_nodeinfo){

	finger_table:= make([]int64,finger_table_size)
	finger_table = finger_successorboth_entry(nodeID,finger_table)
	predecessor := find_predecessor_initial(nodeID)
	node_info.Predecessor = predecessor

	node_info.Finger_table = finger_table
	successor := finger_table[0]
	node_info.Successor = successor

	node_dictionary.Storing(nodeID,node_info)

}

func finger_successorboth_entry(nodeID int64, finger_table []int64)([]int64){

	for i:=0;i<finger_table_size;i++{
		raw_finger_entry:= calculation_finger_entry(i,nodeID)
		finger_table[i] = find_successor_rawentry(raw_finger_entry)
	}
	return finger_table
}

func find_successor_rawentry(raw_entry int64) int64 {

	nodes_ring:=node_dictionary.LoadAllKeys()

	for _, node := range nodes_ring {
		if node >= raw_entry {
			return node
		}
	}
	return nodes_ring[0]
}
func find_predecessor_initial(nodeID int64)int64{

	nodes_ring:=node_dictionary.LoadAllKeys()

	for i, node := range nodes_ring{
		if node >= nodeID {
			if i == 0 {
				return nodes_ring[len(nodes_ring)-1]
			}
			return nodes_ring[i-1]
		}
	}
	return nodes_ring[0]
}
