package main

import (
	"sort"
	"time"
	"math/rand"
	"crypto/md5"
	"math"
)

func nodesFromDictionary()([]int64){
	nodes_ring:=make([]int64,0,len(node_dictionary.node_dictionary))
	for k:=range node_dictionary.node_dictionary{
		nodes_ring=append(nodes_ring,k)
	}
	sort.Sort(node_list_order(nodes_ring))

	return nodes_ring
}

func calculation_finger_entry(entry int, nodeID int64)int64{
	return int64((int(nodeID) + int(math.Pow(2, float64(entry)))) % int(math.Pow(2, 32)))
}

func retrieve_copied_entries_finger_table(nodeID int64, successor int64, finger_table []int64, successor_finger_table []int64)[]int64 {

	for i, v := range successor_finger_table {

		raw_entry := int64((int(nodeID) + int(math.Pow(2, float64(i)))) % int(math.Pow(2, 32)))

		if raw_entry <= successor {
			finger_table[i] = successor
		} else {
			finger_table[i] = v
		}

	}
	return finger_table
}

//For random node genereation
var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func rand_string()string{
	checks := make([]rune, 15)
	for i := range checks {
		rand.Seed(time.Now().UTC().UnixNano())
		checks[i] = letters[rand.Intn(52)]
	}
	return string(checks)
}

func gen_node(node string) int64 {
	checknode := hashing(node)
	return hashVal(checknode[0:4])
}

func hashing(checks string) [md5.Size]byte {
	return md5.Sum([]byte(checks))

}

func hashVal(check []byte) int64 {
	return ((int64(check[3]) << 24) | (int64(check[2]) << 16) | (int64(check[1]) << 8) | (int64(check[0]))) }

