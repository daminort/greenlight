package queries

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	values url.Values
}

func New(values url.Values) *Query {
	return &Query{
		values: values,
	}
}

func (q *Query) ReadString(key, defaultValue string) string {
	s := q.values.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (q *Query) ReadInt(key string, defaultValue int) int {
	s := q.values.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	if i < 1 {
		return defaultValue
	}

	return i
}

func (q *Query) ReadStrings(key string, defaultValue []string) []string {
	s := q.values.Get(key)
	if s == "" {
		return defaultValue
	}

	return strings.Split(s, ",")
}
