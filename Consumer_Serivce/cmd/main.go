package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	sensor "github.com/snipep/Assessment/Consumer-Service/Proto"
	"google.golang.org/grpc"
)

// // server is used to implement the SensorService server.
// type server struct {
// 	sensor.UnimplementedSensorServiceServer
// }

// // SendSensorData implements the server-side logic for the streaming RPC.
// // It receives a stream of SensorDataRequest messages from the client.
// func (s *server) SendSensorData(stream sensor.SensorService_SendSensorDataServer) error {
// 	log.Println("Client connected. Starting to receive stream.")
// 	var dataCount int32

// 	// Loop to receive messages from the stream.
// 	for {
// 		// Recv() blocks until a message is received or the stream is closed.
// 		req, err := stream.Recv()

// 		// If the stream is closed by the client, io.EOF is returned.
// 		if err == io.EOF {
// 			log.Printf("Finished receiving data. Total messages: %d", dataCount)
// 			// Send a final response back to the client and close the connection.
// 			return stream.SendAndClose(&sensor.SensorDataResponse{
// 				Status:  "Success",
// 				Message: fmt.Sprintf("Successfully processed %d data points.", dataCount),
// 			})
// 		}
// 		// Handle any other errors during reception.
// 		if err != nil {
// 			log.Printf("Error while receiving stream: %v", err)
// 			return err
// 		}

// 		// Process the received data.
// 		data := req.GetData()
// 		log.Printf("Received: Type=%s, Value=%.2f, ID1=%s, ID2=%d",
// 			data.GetSensorType(), data.GetSensorValue(), data.GetId1(), data.GetId2())
// 		dataCount++
// 	}
// }

func main() {
	// Listen for incoming TCP connections on port 50051.
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	} 
	gRPCPort := os.Getenv("GPRC_PORT")
	if gRPCPort == "" {
		gRPCPort = ":50051" 	// Default port if not set in .env
	}
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server instance.
	s := grpc.NewServer()

	// Register our server implementation with the gRPC server.
	sensor.RegisterSensorServiceServer(s, &Server{})

	log.Println("Server listening at", lis.Addr())

	// Start serving requests. This is a blocking call.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
