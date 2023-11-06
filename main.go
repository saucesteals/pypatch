package main

// #include "dll.h"
import "C"
import (
	"context"
	_ "embed"
	"time"

	"github.com/saucesteals/pypatch/pypatch"
)

//go:embed inject.py
var code string

//export OnProcessAttach
func OnProcessAttach() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	p, err := pypatch.New(ctx)
	cancel()
	if err != nil {
		panic(err)
	}

	if err := p.Inject(code); err != nil {
		panic(err)
	}
}

func main() {}
