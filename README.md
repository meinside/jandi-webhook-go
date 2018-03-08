# jandi-webhook-go

Jandi's webhook library for golang.

## Sample

### Send an incoming webhook

```go
package main

import (
	"log"

	"github.com/meinside/jandi-webhook-go"
)

const (
	webhookURL = "https://wh.jandi.com/connect-api/webhook/000000/abcd1234567890ef"
)

func main() {
	client := jandi.NewIncomingClient(webhookURL)
	client.SetVerbose(true)

	if txt, err := client.SendIncoming(
		"Some text",
		"#FF0000",
		jandi.ConnectInfoFrom("Sample program", "Sent from this sample program.", ""),
	); err == nil {
		log.Printf(">>> success: %s\n", txt)
	} else {
		log.Printf(">>> failure: %s (%s)\n", txt, err)
	}
}
```

## License

MIT

