package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {

	sema := make(chan struct{}, pool)

	var waitG = sync.WaitGroup{}

	var DB struct {
		users []user
		mute  sync.Mutex
	}
	for i := 0; i < int(n); i++ {
		waitG.Add(1)
		go func(iD int64) {
			sema <- struct{}{}
			U := getOne(iD)
			<-sema
			DB.mute.Lock()
			DB.users = append(DB.users, U)
			DB.mute.Unlock()
			waitG.Done()
		}(int64(i))
	}
	waitG.Wait()
	return DB.users
}
