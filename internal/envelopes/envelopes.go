package envelopes

type Envelope map[string]any

func New(name string, payload any) *Envelope {
	return &Envelope{
		name: payload,
	}
}

func NewPack(payload map[string]any) *Envelope {
	e := Envelope{}
	for k, v := range payload {
		e[k] = v
	}

	return &e
}
