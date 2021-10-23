package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const USER_COUNT = 10
const MAX_FILE_LOGS_CONT = 10
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

func main() {
	t := time.Now()
	rand.Seed(t.Unix())

	var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}
	SaveUserInfo(GenerateUsers(USER_COUNT, &actions))
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(t).Seconds())
}

// ------------ user

func GenerateUsers(count int, actions *[]string) *[]User {
	users := make([]User, count)
	logs := generateUsersLogs(count, MAX_FILE_LOGS_CONT, actions)
	jobs := make(chan int, count)
	results := make(chan *User, count)

	for w := 1; w <= count; w++ {
		go generateUsersWorker(jobs, results, (*logs)[w-1])
	}

	for j := 1; j <= count; j++ {
		jobs <- j - 1
	}
	close(jobs)

	for a := 1; a <= count; a++ {
		users[a-1] = *<-results
	}

	return &users
}

func generateUsersLogs(userCount, maxLogsCount int, actions *[]string) *[][]logItem {
	logsArr := make([][]logItem, userCount)
	for i := 0; i < userCount; i++ {
		logsCount := rand.Intn(maxLogsCount)
		logsArr[i] = make([]logItem, logsCount)

		for j := 0; j < logsCount; j++ {
			logsArr[i][j] = logItem{
				action:    (*actions)[rand.Intn(len((*actions))-1)],
				timestamp: time.Now(),
			}
		}
	}
	return &logsArr
}

func generateUsersWorker(jobs <-chan int, results chan<- *User, logs []logItem) {
	for j := range jobs {
		results <- generateUser(j, logs)
	}
}

func (u *User) getActivityInfo() *bytes.Buffer {
	buffer := new(bytes.Buffer)

	buffer.WriteString("UID: " + strconv.Itoa(u.id) + " Email: " + u.email + ";\nActivity Log:\n")
	for index, item := range u.logs {
		buffer.WriteString(strconv.Itoa(index) + ". [" + item.action + "] at " + item.timestamp.Format(time.RFC3339) + "\n")
	}

	return buffer
}

func generateUser(index int, logs []logItem) *User {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("generated user", index+1)

	return &User{
		id:    index + 1,
		email: "user" + strconv.Itoa(index+1) + "@company.com",
		logs:  logs,
	}
}

// -------- file write

func SaveUserInfo(users *[]User) {
	count := len(*users)
	jobs := make(chan *User, count)
	results := make(chan int, count)

	for w := 1; w <= count; w++ {
		go saveUserInfoWorker(jobs, results)
	}

	for j := 1; j <= count; j++ {
		jobs <- &(*users)[j-1]
	}
	close(jobs)

	for a := 1; a <= count; a++ {
		<-results
	}
}

func saveUser(user *User) int {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	file, err := os.OpenFile(USER_FILE_FOLDER+USER_FILE_NAME+strconv.Itoa(user.id)+USER_FILE_EXT, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = file.Write(user.getActivityInfo().Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Second)
	return user.id
}

func saveUserInfoWorker(jobs <-chan *User, results chan<- int) {
	for j := range jobs {
		results <- saveUser(j)
	}
}
