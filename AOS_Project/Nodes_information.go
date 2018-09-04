package main

import (
	"fmt"

)

func information_ring_nodes(){
	fmt.Println("LIST OF NODES")
	for k,node := range node_dictionary.node_dictionary{
		fmt.Println("START OF NODE INFORMATION")
		fmt.Printf("\n Node %d present in ring\n",k)
		fmt.Printf("Contents of Nodes %d",node.ChannelID)
		if node.Successor != -1 {
			fmt.Printf("Successor Id: %d\n", node.Successor)
		}else{
			fmt.Printf("Successor Id: nil\n")
		}

		if node.Predecessor != -1 {
			fmt.Printf("Predecessor Id: %d\n", node.Predecessor)
		}else{
			fmt.Printf("Predecessor Id: nil\n")
		}
		if node.Finger_table != nil {
			for node_id, node_entry := range node.Finger_table {
				if node_entry != -1 {
					fmt.Printf("Finger Table at %d is %d\n", node_id, node_entry)
				}
			}
		}else{
			fmt.Println("Finger Table Empty")
		}

		if node.Bucket!= nil{
			for k,v:= range node.Bucket{
				fmt.Println("Buckets: Data Key:",k,"&","Data Value:",v)
			}
		}else {
			fmt.Println("Bucket Empty")
		}
		fmt.Println("END OF NODE INFORMATION")

	}
}
