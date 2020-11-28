package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"

	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

type chank struct {
	Chunk              byte
	Id                 int32
	NombreLibro        string
	Cantidad           int32
	TipoExclusionMutua string
}

func (s *server) SubirLibro(ctx context.Context, request *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("recibi la wea")

	//creo la carpeta para guardar chunks del libro
	if _, err := os.Stat(request.NombreLibro); os.IsNotExist(err) {
		err = os.Mkdir(request.NombreLibro, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Guardo los chunks en la carperta recien creada
	for i := 0; i < int(request.Cantidad); i++ {
		chunksito := chank{
			Chunk:              request.Chunk,
			Id:                 request.Id,
			NombreLibro:        request.NombreLibro,
			Cantidad:           request.Cantidad,
			TipoExclusionMutua: request.TipoExclusionMutua,
		}

		dst, err := os.Create(filepath.Join(request.NombreLibro, filepath.Base(chunksito.NombreLibro+"-"+i))) // dir is directory where you want to save file.
		if err != nil {
			checkErr(err)
		}
		defer dst.Close()
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
