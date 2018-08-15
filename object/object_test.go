package object

import "testing"

func TestStringHashKey(t *testing.T){
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is jonny"}
	diff2 := &String{Value: "My name is jonny"}

	if hello1.HashKey() != hello2.HashKey(){
		t.Errorf("strings with same content have diffrent hash key")
	}
	if diff1.HashKey() != diff2.HashKey(){
		t.Errorf("strings with same content have diffrent hash key")
	}
	if hello1.HashKey() == diff1.HashKey(){
		t.Errorf("strings with diffrent content have some hash keys")
	}

}