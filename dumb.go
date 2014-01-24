package main

import (
	"github.com/jgrahamc/go-openssl/sha1"
	//"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	//"sync"
	"encoding/hex"
	"runtime"
)

//func gitMoney(difficulty string,in []byte, w *sync.WaitGroup) {
func gitMoney(difficulty string, in []byte, w chan bool) {
	hashes := 0
	now := time.Now()
	start_second := now.Truncate(time.Second)
	i := 0
	debug := false
	for {
		//text := fmt.Sprintf("tree %s \n parent %s \n author CTF user <me@example.com> %s +0000 \n committer CTF user <me@example.com> %s +0000 \n Give me a Gitcoin\n $d", "tree", "parent", "time", counter)

		//dumb
		t := strconv.Itoa(i)
		counter := []byte(t)

		h := sha1.New()
		body := append(in, counter...)
		fmt.Fprintf(h, "commit %d\x00", len(body))
		h.Write(body)
		sum := h.Sum(nil)
		//cs := fmt.Sprintf("%x\n", h.Sum(nil))
		cs := hex.EncodeToString(sum[:])
		if cs < difficulty {
			fmt.Printf("%s%s", in, t)
			fmt.Fprintln(os.Stderr, "\nHash:", cs)
			break
		}
		hashes++
		if debug == true {
			now := time.Now()
			end_second := now.Truncate(time.Second)
			if end_second.After(start_second) {
				fmt.Fprintln(os.Stderr, "hashes per second: %d \n", hashes)
				start_second = end_second
				hashes = 0
			}
		}
		i++
	}
	w <- true

}

func main() {

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println(err, string(in))
	}
	difficulty := os.Args[1]
	fmt.Fprintln(os.Stderr, "Called with difficulty:", difficulty)
	quit := make(chan bool)
	runtime.GOMAXPROCS(runtime.NumCPU())
	cores := runtime.NumCPU()
	//	var wg sync.WaitGroup
	//	wg.Add(1)
	for i:=0; i<=cores; i++{
	go gitMoney(difficulty, in, quit)
	}
	for {
		select {
		case <-quit:
			return
		}
	}
	//	wg.Wait()
}
