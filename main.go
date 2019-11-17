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
)

type user struct {
	Id       int    `json:"Id"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	UserType int    `json:"UserType"`
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
	Id          int    `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	CookingTime int    `json:"CookingTime"`
	RecipeMakes int    `json:"RecipeMakes"`
	Calories    int    `json:"Calories"`
	Carbs       int    `json:"Carbs"`
	Fat         int    `json:"Fat"`
	Protein     int    `json:"Protein"`
	Directions  string `json:"Directions"`
	fkUserId    int    `json:"fkUserId"`
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
	router.HandleFunc("/dishingredients", getAlldishingredients).Methods("GET")
	router.HandleFunc("/dishingredients/{id}", getOnedishingredient).Methods("GET")
	router.HandleFunc("/dishingredients/{id}", updatedishingredient).Methods("PATCH")
	router.HandleFunc("/dishingredients/{id}", deletedishingredient).Methods("DELETE")

	router.HandleFunc("/dishes", createdish).Methods("POST")
	router.HandleFunc("/dishes", getAlldishes).Methods("GET")
	router.HandleFunc("/dishes/{id}", getOnedish).Methods("GET")
	router.HandleFunc("/dishes/{id}", updatedish).Methods("PATCH")
	router.HandleFunc("/dishes/{id}", deletedish).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createcomment(w http.ResponseWriter, r *http.Request) {
	var comments comment
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &comments)
	dt := time.Now()
	str := fmt.Sprintf(dt.Format("2006-01-02 15:04:05"))
	comments.Date = str
	comments.Text = r.FormValue("Text")
	u2, err := strconv.Atoi(r.FormValue("fkDishId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	comments.fkDishId = u2
	u3, err := strconv.Atoi(r.FormValue("fkUserId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	comments.fkUserId = u3
	if len(comments.Text) < 1 {
		fmt.Fprintf(w, "Text can't be less than 1 character long")
	} else if comments.fkDishId == 0 {
		fmt.Fprintf(w, "Wrong Dish ID")
	} else if comments.fkUserId == 0 {
		fmt.Fprintf(w, "Wrong Dish ID")
	} else {
		str = fmt.Sprintf("INSERT INTO `myfridgefood`.`comments` (`Date`, `Text`, `fk_DishId`, `fk_UserId`) VALUES ( '%v' , '%v' , '%v', '%v' );", comments.Date, comments.Text, comments.fkDishId, comments.fkUserId)
		fmt.Println(str)
		res, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Uploaded Sucessfully")
			fmt.Println(res)
		}
	}

}

func createingredient(w http.ResponseWriter, r *http.Request) {
	var ingredient ingredient
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &ingredient)

	ingredient.Name = r.FormValue("Name")
	if len(ingredient.Name) < 1 {
		fmt.Fprintf(w, "Please enter a name of the product")
	} else {
		str := fmt.Sprintf("INSERT INTO `myfridgefood`.`ingredients` (`Name`) VALUES ( '%v' )", ingredient.Name)
		fmt.Println(str)
		res, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Uploaded Sucessfully")
			fmt.Println(res)
		}
	}

}

func createuser(w http.ResponseWriter, r *http.Request) {
	var user user
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &user)

	user.Email = r.FormValue("Email")
	user.Password = r.FormValue("Password")

	u2, err := strconv.Atoi(r.FormValue("UserType"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	user.UserType = u2
	if len(user.Email) < 1 {
		fmt.Fprintf(w, "please enter an user name")
	} else if len(user.Password) < 6 {
		fmt.Fprintf(w, "your password canot be less than, 6 characters long")
	} else {
		str := fmt.Sprintf("INSERT INTO `myfridgefood`.`users` (`Email`, `Password`, `UserType`) VALUES ( '%v' , '%v' , '%v' );", user.Email, user.Password, user.UserType)
		fmt.Println(str)
		res, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Uploaded Sucessfully")
			fmt.Println(res)
		}
	}
}

func createdish(w http.ResponseWriter, r *http.Request) {
	var dish dish
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &dish)

	dish.Name = r.FormValue("Name")
	dish.Description = r.FormValue("Description")

	u2, err := strconv.Atoi(r.FormValue("CookingTime"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.CookingTime = u2
	u2, err = strconv.Atoi(r.FormValue("RecipeMakes"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.RecipeMakes = u2

	dish.Calories, err = strconv.Atoi(r.FormValue("Calories"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Carbs, err = strconv.Atoi(r.FormValue("Carbs"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Fat, err = strconv.Atoi(r.FormValue("Fat"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Protein, err = strconv.Atoi(r.FormValue("Protein"))
	if err != nil {
		fmt.Println("something went wrong")
	}

	dish.Directions = r.FormValue("Directions")

	u2, err = strconv.Atoi(r.FormValue("fkUserId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.fkUserId = u2
	if len(dish.Name) < 1 {
		fmt.Println("enter dish name")
	} else if len(dish.Description) < 1 {
		fmt.Println("enter dish description")
	} else if dish.CookingTime < 1 {
		fmt.Println("enter dish cookingTime")
	} else if dish.RecipeMakes < 1 {
		fmt.Println("enter dish recipeMakes")
	} else if len(dish.Directions) < 1 {
		fmt.Println("enter dish Directions")
	} else if dish.Calories < 1 {
		fmt.Println("enter dish calories cannot be negative")
	} else if dish.fkUserId == 0 {
		fmt.Println("user not assigned")
	} else {

		str := fmt.Sprintf("INSERT INTO `myfridgefood`.`dishes` (`Name`, `Description`, `CookingTime`, `RecipeMakes`, `Calories`, `Carbs`, `Fat`, `Protein`, `Directions`, `fk_UserId`) VALUES ( '%v' , '%v' , '%v', '%v' , '%v' , '%v', '%v' , '%v' , '%v' , '%v');", dish.Name, dish.Description, dish.CookingTime, dish.RecipeMakes, dish.Calories, dish.Carbs, dish.Fat, dish.Protein, dish.Directions, dish.fkUserId)
		fmt.Println(str)
		res, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Uploaded Sucessfully")
			fmt.Println(res)
		}
	}

}

func createdishingredient(w http.ResponseWriter, r *http.Request) {
	var dishingredient dishingredient
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &dishingredient)

	dishingredient.Amount = r.FormValue("Amount")
	u2, err := strconv.Atoi(r.FormValue("fkDishId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dishingredient.fkDishId = u2

	u2, err = strconv.Atoi(r.FormValue("fkIngredientId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dishingredient.fkIngredientId = u2

	str := fmt.Sprintf("INSERT INTO `myfridgefood`.`dishingredients` (`Amount`, `fk_DishId`, `fk_IngredientId`) VALUES ( '%v' , '%v' , '%v' );", dishingredient.Amount, dishingredient.fkDishId, dishingredient.fkIngredientId)
	fmt.Println(str)
	res, err := db.Query(str)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Uploaded Sucessfully")
		fmt.Println(res)
	}

}

func getOneuser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	results, err := db.Query("Select Email, Id, Password, UserType FROM users WHERE id = ?", userID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var user user

		err = results.Scan(&user.Email, &user.Id, &user.Password, &user.UserType)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(user)
	}
}

func getOneingredient(w http.ResponseWriter, r *http.Request) {
	ingredientID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	results, err := db.Query("Select * FROM ingredients WHERE id = ?", ingredientID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var ingredient ingredient

		err = results.Scan(&ingredient.Name, &ingredient.Id)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(ingredient)
	}
}

func getOnecomment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	results, err := db.Query("Select Date, Text, Id, fk_DishId, fk_UserId FROM comments WHERE Id = ?", commentID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var comments comment
		err = results.Scan(&comments.Date, &comments.Text, &comments.Id, &comments.fkDishId, &comments.fkUserId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(comments)
	}
}

func getOnedishingredient(w http.ResponseWriter, r *http.Request) {
	dishingredientID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	results, err := db.Query("Select Amount, Id, fk_DishId AS fkDishId, fk_IngredientId AS fkIngredientId FROM dishingredients WHERE Id = ?;", dishingredientID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var dishingredient dishingredient

		err = results.Scan(&dishingredient.Amount, &dishingredient.Id, &dishingredient.fkDishId, &dishingredient.fkIngredientId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(dishingredient)
	}
}

func getOnedish(w http.ResponseWriter, r *http.Request) {
	dishID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

	results, err := db.Query("Select * FROM dishes WHERE Id = ?;", dishID)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var dish dish

		err = results.Scan(&dish.Name, &dish.Description, &dish.CookingTime, &dish.RecipeMakes, &dish.Calories, &dish.Carbs, &dish.Fat, &dish.Protein, &dish.Directions, &dish.Id, &dish.fkUserId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(dish)
	}
}

func getAllusers(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("Select Email, Id, Password, UserType FROM users")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var user user

		err = results.Scan(&user.Email, &user.Id, &user.Password, &user.UserType)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(user)
	}
}

func getAllingredients(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("Select Name, Id FROM ingredients")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var ingredient ingredient

		err = results.Scan(&ingredient.Name, &ingredient.Id)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(ingredient)
	}
}

func getAllcomments(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("Select * FROM comments")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var comments comment

		err = results.Scan(&comments.Date, &comments.Text, &comments.Id, &comments.fkDishId, &comments.fkUserId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(comments)
	}
}

func getAlldishes(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("Select * FROM dishes")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var dish dish

		err = results.Scan(&dish.Name, &dish.Description, &dish.CookingTime, &dish.RecipeMakes, &dish.Calories, &dish.Carbs, &dish.Fat, &dish.Protein, &dish.Directions, &dish.Id, &dish.fkUserId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(dish)
	}
}

func getAlldishingredients(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("Select * FROM dishingredients")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var dishingredient dishingredient

		err = results.Scan(&dishingredient.Id, &dishingredient.Amount, &dishingredient.fkIngredientId, &dishingredient.fkDishId)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(dishingredient)
	}
}

func updateuser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	userID := mux.Vars(r)["id"]
	var user user
	fmt.Println("atejo")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "bad data")
	}
	json.Unmarshal(reqBody, &user)

	user.Email = r.FormValue("Email")
	user.Password = r.FormValue("Password")

	u2, err := strconv.Atoi(r.FormValue("UserType"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	user.UserType = u2
	if user.UserType > 2 || user.UserType < 1 {
		user.UserType = 1
	}
	if len(user.Email) < 1 {
		fmt.Fprintf(w, "Email can't be less than 1 character long")
	} else if len(user.Password) < 6 {
		fmt.Fprintf(w, "Password can't be less than 6 character long")
	} else {

		str := fmt.Sprintf("UPDATE `myfridgefood`.`users` SET `Password`='%v', `Email`= '%v', `UserType`= '%v' WHERE  `Id`=%v;", user.Password, user.Email, user.UserType, userID)
		fmt.Println(str)
		results, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Updated Sucessfully")
			fmt.Println(results)
		}
	}
}

func updateingredient(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	ingredientID := mux.Vars(r)["id"]
	var ingredient ingredient
	fmt.Println("atejo")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "bad data")
	}
	json.Unmarshal(reqBody, &ingredient)

	ingredient.Name = r.FormValue("Name")
	if len(ingredient.Name) < 1 {
		fmt.Fprintf(w, "Name can't be less than 1 character long")
	} else {
		str := fmt.Sprintf("UPDATE `myfridgefood`.`ingredients` SET `Name`='%v' WHERE  `Id`=%v;", ingredient.Name, ingredientID)
		fmt.Println(str)
		results, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Updated Sucessfully")
			fmt.Println(results)
		}
	}
}

func updatecomment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	commentID := mux.Vars(r)["id"]
	var comment comment
	fmt.Println("atejo")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "bad data")
	}
	json.Unmarshal(reqBody, &comment)

	comment.Text = r.FormValue("Text")
	if len(comment.Text) < 1 {
		fmt.Fprintf(w, "Text can't be less than 1 character long")
	} else {
		str := fmt.Sprintf("UPDATE `myfridgefood`.`comments` SET `Text`= '%v' WHERE  `Id`=%v;", comment.Text, commentID)
		fmt.Println(str)
		results, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Updated Sucessfully")
			fmt.Println(results)
		}
	}
}

func updatedishingredient(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	dishingredientID := mux.Vars(r)["id"]
	var dishingredient dishingredient
	fmt.Println("atejo")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "bad data")
	}
	json.Unmarshal(reqBody, &dishingredient)

	dishingredient.Amount = r.FormValue("Amount")
	if len(dishingredient.Amount) < 1 {
		fmt.Fprintf(w, "Amount can't be less than 1 character long")
	} else {
		str := fmt.Sprintf("UPDATE `myfridgefood`.`dishingredients` SET `Amount`= '%v' WHERE  `Id`=%v;", dishingredient.Amount, dishingredientID)
		fmt.Println(str)
		results, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Updated Sucessfully")
			fmt.Println(results)
		}
	}
}

