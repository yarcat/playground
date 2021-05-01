package protocol

type Reader struct{ p *Protocol }

func (r Reader) ConsumeMessage()                  { r.p.ConsumeMessage() }
func (r Reader) Consume(m Msg)                    { r.p.Consume(m) }
func (r Reader) Msg(f func(Msg))                  { r.p.Msg(f) }
func (r Reader) ReadLine(f func([]byte))          { r.p.ReadLine(f) }
func (r Reader) ReadInt(f func(int))              { r.p.ReadInt(f) }
func (r Reader) ReadBulk(f func(int, BulkReader)) { r.p.ReadBulk(f) }

type messageConsumer struct{}

func (messageConsumer) Receive(r Reader) { r.ConsumeMessage() }

func IgnoreOutput() ResFunc { return messageConsumer{}.Receive }

type SimpleStrFunc func(line []byte, ok bool)

var (
	MsgTypeUnexpected = []byte("unexpected message type")
	MsgKeyNotFound    = []byte("key not found")
)

func (ssfn SimpleStrFunc) Receive(r Reader) {
	r.Msg(func(m Msg) {
		if m != MsgErr && m != MsgSimpleStr {
			r.Consume(m)
			ssfn(MsgTypeUnexpected, true /*err*/)
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
			ifn(0, MsgTypeUnexpected)
		case m == MsgInt:
			r.ReadInt(func(n int) { ifn(n, nil) })
		}
	})
}

type BulkStrFunc func([]byte, int, BulkReader)

func (bsf BulkStrFunc) Receive(r Reader) {
	r.Msg(func(m Msg) {
		switch {
		case m == MsgErr:
			r.ReadLine(func(data []byte) { bsf(data, 0, BulkReader{}) })
			return
		case m != MsgBulkStr:
			r.Consume(m)
			bsf(MsgTypeUnexpected, 0, BulkReader{})
			return
		}
		r.ReadBulk(func(n int, r BulkReader) {
			if n < 0 {
				bsf(MsgKeyNotFound, 0, BulkReader{})
				return
			}
			bsf(nil, n, r)
		})
	})
}

func AcceptStatus(f SimpleStrFunc) ResFunc { return f.Receive }
func AcceptInt(f IntFunc) ResFunc          { return f.Receive }
func AcceptBulk(f BulkStrFunc) ResFunc     { return f.Receive }
