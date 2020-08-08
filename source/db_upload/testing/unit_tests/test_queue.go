package main

import "fmt"

type queue struct {
	Stack1 []string // items pushed to queue / input
	Stack2 []string // items popped from queue / output
	Index1 int      // next index after top of stack (top == index - 1)
	Index2 int      // next index after top of stack
}

func initIDQueue() queue {
	q := queue{Stack1: make([]string, 1000), Stack2: make([]string, 1000), Index1: 0, Index2: 0}
	return q
}

func initTxQueue() queue {
	q := queue{Stack1: make([]string, 25), Stack2: make([]string, 25), Index1: 0, Index2: 0}
	return q
}

// Add single item to queue
func (q *queue) Push(item string) {
	q.Stack1[q.Index1] = item
	q.Index1++
}

// Add list of items to queue
func (q *queue) PushMulti(items []string) {
	for _, item := range items {
		q.Push(item)
	}
}

// Remove single item from queue
func (q *queue) Pop() string {
	if q.Index1 == 0 && q.Index2 == 0 { // popping from empty queue
		return ""
	}

	if q.Index2 == 0 {
		// reassign each item in Stack1 to Stack2 in reverse order
		for i := (q.Index1 - 1); i >= 0; i-- {
			q.Stack2[q.Index2] = q.Stack1[i]
			q.Stack1[i] = ""
			q.Index1--
			q.Index2++
		}
	}

	// return last item in Stack2
	item := q.Stack2[q.Index2-1]
	q.Stack2[q.Index2-1] = ""
	q.Index2--
	return item
}

// Remove and return multiple items from queue
func (q *queue) MultiPop(num int) []string {
	items := []string{}
	for i := 0; i < num; i++ {
		items = append(items, q.Pop())
		if q.Index2 == 0 {
			break
		}
	}
	return items
}

func main() {
	shortList := []string{"1", "2", "3", "4", "5"}
	// list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	q1 := initTxQueue() // cap == 25
	// q2 := initIDQueue() // cap == 1000

	fmt.Println("----- SINGLE ITEM TEST -----")
	for _, num := range shortList {
		q1.Push(num)
	}
	fmt.Println("q1 - each item pushed individually")
	fmt.Println(q1)

	item1 := q1.Pop()
	item2 := q1.Pop()
	fmt.Println(q1)
	// fmt.Println("len(Stack1): ", len(q1.Stack1))
	// fmt.Println("cap(Stack1): ", cap(q1.Stack1))

	fmt.Println("1st Pop - item1: ", item1)
	fmt.Println("2nd Pop - item2: ", item2)

	q1.Push("6")
	q1.Push("7")
	q1.Push("8")
	fmt.Println(q1)
	q1.Pop()
	q1.Pop()
	q1.Pop()
	fmt.Println(q1)
	last := q1.Pop()
	fmt.Println(last)
	fmt.Println(q1)

	/* fmt.Println("----- MULTI ITEMS TEST -----")

	q2.PushMulti(list)
	fmt.Println("q2 - multi push")
	fmt.Println(q2)
	fmt.Println("Index1: ", q2.Index1)
	fmt.Println("Index2: ", q2.Index2)

	items := q2.MultiPop(10)
	fmt.Println("MultiPop - items: ", items)
	fmt.Println(q2)
	items2 := q2.MultiPop(10)
	fmt.Println("items2: ", items2)
	fmt.Println("Index1, Index2: ", q2.Index1, q2.Index2)
	fmt.Println("items3 :", q2.MultiPop(10))
	fmt.Println("Index1, Index2: ", q2.Index1, q2.Index2) */

}
