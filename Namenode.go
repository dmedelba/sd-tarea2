package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"./propu"
	"./uploader"
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

func stringToList(texto string) []int {
	lista := strings.Split(texto, ",")
	listaInt := make([]int, len(lista)-1)
	//convertir a int
	for i, s := range lista {
		if i == len(lista)-1 {
			break
		}
		listaInt[i], _ = strconv.Atoi(s)
	}
	return listaInt
}

func borrarMaquina(propuesta []int, value int) ([]int, int) {
	var cant int
	cant = 0
	for i := 0; i < len(propuesta); i++ {
		if value == propuesta[i] {
			cant = cant + 1
			copy(propuesta[i:], propuesta[i+1:])
			propuesta[len(propuesta)-1] = 0
			intSlice := propuesta[:len(propuesta)-1]
			fmt.Printf(intSlice)
		}
	}
	return propuesta, cant
}

func evaluarPropuesta(propuesta string) {
	//pasar propuesta a lista
	propuestita := stringToList(propuesta)
	maquinitas := []int{70, 71, 72}
	var cant int
	total := 0
	//recorro la lista de maquinas para verificar nodos caidos

	var conn *grpc.ClientConn

	for i := 0; i < len(maquinitas); i++ {
		numeroMaquina := strconv.Itoa(maquinitas[i])
		conn, _ := grpc.Dial("dist"+numeroMaquina+":6009", grpc.WithInsecure())
		defer conn.Close()

		c := uploader.NewUploaderClient(conn)
		conexion, error := c.EstadoMaquina(context.Background(), &uploader.Solicitud_EstadoMaquina{
			EstadoMaquina: "1",
		})

		if error != nil {
			log.Printf("dist" + numeroMaquina + ":6009, Maquina caida")
			propuestita, cant = borrarMaquina(propuestita, maquinitas[i])
			total = cant + total
		}

	}

	//verificar maquinas caidas
}

func main() {
	log.Printf("[Namenode]")
	lis, err := net.Listen("tcp", ":6006")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	propu.RegisterPropuServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
