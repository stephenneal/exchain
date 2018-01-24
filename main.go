package main

import (
    "flag"
    //"time"

    "github.com/stephenneal/exchain/microservice"
)

func main() {
    var (
        listen = flag.String("listen", ":8080", "HTTP listen address")
        //proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
    )
    flag.Parse()

    microservice.Start(listen)
	//ticker := time.NewTicker(time.Millisecond * 1000)
    //go func() {
    //    for range ticker.C {
    //    }
    //}()
    //time.Sleep(time.Millisecond * 3000)
    //ticker.Stop()
}
