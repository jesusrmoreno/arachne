package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn, cleanup func()) {
	defer cleanup()
	for {
		// never going to care just need to handle
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	cleaup := GraphState.Subscribe(func(v Graph) {
		ws.WriteJSON(v)
	})

	reader(ws, cleaup)
}

func servePlugins(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	cleaup := GraphState.Subscribe(func(v Graph) {
		ws.WriteJSON(v)
	})

	reader(ws, cleaup)
}

func graphEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(GraphState.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func StartServer(port string, root string) {
	mux := http.NewServeMux()

	// Notes / Media
	contentServer := http.FileServer(http.Dir(root))
	mux.Handle("/content/", http.StripPrefix("/content", contentServer))

	// Websockets
	mux.HandleFunc("/ws", wsEndpoint)
	mux.HandleFunc("/graph", graphEndpoint)

	log.Println("Starting server on :" + port)
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
