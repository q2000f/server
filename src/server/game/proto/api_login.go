package proto

type ILogin struct {
	AID      string
	Platform string
}

type OLogin struct {
	PID string
	SID uint64
}

type IGetPlayer struct {
}

type OGetPlayer struct {
	Players []Player
}

type ICreatePlayer struct {
	Name string
}

type OCreatePlayer struct {
	PlayerID uint64
	Name     string
}
