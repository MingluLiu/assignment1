package main

import (
    "fmt"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "os"
    "strings"
    "math/rand"
    "strconv"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

var RandomID int

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


type StockDetail struct{}

var hmap map[int]string
var umap map[int]float32

func GetStockPrice(stockKey string) string{
    UrlLeft := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quote%20where%20symbol%20%3D%20%22"
    UrlRight := "%22&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
    url := UrlLeft+stockKey+UrlRight

    r, _ := http.Get(url)

    var quote map[string]interface{}
    body, _ := ioutil.ReadAll(r.Body)
    byt := []byte(body)
    if err := json.Unmarshal(byt, &quote); err!=nil{
        panic(err)
    }

    link1 := quote["query"].(map[string]interface{})
    link2 := link1["results"].(map[string]interface{})
    link3 := link2["quote"].(map[string]interface{})
    link4 := link3["LastTradePriceOnly"].(string)

    return link4
}

func GenerateID() int{
    if RandomID == 0{
        for RandomID == 0{
            RandomID = rand.Intn(100)
        }
    }else{
        RandomID = RandomID + 1
    }
    return RandomID
}


func (t *StockDetail) Stocks(args *Request, reply *Response) error {
    WholeString := args.StringInput
    budget := args.Budget
    replyString := ""
    var leftoverMoney float32

    if strings.Contains(WholeString, ","){
        StringArray := strings.Split(WholeString, ",")
        for i := 0; i < len(StringArray); i++{
            firstString := strings.Split(StringArray[i], ":")
            stockPrice := GetStockPrice(firstString[0])
            priceParse, _ := strconv.ParseFloat(stockPrice, 64)
            price := (float32(priceParse))

            StockArray := strings.Split(firstString[1], "%")
            goParse := StockArray[0]
            percent1, _ := strconv.ParseFloat(goParse, 64)
            percentage1 := (float32(percent1))
            budgetStock := budget * (percentage1 /100)
            if i == 0{
                leftoverMoney = budget - budgetStock
            }else{
                leftoverMoney = leftoverMoney - budgetStock
                if leftoverMoney < 0{
                    os.Exit(1)
                }
            }

            var result1 float32
            result1 = budgetStock / price
            var stockshare int = int(result1)
            leftoverMoney = budgetStock - (price * float32(stockshare)) + leftoverMoney

            fmt.Println(stockshare)
            fmt.Println(leftoverMoney)

            stockShare1 := strconv.Itoa(stockshare)

            reply.TradeId = GenerateID()

            if i == 0{
                replyString = firstString[0]+":"+stockShare1+":$"+stockPrice
            }else{
                replyString = replyString+","+firstString[0]+":"+stockShare1+":$"+stockPrice
            }
            
            reply.LeftOver = leftoverMoney
        }
        hmap[reply.TradeId] = replyString
        umap[reply.TradeId] = reply.LeftOver
    }else{
        oneStockString := strings.Split(WholeString, ":")
        stockP := GetStockPrice(oneStockString[0])
        oneStockPrice, _ := strconv.ParseFloat(stockP, 64)
        stockPrice := (float32(oneStockPrice))
        
        PersentageString := strings.Split(oneStockString[1],"%")
        PercentageParse := PersentageString[0]
        percent, _ := strconv.ParseFloat(PercentageParse, 64)
        percentage := (float32(percent))
        
        budget_stock := budget * (percentage / 100)
        residualValue := budget - budget_stock
        var result float32 
        result = budget_stock / stockPrice
        var share int = int(result) 
        residualValue = budget_stock - (stockPrice * float32(share)) + residualValue

        fmt.Println(share)
        fmt.Println(residualValue)

        stockShare := strconv.Itoa(share)
        replyString = oneStockString[0]+":"+stockShare+":"+"$"+stockP
        fmt.Println(replyString)

        reply.LeftOver = residualValue
        reply.TradeId = GenerateID()

        hmap[reply.TradeId] = replyString
        umap[reply.TradeId] = reply.LeftOver
    }
   reply.Stocks = replyString
    
    return nil
}

func (t *StockDetail) CheckPortfolio(args *RequestCheck, reply *ReplyCheck) error{
    getContent := ""
    stringReply := ""
    var currentMarketP float32
    var total float32
    currentMarketP = 0
    total = 0
    checkID := args.TradeId
    for key, value := range hmap{
        if key == checkID{
            getContent = value
            if strings.Contains(getContent, ","){
                StringArrayC := strings.Split(getContent, ",")
                for i := 0; i < len(StringArrayC); i++{
                    firstStringC := strings.Split(StringArrayC[i], ":")
                    stockName := firstStringC[0]
                    stockShareC := firstStringC[1]
                    intShare, _ := strconv.Atoi(stockShareC)
                    purchasePrice := firstStringC[2]
                    purchasePrice1 := strings.Split(purchasePrice, "$")
                    purchasePrice64, _ := strconv.ParseFloat(purchasePrice1[1], 64)
                
                    currentPrice32 := (float32(purchasePrice64))
                    currentPrice := GetStockPrice(stockName)
                    currentPriceC, _ := strconv.ParseFloat(currentPrice, 64)
                    currentQuote := (float32(currentPriceC))

                    currentMarketP = currentQuote * (float32(intShare))
                    total = total + currentMarketP
                    compare := currentQuote - currentPrice32

                    if i == 0{
                        if compare > 0{
                            stringReply = stockName+":"+stockShareC+":+$"+currentPrice
                        }else if compare < 0{
                            stringReply = stockName+":"+stockShareC+":-$"+currentPrice
                        }else if compare == 0{
                            stringReply = stockName+":"+stockShareC+":$"+currentPrice
                        }
                    }else{
                        if compare > 0{
                            stringReply = stringReply + "," + stockName+":"+stockShareC+":+$"+currentPrice
                        }else if compare < 0{
                            stringReply = stringReply + "," + stockName+":"+stockShareC+":-$"+currentPrice
                        }else if compare == 0{
                            stringReply = stringReply + "," + stockName+":"+stockShareC+":$"+currentPrice
                        }
                    }
                }
                reply.StocksC = stringReply
                reply.CurrentMarketValue = total
            }else{
                singleStockString := strings.Split(getContent,":")
                singleStockName := singleStockString[0]
                SingleStockShare := singleStockString[1]
                intShare, _ := strconv.Atoi(SingleStockShare)
                purchasePrice := singleStockString[2]
                purchasePrice1 := strings.Split(purchasePrice, "$")
                fpurchasePrice, _ := strconv.ParseFloat(purchasePrice1[1], 64)
                singleCurrentPrice := (float32(fpurchasePrice))

                currentPrice := GetStockPrice(singleStockName)
                fcurrentPrice, _ := strconv.ParseFloat(currentPrice, 64)
                currentQuote := (float32(fcurrentPrice))

                calculate := currentQuote - singleCurrentPrice
                
                currentMarketP = currentQuote * (float32(intShare))
                total = total + currentMarketP

                if calculate > 0{
                    stringReply = singleStockName+":"+SingleStockShare+":+$"+currentPrice
                }else if calculate < 0{
                    stringReply  = singleStockName+":"+SingleStockShare+":-$"+currentPrice
                }else{
                    stringReply  = singleStockName+":"+SingleStockShare+":$"+currentPrice
                }
                reply.StocksC = stringReply
                reply.CurrentMarketValue = total
            }
        }else{
            os.Exit(1)
        }
    }

    for key1, value1 := range umap{
        if(key1 == checkID){
            reply.LeftAmount = value1
        }else{
            os.Exit(1)
        }
    }
    return nil
}







func main() {
    RandomID = 0
    hmap = make(map[int]string)
    umap = make(map[int]float32)
    StockDetail := *(new(StockDetail))
    rpc.Register(&StockDetail)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }

}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
