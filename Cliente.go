package main

import (
	"log"
	"fmt"
)

func main(){
	//crear conexion
	var seleccion int
	for true{
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
				//centralizada
			case 2:
				//distribuida	
			}
		case 2:
			//descargar libro
		case 3:
			//ver biblioteca
		case 4:
			//finalizar
			break
		}
	}
}