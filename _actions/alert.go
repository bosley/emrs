package alert

import (
  "fmt"
  "errors"
)

func Sms(origin string, route []string, data []byte) error {

  fmt.Println("alert::sms", "from:", origin)

  if len(route) < 2 {
    fmt.Println("no specific routes given, dumping to stdout")
    fmt.Println("alert::data", string(data))
    return nil
  }

  level := route[0]
  provider := route[1]

  fmt.Println("using", level, "as the alert level, and", provider, "as the provider")

  switch level {
  case "info":
  case "warn":
  case "error":
  case "debug":
    break
  default:
    return errors.New(fmt.Sprintf("unknown level given for alert request:%s", level))
  }

  switch provider {
  case "twilio":
    return alertViaTwilio(origin, level, data)
  }

  return errors.New(fmt.Sprintf("unknown provider given for alert request:%s", provider))
}

func alertViaTwilio(origin string, level string, data []byte) error {

  fmt.Println("actions::alert::TODO: Fill in alertViaTwilio\n", origin, level, string(data))

  return errors.New("not yet implemented")
}
