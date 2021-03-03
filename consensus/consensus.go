package consensus

type Consensus interface {
FindNonce() int64
}

func NewPow() Consensus {


}
