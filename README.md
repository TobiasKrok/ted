# ted - toby's editor

just to learn data structures and such

Plan:

- add, delete and replace
- search?
- highlighting module?
- gap buffer implementation
- search

References:

- https://www.gnu.org/software/emacs/manual/html_node/elisp/Buffer-Gap.html
- https://doc.cat-v.org/plan_9/4th_edition/papers/sam/
- https://workdad.dev/posts/building-a-text-editor-the-gap-buffer/
- https://viewsourcecode.org/snaptoken/kilo/index.html
- https://www.man7.org/linux/man-pages/man3/termios.3.html
- https://cs.opensource.google/go/x/sys/+/refs/tags/v0.41.0:unix/ioctl_unsigned.go;l=48

https://vt100.net/docs/vt510-rm/DECSCUSR.html

CURSOR COMMANDS
https://en.wikipedia.org/wiki/ANSI_escape_code#Control_Sequence_Introducer_commands

## gap buffer

for line numbers, maybe we can have an array of lines and in normal mode they are normal arrays, but in insert mode we make it a gap buffer
We can for example then make the x lines under and above gap buffers as it would be plausable that the user goes wants ot edit the closst line more often
