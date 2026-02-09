package documentstore

import (
	"reflect"
	"testing"
)

func NewValidDocument() *Document {
	return &Document{
		Fields: map[string]DocumentField{
			"Id": {
				Type:  DocumentFieldTypeString,
				Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
			},
			"Name": {
				Type:  DocumentFieldTypeString,
				Value: "John",
			},
		},
	}
}

func TestMarshalDocument(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    *Document
		wantErr bool
	}{
		{
			name: "doc created is valid",
			args: args{
				input: struct {
					Id   string
					Name string
				}{
					Id:   "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
					Name: "John",
				},
			},
			want:    NewValidDocument(),
			wantErr: false,
		},
		{
			name: "doc created is invalid",
			args: args{
				input: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalDocument(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalDocument(t *testing.T) {
	type args struct {
		doc    *Document
		output any
	}
	tests := []struct {
		name       string
		args       args
		wantOutput any
		wantErr    bool
	}{
		{
			name: "unmarshal to struct",
			args: args{
				doc: NewValidDocument(),
				output: &struct {
					Id   string
					Name string
				}{},
			},
			wantOutput: struct {
				Id   string
				Name string
			}{
				Id:   "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
				Name: "John",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalDocument(tt.args.doc, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
