package main

import (
	"authentication/internal/apple"
	"fmt"
)

func main() {
	s, err := apple.New("./conf/apple/AuthKey.key", "", "", "", "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := s.AuthTokenWithApp("")
	if err != nil {
		return
	}

	fmt.Println(resp)
}
