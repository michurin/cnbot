package processors

type MessageProcessor func(text string) ([]string, error)
