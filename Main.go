package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"time"
	"net/smtp"
	"log"
	"fmt"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetAllUser())
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(ctx iris.Context) {
	var person Person
	_ = ctx.ReadJSON(&person)
	fmt.Print("Person's username data: "+ person.Username)
	address := person.Address.City + " " + person.Address.State
	t,_ := time.Parse("2006-01-02","2017-02-02")
	insertUser(person.Username, person.Password, person.Firstname, person.Lastname, address, t.Format("2006-01-02"))
	sendEmailVerification()
	ctx.JSON("true")
}

func sendEmailVerification(){
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		"erwin@excite.co.id",
		"jmb_Ultima[1]",
		"smtp.gmail.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"button.setonclicklistener@gmail.com",
		[]string{"ultima51@yahoo.com"},
		[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// authenticate an user
func AuthenticateUser(ctx iris.Context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	ctx.JSON(UserAuthentication(username, password))
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

// our main function
func main() {
	/*
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/registration", CreatePerson).Methods("POST").HeadersRegexp("Content-Type", "application/(text|json)")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
	*/

	//lets try to implement iris
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal
	app.Use(recover.New())
	app.Use(logger.New())


	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/people", func(ctx iris.Context) {
		ctx.JSON(GetAllUser())
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	app.Post("/registration", CreatePerson)
	app.Post("/authentication", AuthenticateUser)

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8000"), iris.WithoutServerError(iris.ErrServerClosed))
}

// Game contains the state of a bowling game.
type Game struct {
	rolls   []int
	current int
}

// NewGame allocates and starts a new game of bowling.
func NewGame() *Game {
	game := new(Game)
	game.rolls = make([]int, maxThrowsPerGame)
	return game
}

// Roll rolls the ball and knocks down the number of pins specified by pins.
func (self *Game) Roll(pins int) {
	self.rolls[self.current] = pins
	self.current++
}

// Score calculates and returns the player's current score.
func (self *Game) Score() (sum int) {
	for throw, frame := 0, 0; frame < framesPerGame; frame++ {
		if self.isStrike(throw) {
			sum += self.strikeBonusFor(throw)
			throw += 1
		} else if self.isSpare(throw) {
			sum += self.spareBonusFor(throw)
			throw += 2
		} else {
			sum += self.framePointsAt(throw)
			throw += 2
		}
	}
	return sum
}

// isStrike determines if a given throw is a strike or not. A strike is knocking
// down all pins in one throw.
func (self *Game) isStrike(throw int) bool {
	return self.rolls[throw] == allPins
}

// strikeBonusFor calculates and returns the strike bonus for a throw.
func (self *Game) strikeBonusFor(throw int) int {
	return allPins + self.framePointsAt(throw+1)
}

// isSpare determines if a given frame is a spare or not. A spare is knocking
// down all pins in one frame with two throws.
func (self *Game) isSpare(throw int) bool {
	return self.framePointsAt(throw) == allPins
}

// spareBonusFor calculates and returns the spare bonus for a throw.
func (self *Game) spareBonusFor(throw int) int {
	return allPins + self.rolls[throw+2]
}

// framePointsAt computes and returns the score in a frame specified by throw.
func (self *Game) framePointsAt(throw int) int {
	return self.rolls[throw] + self.rolls[throw+1]
}

// testing utilities:

func (self *Game) rollMany(times, pins int) {
	for x := 0; x < times; x++ {
		self.Roll(pins)
	}
}
func (self *Game) rollSpare() {
	self.Roll(5)
	self.Roll(5)
}
func (self *Game) rollStrike() {
	self.Roll(10)
}

const (
	// allPins is the number of pins allocated per fresh throw.
	allPins = 10

	// framesPerGame is the number of frames per bowling game.
	framesPerGame = 10

	// maxThrowsPerGame is the maximum number of throws possible in a single game.
	maxThrowsPerGame = 21
)

