package main

import (
    "fmt"
    "log"
    "net/rpc/jsonrpc"
    "os"
)

type Request struct {
    StringInput string
    Budget float32
}

type Response struct {
    TradeId int
    Stocks string
    LeftOver float32
}

type RequestCheck struct {
     TradeId int
}

type ReplyCheck struct {
     StocksC string
     CurrentMarketValue float32
     LeftAmount float32
}


func main() {
    var i int
    fmt.Println("Put 1 to purchase stocks, 2 to check portfolio, anything else to exit")
    fmt.Scanf("%d", &i)
    switch i{
case 1:
    fmt.Println("Enter Stock Key & Stock Percentage(E.g. “GOOG:50%,YHOO:50%”)")
    var input string
    fmt.Scanln(&input)

    fmt.Println("Enter Your Budget")
    var budget float32
    fmt.Scanln(&budget)



    client, err := jsonrpc.Dial("tcp", ":1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // Synchronous call
    args := Request{input, budget}
    var reply Response
    err = client.Call("StockDetail.Stocks", &args, &reply)
    if err != nil {
        log.Fatal("stock error", err)
    }
    fmt.Println(reply.TradeId)
    fmt.Println(reply.Stocks)
    fmt.Println(reply.LeftOver)
break;


          case 2:
                fmt.Println("Enter TradeID")
                var id int
                fmt.Scanln(&id)

                fmt.Println(id)

                client, err := jsonrpc.Dial("tcp", ":1234")
                if err != nil{
                    log.Fatal("dialing: ", err)
                }
                args := RequestCheck{id}
                var reply ReplyCheck
                err = client.Call("StockDetail.CheckPortfolio", &args, &reply)
                if err != nil{
                    log.Fatal("check Error: ", err)
                }
                fmt.Println(reply.StocksC)
                fmt.Println(reply.CurrentMarketValue)
                fmt.Println(reply.LeftAmount)
break;
        default:
            os.Exit(0)

        }

}
