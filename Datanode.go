package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"./propu"
	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

//transformamos la lista de propuesta a string para enviar por protobuffer
func ListToString(lista []int) string {
	var propuestaString = ""
	for i := 0; i < len(lista); i++ {
		maquina := lista[i]
		maquinaStr := strconv.Itoa(maquina)
		propuestaString += maquinaStr + ","
	}
	return propuestaString
}

//el string a la lista de propuesta
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

//se crea una propuesta inicial simple, asignando cantidades equitativas a cada maquina.
func crearPropuestaInicial(cantidadChunks int) []int {
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

//se envia la propuesta al namenode (centralizado) , tipoExclusion indica que tipo de exclusion hará
func enviarPropuesta(propuesta string, tipoExclusion string, conn *grpc.ClientConn, NombreLibro string) string {
	//enviar propuesta
	var propuestaDistribucion string
	if tipoExclusion == "1" {
		//es centralizada, preguntar al name node

		log.Printf("Propuesta a enviar:")
		log.Printf(propuesta)
		c := propu.NewPropuClient(conn)
		respuestita, err := c.EnviarPropuesta(context.Background(), &propu.Propuesta_Generada{
			ListaPropuesta: propuesta,
			NombreLibro:    NombreLibro,
		})

		if err != nil {
			log.Fatalf("Error de envio de mensaje %s", err)
		}
		//aprobado o rechazo
		log.Printf("Decision:")
		log.Printf(respuestita.Respuesta)
		propuestaDistribucion = respuestita.Respuesta
	}
	return propuestaDistribucion
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

//leemos el chunk para ver su contenido y almacenar
func leerChunk(nombreLibro string, indice int) []byte {
	indiceStr := strconv.Itoa(indice)
	file, err := os.Open("./libros_subidos/" + nombreLibro + "-" + indiceStr)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

//enviamos a los otros datanode los chunks , depende de la propuesta.
func distribuirChunks(distribucion string, nombreLibro string) {
	//recorrer la lista y enviar a los chunks correspondientes.
	listaChunks := stringToList(distribucion)
	for i := 0; i < len(listaChunks); i++ {
		//tengo que enviar el id del chunk, el nombre y el contenido.
		contenidoChunk := leerChunk(nombreLibro, i)
		maquinaStr := strconv.Itoa(listaChunks[i])
		//crear la conexion con la maquina en cuestion
		conn, err := grpc.Dial("dist"+maquinaStr+":5050", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("No se pudo conectar al datanode para distribuir: %s", err)
		}
		defer conn.Close()
		c := uploader.NewUploaderClient(conn)
		listo, _ := c.Distribuir(context.Background(), &uploader.Solicitud_Distribucion{
			IdChunk:        int32(i),
			NombreLibro:    nombreLibro,
			ContenidoChunk: contenidoChunk,
		})
		if listo.Respuesta == "1" {
			log.Printf("Se ha recibido el chunk en la maquina" + maquinaStr)
		}

	}
}

//recibo los libros desde el cliente, los almaceno, genero propuesta y envio segun tipo de exclusión.
func (s *server) SubirLibro(ctx context.Context, in *uploader.Solicitud_SubirLibro) (*uploader.Respuesta_SubirLibro, error) {
	log.Printf("[Datanode] Chunks recibidos por parte del cliente")

	//Recibimos los chunks desde el cliente
	//creo la carpeta para guardar chunks del libro
	idChunk := strconv.Itoa(int(in.Id))
	fileName := "./libros_subidos/" + in.NombreLibro + "-" + idChunk
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//guardo el contenido del chunk
	ioutil.WriteFile(fileName, in.Chunk, os.ModeAppend)
	//una vez creado los chunks por el nodo, creo la propuesta
	if int(in.Id) == int(in.Cantidad)-1 {
		fmt.Printf("Se crearon todos los chunks")
		//a la funcion pasar el tipo de exlusión mutua
		propuestaInicial := crearPropuestaInicial(int(in.Cantidad))
		propuestaInicialString := ListToString(propuestaInicial)
		//conexion
		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist69:5050", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Error al conectarse con la maquina 69 [Name node]. %s", err)
		}
		defer conn.Close()
		log.Printf("Propuesta enviada")

		distribucion := enviarPropuesta(propuestaInicialString, in.TipoExclusionMutua, conn, in.NombreLibro)
		distribuirChunks(distribucion, in.NombreLibro)

		return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
	}

	return &uploader.Respuesta_SubirLibro{Respuesta: int32(0)}, nil
}

//respondemos cual es el estado de la maquina
func (s *server) EstadoMaquina(ctx context.Context, respuesta *uploader.Solicitud_EstadoMaquina) (*uploader.Respuesa_EstadoMaquina, error) {
	return &uploader.Respuesa_EstadoMaquina{EstadoMaquina: "1"}, nil
}

//funcion que recibe los chunks luego de enviada la distribución (propuesta aceptada)
func (s *server) Distribuir(ctx context.Context, respuesta *uploader.Solicitud_Distribucion) (*uploader.Respuesta_Distribucion, error) {
	log.Printf("Guardando chunk en la maquina correspondiente:")
	//Recibimos el chunk correspondiente desde el nodo distribución
	idChunk := strconv.Itoa(int(respuesta.IdChunk))
	fileName := "./mis_chunks/" + respuesta.NombreLibro + "-" + idChunk
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//guardo el contenido del chunk
	ioutil.WriteFile(fileName, respuesta.ContenidoChunk, os.ModeAppend)
	return &uploader.Respuesta_Distribucion{Respuesta: "1"}, nil
}

func main() {
	log.Printf("[Datanode]")
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	uploader.RegisterUploaderServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
