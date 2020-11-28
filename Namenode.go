package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) EnviarPropuesta(ctx context.Context, request *uploader.Propuesta_Generada) (*uploader.Respuesta_Propuesta, error) {
	listaPropuesta := request.ListaPropuesta
	fmt.Println(listaPropuesta)
}

func main() {
	log.Printf("[Namenode]")
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	uploader.RegisterUploaderServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
