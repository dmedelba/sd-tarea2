package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"./propu"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) EnviarPropuesta(ctx context.Context, in *propu.Propuesta_Generada) (*propu.Respuesta_Propuesta, error) {
	listaPropuesta := in.ListaPropuesta
	fmt.Printf("recibi algo")
	fmt.Println(listaPropuesta)
	return &propu.Respuesta_Propuesta{Respuesta: "PROPUESTA_RECIBIDA POR EL NAMENODE"}, nil
}

func main() {
	log.Printf("[Namenode]")
	lis, err := net.Listen("tcp", ":6002")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	propu.RegisterPropuServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
