package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"./uploader"

	"./propu"
	//"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
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
func crearPropuestaInicial(nombreLibro string, cantidadChunks int) []int {
	//creamos la propuesta inicial simple
	propuestaMaquinas := make([]int, cantidadChunks)
	var indice = 0
	var maquina = 70
	for i := 0; i < cantidadChunks; i++ {

		propuestaMaquinas[i] = int(maquina)
		indice++
		maquina++
		if indice == 3 {
			indice = 0
			maquina = 70
		}
	}
	return propuestaMaquinas
}
func enviarPropuesta(propuesta string, tipoExclusion string, conn *grpc.ClientConn) {
	//enviar propuesta
	if tipoExclusion == "1" {
		//es centralizada, preguntar al name node

		log.Printf("Propuesta a enviar:")
		log.Printf(propuesta)
		c := propu.NewPropuClient(conn)
		decision, error := c.EnviarPropuesta(context.Background(), &propu.Propuesta_Generada{
			ListaPropuesta: propuesta,
		})

		if error != nil {
			log.Fatalf("Error al tratar de enviar la propuesta al [Name node]. %s", error)
		}

		//aprobado o rechazo
		log.Printf("Decision:")
		log.Printf(decision.Respuesta)
	}
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

func generarNuevaPropuesta(propuestaMaquinas []int32) []int32 {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(propuestaMaquinas), func(i, j int) {
		propuestaMaquinas[i], propuestaMaquinas[j] = propuestaMaquinas[j], propuestaMaquinas[i]
	})
	return propuestaMaquinas
}
func propuestaToString(propuestaMaquinas []int32, nombreLibro string) string {
	cantidadChunks := len(propuestaMaquinas)
	cChunksStr := strconv.Itoa(cantidadChunks)
	propuesta := nombreLibro + " " + cChunksStr + "\n"

	for i := 0; i < cantidadChunks; i++ {
		chunk := strconv.Itoa(i)
		maquina := propuestaMaquinas[i]
		maquinaStr := strconv.Itoa(int(maquina))
		propuesta += nombreLibro + "-" + chunk + " dist" + maquinaStr + "\n"
	}
	return propuesta
}
func (s *server) SubirLibro(ctx context.Context, in *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("recibi la wea")

	//creo la carpeta para guardar chunks del libro
	idChunk := strconv.Itoa(int(in.Id))
	fileName := "./libros_subidos/" + in.NombreLibro + "-" + idChunk
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//una vez creado los chunks por el nodo, creo la propuesta
	if int(in.Id) == int(in.Cantidad)-1 {
		fmt.Printf("Se crearon todos los chunks")
		//a la funcion pasar el tipo de exlusión mutua
		propuestaInicial := crearPropuestaInicial(in.NombreLibro, int(in.Cantidad))
		propuestaInicialString := ListToString(propuestaInicial)
		//conec
		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist69:6006", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error al conectarse con la maquina 69 [Name node]. %s", err)
		}
		defer conn.Close()
		enviarPropuesta(propuestaInicialString, in.TipoExclusionMutua, conn)

		log.Printf("Propuesta enviada")
		return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
	}

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
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
