package account

type mockRedisListReader struct {
	KeyGot string
	Err    error
	List   []string
}

func (lr *mockRedisListReader) LRange(key string, start int64, stop int64) ([]string, error) {
	lr.KeyGot = key
	if lr.Err != nil {
		return []string{}, lr.Err
	}
	if (start == 0) && (stop == -1) {
		return lr.List, nil
	}
	if (start > int64(len(lr.List))) || ((stop + 1) > int64(len(lr.List))) {
		return []string{}, nil
	}
	return lr.List[start : stop+1], nil
}

func (lr mockRedisListReader) Quit() (string, error) {
	return "", nil
}
