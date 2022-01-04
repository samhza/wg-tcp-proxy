# wg-tcp-proxy
`go install samhza.com/wg-tcp-proxy@latest`

---
```
Usage of wg-tcp-proxy:
  -addr value
    	address to listen on
  -endpoint string
    	endpoint of peer
  -privkey string
    	our private key
  -pubkey string
    	public key of peer
  -target value
    	address to reverse proxy to
  -v	verbose output
```
Example:
```
wg-tcp-proxy \
  -addr 10.0.0.1:80 -target 0.0.0.0:80
  -privkey yJnLs2Kznd1hmun71Z7oMsvvFa9KyaYOW1BP2MHzjnk= \
  -pubkey W4lY5vXro39JG9NJz8GthHrvOsFKgY1Uf1VGWWzxaSw= \
  -endpoint 172.16.1.1:51820
```
