package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"

	"./uploader"

	//"path/filepath"
	"google.golang.org/grpc"
)

const (
	puerto = "dist70:6000"
)

func subirLibro(conn *grpc.ClientConn, tipo string) {
	//buscamos libro, se selecciona y se descompone
	//conexion con el datanode
	nombreLibroSeleccionado := mostrarLibros() //se muestran los libros y se selecciona el libro a subir
	cantidadChunks := generarChunks(nombreLibroSeleccionado)

	c := uploader.NewUploaderClient(conn)
	for i := 0; i < cantidadChunks; i++ {
		contenidoChunk := abrirChunk(nombreLibroSeleccionado, i)
		c.SubirLibro(context.Background(), &uploader.Solicitud_SubirLibro{
			Chunk:              contenidoChunk,
			Id:                 int32(i),
			NombreLibro:        nombreLibroSeleccionado,
			Cantidad:           int32(cantidadChunks),
			TipoExclusionMutua: tipo,
		})
	}
	log.Printf("OK?")

}

/*
func subirLibroDistribuido(conn *grpc.ClientConn) {
	//buscamos libro, se selecciona y se descompone
	var libroSeleccionado int
	log.Printf("----------------------------------")
	mostrarLibros()
	log.Printf("----------------------------------")
	log.Printf("Seleccione un libro a descargar.")
	log.Printf("----------------------------------")
	fmt.Scanln(&libroSeleccionado)
}
*/
func abrirChunk(nombreLibro string, indice int) []byte {
	indiceStr := strconv.Itoa(indice)
	file, err := os.Open("./" + nombreLibro + "-" + indiceStr)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func generarChunks(nombreLibroSeleccionado string) int {
	//funcion obtenida de https://www.socketloop.com/tutorials/golang-recombine-chunked-files-example
	fileToBeChunked := "./libros/" + nombreLibroSeleccionado
	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()
	const fileChunk = 250000 //250kbytes

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Dividiendo el libro en %d partes.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := nombreLibroSeleccionado + "-" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to : ", fileName)
	}
	//devuelve la cantidad de partes que se dividio
	return int(totalPartsNum)
}

func mostrarLibros() string {
	//presentamos al usuario los libros para que seleccione
	var nombreLibroSeleccionado string
	log.Printf("----------------------------------")
	log.Printf("Seleccione un libro a descargar.")
	log.Printf("----------------------------------")
	path := "./libros/"
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	i := 1
	for _, val := range lst {
		if val.IsDir() {
			log.Printf("[%s]\n", val.Name())
		} else {
			s := strconv.Itoa(i)
			log.Println(s + ". " + val.Name())
			i++
		}
	}
	var libroSeleccionado int
	log.Printf("----------------------------------")
	fmt.Scanln(&libroSeleccionado)

	//encontramos el libro seleccionado
	indice := 1
	for _, val := range lst {
		if val.IsDir() {
			log.Printf("[%s]\n", val.Name())
		} else {
			if libroSeleccionado == indice {
				nombreLibroSeleccionado = val.Name()
			}
			indice++
		}
	}
	return nombreLibroSeleccionado
}

//Establecemos conexión con logisitica dist70:6000

func main() {
	//crear conexion
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo establecer la conexión. ERROR: %v", err)
	}

	defer conn.Close()
	//seleccionamos que quiere hacer el cliente
	var seleccion int
	flag := true
	for flag {
		log.Printf("[Cliente] Ingrese opición a realizar:")
		log.Printf("1. Subir libro")
		log.Printf("2. Descargar libro")
		log.Printf("3. Ver biblioteca (libros descargados)")
		log.Printf("4. Finalizar")
		fmt.Scanln(&seleccion)

		var tipoSubida int
		switch seleccion {
		case 1:
			//subir libro, puede subir centralizado o distribuido
			log.Printf("[Cliente] Seleccione el tipo de subida a realizar:")
			log.Printf("1. Subida Centralizada")
			log.Printf("2. Subida Distribuida")
			fmt.Scanln(&tipoSubida)
			switch tipoSubida {
			case 1:
				subirLibro(conn, "1")
				//centralizada
			case 2:
				//distribuida
				//subirLibro(conn,"2")
			}
		case 2:
			//descargar libro, conectarse al name node (69)
		case 3:
			//ver biblioteca

		case 4:
			//finalizar
			log.Printf("Sesión finalizada. Muchas gracias!")
			flag = false
		}
	}
}

//Cliente y datanode es con un datanode aleatorio. [randint=>0-3]
//Si siempre es el mismo nodo, nunca va a existir condición de carrera en el log del name node.
