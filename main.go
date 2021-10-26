package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const USER_COUNT = 100
const MAX_FILE_LOGS_CONT = 1000
const USER_FILE_FOLDER = "users/"
const USER_FILE_NAME = "uid"
const USER_FILE_EXT = ".txt"

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() *bytes.Buffer {
	buffer := new(bytes.Buffer)

	buffer.WriteString("UID: " + strconv.Itoa(u.id) + " Email: " + u.email + ";\nActivity Log:\n")
	for index, item := range u.logs {
		buffer.WriteString(strconv.Itoa(index) + ". [" + item.action + "] at " + item.timestamp.Format(time.RFC3339) + "\n")
	}

	return buffer
}

func main() {
	t := time.Now()
	rand.Seed(t.Unix())

	var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}
	RunPipeline(USER_COUNT, &actions)

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(t).Seconds())
}

// RunPipeline - user processing pipeline
func RunPipeline(count int, actions *[]string) {
	data := make(chan int, count)          // data for user creating
	users := make(chan *User, count)       // users
	saveResults := make(chan error, count) // errors
	wg := &sync.WaitGroup{}

	// creating user pipeline
	n := 0
	for n < count {
		wg.Add(1)
		go generateUserPipe(data, users, MAX_FILE_LOGS_CONT, actions)
		data <- n
		n += 1
	}

	// saving user pipeline
	n = 0
	for n < count {
		go saveUserPipe(wg, users, saveResults)
		n += 1
	}

	// catch errors
	n = 0
	for n < count {
		if <-saveResults != nil {
			fmt.Println((<-saveResults).Error())
		}
		n += 1
	}

	wg.Wait()
	close(data)
	close(users)
	close(saveResults)
}

// generateUserPipe - getting user id from one channel and sending created user to other channel
func generateUserPipe(data <-chan int, users chan<- *User, maxLogsCount int, actions *[]string) {
	for userId := range data {
		users <- generateUser(userId, maxLogsCount, actions)
	}
}

// generateUser - create user
func generateUser(index, maxLogsCount int, actions *[]string) *User {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("generated user", index+1)

	return &User{
		id:    index + 1,
		email: "user" + strconv.Itoa(index+1) + "@company.com",
		logs:  *generateLogs(maxLogsCount, actions),
	}
}

// generateLogs - create logs
func generateLogs(maxLogsCount int, actions *[]string) *[]logItem {
	logsCount := rand.Intn(maxLogsCount)
	logsArr := make([]logItem, logsCount)

	for i := 0; i < logsCount; i++ {
		logsArr[i] = logItem{
			action:    (*actions)[rand.Intn(len((*actions))-1)],
			timestamp: time.Now(),
		}
	}
	return &logsArr
}

// saveUserPipe - getting users from one channel and sending to other channel
func saveUserPipe(wg *sync.WaitGroup, users <-chan *User, saveResults chan<- error) {
	for user := range users {
		saveResults <- saveUser(user)
		wg.Done()
	}
}

// saveUser - write user info into file with store errors
func saveUser(user *User) error {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	file, err := os.OpenFile(USER_FILE_FOLDER+USER_FILE_NAME+strconv.Itoa(user.id)+USER_FILE_EXT, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(user.getActivityInfo().Bytes())
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return nil
}
