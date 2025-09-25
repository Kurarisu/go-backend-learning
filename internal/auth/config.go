package auth

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GetJWTSecret() []byte {
	return jwtSecret
}
