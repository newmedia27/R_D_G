package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"tcp/internal/documentstore"
)

func setupTestCollection(t *testing.T) *documentstore.Collection {
	store := documentstore.NewStore()

	coll, err := store.CreateCollection("users", &documentstore.CollectionConfig{
		PrimaryKey: "id",
	})
	require.NoError(t, err, "should create test collection")
	return coll
}

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) Put(doc documentstore.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockCollection) Get(key string) (*documentstore.Document, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*documentstore.Document), args.Error(1)
}

func (m *MockCollection) Delete(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockCollection) List() []documentstore.Document {
	args := m.Called()
	return args.Get(0).([]documentstore.Document)
}

func (m *MockCollection) Query(fieldName string, params documentstore.QueryParams) ([]documentstore.Document, error) {
	args := m.Called()
	return args.Get(0).([]documentstore.Document), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {

	type args struct {
		u UserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "create user success",
			args: args{
				u: UserRequest{
					Name: "John Doe",
					Age:  "30",
				},
			},
			wantErr: nil,
		},
		{
			name: "create user with empty fields",
			args: args{
				u: UserRequest{
					Name: "",
					Age:  "",
				},
			},
			wantErr: nil, //валідацію не писав!! TODO:
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := setupTestCollection(t)

			service := NewUserService(collection)

			usr, err := service.CreateUser(tt.args.u)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, usr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, usr)
			assert.Equal(t, tt.args.u.Name, usr.Name)
			assert.Equal(t, tt.args.u.Age, usr.Age)
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	type fields struct {
		coll *documentstore.Collection
	}
	tests := []struct {
		name    string
		mocks   []documentstore.Document
		want    []User
		wantErr error
	}{
		{
			name: "get user list",
			mocks: []documentstore.Document{
				{
					Fields: map[string]documentstore.DocumentField{
						"id": {
							Type:  documentstore.DocumentFieldTypeString,
							Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
						},
						"name": {
							Type:  documentstore.DocumentFieldTypeString,
							Value: "John",
						},
						"age": {
							Type:  documentstore.DocumentFieldTypeString,
							Value: "42",
						},
					},
				},
			},
			want: []User{
				{
					ID: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
					UserRequest: UserRequest{
						Name: "John",
						Age:  "42",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "get user empty list",
			mocks:   make([]documentstore.Document, 0),
			want:    make([]User, 0),
			wantErr: nil,
		},
		{
			name:    "get user empty list",
			mocks:   make([]documentstore.Document, 0),
			want:    make([]User, 0),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockColl := new(MockCollection)
			mockColl.On("List").Return(tt.mocks).Once()
			s := NewUserService(mockColl)

			got, err := s.ListUsers()

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				mockColl.AssertExpectations(t)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		mocks   *documentstore.Document
		args    args
		want    *User
		wantErr error
		mockErr error
	}{
		{
			name: "get user success",
			mocks: &documentstore.Document{
				Fields: map[string]documentstore.DocumentField{
					"id": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
					},
					"name": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "John",
					},
					"age": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "42",
					},
				},
			},
			args: args{
				userID: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
			},
			want: &User{
				ID: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
				UserRequest: UserRequest{
					Name: "John",
					Age:  "42",
				},
			},
		},
		{
			name:    "get user not found",
			mocks:   nil,
			mockErr: documentstore.ErrDocumentNotFound,
			args: args{
				userID: "non-existing-id",
			},
			want:    nil,
			wantErr: documentstore.ErrDocumentNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockColl := new(MockCollection)
			mockColl.On("Get", tt.args.userID).Return(tt.mocks, tt.mockErr).Once()
			s := NewUserService(mockColl)

			got, err := s.GetUser(tt.args.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				mockColl.AssertExpectations(t)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)

		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type fields struct {
		coll *documentstore.Collection
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		mocks   bool
		args    args
		wantErr error
	}{
		{
			name:  "delete existing user - success",
			mocks: true,
			args: args{
				userID: "30b8ac02-7a33-4fd4-8a1d-7ce93627f3a0",
			},
			wantErr: nil,
		},
		{
			name:  "delete non-existing user - returns error",
			mocks: false,
			args: args{
				userID: "some_id",
			},
			wantErr: ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockColl := new(MockCollection)
			mockColl.On("Delete", tt.args.userID).Return(tt.mocks).Once()
			s := NewUserService(mockColl)

			err := s.DeleteUser(tt.args.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			mockColl.AssertExpectations(t)
		})
	}
}
