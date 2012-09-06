package count

type TokensMap map[string]int

func (tokensMap *TokensMap) Add(token string) {
	tm := *tokensMap
	tm[token] = tm[token] + 1
	*tokensMap = tm
}
