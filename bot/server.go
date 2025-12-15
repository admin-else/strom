package bot

import "net"

func ServeClient(cNet net.Conn, factory func(c *Conn) (h HandlerInst, err error)) (err error) {
	c := Servee(cNet)
	h, err := factory(c)
	if err != nil {
		return
	}
	err = c.Start(h)
	return
}

func StartServer(listenAddr string, factory func(c *Conn) (h HandlerInst, err error)) (err error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	var cNet net.Conn
	for {
		cNet, err = l.Accept()
		if err != nil {
			return
		}
		go ServeClient(cNet, factory)
	}
}
