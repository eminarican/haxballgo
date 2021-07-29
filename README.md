<h1 align="center">Haxballgo</h1>
<p align="center"><strong>Haxball headless API wrapper for Go</strong></p>

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

```go
func main() {
  r := room.New()
  defer r.Shutdown()
  
  l := r.Logger()
  s := r.Scheduler()

  r.OnPlayerJoin(func(p *room.Player) {
	  l.Info("A player joined!")
  })

  r.OnPlayerLeave(func(p *room.Player) {
	  l.Info("A player leaved!")
  })

  s.Repeating(time.Second, func(){
    r.Announce("Test message")
  })
}
```

```json
# auto-generated config.json
{
  "Bot": {
    "Active": true,
    "Name": "Bot"
  },
  "General": {
    "Name": "My Room",
    "Token": "",
    "Debug": false,
    "MaxPlayer": 16
  },
  "Security": {
    "Public": false,
    "Password": ""
  }
}
```
