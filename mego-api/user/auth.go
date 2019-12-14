package user

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"log"
	"net/http"
	"time"
)

// hence userDB itself is memory-based,
// then no value of making the jwtSecret a const value that span multiple server restarts
var jwtSecret []byte

func init() {
	jwtSecret = make([]byte, 32)
	_, err := rand.Read(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}
}

var usersDB = make(map[string]string) //username->password

type token string

func login(u *user) (token, bool) {

	ewsClient := commons.NewEWSClient(u.Username, u.Password)
	_, err := ewsutil.GetUserPhotoBase64(ewsClient, u.Username) // check ews auth

	herr, ok := err.(*ews.HTTPError)
	if ok && herr.StatusCode == http.StatusUnauthorized {
		return "", false
	}

	t, err := createToken(u.Username)
	if err != nil {
		log.Println(err)
		return "", false
	}

	// insert user into our database
	enc, err := encrypt(u.Password)
	if err != nil {
		log.Println(err)
		return "", false
	}
	usersDB[u.Username] = enc

	return t, true
}

func GetUser(t string) (*user, error) {

	if len(t) == 0 {
		return nil, errors.New("missing token")
	}

	username, err := getUsernameFromToken(token(t))
	if err != nil {
		return nil, err
	}

	p, found := usersDB[username]
	if !found {
		return nil, errors.New("username not found")
	}
	dec, err := decrypt(p)
	if err != nil {
		return nil, err
	}

	return &user{
		Username: username,
		Password: dec,
	}, nil
}

func logout(u *user) {
	delete(usersDB, u.Username)
}

func createToken(username string) (token, error) {
	newJwt := jwt.New(jwt.SigningMethodHS256)
	claims := newJwt.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24 * 10).Unix() // 10 days
	t, err := newJwt.SignedString([]byte(jwtSecret))
	return token(t), err
}

func getUsernameFromToken(t token) (string, error) {
	tt, err := jwt.Parse(string(t), func(tt *jwt.Token) (i interface{}, e error) {
		if _, ok := tt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tt.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := tt.Claims.(jwt.MapClaims)
	if !ok || !tt.Valid {
		return "", errors.New("invalid token")
	}

	u := claims["username"].(string)
	return u, nil
}
