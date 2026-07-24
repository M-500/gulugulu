package xcode

type XCode interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
}
