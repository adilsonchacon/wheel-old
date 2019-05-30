package handler

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
	"wheel.smart26.com/app/session"
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/app/view"
	"wheel.smart26.com/commons/locale"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/commons/mailer"
	"wheel.smart26.com/config"
	"wheel.smart26.com/db/entity"
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
	log.Info.Println("Handler: SessionSingIn")
	w.Header().Set("Content-Type", "application/json")

	userAuth, err := user.Authenticate(r.FormValue("email"), r.FormValue("password"))

	if !user.IsNil(&userAuth) {
		json.NewEncoder(w).Encode(session.SignInSuccessMessage("notice", "signed in successfully", sessionGenerateToken(userAuth, r.RemoteAddr)))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not sign in", []error{err}))
	}
}

func SessionSignOut(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: SessionSignOut")
	w.Header().Set("Content-Type", "application/json")

	authToken, _ := sessionAuthToken(r.Header.Get("token"))

	claims, ok := authToken.Claims.(*SessionClaims)
	if !ok || !authToken.Valid {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	sessionSingOut, errorFindByJti := session.FindByJti(claims.Jti)
	if errorFindByJti != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	if deactivated, _ := session.Deactivate(&sessionSingOut); deactivated {
		json.NewEncoder(w).Encode(session.SignOutSuccessMessage("notice", "signed out successfully"))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
	}
}

func SessionRefresh(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: SessionRefresh")
	w.Header().Set("Content-Type", "application/json")

	authToken, _ := sessionAuthToken(r.Header.Get("token"))

	claims, ok := authToken.Claims.(*SessionClaims)
	if !ok || !authToken.Valid {
		log.Error.Println("invalid token")
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	currentSession, errorFindByJti := session.FindByJti(claims.Jti)
	if errorFindByJti != nil {
		log.Error.Printf("could not find session by token %s", claims.Jti)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	valid, _ := session.Deactivate(&currentSession)
	if !valid {
		log.Error.Println("could not deactivate session")
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	userSession, errorUserNotFound := user.Find(currentSession.UserID)
	if errorUserNotFound != nil {
		log.Error.Printf("could not find user by %d\n", currentSession.UserID)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	json.NewEncoder(w).Encode(session.RefreshSuccessMessage("notice", "session was successfully refreshed", sessionGenerateToken(userSession, r.RemoteAddr)))
}

func SessionSignUp(w http.ResponseWriter, r *http.Request) {
	var newUser = entity.User{}

	log.Info.Println("Handler: SessionSignUp")
	w.Header().Set("Content-Type", "application/json")

	userSetParams(&newUser, r)
	newUser.Admin = false

	if valid, errs := user.Save(&newUser); valid {
		locale.Load(newUser.Locale)

		mailer.AddTo(newUser.Name, newUser.Email)
		subject := locale.Welcome() + " " + user.FirstName(&newUser)
		body := session.SignUpMailer(&newUser)
		go mailer.Send(subject, body, true)

		json.NewEncoder(w).Encode(session.SignUpSuccessMessage("notice", "user was successfully created", sessionGenerateToken(newUser, r.RemoteAddr)))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not created", errs))
	}
}

func SessionPassword(w http.ResponseWriter, r *http.Request) {
	var currentUser, _ = user.FindByEmail(r.FormValue("email"))

	log.Info.Println("Handler: SessionPassword")
	w.Header().Set("Content-Type", "application/json")

	if user.Exists(&currentUser) {
		locale.Load(currentUser.Locale)

		token, _ := user.SetRecovery(&currentUser)
		mailer.AddTo(currentUser.Name, currentUser.Email)
		subject := locale.PasswordRecoveryInstructions()
		body := session.PasswordRecoveryInstructionsMailer(&currentUser, token)
		go mailer.Send(subject, body, true)
	}

	json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user password recovery instructions was successfully sent"))
}

func SessionRecovery(w http.ResponseWriter, r *http.Request) {
	var errs []error
	var valid bool

	log.Info.Println("Handler: SessionRecovery")
	w.Header().Set("Content-Type", "application/json")

	currentUser, _ := user.FindByResetPasswordToken(r.Header.Get("token"))
	currentUser.Password = r.FormValue("new_password")

	if !user.Exists(&currentUser) {
		errs = append(errs, errors.New("invalid reset password token"))
	} else if r.FormValue("new_password") != r.FormValue("password_confirmation") {
		errs = append(errs, errors.New("password confirmation does not match new password"))
	} else if valid, errs = user.Save(&currentUser); valid {
		user.ClearRecovery(&currentUser)
		json.NewEncoder(w).Encode(session.SignInSuccessMessage("notice", "password was successfully changed", sessionGenerateToken(currentUser, r.RemoteAddr)))
	}

	if len(errs) > 0 {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "password could not be changed", errs))
	}
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

func sessionBuildClaims(jti string, userSession *entity.User) (jwt.MapClaims, time.Time) {
	expiresAt := time.Now().Add(time.Second * time.Duration(config.TokenExpirationSeconds()))

	claims := make(jwt.MapClaims)
	claims["iss"] = config.AppName()
	claims["exp"] = expiresAt.Unix()
	claims["jti"] = jti
	claims["uid"] = userSession.ID

	return claims, expiresAt
}

func sessionUpdateStats(jti string) bool {
	currentSession, _ := session.FindByJti(jti)

	if !currentSession.Active {
		return false
	} else {
		session.IncrementStats(&currentSession)
		return true
	}
}

func sessionAuthToken(token string) (*jwt.Token, error) {
	var publicBytes []byte
	var publicKey *rsa.PublicKey
	var errorReadFile error
	var errorParseRsa error
	var err error
	var authToken *jwt.Token

	if token == "" {
		return authToken, errors.New("invalid token")
	}

	publicBytes, errorReadFile = ioutil.ReadFile(publicKeyPath)
	if errorReadFile != nil {
		log.Error.Println(errorReadFile)
		return authToken, errorReadFile
	}

	publicKey, errorParseRsa = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if errorParseRsa != nil {
		log.Error.Println(errorParseRsa)
		return authToken, errorParseRsa
	}

	authToken, err = jwt.ParseWithClaims(token, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	return authToken, err
}

func sessionGenerateToken(userSession entity.User, remoteAddr string) string {
	var privateBytes []byte
	var privateKey *rsa.PrivateKey
	var err error
	var expiresAt time.Time
	var sessionNew entity.Session

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
	signer.Claims, expiresAt = sessionBuildClaims(jti, &userSession)

	token, err := signer.SignedString(privateKey)
	if err != nil {
		log.Error.Println(err)
	} else {
		log.Info.Println("Token was successfully created for user " + userSession.Email)
	}

	t := time.Now()
	ip, _, _ := net.SplitHostPort(remoteAddr)
	sessionNew = entity.Session{Jti: jti, App: "Default", Requests: 0, LastRequestAt: &t, UserID: userSession.ID, Address: ip, ExpiresIn: config.TokenExpirationSeconds(), ExpiresAt: expiresAt}
	session.Save(&sessionNew)

	return token
}