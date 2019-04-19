package proto

type ILogin struct {
	AID string
	Platform string
}

type OLogin struct {
	PID string
	SID string
}

type IGetPlayer struct {
	PID string
	SID string
}

type OGetPlayer struct {
	Players []Player
}
