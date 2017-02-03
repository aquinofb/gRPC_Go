package main

import (
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "../customer"
)

const (
	address = "localhost:50051"
)

func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "Felipe Aquino",
		Email: "aquinofb@gmail.com",
		Phone: "555-555555",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "27 Kilmainham Orchard",
				City:              "Dublin",
				State:             "Dublin",
				Zip:               "D8",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "Beira Mar",
				City:              "Fortaleza",
				State:             "CE",
				Zip:               "123123",
				IsShippingAddress: true,
			},
		},
	}

	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    102,
		Name:  "Maria Olivia",
		Email: "mariaolivia@gmail.com",
		Phone: "555-555555",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Av Beira Mar",
				City:              "Fortaleza",
				State:             "CE",
				Zip:               "010101",
				IsShippingAddress: true,
			},
		},
	}

	createCustomer(client, customer)
	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)
}
