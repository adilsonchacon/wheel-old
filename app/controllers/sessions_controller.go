package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/app/views"
	"wheel.smart26.com/commons/locale"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/commons/mailer"
	"wheel.smart26.com/config"
)

type SessionClaims struct {
	Uid uint   `json:"uid"`
	Jti string `json:"jti"`
	jwt.StandardClaims
}

const (
	privateKeyPath = "config/keys/app.key.rsa"
	publicKeyPath  = "config/keys/app.key.rsa.pub"
)

func SessionSignIn(w http.ResponseWriter, r *http.Request) {
	var user = models.UserAuthenticate(r.FormValue("email"), r.FormValue("password"))
	var errors []string

	log.Info.Println("controllers: SessionSingIn")
	w.Header().Set("Content-Type", "application/json")

	if !models.UserNil(user) {
		json.NewEncoder(w).Encode(views.SessionSignInSuccessMessage("notice", "signed in successfully", sessionGenerateToken(user, r.RemoteAddr)))
	} else {
		errors = append(errors, "invalid credentials")
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "could not sign in", errors))
	}
}

func SessionSignOut(w http.ResponseWriter, r *http.Request) {
	var session models.Session
	var errors []string

	authToken, _ := sessionAuthToken(r.Header.Get("token"))
	claims, ok := authToken.Claims.(*SessionClaims)
	session = models.SessionFindByToken(claims.Jti)

	if ok && authToken.Valid && models.SessionDeactivate(&session) {
		json.NewEncoder(w).Encode(views.SessionSignOutSuccessMessage("notice", "signed out successfully"))
	} else {
		errors = append(errors, "invalid token")
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "access denied", errors))
	}
}

func SessionSignUp(w http.ResponseWriter, r *http.Request) {
	var user = models.User{Name: r.FormValue("name"), Email: r.FormValue("email"), Password: r.FormValue("password"), Locale: r.FormValue("locale"), Admin: false}

	log.Info.Println("controllers: SignUp")
	w.Header().Set("Content-Type", "application/json")

	if models.UserSave(&user) {
		locale.Load(user.Locale)

		mailer.AddTo(user.Name, user.Email)
		subject := locale.Welcome() + " " + models.UserFirstName(&user)
		body := views.SessionSignUpMailer(&user)
		go mailer.Send(subject, body, true)

		json.NewEncoder(w).Encode(views.SessionSignUpSuccessMessage("notice", "user was successfully created", sessionGenerateToken(user, r.RemoteAddr)))
	} else {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "user was not created", models.Errors))
	}
}

func SessionPassword(w http.ResponseWriter, r *http.Request) {
	var user = models.UserFindByEmail(r.FormValue("email"))

	log.Info.Println("controllers: SessionPassword")

	if models.UserExists(user) {
		locale.Load(user.Locale)

		token := models.UserSetRecovery(&user)
		mailer.AddTo(user.Name, user.Email)
		subject := locale.PasswordRecoveryInstructions()
		body := views.SessionPasswordRecoveryInstructionsMailer(&user, token)
		go mailer.Send(subject, body, true)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(views.SetSystemMessage("notice", "user password recovery instructions was successfully sent"))
}

func SessionRecovery(w http.ResponseWriter, r *http.Request) {
	var errors []string
	user := models.UserFindByResetPasswordToken(r.Header.Get("token"))

	log.Info.Println("controllers: SessionRecovery")

	w.Header().Set("Content-Type", "application/json")

	if !models.UserExists(user) {
		errors = append(errors, "invalid reset password token")
	} else if r.FormValue("new_password") != r.FormValue("password_confirmation") {
		errors = append(errors, "password confirmation does not match new password")
	} else if user.Password = r.FormValue("new_password"); models.UserSave(&user) {
		models.UserClearRecovery(&user)
		json.NewEncoder(w).Encode(views.SessionSignInSuccessMessage("notice", "password was successfully changed", sessionGenerateToken(user, r.RemoteAddr)))
	} else {
		errors = models.Errors
	}

	if len(errors) > 0 {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "password could not be changed", errors))
	}
}

func SessionRefresh(w http.ResponseWriter, r *http.Request) {
	authToken, _ := sessionAuthToken(r.Header.Get("token"))
	claims, _ := authToken.Claims.(*SessionClaims)
	session := models.SessionFindByToken(claims.Jti)
	models.SessionDeactivate(&session)

	log.Info.Println("controllers: SessionRefresh")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(views.SessionRefreshSuccessMessage("notice", "session was successfully refreshed", sessionGenerateToken(models.UserCurrent, r.RemoteAddr)))
}

func SessionCheck(token string) (uint, error) {
	authToken, err := sessionAuthToken(token)

	if err != nil {
		return 0, errors.New("invalid token")
	}

	claims, ok := authToken.Claims.(*SessionClaims)

	if ok && authToken.Valid && sessionUpdateStats(claims.Jti) {
		return claims.Uid, nil
	} else {
		return 0, errors.New("invalid token")
	}
}

func sessionBuildClaims(jti string, user *models.User) (jwt.MapClaims, time.Time) {
	expiresAt := time.Now().Add(time.Second * time.Duration(config.TokenExpirationSeconds()))

	claims := make(jwt.MapClaims)
	claims["iss"] = config.AppName()
	claims["exp"] = expiresAt.Unix()
	claims["jti"] = jti
	claims["uid"] = user.ID

	return claims, expiresAt
}

func sessionUpdateStats(jti string) bool {
	session := models.SessionFindByToken(jti)

	if !session.Active {
		return false
	} else {
		models.SessionIncrementStats(&session)
		return true
	}
}

func sessionAuthToken(token string) (*jwt.Token, error) {
	var publicBytes []byte
	var publicKey *rsa.PublicKey
	var err_read_file error
	var err_parse_rsa error
	var err error
	var authToken *jwt.Token

	if token == "" {
		return authToken, errors.New("invalid token")
	}

	publicBytes, err_read_file = ioutil.ReadFile(publicKeyPath)
	if err_read_file != nil {
		log.Fatal.Println(err_read_file)
	}

	publicKey, err_parse_rsa = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err_parse_rsa != nil {
		log.Fatal.Println(err_parse_rsa)
	}

	authToken, err = jwt.ParseWithClaims(token, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	return authToken, err
}

func sessionGenerateToken(user models.User, remoteAddr string) string {
	var privateBytes []byte
	var privateKey *rsa.PrivateKey
	var err error
	var expiresAt time.Time

	privateBytes, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal.Println(err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal.Println(err)
	}

	jti := uuid.Must(uuid.NewV4()).String()

	signer := jwt.New(jwt.SigningMethodRS256)
	signer.Claims, expiresAt = sessionBuildClaims(jti, &user)

	token, err := signer.SignedString(privateKey)
	if err != nil {
		log.Error.Println(err)
	} else {
		log.Info.Println("Token was succesfully created for user " + user.Email)
	}

	t := time.Now()
	ip, _, _ := net.SplitHostPort(remoteAddr)
	var session = models.Session{Jti: jti, App: "Default", Requests: 0, LastRequestAt: &t, UserID: user.ID, Address: ip, ExpiresIn: config.TokenExpirationSeconds(), ExpiresAt: expiresAt}
	models.SessionSave(&session)

	return token
}
