package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

//SyncToRCenter - Sync events store data to reseach center
//First aid improrement global service
//Health check report integrations & updating first aid data & relavants services if needed
//Monitoring pmtrs which will be impacted once treatment started
//Drill & maintainance

var wg sync.WaitGroup

const sleepTs = 2

var events map[int]time.Time = make(map[int]time.Time, 10)
var eventsCnt int

const PositiveSignalInterval int = 5

var firstAidDB map[string]string = map[string]string{
	"Hattack": "1. Aspirin 15 mg\n2. Nitroglyceren 10 mg\n3. Beta blockers 50 ml",
	"Harrest": "1. apply AED \n2. keep applying CPR 100-120 compressions/minute with same intensity",
}

var signalDB map[int]string = map[int]string{
	1: "Hattack",
	2: "Harrest",
}

var eventsStore map[string][]int = make(map[string][]int, 0)
var mReports map[string]string = make(map[string]string, 0)

func classifySignal(signal int) int {
	res := harrestSignal(signal)
	if res != 0 {
		return res
	}
	res = haSignal(signal)
	return res
}

func haSignal(signal int) int {
	if signal%25 == 0 {
		return 1
	}
	return 0
}

func harrestSignal(signal int) int {
	if signal%23 == 0 {
		return 2
	}
	return 0
}

func main() {
	fmt.Println("Welcome!, MeAhe is watching you")

	for i := 0; i < 10000000; i++ {
		signal := rand.Int()
		sg, ok := signalDB[classifySignal(signal)]

		if ok {
			SendSignal(sg, signal)
		}
		time.Sleep(time.Microsecond * sleepTs)
	}
}

func SendSignal(sgType string, intensity int) {
	if eventsCnt > 1 {
		cur := events[eventsCnt]
		diffMin := cur.Sub(events[eventsCnt-1]).Minutes()
		if int(diffMin) < PositiveSignalInterval && sgType != "Harrest" {
			//fmt.Printf("Received one more signal within %d minutes\n", PositiveSignalInterval)
			return
		}
		eventsCnt++
		events[eventsCnt] = time.Now()
	} else {
		eventsCnt++
		events[eventsCnt] = time.Now()
	}

	action(sgType, intensity)

	if eventsCnt == 0 {
		fmt.Printf("Today was great day, Hoping healthy tomorrow!\nDon't worry i am allready watching tomorrow")
		os.Exit(1)
	}
}

func action(sg string, intensity int) {

	fmt.Printf("EMERGENCY::Cardiac event %s detected......\n", sg)
	wg.Add(3)

	go sendAlarm()
	go sendFirstAid("Hattack")
	go sendCalls()
	wg.Wait()
	saveEvent(sg, intensity)
}
func saveEvent(sg string, intensity int) {
	ev, ok := eventsStore[sg]
	if ok {
		ev = append(ev, intensity)
		eventsStore[sg] = ev
	} else {
		var ev []int
		ev = append(ev, intensity)
		eventsStore[sg] = ev
	}

}
func sendAlarm() {
	fmt.Println("Buzzing alarm ******")
	wg.Done()
}

func sendFirstAid(signal string) {
	fmt.Println(firstAidDB[signal])
	wg.Done()
}

func sendCalls() {
	//Keep calling until response received from emergecy call list
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond * 100)
		if (rand.Intn(100) % 50) == 0 {
			fmt.Println("Call received")
			break
		}
	}
	wg.Done()
}

func getReports() {
	m := ReadReports()
	for k, v := range m {
		m[k] = v
	}
}

func ReadReports() map[string]string {
	m := map[string]string{}
	m["allergy"] = "nitroglycerene"
	m["monitor"] = "cholesterol"
	m["deficiency"] = "B12"

	return m
}
