package logger

import (
	"fmt"
	"strings"
	"time"
)

func Log(origin string, route []string, data []byte) error {

	print(fmt.Sprintf(`

logger: %s

    origin: %s
     route: %s
      data: %s


`, time.Now().String(), origin, strings.Join(route, "."), string(data)))

	return nil
}
