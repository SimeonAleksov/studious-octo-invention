package main

import (
  "log"
  "net/http"
  "github.com/rs/cors"
  "github.com/jackc/pgx/v4"
  socketio "github.com/googollee/go-socket.io"
)

func main() {
     c := cors.New(cors.Options{
         AllowedOrigins: []string{"*"},
         AllowCredentials: true,
   })

  server := socketio.NewServer(nil)

  server.OnConnect("/", func(s socketio.Conn) error {
  s.SetContext("")
  log.Println("connected:", s.ID())
  return nil
  })

  server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
  log.Println("notice:", msg)
  s.Emit("reply", "have "+msg)
  })

  server.OnEvent("/", "bye", func(s socketio.Conn) string {
  last := s.Context().(string)
  s.Emit("bye", last)
  s.Close()
    return last
  })

  server.OnError("/", func(s socketio.Conn, e error) {
    log.Println("meet error:", e)
  })

  server.OnDisconnect("/", func(s socketio.Conn, reason string) {
    log.Println("closed", reason)
  })

  go func() {
    if err := server.Serve(); err != nil {
      log.Fatalf("socketio listen error: %s\n", err)
    }
  }()
  defer server.Close()

  http.Handle("/socket.io/", c.Handler(server))
  http.Handle("/", http.FileServer(http.Dir("../asset")))

  log.Println("Serving at localhost:5000...")
  log.Fatal(http.ListenAndServe(":5000", nil))
}
