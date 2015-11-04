package customListener

import(
	"net"
	"errors"
	"time"
)

type StoppableListener struct{
	*net.TCPListener
	stop  				chan int
}

func New(l net.Listener) (*StoppableListener, error){
	tcpl, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("Canot wrap listenner")
	}

	retval := &StoppableListener{}
	retval.TCPListener = tcpl
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener Stoped")

func (sl *StoppableListener) Accept() (net.Conn, error){
	for {
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.TCPListener.Accept()

		select {
		case <- sl.stop:
			return nil, StoppedError
			default:
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			if ok && netErr.Timeout() && netErr.Temporary(){
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}