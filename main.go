package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type message string

type Worker struct {
	ID   int
	WG   *sync.WaitGroup
	Ch   chan message
	Next *Worker
}

func CreateWorkersList(n int) *Worker {
	var head *Worker
	var currentNode *Worker
	for i := 1; i <= n; i++ {
		newNode := &Worker{ID: i}
		if head == nil {
			head = newNode
			currentNode = newNode
		} else {
			currentNode.Next = newNode
			currentNode = newNode
		}
	}
	return head
}

func (w *Worker) Logger() {
	for msg := range w.Ch {
		log.Printf("logging %d", w.ID)
		fmt.Println(msg)
		time.Sleep(1 * time.Second)
		if w.Next != nil {
			ch := make(chan message)
			w.Next.Ch = ch
			w.Next.WG = w.WG
			w.WG.Add(1)
			go w.Next.Logger()
			w.Next.Ch <- msg
			go func() {
				defer w.WG.Done()
				close(ch)
			}()
		} else {
			defer w.WG.Done()
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}
	msgs := []message{
		"hello 1",
		"hello 2",
		"hello 3",
		"hello 4",
		"hello 5",
	}
	ch := make(chan message)
	workerList := CreateWorkersList(5)
	workerList.Ch = ch
	workerList.WG = wg

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go workerList.Logger()
	}

	for _, m := range msgs {
		ch <- m
	}

	go func() {
		defer wg.Done()
		close(ch)
	}()

	wg.Wait()

}
