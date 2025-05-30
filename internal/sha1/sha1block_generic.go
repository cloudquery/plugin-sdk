package sha1

func block(dig *digest, p []byte) {
	blockGeneric(dig, p)
}
