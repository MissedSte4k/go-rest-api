package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	jwt "github.com/jwt-go"
)

func createcomment(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {
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
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
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

		i, err := strconv.Atoi(claims.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		comments.fkUserId = i
		if len(comments.Text) < 1 {
			fmt.Fprintf(w, "Text can't be less than 1 character long")
			w.WriteHeader(http.StatusBadRequest)
		} else if comments.fkDishId == 0 {
			fmt.Fprintf(w, "Wrong Dish ID")
			w.WriteHeader(http.StatusBadRequest)
		} else if comments.fkUserId == 0 {
			fmt.Fprintf(w, "Wrong Dish ID")
			w.WriteHeader(http.StatusBadRequest)
		} else {
			str = fmt.Sprintf("INSERT INTO `myfridgefood`.`comments` (`Date`, `Text`, `fk_DishId`, `fk_UserId`) VALUES ( `'%v'` , `'%v'` , `'%v'`, `'%v'` );", comments.Date, comments.Text, comments.fkDishId, comments.fkUserId)
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
}

func createingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
			w.WriteHeader(http.StatusBadRequest)
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
}

func createuser(w http.ResponseWriter, r *http.Request) {
	var user, user2 user

	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
	}
	json.Unmarshal(reqBody, &user)

	user.Email = r.FormValue("Email")
	user.Password = r.FormValue("Password")
	user.Username = r.FormValue("Username")

	results, err := db.Query("Select * FROM users WHERE Username = ?", user.Username)

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {

		err = results.Scan(&user2.Email, &user2.Password, &user2.Id, &user2.UserType, &user2.Username)
		if err != nil {
			panic(err.Error())
		}
		if user2.Email == user.Email {
			http.Error(w, "user with this email already registered", http.StatusForbidden)
			return
		}
		if user2.Username == user.Username {
			http.Error(w, "user with this username already registered", http.StatusForbidden)
			return
		}
	}

	u2, err := strconv.Atoi(r.FormValue("UserType"))
	if err != nil {
		fmt.Println("something went wrong")
	}
	user.UserType = u2
	if len(user.Email) < 1 {
		fmt.Fprintf(w, "please enter an email adress")
		w.WriteHeader(http.StatusBadRequest)
	} else if len(user.Email) < 1 {
		fmt.Fprintf(w, "please enter an username")
		w.WriteHeader(http.StatusBadRequest)
	} else if len(user.Password) < 6 {
		fmt.Fprintf(w, "your password canot be less than, 6 characters long")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		str := fmt.Sprintf("INSERT INTO `myfridgefood`.`users` (`Email`, `Password`, `UserType`, `Username`) VALUES ( '%v' , '%v' , '%v', '%v' );", user.Email, user.Password, user.UserType, user.Username)
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
	if Refresh(w, r) && UserType(w, r) >= 1 {
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

		dish.Calories, err = strconv.ParseFloat(r.FormValue("Calories"), 64)
		if err != nil {
			fmt.Println("something went wrong")
		}
		dish.Carbs, err = strconv.ParseFloat(r.FormValue("Carbs"), 64)
		if err != nil {
			fmt.Println("something went wrong")
		}
		dish.Fat, err = strconv.ParseFloat(r.FormValue("Fat"), 64)
		if err != nil {
			fmt.Println("something went wrong")
		}
		dish.Protein, err = strconv.ParseFloat(r.FormValue("Protein"), 64)
		if err != nil {
			fmt.Println("something went wrong")
		}

		dish.Directions = r.FormValue("Directions")

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
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

		i, err := strconv.Atoi(claims.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		dish.fkUserId = i
		if len(dish.Name) < 1 {
			fmt.Println("enter dish name")
			w.WriteHeader(http.StatusBadRequest)
		} else if len(dish.Description) < 1 {
			fmt.Println("enter dish description")
			w.WriteHeader(http.StatusBadRequest)
		} else if dish.CookingTime < 1 {
			fmt.Println("enter dish cookingTime")
			w.WriteHeader(http.StatusBadRequest)
		} else if dish.RecipeMakes < 1 {
			fmt.Println("enter dish recipeMakes")
			w.WriteHeader(http.StatusBadRequest)
		} else if len(dish.Directions) < 1 {
			fmt.Println("enter dish Directions")
			w.WriteHeader(http.StatusBadRequest)
		} else if dish.Calories < 1 {
			fmt.Println("enter dish calories cannot be negative")
			w.WriteHeader(http.StatusBadRequest)
		} else if dish.fkUserId == 0 {
			fmt.Println("user not assigned")
			w.WriteHeader(http.StatusBadRequest)
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
}

func createdishingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {
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
		if len(dishingredient.Amount) < 1 {
			fmt.Println("Amount not entered")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if dishingredient.fkDishId < 1 {
			fmt.Println("Dish id not entered")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if dishingredient.fkIngredientId < 1 {
			fmt.Println("ingredient id not entered")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
}

func getOneuser(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) {
		userID := mux.Vars(r)["id"]
		if UserType(w, r) >= 1 {
			if UserType(w, r) == 1 && userID != TokenIdGet(w, r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

			results, err := db.Query("Select * FROM users WHERE id = ?", userID)

			if err != nil {
				panic(err.Error())
			}

			for results.Next() {
				var user user

				err = results.Scan(&user.Email, &user.Password, &user.Id, &user.UserType, &user.Username)
				if err != nil {
					panic(err.Error())
				}
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func getOneingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
}

func getOnecomment(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {

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
			if UserType(w, r) == 1 && strconv.Itoa(comments.fkUserId) != TokenIdGet(w, r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			json.NewEncoder(w).Encode(comments)
		}
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
	if Refresh(w, r) && UserType(w, r) > 1 {
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
		results, err := db.Query("Select * FROM users")

		if err != nil {
			panic(err.Error())
		}

		for results.Next() {
			var user user

			err = results.Scan(&user.Email, &user.Password, &user.Id, &user.UserType, &user.Username)
			if err != nil {
				panic(err.Error())
			}
			json.NewEncoder(w).Encode(user)
		}
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
	if Refresh(w, r) && UserType(w, r) >= 1 {
		userID := mux.Vars(r)["id"]

		if UserType(w, r) == 1 && userID != TokenIdGet(w, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
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
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if len(user.Password) < 6 {
			fmt.Fprintf(w, "Password can't be less than 6 character long")
			w.WriteHeader(http.StatusBadRequest)
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
}

func updateingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
}

func updatecomment(w http.ResponseWriter, r *http.Request) {

	if Refresh(w, r) && UserType(w, r) >= 1 {
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
		commentID := mux.Vars(r)["id"]

		if err != nil {
			fmt.Fprintf(w, "bad data")
		}
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
			if UserType(w, r) == 1 && strconv.Itoa(comments.fkUserId) != TokenIdGet(w, r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		var comment comment
		reqBody, err := ioutil.ReadAll(r.Body)

		json.Unmarshal(reqBody, &comment)

		comment.Text = r.FormValue("Text")
		if len(comment.Text) < 1 {
			fmt.Fprintf(w, "Text can't be less than 1 character long")
			w.WriteHeader(http.StatusBadRequest)
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
}

func updatedishingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) {
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
			w.WriteHeader(http.StatusBadRequest)
			return
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
}

func updatedish(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
		dishID := mux.Vars(r)["id"]

		results, err := db.Query("Select * FROM dishes WHERE Id = ?;", dishID)

		if err != nil {
			panic(err.Error())
		}

		if results.Next() {

		} else {

		}

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dish.CookingTime = u2
		u2, err = strconv.Atoi(r.FormValue("RecipeMakes"))
		if err != nil {
			fmt.Println("something went wrong")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dish.RecipeMakes = u2

		dish.Calories, err = strconv.ParseFloat(r.FormValue("Calories"), 64)
		if err != nil {
			fmt.Println("something went wrong")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dish.Carbs, err = strconv.ParseFloat(r.FormValue("Carbs"), 64)
		if err != nil {
			fmt.Println("something went wrong")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dish.Fat, err = strconv.ParseFloat(r.FormValue("Fat"), 64)
		if err != nil {
			fmt.Println("something went wrong")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dish.Protein, err = strconv.ParseFloat(r.FormValue("Protein"), 64)
		if err != nil {
			fmt.Println("something went wrong")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dish.Directions = r.FormValue("Directions")

		if len(dish.Name) < 1 {
			fmt.Println("enter dish name")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if len(dish.Description) < 1 {
			fmt.Println("enter dish description")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if dish.CookingTime < 1 {
			fmt.Println("enter dish cookingTime")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if dish.RecipeMakes < 1 {
			fmt.Println("enter dish recipeMakes")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if len(dish.Directions) < 1 {
			fmt.Println("enter dish Directions")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if dish.Calories < 1 {
			fmt.Println("enter dish calories cannot be negative")
			w.WriteHeader(http.StatusBadRequest)
			return
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
}

func deleteuser(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
}

func deleteingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
}

func deletecomment(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")
		commentID := mux.Vars(r)["id"]
		if UserType(w, r) == 1 {
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
				if UserType(w, r) == 1 && strconv.Itoa(comments.fkUserId) != TokenIdGet(w, r) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}

		results, err := db.Query("DELETE FROM `myfridgefood`.`comments` WHERE `id`=?", commentID)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Deleted Sucessfully")
			fmt.Println(results)
		}
	}
}

func deletedishingredient(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) > 1 {
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
}

func deletedish(w http.ResponseWriter, r *http.Request) {
	if Refresh(w, r) && UserType(w, r) >= 1 {
		dishID := mux.Vars(r)["id"]
		db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1)/myfridgefood")

		if UserType(w, r) == 1 {
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
				if UserType(w, r) == 1 && strconv.Itoa(dish.fkUserId) != TokenIdGet(w, r) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}

		results, err := db.Query("DELETE FROM `myfridgefood`.`dishes` WHERE `id`=?", dishID)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "Deleted Sucessfully")
			fmt.Println(results)
		}
	}
}
