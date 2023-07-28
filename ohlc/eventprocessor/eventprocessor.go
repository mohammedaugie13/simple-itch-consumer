package eventprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"io"
	"log"
	"ohlc/models"
	"strings"
)

func ProcessEvent(event *kgo.Record) {
	d := json.NewDecoder(strings.NewReader(string(event.Value)))
	for {
		var v models.EventMessage
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		fmt.Println(v)
	}
}
