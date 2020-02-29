package poller

type MessageFilter interface {
	IsAllowed(update Update) bool
}

type FilterMessageByUser struct {
	Users []int
}

func (f FilterMessageByUser) IsAllowed(update Update) bool {
	u := update.Message.From.ID
	for _, o := range f.Users {
		if o == u {
			return true
		}
	}
	return false
}