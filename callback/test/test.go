package main

import (
	"fmt"
	"strings"
)

func main() {

	var a map[string][]string = make(map[string][]string)

	a["map[username]"] = []string{"root"}
	//key := "username"
	for k , v := range a{

		i := strings.IndexByte(k, '[')
		fmt.Println(k  ,v)
		fmt.Println(k[0:i])
		fmt.Println(i)


	}
}
