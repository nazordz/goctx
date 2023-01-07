package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.WithValue(context.Background(), "key", "value")
	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result: ", val)
	fmt.Println("took: ", time.Since(start))
}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200) // wait for 200 milisecond
	defer cancel()
	respch := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		respch <- Response{
			value: val,
			err:   err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetch data form third party took to long")
		case resp := <-respch:
			return resp.value, resp.err
		}

	}

}

func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	// time.Sleep(time.Millisecond * 150)
	time.Sleep(time.Millisecond * 150) // wait for 150 milisecon and that will no error
	return 666, nil
}
