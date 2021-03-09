package main

import "net"

type Interface struct{
	net.Interface
}



func main() {
	var inter Interface
	inter.Name = "eth0"
	inter.Flags = 1
	inter.Index = 1


}