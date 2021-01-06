package controller

import (
	"cnxbd"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//Register função para registrar usuario
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		firstName := r.FormValue("pname")
		lastName := r.FormValue("sname")
		password := r.FormValue("password")
    //Conexão com o Banco de dados MYSQL
		db := cnxbd.Conectabd()
    
    //Função da Biblioteca bcrypt para cryptografar a senha criada pelo usuario no form, a senha é cryptografada com sha256 e salt.
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err.Error())
		}
		hashedPasswordString := string(hashedPasswordBytes)

		inserir, err := db.Query("INSERT INTO user(username, email, firstName, lastName, password) VALUES(?,?,?,?,?)", username, email, firstName, lastName, hashedPasswordString)
		if err != nil {
			panic(err.Error())
		}
		defer inserir.Close()

		http.Redirect(w, r, "/", 301)
	} else {
		fmt.Println("Erro para registrar o usuario")
	}

}
