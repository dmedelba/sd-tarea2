package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"

	"./propu"
	"./uploader"
	"google.golang.org/grpc"
)

type server struct {
}

//candado para regular el acceso al namenode.
var ocupado bool = false

//recibimos la propuesta , aceptamos o rechazamos y respondemos con la nueva propuesta, puede ser la misma.
func (s *server) EnviarPropuesta(ctx context.Context, in *propu.Propuesta_Generada) (*propu.Respuesta_Propuesta, error) {
	for ocupado {
		log.Printf("[Name node] Se está procesando otra solicitud en estos momentos. Esperar")
	}
	ocupado = true
	listaPropuesta := in.ListaPropuesta
	nombreLibro := in.NombreLibro
	fmt.Printf("Propuesta recibida, a evaluar")
	fmt.Printf(listaPropuesta)
	//evaluamos la propuesta, si hay una maquina que no funcione el namenode genera una nueva propuesta con las maquinas activas.
	nuevaPropuesta := evaluarPropuesta(listaPropuesta)
	//si cambio, entregara la nueva propuesta, si no, entregará la misma.
	//Escribir en el log ya que es una propuesta aceptada

	textoPropuesta := propuestaToString(stringToList(nuevaPropuesta), nombreLibro)
	file, err := os.OpenFile("./log.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo log: %s", err)
	}
	defer file.Close()
	//escribimos en el archivo
	file.WriteString(textoPropuesta)
	ocupado = false
	return &propu.Respuesta_Propuesta{Respuesta: nuevaPropuesta}, nil
}

//respondemos los libros disponibles para que el cliente seleccione el libro a descargar
func (s *server) VerLibros(ctx context.Context, in *propu.Solicitud_VerLibros) (*propu.Respuesta_VerLibros, error) {
	//leer el LOG.txt y enviar los nombres de los libros disponibles.
	listadoLibros := ""
	file, _ := os.Open("./log.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineas := strings.Split(scanner.Text(), " ")
		if isInt(lineas[1]) {
			listadoLibros += lineas[0] + ","
		}
	}
	return &propu.Respuesta_VerLibros{LibrosDisponibles: listadoLibros}, nil
}

//enviamos la ubicacion del libro
func (s *server) VerUbicaciones(ctx context.Context, in *propu.Solicitud_Ubicaciones) (*propu.Respuesta_Ubicaciones, error) {
	//enviar las ubicaciones, leer el archivo LOG.TXT y enviar
	nombreLibro := in.NombreLibro
	listadoMaquinas := ""
	file, _ := os.Open("./log.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineas := strings.Split(scanner.Text(), " ")
		//primeras lineas del archivo
		if lineas[0] == nombreLibro {
			cantChunks, _ := strconv.Atoi(lineas[1])
			for i := 0; i < cantChunks; i++ {
				scanner.Scan()
				listadoMaquinas += scanner.Text() + ","
			}
			break
		}
	}
	return &propu.Respuesta_Ubicaciones{Ubicaciones: listadoMaquinas}, nil
}

//transformamos la propuest apara escribirla en el lOG.txt
func propuestaToString(propuestaMaquinas []int, nombreLibro string) string {
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

//eliminamos la maquina caida y reemplazamos por una que sirva en la propuesta.
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
	for i := 0; i < len(propuesta); i++ {
		if value == propuesta[i] {
			maquinaElegida := rand.Intn(len(maquinas))
			propuesta[i] = maquinas[maquinaElegida]
		}
	}
	return propuesta
}

//verificamos si la propuesta es aceptada o no. Maquinas caida o no.
func evaluarPropuesta(propuesta string) string {
	//pasar propuesta a lista
	propuestita := stringToList(propuesta)
	maquinitas := []int{70, 71, 72}
	//recorro la lista de maquinas para verificar nodos caidos
	for i := 0; i < len(maquinitas); i++ {
		numeroMaquina := strconv.Itoa(maquinitas[i])
		log.Printf(numeroMaquina)

		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist"+numeroMaquina+":5050", grpc.WithInsecure())

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
			log.Printf("dist" + numeroMaquina + ":5050, Maquina caida")
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

//String es un entero ? https://stackoverflow.com/questions/22593259/check-if-string-is-int
func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
func main() {
	log.Printf("[Namenode]")
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Error al tratar de escuchar: %v", err)
	}
	s := grpc.NewServer()
	propu.RegisterPropuServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
