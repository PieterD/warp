package ecs_test

import (
	"fmt"
	"testing"

	"github.com/PieterD/warp/pkg/ecs"
)

func TestRepository_Basic(t *testing.T) {
	type TestComponent struct {
		CoolInt  int
		CoolBool bool
	}
	want := TestComponent{
		CoolInt:  542,
		CoolBool: true,
	}
	var got TestComponent

	repo := ecs.NewRepository()
	e := repo.NewEntity()
	if err := e.GetComponent(&got); !ecs.IsMissingComponentError(err) {
		t.Fatalf("getting component before setting: expected missing component error: %v", err)
	}
	if err := e.DelComponent(got); !ecs.IsMissingComponentError(err) {
		t.Fatalf("deleting component before setting: expected missing component error: %v", err)
	}
	if err := e.SetComponent(want); err != nil {
		t.Fatalf("setting component: %v", err)
	}
	if err := e.GetComponent(&got); err != nil {
		t.Fatalf("getting component after setting: %v", err)
	}
	if want != got {
		t.Fatalf("data mismatch: get does not match set, got: %#v", got)
	}
	{
		got = TestComponent{}
		err := repo.Entities(func(vEntity ecs.Entity) error {
			if vEntity.ID() != e.ID() {
				return fmt.Errorf("visiting entity id != only entity id")
			}
			if got != want {
				return fmt.Errorf("visiting component does not match, got: %#v", got)
			}
			return nil
		}, &got)
		if err != nil {
			t.Fatalf("entities: %v", err)
		}
	}
	if err := e.DelComponent(got); err != nil {
		t.Fatalf("deleting component after setting: %v", err)
	}
	if err := e.GetComponent(&got); !ecs.IsMissingComponentError(err) {
		t.Fatalf("getting component after deleting: expected missing component error: %v", err)
	}
}
