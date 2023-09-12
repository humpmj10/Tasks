package db

import (
	"Tasks/models"
	"encoding/json"
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
	"sync"
)

var (
	db     *badger.DB
	taskID int
	mu     sync.Mutex
)

func InitDB() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic(err)
	}
}

func GetTasks() ([]models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	var tasks []models.Task
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var task models.Task
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &task)
			}); err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func CreateTask(newTask models.Task) error {
	mu.Lock()
	defer mu.Unlock()

	err := db.Update(func(txn *badger.Txn) error {
		taskID++
		newTask.ID = taskID

		taskBytes, err := json.Marshal(newTask)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(newTask.Title), taskBytes)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(taskID int) error {
	mu.Lock()
	defer mu.Unlock()

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(fmt.Sprintf("%d", taskID)))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
