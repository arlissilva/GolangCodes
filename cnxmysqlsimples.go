package cnxbd

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//Conectabd função para conexão com banco de dados
func Conectabd() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "user"
	dbPass := "password"
	dbName := "name Database
  
  //Abre a conexão com o Banco de dados
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(IP:PORT)/"+dbName)
	if err != nil {
		fmt.Println("Erro no banco", err)
	}

	return db
}
