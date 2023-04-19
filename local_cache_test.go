package go_cachec

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBuildInMapCache_Close(t *testing.T) {
	type fields struct {
		data  map[string]*Item
		mu    sync.RWMutex
		close chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildInMapCache{
				data:  tt.fields.data,
				mu:    tt.fields.mu,
				close: tt.fields.close,
			}
			if err := b.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildInMapCache_Del(t *testing.T) {
	type fields struct {
		data  map[string]*Item
		mu    sync.RWMutex
		close chan struct{}
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildInMapCache{
				data:  tt.fields.data,
				mu:    tt.fields.mu,
				close: tt.fields.close,
			}
			if err := b.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildInMapCache_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		key     string
		val     any
		s       time.Duration
		want    any
		wantErr bool
	}{
		{
			name: "set get",
			key:  "name",
			val:  "liu",
			s:    time.Second * 10,
			want: "liu",
		},
		{
			name: "set ex",
			key:  "name",
			val:  "liu",
			s:    time.Second * 3,
			want: "liu",
		},
	}

	//b := NewBuildInMapCache(10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//err := b.Set(context.Background(), tt.key, tt.val, tt.s)
			//if err != nil {
			//	t.Error(err)
			//	return
			//}
			//time.Sleep(tt.s + 1)
			//get, err := b.Get(context.Background(), tt.key)
			//if err != nil {
			//	t.Error(err)
			//	return
			//}
			//assert.Equal(t, get, tt.val)
		})
	}
}

func TestBuildInMapCache_Set(t *testing.T) {
	type fields struct {
		data  map[string]*Item
		mu    sync.RWMutex
		close chan struct{}
	}
	type args struct {
		ctx        context.Context
		key        string
		val        any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildInMapCache{
				data:  tt.fields.data,
				mu:    tt.fields.mu,
				close: tt.fields.close,
			}
			if err := b.Set(tt.args.ctx, tt.args.key, tt.args.val, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewBuildInMapCache(t *testing.T) {
	type args struct {
		size int32
	}
	tests := []struct {
		name string
		args args
		want *BuildInMapCache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := NewBuildInMapCache(tt.args.size); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewBuildInMapCache() = %v, want %v", got, tt.want)
			//}
		})
	}
}
