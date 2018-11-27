package main

import (
	"net/http"
	"io/ioutil"
	"sync"
	"time"
	"fmt"
)


type ChunkAttack struct {
	success bool
	err error
	msg string
}

func NewChunkAttack(t *TestDef) (Attack) {
	return &ChunkAttack{}
}

func (c *ChunkAttack) Run(r chan Attack, d chan bool, wg *sync.WaitGroup) {

Loop:
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://localhost:8000/chunk/1")
		if err != nil {
			c.msg = "Failed"
		} else {
			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				c.msg = "Foiled"
			} else {
				c.msg = "success"
			}		
		}

		r <- c

		select {
			case _ = <- d:
				break Loop
			case <-time.After(100 * time.Millisecond):
				continue
		}
	}
	wg.Done()
}

func(c *ChunkAttack) Msg() (string) {
	return c.msg
}


type Attack interface {
	// Attack needs to change
	Run(r chan Attack, d chan bool, wg *sync.WaitGroup)
	Msg() (string)
}

type AttackRunner = func(t *TestDef) (Attack)


func collect(r chan Attack) {
	for {
		select {
			case attack := <- r:
				fmt.Println(attack.Msg())
		}
	}
}

type TestDef struct {
	resultsChannel chan Attack
	doneChannel chan bool
	wg sync.WaitGroup
	ar AttackRunner
}

func NewTest(ar AttackRunner) (*TestDef) {
	return &TestDef{
		resultsChannel: make(chan Attack, 10000),
		doneChannel: make(chan bool, 10000),
		wg: sync.WaitGroup{},
		ar: ar,
	}
}

func (t *TestDef) rampup(ms time.Duration, usercount int) {
	sleeptime := ms / time.Duration(usercount)
	for i := 0; i < usercount; i++ {
		t.wg.Add(1)
		a := t.ar(t)
		go a.Run(t.resultsChannel, t.doneChannel, &t.wg)
		time.Sleep(sleeptime)
	}
}

func (t *TestDef) rampdown(ms time.Duration, usercount int) {
	sleeptime := ms / time.Duration(usercount)
	for i := 0; i < usercount; i++ {
		t.doneChannel <- true
		time.Sleep(sleeptime)
	}	
}

func (t *TestDef) spawnAttack() {
	go collect(t.resultsChannel)
	t.rampup(2000, 2) 
	t.rampup(2000, 2)
	time.Sleep(5*time.Second)
	t.rampdown(2000, 4) 
	t.wg.Wait()
}

func main() {
	t := NewTest(NewChunkAttack)
	t.spawnAttack()
}