# Laboratorio 2 - Sistemas Distribuidos

## Integrantes : 
- David Medel B. 201573548-4
- Macarana Hidalgo A. 201473608-8

Se definen las siguientes maquinas:

	 dist69--> NameNode 
	 dist70--> DataNode
	 dist71--> DataNode
	 dist72--> DataNode

Para ejecutar los nodos correspondiente se debe correr en las maquinas indicadas 
	 make namenode -> ejecutar en la dist69.
	 make datanode -> ejecutar en dist70/71/72.
	 make cliente -> Para ejecutar al cliente en cualquiera de las maquinas.
		
## Consideraciones : 
- Los puertos de las maquinas se definió como 5050, todos deben estar en el mismo puerto. En caso de fallar algun puerto por problemas con las maquinas, establecer todos un numero igual.
- Si se quiere correr desde 0, se pide porfavor eliminar el contenido de LOG.txt. En la entrega viene vacio.
- No se implementó un limpiado de archivos, por lo que si se sube 2 veces el mismo libro, se registrará 2 veces, sin embargo, cuando se solicita la ubicación de estos al namenode, solo entrega la primera que encuentra en el LOG.txt.
- Se supone que el usuario no va a ingresar otra opción por consola. Por ejemplo, si se le pide que ingrese una opción  0 o 1, no ingresará la opción 2.
- Más detalles en el informe.
- Implementado algoritmo centralizado 100%, distribuido no.