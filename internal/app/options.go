package app

type TlsOptions struct {
	certFile string
	keyFile  string
}

type ConnOptions struct {
	Port int
}

type Options struct {
	tls  TlsOptions
	tcp  *ConnOptions
	quic *ConnOptions
}

func defaultOptions() Options {
	return Options{}
}

type OptionModifier func(*Options)

// OptionModifiers
func WithTls(certFile string, keyFile string) OptionModifier {
	return func(o *Options) {
		o.tls = TlsOptions{
			certFile: certFile,
			keyFile:  keyFile,
		}
	}
}

func WithTcp(port int) OptionModifier {
	return func(o *Options) {
		o.tcp = &ConnOptions{
			Port: port,
		}
	}
}

func WithQuic(port int) OptionModifier {
	return func(o *Options) {
		o.quic = &ConnOptions{
			Port: port,
		}
	}
}
