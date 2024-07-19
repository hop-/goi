package app

type TlsOptions struct {
	certFile string
	keyFile  string
}

type ConnOptions struct {
	Port int
}

type StorageOptions struct {
	storageType string
	uri         string
}

type Options struct {
	tls     TlsOptions
	tcp     *ConnOptions
	quic    *ConnOptions
	storage StorageOptions
}

func defaultOptions() Options {
	return Options{
		// TODO
	}
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

func WithStorage(storageType string, storageUri string) OptionModifier {
	return func(o *Options) {
		o.storage = StorageOptions{
			storageType: storageType,
			uri:         storageUri,
		}
	}
}
