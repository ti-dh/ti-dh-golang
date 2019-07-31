# ti-dh-golang
## DEMO
- go run dh-server.go（go依赖的包只有一个redis，自己解决下...）
- cd example && php client.php（使用php充当客户端）

## API说明
### Init() ( string, string, string, string )
返回p、g、serverNumber、processedServerNumber

### ComputeShareKey( string, string, string ) ( string )
入参p、serverNumber、clientNumber
返回key，也就是最终协商出来对称密钥
