package term

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	original *unix.Termios
	termios  *unix.Termios
	fd       uintptr
	Size     terminalSize
}
type terminalSize struct {
	Cols int
	Rows int
}

func NewTerminal() (*Terminal, error) {

	fd := os.Stdin.Fd()
	// unix.TIOCGETA is mac only?
	term, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	if err != nil {
		return nil, err
	}
	original := *term
	ws, err := getWindowSize(int(fd))
	if err != nil {
		return nil, err
	}
	return &Terminal{
		original: &original,
		termios:  term,
		fd:       fd,
		Size: terminalSize{
			Cols: ws.Cols,
			Rows: ws.Rows,
		},
	}, nil

}

func (t *Terminal) EnableRawMode() error {

	// original := *term
	t.termios.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	t.termios.Oflag &^= unix.OPOST
	t.termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.termios.Cflag &^= unix.CSIZE | unix.PARENB
	t.termios.Cflag |= unix.CS8

	if err := unix.IoctlSetTermios(int(t.fd), unix.TIOCSETA, t.termios); err != nil {
		return err
	}

	return nil
}

func (t *Terminal) DisableRawMode() error {
	err := unix.IoctlSetTermios(int(t.fd), unix.TCIOFLUSH, t.original)
	return err
}

func (t *Terminal) Clear() {
	fmt.Print("\x1b[2J\x1b[H") // refresh and put the cursor at home
}

func getWindowSize(fd int) (terminalSize, error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return terminalSize{}, err
	}
	return terminalSize{
		Cols: int(ws.Col),
		Rows: int(ws.Row),
	}, nil
}
