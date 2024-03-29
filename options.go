package memcache

type Option func(l *Client)

func WithNamespace(namespace string) Option {
	return func(l *Client) {
		l.namespace = namespace
	}
}

func WithEncoding(encoder Encoder, decoder Decoder) Option {
	return func(l *Client) {
		l.encoder = encoder
		l.decoder = decoder
	}
}

func WithAuth(username, password string) Option {
	return func(l *Client) {
		l.username = username
		l.password = password
	}
}

// WithConfig sets extra configs, use mc.DefaultConfig() as a starting point
func WithConfig(config *Config) Option {
	return func(l *Client) {
		l.config = config
	}
}
