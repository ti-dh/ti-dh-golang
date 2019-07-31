package Dh
import (
  "fmt"
  "time"
  "math/big"
  "math/rand"
  //"encoding/json"
)

/*
 * @desc : 定义p、g、ProcessedServerNumber结构体.
 */
/*
type DHInitStruct struct {
  P string `json:p`
  G string `json:g`
  ProcessedServerNumber string `json:processedServerNumber`
}
*/

/*
 * @desc : 初始化
 */
func Init() ( string, string, string, string ) {
  // 计算出p
  var pBaseHex string = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200CBBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
  p         := new( big.Int )
  p, _      = p.SetString( pBaseHex, 16 )
  pBase     := new( big.Int )
  pBase, _  = pBase.SetString( pBaseHex, 16 )
  p.Sub( p, big.NewInt( 1 ) )
  // 计算g
  g := new( big.Int )
  for {
    g.Rand( rand.New( rand.NewSource( time.Now().Unix() ) ), p )
    gFlag := new( big.Int )
    gFlag.Exp( g, p, pBase )
    if 0 == gFlag.Cmp( big.NewInt( 1 ) ) {
      break
    }
  }
  // 随机出server_number
  serverNumber  := new( big.Int )
  serverRandMax := big.NewInt( 10000000 )
  serverNumber.Rand( rand.New( rand.NewSource( time.Now().Unix() ) ), serverRandMax )
  // 最终的processed_server_number
  processedServerNumber  := new( big.Int )
  processedServerNumber   = processedServerNumber.Exp( g, serverNumber, pBase )
  // 返回给客户端数据
  return pBase.String(), g.String(), serverNumber.String(), processedServerNumber.String()
}

/*
 * @desc : 计算共享对称密钥.
 */
func ComputeShareKey( sP string, sServerNumber string, sClientNumber string ) ( string ) {
  // 将clientNumber serverNumber p 转为big.Int
  iClientNumber := new( big.Int )
  iClientNumber.SetString( sClientNumber, 10 )
  iServerNumber := new( big.Int )
  iServerNumber.SetString( sServerNumber, 10 )
  iP := new( big.Int )
  iP.SetString( sP, 10 )
  // 接受client_number（实际上是经过了client客户端处理过后的client_number）
  // 利用client_number,server_number和p计算出公共密钥key
  iKey := new( big.Int )
  iKey.Exp( iClientNumber, iServerNumber, iP )
  return iKey.String()
}
