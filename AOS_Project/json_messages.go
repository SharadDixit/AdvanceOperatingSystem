package main

import "encoding/json"
import (
	"./main_struct"

	//"math/rand"
	"math/rand"
)


func successor_message(sponsor_node int64, target_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:        "find-ring-successor",
		Respond_to: sponsor_node,
		Target_node:  target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)

	return string(instruction)
}

func Updating_Successor_CallInitfingers_message(successor int64,node int64) string {

	instruction_struct := &main_struct.File_info{
		Do:        "update-successor",
		Successor: successor,
		Target_node:node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)

}
func initialize_ring_fingers_message() string {
	instruction_struct := &main_struct.File_info{
		Do: "init-ring-fingers",
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}

func get_ring_fingers_message(target_node int64) string {

	instruction_struct := &main_struct.File_info{
		Do:        "get-ring-fingers",
		Respond_to: target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)

}

func update_finger_entries_message(finger_table []int64) string {
	instruction_struct := &main_struct.File_info{
		Do:          "fingers-updated-from-successor",
		FingerTable: finger_table,
	}
	notifyNodeAboutSuccessorMsg, _ := json.Marshal(instruction_struct)
	return string(notifyNodeAboutSuccessorMsg)
}

func predecessor_message(sponsor_node int64, target_node int64)string{
	instruction_struct := &main_struct.File_info{
		Do:        "find-ring-predecessor",
		Respond_to: sponsor_node,
		Target_node:  target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}

func update_predecessor(predecessor int64,node int64) string {

	instruction_struct := &main_struct.File_info{
		Do:          "update-predecessor",
		Predecessor: predecessor,
		Target_node: node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)

}
func getting_bucket_message(nodeID int64) string {
	instruction_struct := &main_struct.File_info{
		Do:  "get-bucket-join-ring",
		Node: nodeID,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func copy_bucket_message(successor int64, target_node int64, bucket map[int64]string) string {
	copy_bucket := make(map[int64]string)
	for check, v := range bucket {
		if check >= target_node && check < successor {
			copy_bucket[check] = v
		}
	}
	instruction_struct := &main_struct.File_info{
		Do:     "copy-bucket-join-ring",
		Bucket: copy_bucket,
		Target_node: target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func update_bucket_predecessor_message(bucket map[int64]string, predecessor int64,successor int64) string {

	instruction_struct := &main_struct.File_info{
		Do:          "update-predecessor-and-bucket",
		Predecessor: predecessor,
		Successor:successor,
		Bucket:  bucket,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)

}
func notify_successor_message(target_node int64, successor int64)string{
	instruction_struct := &main_struct.File_info{
		Do: "notify",
		Target_node:target_node,
		Successor:successor,
	}
	instruction,_:= json.Marshal(instruction_struct)
	return string(instruction)

}
func stabilize_message (target_node int64, predecessor int64)string{
	instruction_struct := &main_struct.File_info{
		Do: "stabilize",
		Target_node:target_node,
		Predecessor:predecessor,
	}
	instruction,_:= json.Marshal(instruction_struct)
	return string(instruction)

}
func fix_fingers_message(target_node int64, predecessor int64)string{

	nodes_ring := nodesFromDictionary()
	random_channel_node := nodes_ring[rand.Intn(len(nodes_ring))]
	instruction_struct := &main_struct.File_info{
		Do: "fix-fingers",
		Target_node:random_channel_node,
		Predecessor:predecessor,
	}
	instruction,_:= json.Marshal(instruction_struct)
	return string(instruction)
}
func find_successor_fix_finger_message(finger_entry_raw int64,sponsor_node int64,finger_table_sponsor []int64)string{

	instruction_struct := &main_struct.File_info{
		Do: "find-successor-fix-finger",
		Target_node:finger_entry_raw,
		Sponsoring_node:sponsor_node,
		FingerTable:finger_table_sponsor,
	}
	instruction,_:= json.Marshal(instruction_struct)
	return string(instruction)
}
func find_bucket_succesor_message(sponsor_node int64, data_key int64, data_value string) string {
	instruction_struct := &main_struct.File_info{
		Do:      "put-get-bucket-successor",
		Data_Key:  data_key   ,
		Data_Value:   data_value,
		Sponsoring_node: sponsor_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func put_bucket_message(data_key int64, data_value string,target_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:    "put-bucket-successor",
		Data_Key:   data_key,
		Data_Value: data_value,
		Target_node: target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func get_bucket_data_successor_message(data_key int64, sponsor_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:          "get-data-successor",
		Data_Key:         data_key,
		Sponsoring_node:     sponsor_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func get_bucket_data_message(data_key int64,target_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:  "get-bucket-data-successor",
		Data_Key: data_key,
		Target_node:target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func remove_bucket_data_successor_message(data_key int64, sponsoring_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:      "remove-data-successor",
		Data_Key:     data_key,
		Sponsoring_node: sponsoring_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func remove_bucket_data_message(data_key int64,target_node int64) string {
	instruction_struct := &main_struct.File_info{
		Do:  "remove-data",
		Data_Key: data_key,
		Target_node:target_node,
	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}
func print_message()string{
	instruction_struct := &main_struct.File_info{
		Do:  "print-ring",

	}
	instruction, _ := json.Marshal(instruction_struct)
	return string(instruction)
}