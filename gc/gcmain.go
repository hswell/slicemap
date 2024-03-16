package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"slicemap"
	"time"
)

type User struct {
	Id      int64
	Name    string
	Age     int
	Sex     string
	Address string
	Phone   string
	Email   string
}

func (u User) getKey() int64 {
	return u.Id
}
func newUser(id int64, name string, age int, sex string, address string, phone string, email string) *User {
	return &User{
		Id:      id,
		Name:    name,
		Age:     age,
		Sex:     sex,
		Address: address,
		Phone:   phone,
		Email:   email,
	}
}

func newUserObj(id int64, name string, age int, sex string, address string, phone string, email string) User {
	return User{
		Id:      id,
		Name:    name,
		Age:     age,
		Sex:     sex,
		Address: address,
		Phone:   phone,
		Email:   email,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func showGCStats(info debug.GCStats) string {
	return fmt.Sprintf("GC NumGC %v PauseTotal %v.", info.NumGC, info.PauseTotal)
}

func sliceMapGCTest() {
	sm := slicemap.NewSliceMap[int64]()
	debug.SetGCPercent(0)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 500000; i++ {
		randNum := rand.Intn(10)
		randStr := randString(randNum)
		sm.Add(&slicemap.KV[int64]{Key: int64(i), Value: newUserObj(int64(i), randStr, i, randStr, randStr, randStr, randStr)})
	}
	fmt.Println("Init slice-map done.")
	t := time.Now()
	debug.SetGCPercent(100)

	for i := 0; i < 10; i++ {
		runtime.GC()
	}
	fmt.Printf("Slice map gc x 10 cost: %v.\n", time.Now().Sub(t))
}

func sliceGCTest() {
	slice := make([]*User, 0)
	debug.SetGCPercent(0)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 500000; i++ {
		randNum := rand.Intn(10)
		randStr := randString(randNum)
		slice = append(slice, newUser(int64(i), randStr, i, randStr, randStr, randStr, randStr))
	}

	fmt.Println("Init slice done.")
	t := time.Now()
	debug.SetGCPercent(100)

	for i := 0; i < 10; i++ {
		runtime.GC()
	}
	fmt.Printf("Slice map gc x 10 cost: %v.\n", time.Now().Sub(t))
}

func mapGCTest() {
	mm := make(map[int]*User)
	debug.SetGCPercent(0)

	for i := 0; i < 500000; i++ {
		randNum := rand.Intn(10)
		randStr := randString(randNum)
		mm[i] = newUser(int64(i), randStr, i, randStr, randStr, randStr, randStr)
	}

	fmt.Println("Init map done.")
	t := time.Now()
	debug.SetGCPercent(100)

	for i := 0; i < 10; i++ {
		runtime.GC()
	}
	fmt.Printf("Go map gc x 10 cost: %v.\n", time.Now().Sub(t))
}

func main() {
	//fmt.Println("Do test using slice-map")
	//sliceMapGCTest()
	fmt.Println("Do test using go map")
	mapGCTest()
	//fmt.Println("Do test using go slice")
	//sliceGCTest()
}

// cpu: m1 pro
// memory: 32G

//Do test using slice-map
//Init slice-map done.
//Slice map gc x 10 cost: 58.349833ms.

//Do test using go map
//Init map done.
//Go map gc x 10 cost: 235.258125ms.

//Do test using go slice
//Init slice done.
//Slice map gc x 10 cost: 43.631375ms.

//Do test using go map
//Init map done.
//Go map gc x 10 cost: 3.796166ms.

//Do test using go slice
//Init slice done.
//Slice map gc x 10 cost: 3.637083ms.
