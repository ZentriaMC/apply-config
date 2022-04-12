package main

type ctx struct {
	Data map[string]string
	Host hostCtx
}

type hostCtx struct {
	NumCPU int
}
