package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/ssarkar/taskMamagement/middleware"
	"github.com/ssarkar/taskMamagement/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

func (server *Server)UserRegisterHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Trying to register user.........")
	inputUserData :=&models.UserDetails{}
	if err:=json.NewDecoder(r.Body).Decode(inputUserData);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	logrus.Info("input for user registration ",inputUserData)
	if err:=ValidateInputStructs(inputUserData);err!=nil{
		logrus.Printf("Unable to validate user input. Username: %s, Password: %s. Provide valid user data.",inputUserData.UserName,inputUserData.Password)
		http.Error(w,err.Error(),http.StatusBadGateway)
		return
	}
	logrus.Info("User details are validated successfully")
	_, err :=server.db.GetUserByUsername(inputUserData.UserName)
	if errors.Is(err,gorm.ErrRecordNotFound){
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputUserData.Password), bcrypt.DefaultCost)
    if err != nil {
         http.Error(w,err.Error(),http.StatusInternalServerError)
				 return
    }
		if newErr:=server.db.UserRegister(inputUserData.UserName,string(hashedPassword));newErr!=nil{
			logrus.Error("Unable to register user")
			http.Error(w,newErr.Error(),http.StatusInternalServerError)
			return
		}
	}else if err!=nil{
		logrus.Error("unable to get user details from the database")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return		
	}else{
		errStr :=fmt.Sprintf("user [%s] is already registered",inputUserData.UserName)
		http.Error(w,errors.New(errStr).Error(),http.StatusInternalServerError)
		return
	}
	responseData :=fmt.Sprintf("User [%s] is registered successfully.",inputUserData.UserName)
	w.Write([]byte(responseData))
}



func ValidateInputStructs(data interface{})error{
	validate :=validator.New()
	return validate.Struct(data)
}

func (server *Server)LoginHandler(w http.ResponseWriter, r *http.Request){
	loginRequest :=&models.UserDetails{}
	if err:=json.NewDecoder(r.Body).Decode(loginRequest);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if err:=ValidateInputStructs(loginRequest);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	user, err:=server.db.GetUserByUsername(loginRequest.UserName)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(loginRequest.Password));err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return		
	}
	token, err:=middleware.GenerateToken(user.Username)
	if err!=nil{
		logrus.Info("Token Generation error.......")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	resonse :=map[string]string{
		"token":token,
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resonse)
}