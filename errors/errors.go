package errors

import "errors"

var (
	// ErrBindJSON은 넘어온 JSON 데이터를 map으로 변환하지 못했을 경우 발생함.
	ErrBindJSON = errors.New("cannot bind json")
	// ErrCreateChat는 DB에서 Chat 데이터를 저장을 하지 못했을 경우 발생함.
	ErrCreateChat = errors.New("cannot chat create")
	// ErrGetChat은 DB에서 Chat 데이터를 가져오지 못했을 경우 발생함.
	ErrGetChat = errors.New("cannot get chats")
)
