package main
import (
  "io"
  "log"
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/garyburd/redigo/redis"
  "Dh"
)

type DHInitStruct struct {
  P string `json:"p"`
  G string `json:"g"`
  ProcessedServerNumber string `json:"processedServerNumber"`
}

func main() {
  http.HandleFunc( "/init", Init )
  http.HandleFunc( "/compute", Compute )
  err := http.ListenAndServe( ":6666", nil )
  if err != nil {
    log.Fatal( "Listen&Serve : ", err )
  }
}

func Init( w http.ResponseWriter, r *http.Request ) {
  p, g, serverNumber, processedServerNumber := Dh.Init()
  // 写入redis数据
  redisConn, redisErr := redis.Dial( "tcp", "127.0.0.1:6379" )
  if redisErr != nil {
    fmt.Println( "Connect to redis error", redisErr )
    return
  }
  defer redisConn.Close()
  _, DoErr := redisConn.Do( "SET", "p", p )
  if DoErr != nil {
    fmt.Println( "redis set failed:", DoErr )
  }
  _, DoErr = redisConn.Do( "SET", "serverNumber", serverNumber )
  if DoErr != nil {
    fmt.Println( "redis set failed:", DoErr )
  }
  // 将数据marshal成json字符串吐给客户端.
  dhInitRawData := DHInitStruct{ p, g, processedServerNumber }
  dhInitData, _ := json.Marshal( dhInitRawData )
  // 返回给客户端数据
  io.WriteString( w, string( dhInitData ) )
}

func Compute( w http.ResponseWriter, r *http.Request ) {
  // 从post中获取client_number
  sClientNumber := r.PostFormValue("clientNumber")
  // 从redis中获取 serverNumber 和 p
  redisConn, redisErr := redis.Dial( "tcp", "127.0.0.1:6379" )
  if redisErr != nil {
    fmt.Println( "Connect to redis error", redisErr )
    return
  }
  defer redisConn.Close()
  sP, doErr := redis.String( redisConn.Do( "GET", "p" ) )
  if doErr != nil {
    fmt.Println( "redis get failed:", doErr )
  }
  sServerNumber, doErr := redis.String( redisConn.Do( "GET", "serverNumber" ) )
  if doErr != nil {
    fmt.Println( "redis get failed:", doErr )
  }
  // 调用Dh包中的方法
  iKey := Dh.ComputeShareKey( sP, sServerNumber, sClientNumber )
  fmt.Println( "key:", iKey )
  // 这个key便是计算出出来的用于对称加解密的公钥
  io.WriteString( w, "compute\n" )
}
