package controller

import (
	"cnxbd"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Attachment struct {
	NameAttachment string
	BodyAttachment []byte
	ExtAttachment  string
}

var (
	fileName    string
	fullURLFile string
)

func ShowAttachmentId(w http.ResponseWriter, r *http.Request) {

  //Recebe do form o id do arquivo que esta salvo no banco
	idURLAttach := r.URL.Query().Get("id")
	db := cnxbd.Conectabd()

	//Localiza o arquivo no banco pelo id
  attachment, err := db.Query("SELECT * FROM Evidencias where idEvidencias = ?", idURLAttach)
	if err != nil {
		fmt.Println("Erro para buscar evidencia no Banco", err)
	}

	for attachment.Next() {
		var idAttach int
		var nameAttach string
		var bodyAttach []byte
		var extAttach string
		var idPlanTest int

		err := attachment.Scan(&idAttach, &bodyAttach, &nameAttach, &extAttach, &idPlanTest)
		if err != nil {
			fmt.Println("Erro no Scan da Evidencia", err)
		}
    
    //Cria o arquivo temporario no servidor, dentro da pasta "attach" com nome e extenção extraido do banco
		MountFile, err := ioutil.TempFile("./attach/", "*"+nameAttach+extAttach)
		if err != nil {
			fmt.Println("Erro para criar o arquivo temporario: ", err)

		}

		//Adia o fechamento do arquivo
    defer os.Remove(MountFile.Name())
    
    //Escreve o conteudo extraido do banco de dados dentro do arquivo temporario.
		err = ioutil.WriteFile(MountFile.Name(), bodyAttach, 0644)
		if err != nil {
			fmt.Println("Deu erro para inserir conteudo no arquivo", err)
		}
		urlstring := "http://localhost/" + MountFile.Name()

		
		timeout := time.Duration(5) * time.Second
		transport := &http.Transport{
			ResponseHeaderTimeout: timeout,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, timeout)
			},
			DisableKeepAlives: true,
		}
		client := &http.Client{
			Transport: transport,
		}

		resp, err := client.Get(urlstring)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+nameAttach+extAttach)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
    
    //realiza o download para o cliente.
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Println("Retorno do Copy: ", err)
		}

	}

}
