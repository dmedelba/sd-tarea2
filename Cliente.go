package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
func mostrar_libros(){
	//presentamos al usuario los libros para que seleccione
	path := "./libros/"
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, val := range lst {
		if val.IsDir() {
			fmt.Printf("[%s]\n", val.Name())
		} else {
			fmt.Println(val.Name())
		}
}

func main(){
	//crear conexion
	//Establecemos conexión con logisitica dist70:6000

	/*
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo establecer la conexión. ERROR: %v", err)
	}
	*/
	defer conn.Close()

	var seleccion int
	flag := true
	for flag{
		log.Printf("[Cliente] Ingrese opición a realizar:")
		log.Printf("1. Subir libro")
		log.Printf("2. Descargar libro")
		log.Printf("3. Ver biblioteca (libros descargados)")
		log.Printf("4. Finalizar")
		fmt.Scanln(&seleccion)
		
		var tipo_subida int
		switch seleccion {
		case 1:
			//subir libro, puede subir centralizado o distribuido
			log.Printf("[Cliente] Seleccione el tipo de subida a realizar:")
			log.Printf("1. Subida Centralizada")
			log.Printf("2. Subida Distribuida")
			fmt.Scanln(&tipo_subida)
			switch tipo_subida {
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
			mostrar_libros()
			//ver biblioteca
		case 4:
			//finalizar
			log.Printf("Sesión finalizada. Muchas gracias!")
			flag = false
		}
	}
}