package protocol

type Reader struct{ p *Protocol }

func (r Reader) ConsumeMessage()       { r.p.ConsumeMessage() }
func (r Reader) Msg(f func(Msg))       { r.p.Msg(f) }
func (r Reader) ReadLn(f func([]byte)) { r.p.ReadLine(f) }

type messageConsumer struct{}

func (messageConsumer) Receive(r Reader) { r.ConsumeMessage() }

func IgnoreOutput() ResFunc { return messageConsumer{}.Receive }

type SimpleStrFunc func(line []byte, ok bool)

var MessageTypeUnexpected = []byte("unexpected message type")

func (sa SimpleStrFunc) Receive(r Reader) {
	r.Msg(func(m Msg) {
		if m != MsgErr && m != MsgSimpleStr {
			sa(MessageTypeUnexpected, true /*err*/)
			return
		}
		r.ReadLn(func(data []byte) { sa(data, m == MsgSimpleStr) })
	})
}

func AcceptStatus(f SimpleStrFunc) ResFunc { return f.Receive }
