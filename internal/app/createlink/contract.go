package createlink

type CreateLinkCommand struct {
	OwnerID   string
	OriginURL string
}

type CreateLinkResult struct {
	Code string
}
