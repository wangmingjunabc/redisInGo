package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

//User is a struct representing newly register users
type User struct {
	UserName string
	Email    string
}

//MarshalBinary encode the struct into a binary blob
//Here I cheat and use regular json:)
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

//UnmarshalBinary decodes the struct into a user
func (u *User) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return nil
}

//Names Some Non-Random name lists used to generate Random Users
var Names []string = []string{"Jasper", "Johan", "Edward", "Niel", "Percy"}

//SirName Some Non-Random name lists used to generate Random Users
var SirNames []string = []string{"Ericsson", "Redisson", "Edisson", "Tesla", "Bolmer"}

//EmailProviders Some Non-Random email lists to generate Random Users
var EmailProviders []string = []string{"Hotmail.com", "Gmail.com", "Awesomeness.com", "Redis.com"}

func main() {
	//Create a new Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "superSecret",
		DB:       0,
	})
	//Ping the redis server and check if any errors occurred
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		//Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	//Generate a new background context that we will use
	ctx := context.Background()
	//Loop and randomly generate users on a random timer
	for {
		//Publish a generate user to the new_user channel
		err := redisClient.Publish(ctx, "new_users", GenerateRandomUser()).Err()
		if err != nil {
			panic(err)
		}
		//sleep random time
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(4)
		time.Sleep(time.Duration(n) * time.Second)
	}
}

//GenerateRandomUser creates a random user, don't care too much about this.
func GenerateRandomUser() *User {
	rand.Seed(time.Now().UnixNano())
	nameMax := len(Names)
	sirNameMax := len(SirNames)
	EmailProviderMax := len(EmailProviders)

	nameIndex := rand.Intn(nameMax-1) + 1
	sirNameIndex := rand.Intn(sirNameMax-1) + 1
	emailIndex := rand.Intn(EmailProviderMax-1) + 1

	return &User{
		UserName: Names[nameIndex] + " " + SirNames[sirNameIndex],
		Email:    Names[nameIndex] + SirNames[sirNameIndex] + "@" + EmailProviders[emailIndex],
	}
}

