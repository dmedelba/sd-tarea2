package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

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
			intSlice = propuesta[:len(propuesta)-1]
		}
	}
	return propuesta, cant
}

func evaluarPropuesta(propuesta string) {
	//pasar propuesta a lista
	propuestita := stringToList(propuesta)
	maquinitas := []int{70, 71, 72}
	var cant int
	var total int
	total = 0

	//recorro la lista de maquinas paraverificar nodos caidos
	var conn *grpc.ClientConn
	for i := 0; i < len(maquinitas); i++ {
		conn, err := grpc.Dial("dist"+maquinitas[i]+":6009", grpc.WithInsecure())
		if err != nil {
			log.Printf("Maquina caida")
			propuestita, cant = borrarMaquina(propuestita, maquinitas[i])
			total = cant + total
			continue
		}
		defer conn.Close()
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
