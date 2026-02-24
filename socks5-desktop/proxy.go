package main

type Proxy struct {
	rhost string
	ruser string
	rpass string
}

func (p *Proxy) Start(host, user, pass string) error {
	return nil
}

func (p *Proxy) Stop() error {
	return nil
}

func (p *Proxy) SetRemote(host, user, pass string) error {
	return nil
}
