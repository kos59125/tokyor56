package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type (
	PageWeight struct {
		Path   string
		Weight float64
	}
	UserAction func()
)

var (
	landing = []PageWeight{
		PageWeight{Path: "/", Weight: 0.2},
		PageWeight{Path: "/page1", Weight: 0.4},
		PageWeight{Path: "/page2", Weight: 0.4},
	}
	transit = map[string][]PageWeight{
		"/": []PageWeight{
			PageWeight{Path: "/page1", Weight: 0.5},
			PageWeight{Path: "/page2", Weight: 0.4},
			PageWeight{Path: "", Weight: 0.1},
		},
		"/page1": []PageWeight{
			PageWeight{Path: "/", Weight: 0.2},
			PageWeight{Path: "/page2", Weight: 0.7},
			PageWeight{Path: "/404", Weight: 0.05},
			PageWeight{Path: "", Weight: 0.05},
		},
		"/page2": []PageWeight{
			PageWeight{Path: "/", Weight: 0.1},
			PageWeight{Path: "/page1", Weight: 0.6},
			PageWeight{Path: "", Weight: 0.3},
		},
		"/form": []PageWeight{
			PageWeight{Path: "/", Weight: 0.1},
			PageWeight{Path: "/form", Weight: 0.2},
			PageWeight{Path: "/503", Weight: 0.01},
			PageWeight{Path: "/cv", Weight: 0.09},
			PageWeight{Path: "", Weight: 0.6},
		},
		"/cv": []PageWeight{
			PageWeight{Path: "", Weight: 1.0},
		},
		"/404": []PageWeight{
			PageWeight{Path: "/", Weight: 0.2},
			PageWeight{Path: "", Weight: 0.8},
		},
		"/503": []PageWeight{
			PageWeight{Path: "/", Weight: 0.1},
			PageWeight{Path: "/503", Weight: 0.4},
			PageWeight{Path: "", Weight: 0.5},
		},
	}
	lineFeed = []byte{0x0A}
)

func generatePage(table []PageWeight) string {
	var page string
	s := 0.0
	w := rand.Float64()
	for i := 0; i < len(table); i++ {
		s = s + table[i].Weight
		page = table[i].Path
		if w <= s {
			break
		}
	}
	return page
}

func generateWeightTime(s, base float64) time.Duration {
	w := rand.ExpFloat64()*s + base
	return time.Duration(int(w)) * time.Millisecond
}

func newUser() *User {
	lp := generatePage(landing)
	return NewUser(lp)
}

func generateUserAction(user *User) UserAction {
	return func() {
		log.Printf("User %s arrived at %s\n", user.ID, user.CurrentPath)
		for user.CurrentPath != "" {
			// サーバーログ
			user.LeaveFootprint()

			// ページ閲覧なう
			wt := generateWeightTime(1000.0, 200.0)
			time.Sleep(wt)

			// ページ遷移
			nextPage := generatePage(transit[user.CurrentPath])
			user.MoveTo(nextPage)
		}
		log.Printf("User %s exited\n", user.ID)
	}
}

type configuration struct {
	parallel int
	seed     int64
	output   string
}

func main() {
	args := new(configuration)
	flag.IntVar(&args.parallel, "p", 1, "並列度")
	flag.Int64Var(&args.seed, "seed", time.Now().UnixNano(), "乱数シード（デフォルト: UNIXTIME）")
	flag.StringVar(&args.output, "o", "access.log", "データ出力ファイルパス")
	flag.Parse()
	log.Printf("Degree of Parallelism: %d, Random Seed: %d, Output: %s\n", args.parallel, args.seed, args.output)

	rand.Seed(args.seed)

	var mutex sync.Mutex
	f, err := os.OpenFile(args.output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	SetFootprint(func(applog *AppLog) {
		j, _ := json.Marshal(applog)

		mutex.Lock()
		defer mutex.Unlock()

		f.Write(j)
		f.Write(lineFeed)
	})

	// var wg sync.WaitGroup
	mailbox := make(chan func(), 100)
	for i := 0; i < args.parallel; i++ {
		// wg.Add(1)
		go func() {
			// defer wg.Done()
			for task := range mailbox {
				task()
			}
		}()
	}

	for {
		// ユーザー行動
		user := newUser()
		action := generateUserAction(user)
		mailbox <- action
		// 次のユーザーが来るまでの待ち時間
		wt := generateWeightTime(5000.0, 1000.0)
		time.Sleep(wt)
	}

	// wg.Wait()
}
