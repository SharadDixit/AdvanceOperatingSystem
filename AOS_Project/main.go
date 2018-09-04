package main

import ("sync"
	"./main_struct"
	"sort"
	"fmt"
	"time"
	"math/rand"
	"encoding/json"

	//"os"
	//"strconv"
	"os"
	"strconv"
)
var finger_table_size int

var wait_group sync.WaitGroup

var coordinator_channel chan string

//var node_dictionary = make(map[int64] main_struct.Global_nodeinfo)

var node_list []int64

var ticker *time.Ticker

type node_list_order []int64

func (n node_list_order) Len() int           { return len(n) }
func (n node_list_order) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n node_list_order) Less(i, j int) bool { return n[i] < n[j] }

func generate_randomID(){

	//Manually entering nodes first
	node_list = []int64{54295823, 123415253, 234153252, 352241563, 431521242, 512451512, 632625325, 1058693285, 1715125251, 2526235153, 4241512517}
	sort.Sort(node_list_order(node_list))
	fmt.Println("Manually entered node list: ")
	fmt.Println("\n")
	fmt.Println(node_list)
}

func nodeAddDictionary(){

	for i:=0; i<len(node_list);i++ {
		var node = main_struct.Global_nodeinfo{ChannelID: node_list[i], Successor: -1, Predecessor: -1, Finger_table: nil}
		node_dictionary.Storing(node_list[i],node)
	}
}

func insertmessages_randomly(){
	ticker = time.NewTicker(8* time.Second)

	go func() {
		for t := range ticker.C {

			coordinator_channel <- generate_random_instructions(t)

		}
	}()
}

func generate_random_instructions(time time.Time)(string) {
	var instruction string
	rand.Seed(time.UTC().UnixNano())
	option := rand.Intn(6)

	switch option {
	case 0:
		{
			msg_struct := main_struct.File_info{

				Do:     "print-ring",

			}
			msg, _ := json.Marshal(&msg_struct)
			instruction = string(msg)
			fmt.Println( "\nPRINT - RING\n")
			//information_ring_nodes()

		}
	case 1:
		{
			node_ring := node_dictionary.LoadAllKeys() 
			sponsor_node := node_ring[rand.Intn(len(node_ring))]

			msg_struct := main_struct.File_info{
				Do:              "join-ring",
				Sponsoring_node: sponsor_node,
			}
			msg, _ := json.Marshal(&msg_struct)
			instruction = string(msg)
			fmt.Println("\n JOIN - RING\n")

		}
	case 2:
		{
			nodes_ring := nodesFromDictionary()
			leaving_node := nodes_ring[rand.Intn(len(nodes_ring))]


			fmt.Printf("Leaving Node: %d", leaving_node)
			choice := rand.Intn(2)
			switch choice {

			case 0:
				{
					mode := "orderly"
					msg_struct := &main_struct.File_info{
						Do:             "leave-ring",
						Recipient_node: 1715125251,
						Mode:           mode,
						Leaving_node:   leaving_node,
					}
					msg, _ := json.Marshal(msg_struct)
					instruction = string(msg)
					fmt.Println("\nLEAVE - RING\n")
				}
			case 1:
				{
					mode := "immediate"
					msg_struct := &main_struct.File_info{
						Do:           "leave-ring",
						Mode:         mode,
						Leaving_node: leaving_node,
					}
					msg, _ := json.Marshal(msg_struct)

					instruction = string(msg)
					fmt.Println("\nLEAVE - RING\n")
				}
			}
		}
	case 3:
		{
			nodes_ring := nodesFromDictionary()
			sponsoring_node := nodes_ring[rand.Intn(len(nodes_ring))]
			//sponsor := nodeList[rand.Intn(len(nodeList))]
			key := gen_node(rand_string())

			msg_struct := &main_struct.File_info{
				Do:        "put",
				Respond_to: sponsoring_node,
				Data_Key:      key,
				Data_Value:  "File Value",
			}
			msg, _ := json.Marshal(msg_struct)
			instruction = string(msg)
			fmt.Println("\n PUT - DATA \n")
		}
	case 4:
		{
			nodes_ring := nodesFromDictionary()
			sponsoring_node := nodes_ring[rand.Intn(len(nodes_ring))]

			get_key := gen_node(rand_string())

			msg_struct := &main_struct.File_info{
				Do:        "get",
				Respond_to: sponsoring_node,
				Data_Key:      get_key,
			}
			msg, _ := json.Marshal(msg_struct)
			instruction = string(msg)
			fmt.Println("\n GET - DATA \n")
		}
	case 5:
		{
			nodes_ring := nodesFromDictionary()
			sponsoring_node := nodes_ring[rand.Intn(len(nodes_ring))]

			data_key := gen_node(rand_string())

			msg_struct := &main_struct.File_info{
				Do:        "remove",
				Respond_to: sponsoring_node,
				Data_Key:      data_key,
			}
			msg, _ := json.Marshal(msg_struct)
			instruction = string(msg)
			fmt.Println("\n REMOVE - DATA\n")
		}

	}


	fmt.Println("\nInstruction Sending: ", instruction)
	fmt.Println("\n")
	return instruction
}

func main(){
	var pro= os.Args[1:]

		if len(pro) != 1 {
			fmt.Println("USAGE: go run main.go N")
			os.Exit(0)
		}

		finger_table_size, _ = strconv.Atoi(pro[0])

	//finger_table_size = 16  Advisable size as 16

	fmt.Println("Finger Table Size Entered as:",finger_table_size)

	coordinator_channel = make(chan  string)

	generate_randomID()

	nodeAddDictionary()

	wait_group.Add(1)

	go coordinator_function()

	insertmessages_randomly()

	wait_group.Wait()
}