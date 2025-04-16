package main

import (
	"context"
	"github.com/joho/godotenv"
	"kc-backend/judge/orchestrator"
	"kc-backend/judge/state"
	"kc-backend/judge/types"
	"log"
	"sync"
)

func JudgeTest() {
	//container := docker.Docker{}
	//err := container.StartContainer(context.Background())
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	state.PendingJobs = &state.PendingJobs_t{
		Jobs:  make(map[string]types.Job),
		Mutex: sync.Mutex{},
	}

	o := orchestrator.Orchestrator{}
	o.Start(context.Background(), state.PendingJobs)
	o.Run()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//yjsServer := yjs_server.InitYjsServer()
	//syncServer := sync_server.InitSyncServer()
	//http.HandleFunc("/", yjsServer.HandleWebSocket)
	//http.HandleFunc("/sync", syncServer.HandleWebSocket)
	//log.Println("Yjs WebSocket server started on :1234")
	//log.Fatal(http.ListenAndServe(":1234", nil))
	JudgeTest()
}
