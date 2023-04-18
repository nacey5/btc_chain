package data_struct

import (
	"fmt"
	"testing"
)

func TestCreateMHTree(t *testing.T) {
	var fileName string
	fmt.Println("input the origin fileName")
	fmt.Scanln(&fileName)
	mhTree1 := CreateMHTree(fileName)
	fmt.Println("please input the compare file")
	fmt.Scanln(&fileName)
	mhTree2 := CreateMHTree(fileName)

	hash1 := mhTree1.GetRootHash()
	hash2 := mhTree2.GetRootHash()
	if hash1 == hash2 {
		fmt.Println("the user not changes the data")
	} else {
		fmt.Println("the user changes the data")
	}
}
