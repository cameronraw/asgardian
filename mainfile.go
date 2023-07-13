package asgardian

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
  app := CreateApplication(CreateConfig(), log.New(os.Stdout, "", log.Ldate|log.Ltime))
  mux := app.CreatePreconfiguredMux()
  
  srv := &http.Server{
    Addr: app.CreateAddr(),
    Handler: mux,
    IdleTimeout: 120 * time.Second,
    ReadTimeout: 1 * time.Second,
  }

  err := srv.ListenAndServe()
  if err != nil {
    log.Fatal(err)
  }
}


