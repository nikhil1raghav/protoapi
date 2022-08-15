package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"protoapi/protobufexample"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RandomServer struct {
	protobufexample.UnimplementedRandomServer
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
func getString(len int64) string {
	temp := ""
	startChar := "!"
	var i int64 = 1
	for {
		// For getting valid ASCII characters
		myRand := random(0, 94)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == len {
			break
		}
		i++
	}
	return temp
}

func (RandomServer) GetDate(ctx context.Context, r *protobufexample.RequestDateTime) (*protobufexample.DateTime, error) {
	currentTime := time.Now()
	response := &protobufexample.DateTime{
		Value: currentTime.String(),
	}
	return response, nil
}

func (RandomServer) GetRandom(ctx context.Context, r *protobufexample.RandomParams) (*protobufexample.RandomInt, error) {
	rand.Seed(r.GetSeed())
	place := r.GetPlace()
	temp := rand.Int()
	for {
		place--
		if place <= 0 {
			break
		}
		temp = rand.Int()
	}
	response := &protobufexample.RandomInt{
		Value: int64(temp),
	}
	return response, nil
}

func (RandomServer) GetRandomPass(ctx context.Context, r *protobufexample.RequestPass) (*protobufexample.RandomPass, error) {
	rand.Seed(r.GetSeed())
	temp := getString(r.GetLength())
	response := &protobufexample.RandomPass{
		Password: temp,
	}
	return response, nil
}

var port = ":8080"

func main() {
	if len(os.Args) != 1 {
		port = os.Args[1]
	}
	server := grpc.NewServer()
	var randomServer RandomServer
	protobufexample.RegisterRandomServer(server, randomServer)
	reflection.Register(server)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Serving Requests....")
	server.Serve(listen)

}
