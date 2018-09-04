package main

//import "time"
import (
	"./main_struct"

	"encoding/json"
	"fmt"
	"sort"

	"time"
)

var successor int64
var predecessor int64
var finger_table []int64
var bucket map[int64]string

func nodes_workers(nodeID int64, check bool) {
	defer wait_group.Done()
	node_channel := channel_dictionary[nodeID] ///This is the channel the node uses

	node_info_NotPresentInRing, _ := node_dictionary.Loading(nodeID)

	//finger_table:= make([]int64,32)
	//bucket := make(map[int64]string)
	//finger_table := node_info.Finger_table
	//predecessor := node_info.Predecessor
	//successor:= node_info.Successor

	if check {
		intialRingSetup(nodeID, node_info_NotPresentInRing)

	} else {
		//print_ring(nodeID)
	}
	targetID := nodeID
	node_info, checkNode := node_dictionary.Loading(nodeID)
	if checkNode == true {

		successor = node_info.Successor
		predecessor = node_info.Predecessor
		finger_table = node_info.Finger_table
		bucket = node_info.Bucket
	} else {

		finger_table = make([]int64, finger_table_size)
		successor = finger_table[0]
		predecessor = int64(0)
		bucket = make(map[int64]string)
		fmt.Println("New Node")
	}

	//fmt.Println(node_info)
	//information_ring_nodes()   ///placed at wrong position just for checking /////////*******************

	for instruction := range node_channel {

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(instruction), &data); err != nil {
			panic(err)
		}
		option := data["Do"]

		switch option {
			case "join-ring":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				sponsor_node := message_received.Sponsoring_node

				//updating bucket list
				joining_node_initial_info:=main_struct.Global_nodeinfo{nodeID,successor,predecessor,finger_table,bucket}
				node_dictionary.Storing(nodeID,joining_node_initial_info)

				join_ring(sponsor_node, targetID)
				fmt.Printf("\nNode has joined ring %d ", targetID)
				fmt.Println("\n")

			}
			case "leave-ring":
			{

				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				leaving_node:= message_received.Leaving_node
				leaving_node_node_info,_:=node_dictionary.Loading(leaving_node)
				leaving_node_successor:= leaving_node_node_info.Successor
				leaving_node_predecessor:= leaving_node_node_info.Predecessor
				leaving_node_bucket:= leaving_node_node_info.Bucket
				if message_received.Mode == "orderly" {
					leaving_ring_orderly(leaving_node_successor, leaving_node_predecessor, leaving_node_bucket, leaving_node)
				}
				leave_ring(leaving_node)
				fmt.Printf("Node %d left the ring - ", leaving_node)
				fmt.Println("\n")
			}
			case "find-ring-successor":
			{

				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)

				respond_to := message_received.Respond_to

				resondto_node_info,_ := node_dictionary.Loading(respond_to)
				finger_table = resondto_node_info.Finger_table

				target_node := message_received.Target_node

				find_successor(respond_to, target_node, finger_table)      //finger table of joining node is getting called checkkkkk //sponsor node fingertable to take
				fmt.Printf("\nFinding ring successor of %d \n", target_node)
				fmt.Println("\n")

			}
			case "update-successor":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				successor = message_received.Successor
				target_node:= message_received.Target_node
				fmt.Println("\nSuccessor:\n",successor)
				//node_dictionary.node_dictionary[nodeID].Successor = successor  /////Putting successor in the node_dictionary nodes ***********
				node_info_succ_update,_:= node_dictionary.Loading(target_node)
				node_info_succ_update.Successor = successor
				node_dictionary.Storing(target_node,node_info_succ_update)  /////Successor updated in node dictionary

				//Notify:  Calling notfiy message to tell the successor found that its predecessor is newly joined node
				channel_dictionary[target_node] <- initialize_ring_fingers_message()

				channel_dictionary[successor]<-notify_successor_message(target_node,successor)


				channel_dictionary[successor] <- getting_bucket_message(target_node)
				//channel_dictionary[successor] <- triggerGetBucktMessage(key)    for leave ring ***********************
			}
			case "init-ring-fingers":
			{
				//fmt.Println("Init-ring fingers successor:",successor)
				//fmt.Println("nodeID:",nodeID)
				nodeID_info,_:= node_dictionary.Loading(nodeID)     ///Loading the joint node from dictionary
				//fmt.Println("node_info successor",nodeID_info.Successor)
				initialize_ring_fingers(nodeID, nodeID_info.Successor)     ////The successor here is new updated succesor of the new joining node
				fmt.Printf("\nInit ring fingers for %d \n", nodeID)
				fmt.Println("\n")

			}
			case "get-ring-fingers":
			{
				get_ring_fingers(instruction, finger_table)
				//The above is successor's finger entry.
			}
			case "fingers-updated-from-successor":
			{
				node_info_finger_table,_:= node_dictionary.Loading(nodeID)
				/////Null Finger table send so as to copy successor's finger table to joining node
				init_finger_table_joiningnode := node_info_finger_table.Finger_table
				update_fingers_from_successor(nodeID, successor, instruction, init_finger_table_joiningnode)  //This finger table should be blank but it is not and it does not matter as we update it

				fmt.Printf("\nJoining node %d finger table:%d\n",nodeID, finger_table)
			}
			case "find-ring-predecessor":
			{

				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				//fmt.Println("find_predecessor finger table:",finger_table)
				respond_to := message_received.Respond_to
				target_node := message_received.Target_node

				node_info_respond_to_fingertable,_:= node_dictionary.Loading(respond_to)
				finger_table_predecessor:= node_info_respond_to_fingertable.Finger_table
				//fmt.Println("Chcecac:",finger_table_predecessor)
				find_predecessor(respond_to, target_node, finger_table_predecessor)
				fmt.Printf("\nFinding ring predecessor of %d \n ", target_node)
				fmt.Println("\n")

			}
			case "update-predecessor":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				predecessor = message_received.Predecessor
				target_node:= message_received.Target_node
				node_info_predecessor_update,_:= node_dictionary.Loading(target_node)
				node_info_predecessor_update.Predecessor = predecessor
				node_dictionary.Storing(target_node,node_info_predecessor_update)
				//node_dictionary.node_dictionary[nodeID].Predecessor = predecessor
				 channel_dictionary[predecessor]<-stabilize_message(target_node,predecessor)
				fmt.Printf("\nUpdating Predecessor of Node %d as Node %d\n", target_node, predecessor)
				fmt.Println("\n")
			}
			case "get-bucket-join-ring":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				n := message_received.Node
				getting_bucket(nodeID, n, bucket)
			}
			case "copy-bucket-join-ring":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				updated_bucket := message_received.Bucket
				target_node := message_received.Target_node
				for i, entry := range updated_bucket {
					bucket[i] = entry
				}
				fmt.Println("\nJOIN-RING,bucket updated\n ", bucket)
				bucket_update_node_info,_:=node_dictionary.Loading(target_node)
				bucket_update_node_info.Bucket = bucket
				node_dictionary.Storing(target_node,bucket_update_node_info)
				fmt.Println("\n")
			}
			case "update-predecessor-and-bucket":
			{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction),&message_received)

				successor_node_info:= message_received.Successor

				bucketreceived_node_info:= message_received.Bucket
				updated_bucket := update_bucket(instruction, bucketreceived_node_info)

				successor_node_info_bucket_update,_:= node_dictionary.Loading(successor_node_info)
				successor_node_info_bucket_update.Bucket = updated_bucket
				node_dictionary.Storing(successor_node_info,successor_node_info_bucket_update)

				predecessor_node_info:= message_received.Predecessor


				//node_info_succesor_predecessortoUpdate,_:= node_dictrionary.Loading(successor_node_info)

				channel_dictionary[successor_node_info] <- update_predecessor(predecessor_node_info,successor_node_info)
			}
			//case "stabilize":{
			//fmt.Printf("Stabilize will now run for Node %d",nodeID)
			//node_info_nodeID,_:= node_dictionary.Loading(nodeID)
			//
			//channel_dictionary[node_info_nodeID.Successor]<- update_bucket_predecessor_message(node_info_nodeID.Bucket,node_info_nodeID.Predecessor,node_info_nodeID.Successor)
			//channel_dictionary[node_info_nodeID.Predecessor] <- Updating_Successor_CallInitfingers_message(successor,predecessor)
			//
			//}
			case "notify":{
			message_received := main_struct.File_info{}
			json.Unmarshal([]byte(instruction),&message_received)

			fmt.Printf("\nNotify running for successor node %d\n",message_received.Successor)
			node_info_successor_notify,_:= node_dictionary.Loading(message_received.Successor)
			node_info_successor_notify.Predecessor = message_received.Target_node
			node_dictionary.Storing(message_received.Successor,node_info_successor_notify)

			fmt.Printf("\nNotify done and Successor's precessor setted up as %d\n",message_received.Target_node)
			}
			case "stabilize":{
			message_received := main_struct.File_info{}
			json.Unmarshal([]byte(instruction),&message_received)
			/////////NOTIFY PREDECESSOR CAN ALSO BE CALLED *****************
			fmt.Printf("\n Stabilize running for predecessor node %d\n",message_received.Predecessor)
			node_info_predecessor_notify,_:= node_dictionary.Loading(message_received.Predecessor)
			node_info_predecessor_notify.Successor = message_received.Target_node
			node_dictionary.Storing(message_received.Predecessor,node_info_predecessor_notify)

			fmt.Printf("\n Stabilize done and predecessor's successor setted up as %d\n",message_received.Target_node)

			fmt.Printf("\n Now Fix-Finger Will run\n")
			time.Sleep(1*time.Second)
			channel_dictionary[message_received.Predecessor]<- fix_fingers_message(message_received.Target_node,message_received.Predecessor)
			}
			case "fix-fingers":{
			message_received:= main_struct.File_info{}
			json.Unmarshal([]byte(instruction),&message_received)
			fix_finger_predecessor_node:= message_received.Predecessor
			fix_finger_target_node:= message_received.Target_node

			fmt.Printf("\nFix finger running for Node %d\n",fix_finger_predecessor_node)
				finger_entry:=0
			 fix_fingers(fix_finger_predecessor_node,fix_finger_target_node,finger_entry)

			}
			case "put":{
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				respond_to := message_received.Respond_to
				data_key := message_received.Data_Key
				data_value:= message_received.Data_Value
				channel_dictionary[respond_to] <- find_bucket_succesor_message(respond_to, data_key, data_value)
			}
			case "put-get-bucket-successor":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					node_info_get_bucket_successor,_:= node_dictionary.Loading(message_received.Sponsoring_node)

					put_find_bucket_successor(message_received.Sponsoring_node, message_received.Data_Key, message_received.Data_Value, node_info_get_bucket_successor.Finger_table)
				}
			case "put-bucket-successor":
			{
				bucket_put:= make(map[int64]string)
				message_received := main_struct.File_info{}
				json.Unmarshal([]byte(instruction), &message_received)
				target_node:= message_received.Target_node
				node_info_put_bucket_successor,_:=node_dictionary.Loading(target_node)
				node_info_put_bucket_successor_bucket:= node_info_put_bucket_successor.Bucket

				write_lock.Lock()
				for k,v := range node_info_put_bucket_successor_bucket {
					bucket_put[k] = v
				}
				write_lock.Unlock()
				write_lock.Lock()
				bucket_put[message_received.Data_Key] = message_received.Data_Value
				write_lock.Unlock()
				node_info_put_bucket_successor.Bucket = bucket_put
				node_dictionary.Storing(target_node,node_info_put_bucket_successor)

				fmt.Println("Put Bucket", bucket_put)
				fmt.Printf("in node %d\n",target_node)
				fmt.Println("\n")
			}
			case "get":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					respond_to := message_received.Respond_to
					data_key := message_received.Data_Key
					channel_dictionary[respond_to] <- get_bucket_data_successor_message(data_key, respond_to)
				}
			case "get-data-successor":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					node_info_get_bucket_successor,_:= node_dictionary.Loading(message_received.Sponsoring_node)

					get_find_bucket_data_successor(message_received.Sponsoring_node,message_received.Data_Key, node_info_get_bucket_successor.Finger_table)
				}
			case "get-bucket-data-successor":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					target_node:= message_received.Target_node
					//fmt.Println("Target Node:",target_node)
					data_key:= message_received.Data_Key
					node_info_get_bucket_data,_:= node_dictionary.Loading(target_node)
					node_info_get_bucket_data_bucket:= node_info_get_bucket_data.Bucket
					if len(node_info_get_bucket_data_bucket)==0 {
						fmt.Printf("Local Bucket of Node %d Empty and hence Data key trying to get not present in the chord ring",target_node)
					}else{
						write_lock.Lock()
						for k ,v := range node_info_get_bucket_data_bucket{
							if k==data_key{
								fmt.Println("\nKey present with the nodes %d local bucket\n",target_node)
								fmt.Println("Data for key - ", k, " = ", v)
							}
						}
						write_lock.Unlock()
					}
				}
			case "remove":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					respond_to := message_received.Respond_to
					data_key := message_received.Data_Key
					channel_dictionary[respond_to] <- remove_bucket_data_successor_message(data_key, respond_to)
				}
			case "remove-data-successor":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					node_info_get_bucket_successor,_:= node_dictionary.Loading(message_received.Sponsoring_node)
					remove_find_bucket_data_successor(message_received.Sponsoring_node, message_received.Data_Key, node_info_get_bucket_successor.Finger_table)
				}
			case "remove-data":
				{
					message_received := main_struct.File_info{}
					json.Unmarshal([]byte(instruction), &message_received)
					target_node:= message_received.Target_node
					data_key:= message_received.Data_Key
					node_info_remove_bucket_data,_:= node_dictionary.Loading(target_node)
					node_info_remove_bucket_data_bucket:= node_info_remove_bucket_data.Bucket

					if len(node_info_remove_bucket_data_bucket)==0 {
						fmt.Printf("\nLocal Bucket of Node %d Empty and hence Data key trying to remove not present in the chord ring\n",target_node)
					}else{
						write_lock.Lock()
						delete(node_info_remove_bucket_data_bucket,data_key)
						write_lock.Unlock()
						fmt.Printf("\nKey present with the node %d local bucket has been removed\n",target_node)
					}

				}
			case "print-ring":{
				information_ring_nodes()
			}
		}
	}
}
func join_ring(sponsor_node int64, target_node int64){

	///First add this to dictionary then proceed with the below


	channel_dictionary[sponsor_node] <- successor_message(sponsor_node,target_node)
	channel_dictionary[sponsor_node] <-predecessor_message(sponsor_node,target_node)

	//Calling Stabilize for both successor and predecessor of newly joint node
	//fmt.Printf("Newly joint node: %d",target_node)
	//node_info_new_node,_:= node_dictionary.Loading(target_node)
	//fmt.Printf("Successor Of new node: %d",node_info_new_node.Successor)
	//fmt.Printf("Predecessor of new node: %d",node_info_new_node.Predecessor)
	////channel_dictionary[sponsor_node]<-

	//JoinsRing(target_node)
}

