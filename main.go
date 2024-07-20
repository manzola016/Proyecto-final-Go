
// Manuel Anzola - 32666091

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	driver := "mysql"
	usuario := "root"
	contraseña := ""
	nombre := "sistema"
	conexion, err := sql.Open(driver, usuario+":"+contraseña+"@tcp(127.0.0.1)/"+nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)
	fmt.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

func borrar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idContacto := r.FormValue("id")
		conexionEstablecida := conexionBD()
		borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		_, err = borrarRegistro.Exec(idContacto)
		if err != nil {
			panic(err.Error())
		}
		defer conexionEstablecida.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

type Contacto struct {
	Id     int
	Nombre string
	Correo string
	Numero string
}

func inicio(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}
	contacto := Contacto{}
	arregloContacto := []Contacto{}
	for registros.Next() {
		var id int
		var nombre, correo, numero string
		err = registros.Scan(&id, &nombre, &correo, &numero)
		if err != nil {
			panic(err.Error())
		}
		contacto.Id = id
		contacto.Nombre = nombre
		contacto.Correo = correo
		contacto.Numero = numero
		arregloContacto = append(arregloContacto, contacto)
	}
	plantillas.ExecuteTemplate(w, "inicio", arregloContacto)
}

func editar(w http.ResponseWriter, r *http.Request) {
	idContacto := r.URL.Query().Get("id")
	fmt.Println(idContacto)
	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idContacto)
	contacto := Contacto{}
	for registro.Next() {
		var id int
		var nombre, correo, numero string
		err = registro.Scan(&id, &nombre, &correo, &numero)
		if err != nil {
			panic(err.Error())
		}
		contacto.Id = id
		contacto.Nombre = nombre
		contacto.Correo = correo
		contacto.Numero = numero
	}
	fmt.Println(contacto)
	plantillas.ExecuteTemplate(w, "editar", contacto)
}

func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		numero := r.FormValue("numero")
		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo, numero) VALUES(?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo, numero)
		http.Redirect(w, r, "/", 301)
	}
}

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		numero := r.FormValue("numero")
		conexionEstablecida := conexionBD()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=?, numero=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, correo, numero, id)
		http.Redirect(w, r, "/", 301)
	}
}
