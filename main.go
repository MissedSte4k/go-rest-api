package main

//go build github.com/MissedSte4k/go-rest-api/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type user struct {
	Id       int    `json:"Id"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	UserType int    `json:"UserType"`
}

type Credentials struct {
	Password string `json:"Password"`
	Username string `json:"Username"`
}

type Claims struct {
	Username string `json:"username"`
	Type     int    `json:"usertype"`
	jwt.StandardClaims
}

type ingredient struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type dishingredient struct {
	Id             int    `json:"Id"`
	Amount         string `json:"Amount"`
	fkDishId       int    `json:"fkDishId"`
	fkIngredientId int    `json:"fkIngredientId"`
}
type comment struct {
	Id       string `json:"Id"`
	Date     string `json:"Date"`
	Text     string `json:"Text"`
	fkDishId int    `json:"fkDishId"`
	fkUserId int    `json:"fkUserId"`
}

type dish struct {
	Id          int     `json:"Id"`
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	CookingTime int     `json:"CookingTime"`
	RecipeMakes int     `json:"RecipeMakes"`
	Calories    float64 `json:"Calories"`
	Carbs       float64 `json:"Carbs"`
	Fat         float64 `json:"Fat"`
	Protein     float64 `json:"Protein"`
	Directions  string  `json:"Directions"`
	fkUserId    int     `json:"fkUserId"`
}
type Data struct {
	dishingredient
	ingredient
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Print("Succesfully conected")
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/users", createuser).Methods("POST")
	router.HandleFunc("/users", getAllusers).Methods("GET")
	router.HandleFunc("/users/{id}", getOneuser).Methods("GET")
	router.HandleFunc("/users/{id}", updateuser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteuser).Methods("DELETE")

	router.HandleFunc("/ingredients", createingredient).Methods("POST")
	router.HandleFunc("/ingredients", getAllingredients).Methods("GET")
	router.HandleFunc("/ingredients/{id}", getOneingredient).Methods("GET")
	router.HandleFunc("/ingredients/{id}", updateingredient).Methods("PATCH")
	router.HandleFunc("/ingredients/{id}", deleteingredient).Methods("DELETE")

	router.HandleFunc("/comments", createcomment).Methods("POST")
	router.HandleFunc("/comments", getAllcomments).Methods("GET")
	router.HandleFunc("/comments/{id}", getOnecomment).Methods("GET")
	router.HandleFunc("/comments/{id}", updatecomment).Methods("PATCH")
	router.HandleFunc("/comments/{id}", deletecomment).Methods("DELETE")

	router.HandleFunc("/dishingredients", createdishingredient).Methods("POST")
	router.HandleFunc("/dishes/{id}/ingredients", getAlldishingredients).Methods("GET")
	router.HandleFunc("/dishingredients/{id}", getOnedishingredient).Methods("GET")
	router.HandleFunc("/dishingredients/{id}", updatedishingredient).Methods("PATCH")
	router.HandleFunc("/dishingredients/{id}", deletedishingredient).Methods("DELETE")

	router.HandleFunc("/dishes", createdish).Methods("POST")
	router.HandleFunc("/dishes", getAlldishes).Methods("GET")
	router.HandleFunc("/dishes/{id}", getOnedish).Methods("GET")
	router.HandleFunc("/dishes/{id}", updatedish).Methods("PATCH")
	router.HandleFunc("/dishes/{id}", deletedish).Methods("DELETE")

	router.HandleFunc("/signin", Signin)
	router.HandleFunc("/welcome", Welcome)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

//Signin makes the signin
func Signin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	var creds Credentials
	var went = false
	// Get the JSON body and decode into credentials

	//err = json.NewDecoder(r.Body).Decode(&creds)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		fmt.Println("something went wrong")
		return
	}

	// Get the expected password from our in memory map
	str := fmt.Sprintf("SELECT `Password`, `UserType`, `Id` FROM users WHERE Username = '%v'", creds.Username)
	fmt.Print(str)
	expectedPassword, ok := db.Query(str)
	if ok != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	for expectedPassword.Next() {
		var user user
		err = expectedPassword.Scan(&user.Password, &user.UserType, &user.Id)
		went = true
		// If a password exists for the given user
		// AND, if it is the same as the password we received, the we can move ahead
		// if NOT, then we return an "Unauthorized" status
		if user.Password != creds.Password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: creds.Username,
			Type:     user.UserType,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
				Id:        strconv.Itoa(user.Id),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
	if !went {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func UserType(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return 0
		}
		w.WriteHeader(http.StatusBadRequest)
		return -1
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return 0
		}
		w.WriteHeader(http.StatusBadRequest)
		return -1
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return 0
	}

	return claims.Type

}

func TokenIdGet(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return "-1"
		}
		w.WriteHeader(http.StatusBadRequest)
		return "-1"
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return "-1"
		}
		w.WriteHeader(http.StatusBadRequest)
		return "-1"
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return "-1"
	}
	return claims.Id

}

func Refresh(w http.ResponseWriter, r *http.Request) bool {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 120 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 120*time.Second {
		//w.WriteHeader(http.StatusBadRequest)
		return true
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return true
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
