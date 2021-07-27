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
r := room.New()

r.OnPlayerJoin(func(p *room.Player) {
	println("a player joined!")
})

r.OnPlayerLeave(func(p *room.Player) {
	println("a player leaved!")
})

println("room link:", r.Link())
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
    "MaxPlayer": 16
  },
  "Security": {
    "Public": false,
    "Password": ""
  }
}
```
