package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	//"os"
	//"path/filepath"
	//"google.golang.org/grpc"
)

const (
	puerto = "dist70:6000"
)

/*
func subirLibroCentralizado(conn *grpc.ClientConn) {
	//buscamos libro, se selecciona y se descompone
	var libroSeleccionado int
	log.Printf("----------------------------------")
	mostrarLibros()
	log.Printf("----------------------------------")
	log.Printf("Seleccione un libro a descargar.")
	log.Printf("----------------------------------")
	fmt.Scanln(&libroSeleccionado)

}

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
func mostrarLibros() {
	//presentamos al usuario los libros para que seleccione
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
				log.Printf(val.Name())
			}
			indice++
		}
	}
}

//Establecemos conexión con logisitica dist70:6000

/*
	var conn *grpc.ClientConn
	conn, e	rr := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo establecer la conexión. ERROR: %v", err)
	}

	defer conn.Close()
*/
func main() {
	//crear conexion
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
				//subirLibroCentralizado(conn)
				//centralizada
			case 2:
				//distribuida
				//subirLibroDistribuido(conn)
			}
		case 2:
			//descargar libro, conectarse al name node (69)
		case 3:
			mostrarLibros()
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
