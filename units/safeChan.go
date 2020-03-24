package units

//用来关闭通知类的通道
func SafeCloseChan(ch chan struct{}) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()

	// assume ch != nil here.
	close(ch) // panic if ch is closed
	return true
}


