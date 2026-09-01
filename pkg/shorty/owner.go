package shorty

import "context"

const ownerHeader = "X-Shorty-Owner"

type ownerContextKey struct{}

// WithOwner associates the owner of managed links with API requests made using
// ctx. Public link resolution does not require an owner.
func WithOwner(ctx context.Context, ownerID string) context.Context {
	return context.WithValue(ctx, ownerContextKey{}, ownerID)
}

func ownerFromContext(ctx context.Context) string {
	ownerID, _ := ctx.Value(ownerContextKey{}).(string)
	return ownerID
}
