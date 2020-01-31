package upload

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"

	"projects/elections/dynamo"
	"projects/elections/persist"
)

/* queue interface uses 2 stacks of fixed size to implement FIFO queue */
type queue struct {
	Stack1 []string // items pushed to queue / input
	Stack2 []string // items popped from queue / output
	Index1 int      // next index after top of stack1 (top == index - 1)
	Index2 int      // next index after top of stack2
	Limit  int
}

func initQueue(size int) queue {
	q := queue{Stack1: make([]string, size), Stack2: make([]string, size), Index1: 0, Index2: 0, Limit: size}
	return q
}

// Add single item to queue
func (q *queue) Push(item string) error {
	if q.Index1+q.Index2 == q.Limit {
		return fmt.Errorf("queue full")
	}
	q.Stack1[q.Index1] = item
	q.Index1++
	return nil
}

// Add list of items to queue
func (q *queue) PushMulti(items []string) ([]string, error) {
	for i, item := range items {
		err := q.Push(item)
		if err != nil { // queue full; return remaining items
			return items[i:], err
		}
	}
	return nil, nil
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
func (q *queue) PopMulti(num int) []string {
	items := []string{}
	for i := 0; i < num; i++ {
		items = append(items, q.Pop())
		if q.Index2 == 0 {
			break
		}
	}
	return items
}

/* End queue operations */

func getYear() (string, error) {
	yearStr := "" // default return value

	for {
		// get flag and check validity by verifying if bucket for specified year exists
		yearFlag := flag.Int("year", 0, "'year' flag defines which election cycle's dataset to process")
		flag.Parse()
		year := *yearFlag
		if year == 0 {
			fmt.Println("'year' flag must be set to valid year")
			continue
		}

		// convert int value to string
		yearStr = strconv.Itoa(year)

		// open db and start view tx
		db, err := bolt.Open("db/offline_db.db", 0644, nil)
		defer db.Close()
		if err != nil {
			fmt.Println("FATAL: getYear failed: ")
			return "", fmt.Errorf("getYear failed: %v", err)
		}

		// check validity by searching for bucket corresponding to given year
		exists := true
		if err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(yearStr))
			if bucket == nil {
				exists = false
			}
			return nil
		}); err != nil {
			fmt.Println("FATAL: getYear failed: ", err)
			return "", fmt.Errorf("getYear failed: %v", err)
		}

		if exists == false {
			fmt.Printf("Invalid year: %d --- No dataset found!\n", year)
			continue
		} else {
			break
		}
	}
	return yearStr, nil
}

// upload Candidate objects
func uploadCandObjs(year, curr string, dbi *dynamo.DbInfo) (string, error) {
	idQueue := initQueue(1000)

	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: uploadCandObjs failed: 'offline_db.db' failed to open")
		return curr, fmt.Errorf("uploadCandObjs failed: 'offline_db.db' failed to open: %v", err)
	}

	end := false // break outermost loop when true
	for {
		// view Tx - loop thru 1000 keys per outer iteration / 1 tx per iteration
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(year)).Bucket([]byte("candidates"))
			c := b.Cursor()

			if curr == "" { // start of bucket
				// iterate thru each item in the given year's dataset
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					// create a queue of 1000 item id's max
					idQueue.Push(string(k))
					if idQueue.Index1 == idQueue.Limit {
						curr = string(k)
						break
					}
				}
				if idQueue.Index1 < idQueue.Limit {
					end = true
				}
			} else { // subsequent iterations
				// iterate thru each item in the given year's dataset
				start, _ := c.Seek([]byte(curr))
				for k, _ := c.Seek([]byte(curr)); k != nil; k, _ = c.Next() {
					// skip first item - processed as last item in previous iteration of outer loop
					if string(k) == string(start) {
						continue
					}
					// create a queue of 1000 item id's max
					idQueue.Push(string(k))
					if idQueue.Index1 == idQueue.Limit {
						curr = string(k)
						break
					}
				}
				if idQueue.Index1 < idQueue.Limit {
					end = true
				}
			}
			return nil
		})
		// end tx logic

		// while id's in queue
		for {
			// pop 25 item ids from queue - tx queue
			txIDs := idQueue.PopMulti(25)

			// for each item, get corresponding obj from database and add to tx queue
			writes := make([]interface{}, 25)
			for i, ID := range txIDs {
				obj, err := persist.GetCandidate(year, ID)
				if err != nil {
					fmt.Println("uploadCandObjs failed: ", err)
					return curr, fmt.Errorf("")
				}
				writes[i] = obj
			}

			// batch write 25 items to corresponding dynamoDB table
			err := dynamo.BatchWriteCreate(dbi.Svc, writes, dbi.Tables["candidates"], dbi.FailConfig)
			if err != nil {
				fmt.Println("uploadCandObjs failed: ", err)
				return curr, fmt.Errorf("uploadCandObjs failed: %v", err)
			}

			// break when tx queue < 25
			if len(txIDs) < 25 {
				break
			}
		}

		// break when id queue < 1000
		if end == true {
			break
		}
	}
	return curr, nil
}
