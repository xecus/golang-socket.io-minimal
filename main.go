package main

import (
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

type customServer struct {
	Server *socketio.Server
}

func (s *customServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	s.Server.ServeHTTP(w, r)
}

func main() {
	ioServer := SocketIoServer()
	wsServer := new(customServer)
	wsServer.Server = ioServer

	http.Handle("/socket.io/", wsServer)

	port := "8080"
	log.Println("[Main] Starting Server Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func SocketIoServer() *socketio.Server {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		println(so.Id() + " joined clients.")
		so.Join("chat")

		so.On("msg", func(msg string) {
			log.Println("[" + msg + "]")
			so.Emit("hoge", msg)
			// log.Println("emit:", so.Emit("chat message", msg))
			// so.BroadcastTo("chat", "chat message", msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server
}
