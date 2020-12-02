package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	//cantChunks := len(listaPropuesta)
	fmt.Printf("Propuesta recibida, a evaluar")
	fmt.Printf(listaPropuesta)
	//evaluamos la propuesta, si hay una maquina que no funcione el namenode genera una nueva propuesta con las maquinas activas.
	nuevaPropuesta := evaluarPropuesta(listaPropuesta)
	//si cambio, entregara la nueva propuesta, si no, entregará la misma.
	//Escribir en el log

	return &propu.Respuesta_Propuesta{Respuesta: nuevaPropuesta}, nil
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
func ListToString(lista []int) string {
	var propuestaString = ""
	for i := 0; i < len(lista); i++ {
		maquina := lista[i]
		maquinaStr := strconv.Itoa(maquina)
		propuestaString += maquinaStr + ","
	}
	return propuestaString
}
func borrarMaquina(propuesta []int, value int) []int {
	maquinas := []int{70, 71, 72}
	//eliminar maquina que no esta funcionando de nuestra lista maquinas
	for i := 0; i < len(maquinas); i++ {
		if maquinas[i] == value {
			copy(maquinas[i:], maquinas[i+1:])
			maquinas[len(maquinas)-1] = 0
			maquinas = maquinas[:len(maquinas)-1]
		}
	}
	//reemplazar la maquina que esta caida con una que no, de manera random.
	maquinaElegida := rand.Intn(len(maquinas))
	for i := 0; i < len(propuesta); i++ {
		if value == propuesta[i] {
			propuesta[i] = maquinas[maquinaElegida]
		}
	}
	return propuesta
}

func evaluarPropuesta(propuesta string) string {
	//pasar propuesta a lista
	propuestita := stringToList(propuesta)
	maquinitas := []int{70, 71, 72}
	//recorro la lista de maquinas para verificar nodos caidos
	for i := 0; i < len(maquinitas); i++ {
		numeroMaquina := strconv.Itoa(maquinitas[i])
		log.Printf(numeroMaquina)

		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist"+numeroMaquina+":6000", grpc.WithInsecure())

		if err != nil {
			log.Fatalf("Error de envio de mensaje %s", err)
		}

		defer conn.Close()

		c := uploader.NewUploaderClient(conn)
		conexion, error := c.EstadoMaquina(context.Background(), &uploader.Solicitud_EstadoMaquina{
			EstadoMaquina: "1",
		})

		if error != nil {
			//log.Printf(conexion.EstadoMaquina)
			log.Printf("dist" + numeroMaquina + ":6000, Maquina caida")
			propuestita = borrarMaquina(propuestita, maquinitas[i])
		} else {
			log.Printf("Maquina funcionando")
			log.Printf(conexion.EstadoMaquina)
		}
	}
	propuestitaString := ListToString(propuestita)
	return propuestitaString
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
