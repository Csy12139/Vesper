package sender

import (
	"fmt"
	"io"
)

func (s *Session) readFile() {
	s.readingStats.Start()
	defer func() {
		s.readingStats.Pause()
		close(s.output)
	}()
	for {
		s.dataBuff = s.dataBuff[:cap(s.dataBuff)]
		n, err := s.stream.Read(s.dataBuff)
		if err != nil {
			if err == io.EOF {
				s.readingStats.Stop()
				return
			}
			return
		}
		s.dataBuff = s.dataBuff[:n]
		s.readingStats.AddBytes(uint64(n))
		s.output <- outputMsg{
			n:    n,
			buff: append([]byte(nil), s.dataBuff...),
		}
	}
}
func (s *Session) onBufferedAmountLow() func() {
	return func() {
		data := <-s.output
		if data.n != 0 {
			s.msgToBeSent = append(s.msgToBeSent, data)
		} else if len(s.msgToBeSent) == 0 && s.dataChannel.BufferedAmount() == 0 {
			s.sess.NetworkStats.Stop()
			s.close(false)
			return
		}
		currentSpeed := s.sess.NetworkStats.Bandwidth()
		fmt.Printf("Transferring at %.2f MB/s\r", currentSpeed)

		for len(s.msgToBeSent) != 0 {
			cur := s.msgToBeSent[0]

			if err := s.dataChannel.Send(cur.buff); err != nil {
				return
			}
			s.sess.NetworkStats.AddBytes(uint64(cur.n))
			s.msgToBeSent = s.msgToBeSent[1:]
		}
	}
}
func (s *Session) writeToNetwork() {
	s.dataChannel.OnBufferedAmountLow(s.onBufferedAmountLow())
	<-s.stopSending
	s.dataChannel.OnBufferedAmountLow(nil)
	s.sess.NetworkStats.Pause()
}
