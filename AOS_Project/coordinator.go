package main

import (
	"encoding/json"

	"math/rand"
)

var channel_dictionary map[int64]chan string

func coordinator_function() {
	defer wait_group.Done()
	channel_dictionary = make(map[int64]chan string)

	nodes_ring:= node_dictionary.LoadAllKeys()

	for i := 0; i < len(nodes_ring); i++ {
		nodeID := nodes_ring[i]
		channel_dictionary[nodeID] = make(chan string, 10)
	}
	for i := 0; i < len(nodes_ring); i++ {
		nodeID := nodes_ring[i]
		wait_group.Add(1)
		go nodes_workers(nodeID, true)
	}

	for instruction := range coordinator_channel {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(instruction), &data); err != nil {
			panic(err)
		}
		if data["Do"] == "leave-ring" {
			channel_dictionary[1715125251] <- instruction
		}
		if data["Do"] == "join-ring" {

			node := gen_node(rand_string())
			channel_dictionary[node] = make(chan string, 10)
			wait_group.Add(1)
			go nodes_workers(node, false)
			channel_dictionary[node] <- instruction

		}
		if data["Do"] == "print-ring"{
			information_ring_nodes()
		}
		if data["Do"] == "put"{

			nodes_ring := nodesFromDictionary()
			random_channel_node := nodes_ring[rand.Intn(len(nodes_ring))]
			channel_dictionary[random_channel_node]<- instruction
		}
		if data["Do"] == "get" {

			nodes_ring := nodesFromDictionary()
			random_channel_node := nodes_ring[rand.Intn(len(nodes_ring))]
			channel_dictionary[random_channel_node] <- instruction
		}
		if data["Do"] == "remove"{

			nodes_ring := nodesFromDictionary()
			random_channel_node := nodes_ring[rand.Intn(len(nodes_ring))]
			channel_dictionary[random_channel_node] <- instruction
		}
	}
}
