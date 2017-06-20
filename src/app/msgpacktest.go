package main

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"fmt"
)

type testStruct struct {
	A 			int
}

func marshal_args(i ...interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}


func msgpack_test() {


	bret, err := marshal_args(1, "aaa", []byte{2, 2, 2}, &testStruct{A: 111})
	if err != nil {
		fmt.Println("marshal error")
	}

	var aret[]interface{}
	err = msgpack.Unmarshal(bret, &aret)

	fmt.Println(bret, aret)
}

func unpack(i, s, sli, st interface{}) {
	fmt.Println(i, s, sli, st)
}

func msgUnpackTest() {

	bret, err := marshal_args(1, "aaa", []byte{2, 2, 2}, &testStruct{A: 111})
	if err != nil {
		fmt.Println("marshal error ")
	}

	fmt.Println("marshal ret bytes " , bret)


	var rret []interface{}
	err = msgpack.Unmarshal(bret, &rret)
	fmt.Println("unmarshal ret is  ", rret)

	//unpack(rret[0], rret[1], rret[2], rret[3])
}

