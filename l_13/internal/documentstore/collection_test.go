package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Put(t *testing.T) {
	type fields struct {
		Documents map[string]*Document
		name      string
		Cfg       CollectionConfig
	}
	type args struct {
		doc Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "put doc success",
			fields: fields{
				Documents: make(map[string]*Document),
				name:      "users",
				Cfg: CollectionConfig{
					PrimaryKey: "id",
				},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeString,
							Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
						},
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "John",
						},
					}},
			},
			wantErr: nil,
		},
		{
			name: "put doc without primary key",
			fields: fields{
				Documents: make(map[string]*Document),
				name:      "users",
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeString,
							Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
						},
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "John",
						},
					}},
			},
			wantErr: ErrEmptyPrimaryKey,
		},
		{
			name: "put doc with non-string primary key",
			fields: fields{
				Documents: make(map[string]*Document),
				name:      "users",
				Cfg:       CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeNumber,
							Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
						},
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "John",
						},
					}},
			},
			wantErr: ErrUnsupportedDocumentField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Documents: tt.fields.Documents,
				name:      tt.fields.name,
				Cfg:       tt.fields.Cfg,
			}

			var id string
			if f, ok := tt.args.doc.Fields[tt.fields.Cfg.PrimaryKey]; ok && f.Type == DocumentFieldTypeString {
				id = f.Value.(string)
			}

			err := c.Put(tt.args.doc)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			d, updated := c.Documents[id]
			assert.True(t, updated)
			assert.Equal(t, tt.args.doc, *d)
			require.NotNil(t, d)

			assert.Equal(t, tt.args.doc.Fields[tt.fields.Cfg.PrimaryKey].Value, d.Fields[tt.fields.Cfg.PrimaryKey].Value)
		})
	}
}

func TestCollection_Get(t *testing.T) {
	type fields struct {
		Documents map[string]*Document
		name      string
		Cfg       CollectionConfig
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Document
		wantErr error
	}{
		{
			name: "get doc success",
			fields: fields{
				Documents: map[string]*Document{
					"30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0": {
						Fields: map[string]DocumentField{
							"id": {
								Type:  DocumentFieldTypeString,
								Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
							},
							"name": {
								Type:  DocumentFieldTypeString,
								Value: "John",
							},
						},
					},
				},
				name: "users",
				Cfg: CollectionConfig{
					PrimaryKey: "id",
				},
			},
			args: args{
				key: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
			},
			want: &Document{
				Fields: map[string]DocumentField{
					"id": {
						Type:  DocumentFieldTypeString,
						Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
					},
					"name": {
						Type:  DocumentFieldTypeString,
						Value: "John",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "get from empty collection",
			fields: fields{
				Documents: make(map[string]*Document),
				name:      "empty",
				Cfg:       CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				key: "any-key",
			},
			want:    nil,
			wantErr: ErrDocumentNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Documents: tt.fields.Documents,
				name:      tt.fields.name,
				Cfg:       tt.fields.Cfg,
			}
			got, err := c.Get(tt.args.key)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	type fields struct {
		Documents map[string]*Document
		name      string
		Cfg       CollectionConfig
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success delete doc",
			fields: fields{
				Documents: map[string]*Document{
					"30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0": {
						Fields: map[string]DocumentField{
							"id": {
								Type:  DocumentFieldTypeString,
								Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
							},
							"name": {
								Type: DocumentFieldTypeString,
							},
						},
					},
				},
				name: "users",
				Cfg: CollectionConfig{
					PrimaryKey: "id",
				},
			},
			args: args{
				key: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
			},
			want: true,
		},
		{
			name: "delete non-existing document",
			fields: fields{
				Documents: make(map[string]*Document),
				name:      "users",
				Cfg:       CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				key: "non-existing-id",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Documents: tt.fields.Documents,
				name:      tt.fields.name,
				Cfg:       tt.fields.Cfg,
			}
			count := len(c.Documents)
			ok := c.Delete(tt.args.key)

			if tt.want {
				require.True(t, ok)
				_, ok = c.Documents[tt.args.key]
				assert.False(t, ok)
				assert.Equal(t, count-1, len(c.Documents))

			} else {
				assert.Equal(t, count, len(c.Documents))
			}

		})
	}
}

func TestCollection_List(t *testing.T) {
	type fields struct {
		Documents map[string]*Document
		name      string
		Cfg       CollectionConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   []Document
	}{
		{
			name: "get list of documents",
			fields: fields{
				Documents: map[string]*Document{
					"30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0": {
						Fields: map[string]DocumentField{
							"id": {
								Type:  DocumentFieldTypeString,
								Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
							},
							"name": {
								Type:  DocumentFieldTypeString,
								Value: "John",
							},
						},
					},
				},
			},
			want: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeString,
							Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
						},
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "John",
						},
					},
				},
			},
		},
		{
			name: "get empty list",
			fields: fields{
				Documents: make(map[string]*Document),
			},
			want: make([]Document, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Documents: tt.fields.Documents,
				name:      tt.fields.name,
				Cfg:       tt.fields.Cfg,
			}
			got := c.List()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
