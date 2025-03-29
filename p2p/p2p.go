package p2p

func SendData(mnAddr string, sourceUUID string, targetUUID string, data []byte) error {
	s := newSession(mnAddr)
	defer s.ctxCancel()
	//sdpOffer, candidates, err := s.initP2PConnection()
	//if err != nil {
	//	return err
	//}
	//TODO exchange SDP and candidate
	//exchangeSDPCandidates(sourceUUID, targetUUID, sdpOffer, candidates)

	return nil
}
func ReceiveData(sourceUUID string, targetUUID string, data []byte) error {
	return nil
}
