package migrate

type Dialect string

func (d Dialect) String() string {
	return string(d)
}
