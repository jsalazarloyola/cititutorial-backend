package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Verifica que existan variables necesarias para este programa
func checkVars() []string {
	vars := []string{
		"GO_REST_ENV",
		"GIN_MODE",
		"JWT_KEY",
		"API_AUTH_URL",
		"API_AUTH_USER",
		"API_AUTH_PASS",
	}
	missing := []string{}
	for _, v := range vars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = append(missing, v)
		}
	}
	return missing
}


func LoadEnv(){
	env := os.Getenv("TASK_ENV")
	if env == "" { env = "development" }

	// Este módulo es "equivalente", al final, a hacer
	// $ VAR=value go run app
	// con múltiples pares VAR=value
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
	if vars := checkVars(); len(vars) != 0 {
		log.Printf("ERROR: Variables de entorno necesarias no definidas: %v", vars)
		panic(fmt.Sprintf("ERROR: Variables de entorno necesarias no definidas: %v", vars))
	}
}
