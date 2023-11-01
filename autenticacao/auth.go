package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Recuperando os valores de login e email do formulário HTML
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")

	// Conectando ao MongoDB
	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Acessando a coleção "newsletter"
	collection := client.Database("webform").Collection("newsletter")

	// Criando um filtro com base no login e email fornecidos
	filter := bson.M{"name": name, "email": email}

	// Executando uma consulta para verificar se os dados existem na coleção
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Dados de login e email não encontrados
		fmt.Fprintln(w, "Falha no login. Verifique suas credenciais.")
		return
	} else if err != nil {
		log.Fatal(err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Dados de login e email encontrados
	fmt.Fprintln(w, "Logado com sucesso!")
}

// Roteamento para a página do formulário
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "login.html")
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/form", formHandler) // Rota para a página do formulário

	log.Println("Servidor em execução na porta 4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}
