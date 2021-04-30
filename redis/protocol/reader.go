package protocol

type Reader struct{ p *Protocol }

func (r Reader) ConsumeMessage()         { r.p.ConsumeMessage() }
func (r Reader) Consume(m Msg)           { r.p.Consume(m) }
func (r Reader) Msg(f func(Msg))         { r.p.Msg(f) }
func (r Reader) ReadLine(f func([]byte)) { r.p.ReadLine(f) }
func (r Reader) ReadInt(f func(int))     { r.p.ReadInt(f) }

type messageConsumer struct{}

func (messageConsumer) Receive(r Reader) { r.ConsumeMessage() }

func IgnoreOutput() ResFunc { return messageConsumer{}.Receive }

type SimpleStrFunc func(line []byte, ok bool)

var MessageTypeUnexpected = []byte("unexpected message type")

func (ssfn SimpleStrFunc) Receive(r Reader) {
	r.Msg(func(m Msg) {
		if m != MsgErr && m != MsgSimpleStr {
			r.Consume(m)
			ssfn(MessageTypeUnexpected, true /*err*/)
			return
		}
		r.ReadLine(func(data []byte) { ssfn(data, m == MsgSimpleStr) })
	})
}

type IntFunc func(n int, err []byte)

func (ifn IntFunc) Receive(r Reader) {
	r.Msg(func(m Msg) {
		switch {
		case m == MsgErr:
			r.ReadLine(func(data []byte) { ifn(0, data) })
		case m != MsgInt:
			r.Consume(m)
			ifn(0, MessageTypeUnexpected)
		case m == MsgInt:
			r.ReadInt(func(n int) { ifn(n, nil) })
		}
	})
}

func AcceptStatus(f SimpleStrFunc) ResFunc { return f.Receive }
func AcceptInt(f IntFunc) ResFunc          { return f.Receive }
