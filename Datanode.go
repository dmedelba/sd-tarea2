package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"./uploader"

	"./propu"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SubirLibro(ctx context.Context, request *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("recibi la wea")

	//creo la carpeta para guardar chunks del libro
	idChunk := strconv.Itoa(int(request.Id))
	fileName := "./libros_subidos/" + request.NombreLibro + "-" + idChunk
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//a la funcion pasar el tipo de exlusión mutua
	propuestaInicial := crearPropuestaInicial(request.NombreLibro, int(request.Cantidad))
	enviarPropuesta(propuestaInicial, request.TipoExclusionMutua)

	log.Printf("Propuesta enviada")
	return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
}

func crearPropuestaInicial(nombreLibro string, cantidadChunks int) []int32 {
	//creamos la propuesta inicial simple
	propuestaMaquinas := make([]int32, cantidadChunks)
	var indice = 0
	var maquina = 70
	for i := 0; i < cantidadChunks; i++ {
		maquina += indice
		propuestaMaquinas[i] = int32(maquina)
		indice++
		if indice == 3 {
			indice = 0
			maquina = 70
		}
	}
	return propuestaMaquinas
}
func enviarPropuesta(propuesta []int32, tipoExclusion string) {
	//enviar propuesta
	if tipoExclusion == "1" {
		//es centralizada, preguntar al name node
		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist69:6000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error al conectarse con la maquina 69 [Name node]. %s", err)
		}
		defer conn.Close()

		c := propu.NewPropuClient(conn)
		decision, _ := c.EnviarPropuesta(context.Background(), &propu.Propuesta_Generada{
			ListaPropuesta: propuesta,
		})

		//aprobado o rechazo
		log.Printf("hola")
	}
	/*
		else{

			//rechazo, llama a la funcion que crea nueva propuesta y hace recursividad
			propuestaInicial := generarNuevaPropuesta(propuestaInicial)
			decision:= enviarPropuesta( propuestaInicial, request.TipoExclusionMutua)
			if (decision == "1"){
				//propuesta aceptada
				break
			}
		}
	*/

	//falta el else en caso de que sea distribuida

}

func generarNuevaPropuesta(propuestaMaquinas []int32) []int32 {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(propuestaMaquinas), func(i, j int) {
		propuestaMaquinas[i], propuestaMaquinas[j] = propuestaMaquinas[j], propuestaMaquinas[i]
	})
	return propuestaMaquinas
}

func propuestaToString(propuestaMaquinas []int32, nombreLibro string) string {
	cantidadChunks := len(propuestaMaquinas)
	cChunks_str := strconv.Itoa(cantidadChunks)
	propuesta := nombreLibro + " " + cChunks_str + "\n"

	for i := 0; i < cantidadChunks; i++ {
		chunk := strconv.Itoa(i)
		maquina := propuestaMaquinas[i]
		maquina_str := strconv.Itoa(int(maquina))
		propuesta += nombreLibro + "-" + chunk + " dist" + maquina_str + "\n"
	}
	return propuesta
}

func main() {
	log.Printf("[Datanode]")
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	propu.RegisterPropuServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
