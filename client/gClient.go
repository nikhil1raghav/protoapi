package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	pb "protoapi/protobufexample"
	"time"

	"google.golang.org/grpc"
)

var port = ":8080"

func AskingDateTime(ctx context.Context, m pb.RandomClient) (*pb.DateTime, error) {
	request := &pb.RequestDateTime{
		Value: "Doesn't matter what you write here",
	}
	return m.GetDate(ctx, request)
}
func AskPass(ctx context.Context, m pb.RandomClient, seed int, length int) (*pb.RandomPass, error) {
	request := &pb.RequestPass{
		Seed:   int64(seed),
		Length: int64(length),
	}
	return m.GetRandomPass(ctx, request)
}
func AskRandom(ctx context.Context, m pb.RandomClient, seed int, place int) (*pb.RandomInt, error) {
	request := &pb.RandomParams{
		Seed:  int64(seed),
		Place: int64(place),
	}
	return m.GetRandom(ctx, request)
}
func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
func main() {
	if len(os.Args) != 1 {
		port = os.Args[1]
	}
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	rand.Seed(time.Now().Unix())
	seed := rand.Intn(100)

	client := pb.NewRandomClient(conn)
	r, err := AskingDateTime(context.Background(), client)
	handleError(err)

	fmt.Println("Server Data & Time : ", r.Value)

	length := rand.Intn(20)
	place := rand.Intn(10)

	d, err := AskPass(context.Background(), client, seed, length)
	handleError(err)
	fmt.Println("Random pass : ", d.Password)

	num, err := AskRandom(context.Background(), client, seed, place)
	handleError(err)
	fmt.Println("Random number : ", num.Value)

}