func find_successor(sponsor_node int64, target_node int64, finger_table []int64)(){
	//fmt.Println(finger_table)
	//fmt.Println(sponsor_node)
	if target_node > sponsor_node && target_node < finger_table[0] {
		channel_dictionary[target_node] <- Updating_Successor_CallInitfingers_message(finger_table[0],target_node)

	} else {

		closestNode := find_closest_preceeding_node(target_node, finger_table)

		channel_dictionary[closestNode] <- successor_message(closestNode, target_node)
	}

}

func find_closest_preceeding_node(target_node int64, finger_table []int64) int64 {
	//fmt.Println(len(finger_table))
	entries := make([]int64, finger_table_size)

	copy(entries, finger_table)
	sort.Sort(sort.Reverse(node_list_order(entries)))

	for _, node := range entries {
		if node < target_node {
			return node
		}
	}
	return entries[0]
}

func initialize_ring_fingers(nodeID int64, successor int64){
	channel_dictionary[successor] <- get_ring_fingers_message(nodeID)
}

func get_ring_fingers(instruction string, fingerTable []int64) {
	check := main_struct.File_info{}
	json.Unmarshal([]byte(instruction), &check)

	receiving_node := check.Respond_to
	channel_dictionary[receiving_node] <- update_finger_entries_message(fingerTable)   //This is successor finger table
}