func updatedish(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	dishID := mux.Vars(r)["id"]
	var dish dish
	fmt.Println("atejo")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "bad data")
	}
	json.Unmarshal(reqBody, &dish)

	dish.Name = r.FormValue("Name")
	dish.Description = r.FormValue("Description")

	u2, err := strconv.Atoi(r.FormValue("CookingTime"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.CookingTime = u2
	u2, err = strconv.Atoi(r.FormValue("RecipeMakes"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.RecipeMakes = u2

	dish.Calories, err = strconv.Atoi(r.FormValue("Calories"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Carbs, err = strconv.Atoi(r.FormValue("Carbs"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Fat, err = strconv.Atoi(r.FormValue("Fat"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.Protein, err = strconv.Atoi(r.FormValue("Protein"))
	if err != nil {
		fmt.Println("something went wrong")
	}

	dish.Directions = r.FormValue("Directions")

	u2, err = strconv.Atoi(r.FormValue("fkUserId"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	dish.fkUserId = u2
	if len(dish.Name) < 1 {
		fmt.Println("enter dish name")
	} else if len(dish.Description) < 1 {
		fmt.Println("enter dish description")
	} else if dish.CookingTime < 1 {
		fmt.Println("enter dish cookingTime")
	} else if dish.RecipeMakes < 1 {
		fmt.Println("enter dish recipeMakes")
	} else if len(dish.Directions) < 1 {
		fmt.Println("enter dish Directions")
	} else if dish.Calories < 1 {
		fmt.Println("enter dish calories cannot be negative")
	} else if dish.fkUserId == 0 {
		fmt.Println("user not assigned")
	} else {

		str := fmt.Sprintf("UPDATE `myfridgefood`.`dishes` SET `Name` = '%v', `Description` = '%v', `CookingTime` = '%v', `RecipeMakes` = '%v', `Calories` = '%v', `Carbs` = '%v', `Fat` = '%v', `Protein` = '%v', `Directions` = '%v' WHERE `Id`=%v;", dish.Name, dish.Description, dish.CookingTime, dish.RecipeMakes, dish.Calories, dish.Carbs, dish.Fat, dish.Protein, dish.Directions, dishID)
		fmt.Println(str)
		results, err := db.Query(str)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Updated Sucessfully")
			fmt.Println(results)
		}
	}
}

func deleteuser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("DELETE FROM `myfridgefood`.`users` WHERE `id`=?", userID)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Deleted Sucessfully")
		fmt.Println(results)
	}
}

func deleteingredient(w http.ResponseWriter, r *http.Request) {
	ingredientID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("DELETE FROM `myfridgefood`.`ingredients` WHERE `id`=?", ingredientID)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Deleted Sucessfully")
		fmt.Println(results)
	}
}

func deletecomment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("DELETE FROM `myfridgefood`.`comments` WHERE `id`=?", commentID)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Deleted Sucessfully")
		fmt.Println(results)
	}
}

func deletedishingredient(w http.ResponseWriter, r *http.Request) {
	dishingredientID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("DELETE FROM `myfridgefood`.`dishingredients` WHERE `id`=?", dishingredientID)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Deleted Sucessfully")
		fmt.Println(results)
	}
}

func deletedish(w http.ResponseWriter, r *http.Request) {
	dishID := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	results, err := db.Query("DELETE FROM `myfridgefood`.`dishes` WHERE `id`=?", dishID)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(w, "Deleted Sucessfully")
		fmt.Println(results)
	}
}
