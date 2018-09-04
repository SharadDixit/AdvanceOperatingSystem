package main_struct


type Global_nodeinfo struct{
	ChannelID int64
	Successor int64
	Predecessor int64
	Finger_table []int64
	Bucket		map[int64]string
	//Keys map[int64]int64
}
type File_info struct {
	Do              string
	Sponsoring_node int64
	Mode            string
	Respond_to      int64
	Data_Key        int64
	Data_Value		string
	Target_node		int64
	Successor		int64
	FingerTable		[]int64
	Predecessor		int64
	Recipient_node  int64
	Leaving_node  	int64
	Node 			int64
	Bucket		map[int64]string
}
type Data struct {
	Key   int64
	Value string
}
type Communication struct{
	Id int64
}

