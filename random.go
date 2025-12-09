package tools

import (
	"math/rand"
	"time"
)

var (
	rnd     = rand.New(rand.NewSource(time.Now().UnixNano()))
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	numbers = []rune("0123456789")
	ls      = len(letters)
	ns      = 10
)

func RandomAlphaNum(length int) string {
	s := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		s = append(s, letters[rnd.Intn(ls)])
	}
	return string(s)
}

func RandomNumChar(length int) string {
	s := make([]rune, 0, length)
	for i := 0; i < length; i++ {
		s = append(s, numbers[rnd.Intn(ns)])
	}
	return string(s)
}
