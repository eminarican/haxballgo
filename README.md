<p align="center"><img src=".github/assets/logo.png" width="200px" alt="logo"/></p>
<h1 align="center">Haxballgo</h1>
<p align="center"><strong>Haxball room API for Go</strong></p>

<p align="center">
  <a href="https://opensource.org/licenses/gpl-3.0.html">
    <img alt="License" src="https://img.shields.io/github/license/eminarican/haxballgo?color=success&style=for-the-badge">
  </a>

  <a href="https://github.com/eminarican/haxballgo/issues">
    <img alt="GitHub Issues" src="https://img.shields.io/github/issues/eminarican/haxballgo?style=for-the-badge">
  </a>

  <a href="https://github.com/eminarican/haxballgo/stargazers">
    <img alt="GitHub Stars" src="https://img.shields.io/github/stars/eminarican/haxballgo?style=for-the-badge">
  </a>
</p>

## 💡 Simple usage

```
go get github.com/eminarican/haxballgo@0.0.1
```

```go
func main() {
  r := room.New()
  defer r.Shutdown()
  
  l := r.Logger()
  s := r.Scheduler()

  r.OnPlayerJoin(func(p *room.Player) {
	  l.Infof("Player %v joined!", p.Name())
  })

  r.OnPlayerLeave(func(p *room.Player) {
	  l.Infof("Player %v leaved!", p.Name())
  })

  s.Repeating(time.Second, func(stop func()){
    r.Announce("Test message")
  })
}
```

```json
# auto-generated config.json
{
  "Bot": {
    "Active": false,
    "Name": "Bot"
  },
  "General": {
    "Name": "My Room",
    "Token": "",
    "MaxPlayer": 16
  },
  "Security": {
    "Public": true,
    "Password": ""
  },
  "Logging": {
    "Debug": false,
    "Pretty": true
  }
}
```