func update_fingers_from_successor(nodeID int64, successor int64, msg string, finger_table []int64) {
	check := main_struct.File_info{}
	json.Unmarshal([]byte(msg), &check)
	successor_finger_table := check.FingerTable

	/////Putting the initialized finger table in the node's Dictionary
	copied_finger_table := retrieve_copied_entries_finger_table(nodeID, successor, finger_table, successor_finger_table)
	node_info_fingertable_update,_:= node_dictionary.Loading(nodeID)
	node_info_fingertable_update.Finger_table = copied_finger_table

	//node_dictionary.node_dictionary[nodeID].Finger_table = copied_finger_table  //////Putting finger table in the node dictionary nodes *****************
}

func find_predecessor(sponsor_node int64, target_node int64, fingerTable []int64){

	if target_node > sponsor_node && target_node < fingerTable[0] {
		channel_dictionary[target_node] <- update_predecessor(sponsor_node,target_node)

	} else {
		closestNode := find_closest_preceeding_node(target_node, fingerTable)
		if closestNode == target_node {
			channel_dictionary[target_node] <- update_predecessor(closestNode,target_node)
		} else {
			channel_dictionary[closestNode] <- predecessor_message(closestNode, target_node)
		}
	}
}

//func stabilize_ticker(node int64) {
//	ticker = time.NewTicker(5000 * time.Millisecond)
//	go func() {
//		for range ticker.C {
//
//			instruction_struct := &main_struct.File_info{
//				Do: "stabilize",
//			}
//			instruction, _ := json.Marshal(instruction_struct)
//			channel_dictionary[node] <- string(instruction)
//
//			ticker.Stop()
//		}
//	}()
//}

