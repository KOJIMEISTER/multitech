package handlers

import (
	"context"
	"fmt"
	"multitech/pkg/testutils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	containers, err := testutils.SetupContainers(ctx)
	if err != nil {
		fmt.Printf("Failed to setup containers: %v", err)
		os.Exit(1)
	}
	defer containers.Terminate(ctx)
}
