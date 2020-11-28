package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SubirLibro(ctx context.Context, request *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("recibi la wea")

	//creo la carpeta para guardar chunks del libro
	idChunk := strconv.Itoa(int(request.Id))
	fileName := "./libros_subidos/" + request.NombreLibro[0:15] + "_" + idChunk
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf(request.NombreLibro)
	return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
}

/*
func EnviarPropuesta(conn *grpc.ClientConn) {

	if Cantidad%3 == 0 {
		totalpormaquina = Cantidad / 3

		// propuesta para maquina x
		for i := 0; i < totalpormaquina; i++ {
			// sacamos de la carpeta el chunk con nombre = monbre libro
			// y id == a i
			// y le asignamos la maquina x

		}

	}
}
*/

func main() {
	log.Printf("[Datanode]")
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
