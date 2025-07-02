package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"go-template/models"
	"go-template/utils"

	"github.com/gin-gonic/gin"
)

type LoginService struct {
	ServiceURL      string // Dirección del servicio de autenticación
	APIAuthUser     string // Usuario de la API
	APIAuthPassword string // Password de la API
}

// Constructor (?)
func NewLoginService() *LoginService {
	var service LoginService
	service.ServiceURL = os.Getenv("API_AUTH_URL")
	service.APIAuthUser = os.Getenv("API_AUTH_USER")
	service.APIAuthPassword = os.Getenv("API_AUTH_PASS")

	return &service
}

// Autenticación básica en la API
func (service LoginService) GetBasicAuth() string {
	return ("Basic " +
		base64.StdEncoding.EncodeToString([]byte(service.APIAuthUser+
			":"+
			service.APIAuthPassword)))
}

// Función que realiza solicitud de login
func (service LoginService) RequestLogin(username string, password string) (*http.Response, error) {
	// Hashea password para la conexión
	hashPassword := utils.HashPassword(password)

	// Crea  la estructura para luego jsonizarla
	body := models.Login{
		User:     username,
		Password: hashPassword,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// El requerimiento: postea el login
	request, err := http.NewRequest(
		"POST", service.ServiceURL, bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	// Headers para la solicitud
	request.Header.Set("Authorization", service.GetBasicAuth())
	request.Header.Set("Content-Type", "application/json")

	if gin.Mode() == "debug" {
		req, _ := httputil.DumpRequestOut(request, true)
		log.Println(string(req))
	}

	// El requerimiento en sí
	client := &http.Client{}
	response, err := client.Do(request)

	// log.Println(response)
	if err != nil {
		return nil, err
	}

	log.Println(response.StatusCode)
	return response, nil
}

func (service LoginService) DoRequestLogin(username string, password string) (*models.ResponseLogin, error) {
	// uniforma y limpia
	username = strings.ToLower(username)
	username = strings.ReplaceAll(username, "@usach.cl", "")
	// log.Println(username, password)
	// Hace el requerimiento de login a la API de cuentas
	response, err := service.RequestLogin(username, password)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Mensajes de error posibles
	if response.StatusCode != 200 {
		var errorMessage map[string]any
		err = json.Unmarshal(responseBody, &errorMessage)
		if err != nil {
			return nil, errors.New("error al recuperar el error")
		}

		// log.Println(response.StatusCode, errorMessage)

		// El mensaje hay que castearlo, al parecer
		return nil, errors.New(errorMessage["message"].(string))
	}
	// log.Println(string(responseBody), response.StatusCode)

	var responseToLogin models.ResponseLogin
	err = json.Unmarshal(responseBody, &responseToLogin)
	if err != nil {
		return nil, errors.New("error al obtener la respuesta: " + err.Error())
	}

	log.Println(responseToLogin)

	return &responseToLogin, nil
}
