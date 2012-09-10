package count

type TokensMap map[string]int

func (tokensMap *TokensMap) Add(token string) {
	tm := *tokensMap
	tm[token] = tm[token] + 1
	*tokensMap = tm
}

func (tokensMap TokensMap) ToSlice() TokenSlice {
	var ts TokenSlice
	for token, count := range tokensMap {
		ts = append(ts, TokenInfo{token, count})
	}
	return ts
}

type TokenInfo struct {
	Token string
	Count int
}

type TokenSlice []TokenInfo

func (t TokenSlice) Len() int           { return len(t) }
func (t TokenSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TokenSlice) Less(i, j int) bool { return t[i].Count < t[j].Count }

type TokenSliceByCountAsc struct{ TokenSlice }

func (t TokenSliceByCountAsc) Less(i, j int) bool {
	return t.TokenSlice[i].Count < t.TokenSlice[j].Count
}

type TokenSliceByCountDesc struct{ TokenSlice }

func (t TokenSliceByCountDesc) Less(i, j int) bool {
	return t.TokenSlice[i].Count > t.TokenSlice[j].Count
}
