package chronary

import (
	"context"
	"fmt"
	"testing"
)

func TestPageIteratorSinglePage(t *testing.T) {
	iter := newPageIterator(50, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		return &ListResponse[Agent]{
			Data:  []Agent{{ID: "agt_1"}, {ID: "agt_2"}},
			Total: 2,
		}, nil
	})

	items, err := iter.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPageIteratorMultiPage(t *testing.T) {
	iter := newPageIterator(2, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		if offset == 0 {
			return &ListResponse[Agent]{
				Data:  []Agent{{ID: "agt_1"}, {ID: "agt_2"}},
				Total: 3,
			}, nil
		}
		return &ListResponse[Agent]{
			Data:  []Agent{{ID: "agt_3"}},
			Total: 3,
		}, nil
	})

	items, err := iter.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
	if items[2].ID != "agt_3" {
		t.Errorf("expected agt_3, got %s", items[2].ID)
	}
}

func TestPageIteratorEmpty(t *testing.T) {
	iter := newPageIterator(50, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		return &ListResponse[Agent]{Data: []Agent{}, Total: 0}, nil
	})

	items, err := iter.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestPageIteratorError(t *testing.T) {
	iter := newPageIterator(50, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		return nil, fmt.Errorf("network error")
	})

	_, err := iter.Collect(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPageIteratorGetPage(t *testing.T) {
	iter := newPageIterator(50, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		return &ListResponse[Agent]{
			Data:  []Agent{{ID: "agt_1"}},
			Total: 3,
		}, nil
	})

	page, err := iter.GetPage(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if !page.HasMore {
		t.Error("expected HasMore to be true")
	}
	if page.Total != 3 {
		t.Errorf("expected total 3, got %d", page.Total)
	}
}

func TestPageIteratorAllEarlyBreak(t *testing.T) {
	calls := 0
	iter := newPageIterator(2, func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		calls++
		return &ListResponse[Agent]{
			Data:  []Agent{{ID: "agt_1"}, {ID: "agt_2"}},
			Total: 100,
		}, nil
	})

	count := 0
	for _, err := range iter.All(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if count >= 1 {
			break
		}
	}
	if count != 1 {
		t.Errorf("expected 1 item before break, got %d", count)
	}
}
