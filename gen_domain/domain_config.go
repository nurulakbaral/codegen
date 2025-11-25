package domain

type Dir struct {
	Input  string
	Output string
}

type Config struct {
	Root   string
	Domain string
	Force  bool
	Dirs   []Dir
	Data   any
}
