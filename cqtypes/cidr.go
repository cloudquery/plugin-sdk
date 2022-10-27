package cqtypes

type CIDR Inet

func (dst *CIDR) Set(src interface{}) error {
	return (*Inet)(dst).Set(src)
}

func (dst CIDR) Get() interface{} {
	return (Inet)(dst).Get()
}