func leave_ring(leaving_node int64){
	wait_group.Done()		///////////////CHECK **********************
	close(channel_dictionary[leaving_node])
	delete(channel_dictionary,leaving_node)
	node_dictionary.Deleting(leaving_node)
}

func leaving_ring_orderly(successor int64, predecessor int64, bucket map[int64]string, leaving_node int64){

	channel_dictionary[successor] <-  update_bucket_predecessor_message(bucket, predecessor,successor)
	channel_dictionary[predecessor] <- Updating_Successor_CallInitfingers_message(successor,predecessor)  ///since presecessor successor is successor now
}
func getting_bucket(nodeID int64, n int64, bucket map[int64]string) {
	channel_dictionary[n] <- copy_bucket_message(nodeID, n, bucket)
	for check := range bucket {
		if check >= n && check < nodeID {
			delete(bucket, check)
			fmt.Println("join-ring:: Successor's Bucket List updated :", bucket)
			fmt.Println("\n")
		}
	}
}
func update_bucket(message string, bucket map[int64]string) map[int64]string {

	check := main_struct.File_info{}
	json.Unmarshal([]byte(message), &check)
	for k, v := range check.Bucket {
		bucket[k] = v
	}
	return bucket
}
func fix_fingers(node_fix_finger int64, sponsor_node int64,finger_entry int){

	//if finger_entry<32{
	//	finger_fixing_entries(node_fix_finger , sponsor_node ,finger_entry)
	//
	//}else{
	//	fmt.Println("fingers fixed")
	//}


	//fmt.Println("\ncheck1")
	//fmt.Println(sponsor_node)
	fixed_finger_table:= make([]int64,finger_table_size)
	for i:=0;i<finger_table_size;i++ {
		fmt.Println("check",i)
		finger_entry_raw := calculation_finger_entry(i, node_fix_finger)
		fmt.Println("finger_entry_raw",finger_entry_raw)
		fixed_finger_table[i] = find_successor_fix_finger(finger_entry_raw,sponsor_node)
		//fmt.Println("\nfinger_entry_acutal_value:\n",fixed_finger_table[i])
	}
	fmt.Println("fixed fingers",fixed_finger_table)
	node_info_fix_finger_update,_:= node_dictionary.Loading(node_fix_finger)
	node_info_fix_finger_update.Finger_table = fixed_finger_table
	node_dictionary.Storing(node_fix_finger,node_info_fix_finger_update)

}
func finger_fixing_entries(node_fix_finger int64, sponsor_node int64, finger_entry int){

	//node_info_fix_finger,_:= node_dictionary.Loading(node_fix_finger)
	//node_info_fix_finger_fingertable:= node_info_fix_finger.Finger_table

	finger_table_raw := calculation_finger_entry(finger_entry,node_fix_finger)
	//fmt.Println("\nfinger Entry raw", finger_table_raw)
	fmt.Println("\nSponsor node",sponsor_node)
	finger_entry_acutal_value:= find_successor_fix_finger(finger_table_raw,sponsor_node)
	//node_info_fix_finger_fingertable[finger_entry] = finger_entry_acutal_value
	//
	//node_info_fix_finger_fingertable = node_info_fix_finger_fingertable

	//node_dictionary.Storing(node_fix_finger,node_info_fix_finger)
	fmt.Println("\nfinger_entry_acutal_value:\n",finger_entry_acutal_value)
	finger_entry++
	fix_fingers(node_fix_finger , sponsor_node , finger_entry )
}
func find_successor_fix_finger(finger_entry_raw int64 ,sponsor_node int64) (successor int64){
	fmt.Println("sponsor node:",sponsor_node)
	node_info_sponsor_node,_:= node_dictionary.Loading(sponsor_node)
	finger_table_sponsor:= node_info_sponsor_node.Finger_table

	finger_table_sponsor_successor:= node_info_sponsor_node.Successor

	//fmt.Println("Sponsor node fix finger successor:",finger_table_sponsor_successor)
	//fmt.Println("finger table sponsor node",finger_table_sponsor)
	//fmt.Println("finger entry raw",finger_entry_raw)
	//fmt.Println()
	//
	if finger_entry_raw > sponsor_node && finger_entry_raw < finger_table_sponsor_successor {
		successor = finger_table_sponsor_successor
		//channel_dictionary[finger_entry_raw] <- updateSuccessorMessage(fingerTable[0])

	} else {

		closest_node := find_closest_preceeding_node(finger_entry_raw, finger_table_sponsor)

		successor = find_successor_fix_finger(finger_entry_raw,closest_node)
	}
	return successor
}

