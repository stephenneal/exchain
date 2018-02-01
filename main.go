package main

import (
    "os"
    //"time"

    "github.com/stephenneal/exchain/microservice"

    "github.com/go-kit/kit/log"
)

func main() {
    logger := log.NewLogfmtLogger(os.Stderr)

    port := os.Getenv("PORT")
    logger.Log("$PORT", port)
    if port == "" {
        logger.Log("message", "$PORT must be set")
        os.Exit(1)
    } else {
        microservice.Start(port)
    }

	//ticker := time.NewTicker(time.Millisecond * 1000)
    //go func() {
    //    for range ticker.C {
    //    }
    //}()
    //time.Sleep(time.Millisecond * 3000)
    //ticker.Stop()
}
