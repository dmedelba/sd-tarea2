package main

import (
	"context"
	"log"
	"net"

	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SubirLibro(ctx context.Context, request *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("recibi la wea")
	log.Printf(request.NombreLibro)
	return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
}

func main() {
	log.Printf("[Datanode]")
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	uploader.RegisterUploaderServer(s, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
