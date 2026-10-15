package console

import (
	"context"
	"testing"
)

func TestContextFlags(t *testing.T) {
	ctx := context.Background()
	if IsDedicated(ctx) || IsClient(ctx) || IsServer(ctx) || IsCheats(ctx) {
		t.Fatal("flags set on an empty context")
	}
	ctx = AsDedicated(AsServer(ctx))
	if !IsDedicated(ctx) || !IsServer(ctx) || IsClient(ctx) {
		t.Fatal("dedicated server context not detected")
	}
	if !IsCheats(WithCheats(ctx)) {
		t.Fatal("cheats not set")
	}
	if !IsClient(AsClient(context.Background())) {
		t.Fatal("client not set")
	}
}
