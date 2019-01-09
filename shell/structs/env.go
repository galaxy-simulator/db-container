package structs

// Env structure storing the variables used to define what stars are inserted where
type Env struct {
	// general values
	url       string
	data      string
	amount    int64
	treeindex int64
}

func NewEnv(url string, data string, amount int64, treeindex int64) *Env {
	return &Env{url: url, data: data, amount: amount, treeindex: treeindex}
}

// Getters and setters
func (e *Env) Treeindex() int64 {
	return e.treeindex
}
func (e *Env) SetTreeindex(treeindex int64) {
	e.treeindex = treeindex
}

func (e *Env) Amount() int64 {
	return e.amount
}
func (e *Env) SetAmount(amount int64) {
	e.amount = amount
}

func (e *Env) Data() string {
	return e.data
}
func (e *Env) SetData(data string) {
	e.data = data
}

func (e *Env) Url() string {
	return e.url
}
func (e *Env) SetUrl(url string) {
	e.url = url
}
