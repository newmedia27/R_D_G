package documentstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
		cfg  *CollectionConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr error
	}{
		{
			name: "created collection",
			fields: fields{
				Collections: make(map[string]*Collection),
			},
			args: args{
				name: "users",
				cfg:  &CollectionConfig{PrimaryKey: "id"},
			},
			want: &Collection{
				Documents: make(map[string]*Document),
				name:      "users",
				cfg:       CollectionConfig{PrimaryKey: "id"},
			},
			wantErr: nil,
		},
		{
			name: "error collection name",
			fields: fields{
				Collections: make(map[string]*Collection),
			},
			args: args{
				name: "",
				cfg:  &CollectionConfig{PrimaryKey: "id"},
			},
			want:    nil,
			wantErr: ErrEmptyCollectionName,
		},
		{
			name: "collection already exists",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
						name:      "users",
						cfg:       CollectionConfig{PrimaryKey: "id"},
					},
				},
			},
			args: args{
				name: "users",
				cfg:  &CollectionConfig{PrimaryKey: "id"},
			},
			want:    nil,
			wantErr: ErrCollectionAlreadyExists},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Collections: tt.fields.Collections,
			}
			got, err := s.CreateCollection(tt.args.name, tt.args.cfg)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want.name, got.name)
			assert.Equal(t, tt.want.cfg, got.cfg)
			assert.NotNil(t, got.Documents)

			col, ok := s.Collections[tt.args.name]
			assert.True(t, ok, "collection is not exist in store")
			assert.Equal(t, got, col, "collection is not equal to created collection")
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr error
	}{
		{
			name: "collection exists",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
						name:      "users",
						cfg:       CollectionConfig{PrimaryKey: "id"},
					},
				},
			},
			args: args{
				name: "users",
			},
			want: &Collection{
				Documents: make(map[string]*Document),
				name:      "users",
				cfg:       CollectionConfig{PrimaryKey: "id"},
			},
			wantErr: nil,
		},
		{
			name: "collection not exists",
			fields: fields{
				Collections: make(map[string]*Collection),
			},
			args: args{
				name: "some_name",
			},
			want:    nil,
			wantErr: ErrCollectionNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Collections: tt.fields.Collections,
			}
			got, err := s.GetCollection(tt.args.name)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}

		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "delete existing collection",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
						name:      "users",
						cfg:       CollectionConfig{PrimaryKey: "id"},
					},
				},
			},
			args: args{
				name: "users",
			},
			want: true,
		},
		{
			name: "delete not existing collection",
			fields: fields{
				Collections: make(map[string]*Collection),
			},
			args: args{
				name: "users",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Collections: tt.fields.Collections,
			}
			_, isBeforeIsset := s.Collections[tt.args.name]

			got := s.DeleteCollection(tt.args.name)
			assert.Equal(t, tt.want, got)

			if tt.want {
				require.True(t, isBeforeIsset, "collection not exists")
				assert.NotContains(t, s.Collections, tt.args.name, "collection not deleted")
			}

		})
	}
}

// Тут цей тест не потрібен, по факту можу бути помилка лише при конвертації в json
func TestStore_Dump(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "dump normal store",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
					},
					"posts": {
						Documents: make(map[string]*Document),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "dump empty store",
			fields: fields{
				Collections: make(map[string]*Collection),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Collections: tt.fields.Collections,
			}
			got, err := s.Dump()

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.True(t, json.Valid(got))
		})
	}
}

func TestStore_DumpToFIle(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "dump to file success",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
						name:      "users",
						cfg:       CollectionConfig{PrimaryKey: "id"},
					},
				},
			},
			args: args{
				path: "dump.json",
			},
			wantErr: nil,
		},
		{
			name: "dump empty store to file",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {
						Documents: make(map[string]*Document),
					},
				},
			},
			args: args{
				path: "dump.json",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Collections: tt.fields.Collections,
			}
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, tt.args.path)
			err := s.DumpToFIle(path)

			if tt.wantErr != nil {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			info, err := os.Stat(path)
			assert.False(t, info.IsDir(), "path is directory")
			require.NoError(t, err, "file is not exist")

			data, err := os.ReadFile(path)

			require.NoError(t, err, "can't read file")
			require.NotEmpty(t, data, "file is empty")
			require.True(t, json.Valid(data), "file is not valid json")
		})
	}
}
