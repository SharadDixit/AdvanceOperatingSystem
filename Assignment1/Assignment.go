package main
import (
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"io"
	"bufio"
	"strconv"
	"unicode/utf8"
	"encoding/json"
)
type workers struct{
	Filename string
	Start int
	End int
}
type endwork struct{
	Value int64
	Prefix string
	Suffix string
	Notprefixsuffix string
}
var join[] string
var joins[] string
func ReadInts(r io.Reader) ([]int64) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int64
	for scanner.Scan() {

		newInt, _ := strconv.ParseInt(scanner.Text(), 0, 64)

		result = append(result, newInt)
	}
	return result
}
func main() {
	channel:=make(chan string)
	channel1:=make(chan string)
	go coordinator(channel,channel1)
	time.Sleep(time.Second*10)
}
func coordinator(channel chan string,channel1 chan string) {
	fmt.Println("File for input:ASCII FILE.txt")
	time.Sleep(time.Second * 1)
	bs, err := ioutil.ReadFile("ASCII FILE1.txt")
	if err != nil {
		fmt.Println(err)
	}
	s := string(bs)
	ObtainedASCIIString := ReadInts(strings.NewReader(s))
	fmt.Println("String slice After Conversion:", ObtainedASCIIString)
	elementsinobtainedstring := len(ObtainedASCIIString)
	var id [] string
	for i := 0; i < elementsinobtainedstring; i++ {
		a := strconv.FormatInt(ObtainedASCIIString[i], 10)
		id = append(id, a)
	}
	fmt.Println("String Slice after conversion and then after changing it into string format:",id)
	sa := strings.Join(id, " ")
	fmt.Println("One complete String by joining elements of string:",sa)
	fmt.Println("No.of Bytes in the single stream",utf8.RuneCountInString(sa))
	Noofbits := utf8.RuneCountInString(sa)
	var workersnumberaa int
	fmt.Println("Enter Number Of Workers:")
	fmt.Scan(&workersnumberaa)
	fmt.Println("Workers Entered and Also No. of Goroutines after division:", workersnumberaa)
	elementspassing := Noofbits / workersnumberaa
	fmt.Println("No.of bits to each worker:", elementspassing)
	aa := 0
	var finalsum int64
	/*for i:=0;i<workersnumberaa ;i++  {
		go worker(channel,r)                      //Can also do this way **DON'T DELETE IMPORTANT
		time.Sleep(time.Millisecond*100)
	}*/
	w:=make([]workers,workersnumberaa)
	for i := 0; i < workersnumberaa; i++ {
		go worker(channel,sa,channel1)
		time.Sleep(time.Millisecond*100)
		w[i] = workers{"ASCII FILE1", aa, aa + elementspassing-1}   //start and end pos are at index position
		aa = aa + elementspassing
		x:=[]workers{w[i]}
		firstjson, err := json.Marshal(x)
		if err != nil {
			fmt.Println(err)
		}
		channel<-string(firstjson)
		time.Sleep(time.Millisecond*50)
		secondjsonreceive:=<-channel1
		fmt.Println("SecondJSON msg after marshal and through channel1",string(secondjsonreceive))
		unmarsharsecondjason:=[]byte(secondjsonreceive)
		xp1:=[]endwork{}
		err1:=json.Unmarshal([]byte(unmarsharsecondjason),&xp1)
		if err1!=nil {
			fmt.Println(err1)
		}
		fmt.Printf("Second json:%+v \n",xp1)
		var valuea int64
		var prefixa string
		var suffixa string
		var notprefixsuffixa string
		for _,v:=range xp1{
			valuea=(v.Value)
			prefixa=(v.Prefix)
			suffixa=(v.Suffix)
			notprefixsuffixa=(v.Notprefixsuffix)
		}

		finalsum=finalsum+valuea
		fmt.Println("Prefix returned to coordinator:",prefixa)
		fmt.Println("Suffix returned to coordinator",suffixa)
		fmt.Println("Not Prefix Nor Suffix:",notprefixsuffixa)
		fmt.Println("\n")

			join=append(join,suffixa)
			joins=append(joins,prefixa)


		}
	fmt.Println(join)
	fmt.Println(joins)
	fmt.Println("Total Sum:",finalsum)

	for i:=0;i<len(join)-1 ;i++  {
		var k[]string
		k=append(k,join[i])
		c:=append(k,joins[i+1])
		fmt.Println(c)
		d:=strings.Join(c,"")
		fmt.Println(d)
	}
   /* k=append(k,join[0])
    c:=append(k,joins[1])
    fmt.Println(c)
    d:=strings.Join(c,"")
    fmt.Println(d)
*/
}
func worker(channel chan string,r string,channel1 chan string){
	msg:=<-channel
	fmt.Println("JSON msg after marshal and send through channel:",msg)
	a:=[]byte(msg)
	xp:=[]workers{}
	err:=json.Unmarshal([]byte(a),&xp)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Printf("go data after unmarshal: %+v \n",xp)
	starta:=0
	enda:=0
	for _,v:=range xp{
		starta=(v.Start)
		enda=(v.End)+1
	}
	var newarray string
	newarray = r[starta:enda]
	fmt.Println("String after byte division for workers:",newarray)

	lenghtofnewarray:=len(newarray)    //newarray is a string and new array in format is the array after slicing the string
	newarrayinarrayformat := strings.Fields(newarray)    //separates 1_2 like [0]:1   [1]:2   but if string is 33_ [0]:33  [1]:33
	splittednewarrayinformat:=strings.Split(newarray," ")
	lengthofnewarrayinformat := len(newarrayinarrayformat)
	var prefix string
	var suffix string
	var notprefixsuffix string
	if len(newarrayinarrayformat[0])!=lenghtofnewarray {   //For removing 444 case
	if newarray[0]==32{
		prefix=" "
		fmt.Println(suffix)
	}else {
		prefix = splittednewarrayinformat[0]
	}
		fmt.Println("Prefix:", prefix)
		if newarray[lenghtofnewarray-1]==32 {
			suffix:=" "
			fmt.Println("Suffix:",suffix)
		}else {
			suffix=newarrayinarrayformat[lengthofnewarrayinformat-1]
			fmt.Println("Suffix:",suffix)
		}
	}else {
		notprefixsuffix=newarrayinarrayformat[0]
		fmt.Println("Element not prefix and suffix:",notprefixsuffix)
	}
	//converting original string to array for comparison to get prefix and suffix
	var sumofgoroutnies int64
	for i:=0;i<len(newarrayinarrayformat);i++ {
		xx, _ := strconv.ParseInt(newarrayinarrayformat[i], 10, 64)
		sumofgoroutnies=sumofgoroutnies+xx
	}
	fmt.Println("Partial Sum:",sumofgoroutnies)
	time.Sleep(time.Millisecond*50)
	x1:=endwork{sumofgoroutnies,prefix,suffix,notprefixsuffix}
	xx:=[]endwork{x1}
	secondjson,_:=json.Marshal(xx)
	fmt.Println(string(secondjson))
	channel1<-string(secondjson)

}

