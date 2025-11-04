package remote

type RemotePayloadId struct {
	Service string
	Id      string
}

type IRemoteController interface {
	// Add adds a resource to remote storage
	Add(data []byte) (*RemotePayloadId, error)

	// Remove removes a resource from a remote storage
	Remove(rpi *RemotePayloadId) error

	// Ping checks the health of a remote storage
	Ping() error
}
