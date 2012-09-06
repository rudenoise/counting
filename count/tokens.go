package count

type TokensMap map[string]int

func (tm TokensMap) Add(token string) {
	tm[token] = tm[token] + 1
}
