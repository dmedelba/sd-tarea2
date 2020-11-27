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
func subirLibroCentralizado(conn *grpc.ClientConn){
	//buscamos libro, se selecciona y se descompone

}

func subirLibroDistribuido(conn *grpc.ClientConn){
	//buscamos libro, se selecciona y se descompone
}
*/
func mostrarLibros() {
	//presentamos al usuario los libros para que seleccione
	path := "./libros/"
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	i := 1
	for _, val := range lst {
		if val.IsDir() {
			fmt.Printf("[%s]\n", val.Name())
		} else {
			s := strconv.Itoa(i)
			fmt.Println(s + ". " + val.Name())
			i++
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
		fmt.Printf("[Cliente] Ingrese opición a realizar:")
		fmt.Printf("1. Subir libro")
		fmt.Printf("2. Descargar libro")
		fmt.Printf("3. Ver biblioteca (libros descargados)")
		fmt.Printf("4. Finalizar")
		fmt.Scanln(&seleccion)

		var tipoSubida int
		switch seleccion {
		case 1:
			//subir libro, puede subir centralizado o distribuido
			fmt.Printf("[Cliente] Seleccione el tipo de subida a realizar:")
			fmt.Printf("1. Subida Centralizada")
			fmt.Printf("2. Subida Distribuida")
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
