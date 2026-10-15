package noxnet

import (
	"fmt"
	"math/rand"
	"testing"

	_ "github.com/brutusmcforce/opennox/libs/noxnet/discover"
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	_ "github.com/brutusmcforce/opennox/libs/noxnet/netxfer"
)

// Messages decode random data without panics: they come from the network.
func TestDecodeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	seen := map[string]bool{}
	for i := 0; i < 30000; i++ {
		n := 1 + r.Intn(64)
		if r.Intn(10) == 0 {
			n = 1 + r.Intn(600)
		}
		b := make([]byte, n)
		for j := range b {
			switch r.Intn(4) {
			case 0:
				b[j] = 0
			case 1:
				b[j] = 0xff
			default:
				b[j] = byte(r.Intn(256))
			}
		}
		for _, cli := range []bool{false, true} {
			st := &netmsg.State{}
			st.IsClient = cli
			func() {
				defer func() {
					if e := recover(); e != nil {
						k := fmt.Sprint(netmsg.Op(b[0]))
						if !seen[k] {
							seen[k] = true
							t.Errorf("op %v (cli=%v) data=%x: %v", netmsg.Op(b[0]), cli, b, e)
						}
					}
				}()
				data := b
				for len(data) > 0 {
					b = data
					_, m, err := st.DecodeNext(data)
					if err != nil || m <= 0 {
						return
					}
					if m > len(data) {
						t.Errorf("op %v: decoded %d of %d", netmsg.Op(data[0]), m, len(data))
						return
					}
					data = data[m:]
				}
			}()
		}
	}
}
