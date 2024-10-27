package types

type Accounts []Account

func (a Accounts) GetById(in string) Account {
	for _, account := range a {
		if account.Id == in {
			return account
		}
	}
	return Account{}
}

type Account struct {
	Id      string `json:"id"`
	Profile string `json:"profile"`
	Alias   string `json:"alias"`
}