//func FindSuccessor(sponsoringNode *Node, targetID int64) (successor int64) {
//	fmt.Println(targetID, sponsoringNode.ID, sponsoringNode.Successor)
//	if sponsoringNode.ID == sponsoringNode.Successor {
//		successor = sponsoringNode.ID
//		AnsChannels[sponsoringNode.ID] <- successor
//	} else if targetID > sponsoringNode.ID && targetID <= sponsoringNode.Successor {
//		successor = sponsoringNode.Successor
//		AnsChannels[sponsoringNode.ID] <- successor
//	} else {
//		precedingNode := findClosestPrecedingNode(sponsoringNode, targetID)
//		successor = FindSuccessor(NodesInRing[precedingNode], targetID)
//	}
//
//	return successor
//}

func put_find_bucket_successor(sponsor_node int64, data_key int64, data_value string, fingerTable_sponsornode []int64) {
	if data_key > sponsor_node && data_key < fingerTable_sponsornode[0] {
		channel_dictionary[fingerTable_sponsornode[0]] <- put_bucket_message(data_key, data_value,fingerTable_sponsornode[0])

	} else {

		closest_node := find_closest_preceeding_node(data_key, fingerTable_sponsornode)
		channel_dictionary[closest_node] <- find_bucket_succesor_message(closest_node, data_key, data_value)
	}
}
func get_find_bucket_data_successor(sponsor_node int64, data_key int64, fingerTable_sponsornode []int64) {
	if data_key > sponsor_node && data_key < fingerTable_sponsornode[0] {
		channel_dictionary[fingerTable_sponsornode[0]] <- get_bucket_data_message(data_key,fingerTable_sponsornode[0])

	} else {

		closest_node := find_closest_preceeding_node(data_key, fingerTable_sponsornode)
		channel_dictionary[closest_node] <- get_bucket_data_successor_message(data_key, closest_node)
	}
}
func remove_find_bucket_data_successor(sponsor_node int64, data_key int64,fingerTable_sponsornode []int64){
	if data_key > sponsor_node && data_key < fingerTable_sponsornode[0] {
		channel_dictionary[fingerTable_sponsornode[0]] <- remove_bucket_data_message(data_key,fingerTable_sponsornode[0])

	} else {

		closest_node := find_closest_preceeding_node(data_key, fingerTable_sponsornode)
		channel_dictionary[closest_node] <- remove_bucket_data_successor_message(data_key, closest_node)
	}
}
func print_ring(node int64) {
	ticker = time.NewTicker(12000 * time.Millisecond)
	go func() {
		for range ticker.C {

			channel_dictionary[node] <- print_message()
			ticker.Stop()
		}
	}()
}