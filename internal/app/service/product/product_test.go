package sproduct

import (
	"context"
	"errors"
	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/testutil"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/andreyloginov-afk/catalog-service/internal/app/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type createProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

type getByGUIDProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

type updateProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

type deleteProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

type listProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}

}
func TestListProductSuit(t *testing.T) {
	suite.Run(t, new(listProductSuite))

}

func TestCreateProductSuit(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}
	type want struct {
		err error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGuid: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				var nilUUID *uuid.UUID

				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), nilUUID).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(mock.Anything, args.req.CategoryGuid).
					Return(entity.Category{}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGuid: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				var nilUUID *uuid.UUID

				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), nilUUID).
					Return([]entity.Product{{Name: args.req.Name}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGuid: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				var nilUUID *uuid.UUID

				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), nilUUID).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(mock.Anything, args.req.CategoryGuid).
					Return(entity.Category{}, entity.ErrNotFound).
					Once()
			},
		},
		{
			name: "list error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGuid: categoryGUID,
				},
			},
			want: want{err: errors.New("list failed")},
			prepare: func(args args) {
				var nilUUID *uuid.UUID

				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), nilUUID).
					Return(nil, errors.New("list failed")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.svc.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorContains(err, tc.want.err.Error())
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGuid, result.CategoryGuid)
			}
		})
	}
}

func (s *getByGUIDProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDProductSuite))
}

func (s *getByGUIDProductSuite) TestGetByGUID() {
	type args struct {
		guid uuid.UUID
	}
	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				guid: productGUID,
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUID(s.ctx, args.guid).
					Return(entity.Product{
						GUID: args.guid,
						Name: "Test Product",
					}, nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUID(s.ctx, args.guid).
					Return(entity.Product{}, entity.ErrNotFound).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.svc.GetByGUID(s.ctx, tc.args.guid)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.Equal(tc.args.guid, result.GUID)
			}
		})
	}
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuit(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	type args struct {
		guid uuid.UUID
		req  entity.RequestProductUpdate
	}

	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())
	categoryGUID := uuid.Must(uuid.NewV4())

	description := "desc"

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "full update",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "NewName",
					Description:  &description,
					Price:        200,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil)

				s.productRepo.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil)
			},
		},
		{
			name: "partial update - name only",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:  "OnlyName",
					Price: 100,
				},
			},
			want: want{},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil)

				s.productRepo.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil)
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "Name",
					Price:        100,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{}, entity.ErrNotFound)
			},
		},
		{
			name: "duplicate name",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "Duplicate",
					Price:        100,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return([]entity.Product{
						{GUID: uuid.Must(uuid.NewV4())},
					}, nil)
			},
		},
		{
			name: "category not found (ignored by current svc)",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "Name",
					Price:        100,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil)

				s.productRepo.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil)
			},
		},

		// 🔥 ДОБИВАЕМ ПОКРЫТИЕ

		{
			name: "list error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "ErrList",
				},
			},
			want: want{err: errors.New("list error")},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return(nil, errors.New("list error"))
			},
		},
		{
			name: "get by guid error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "ErrGet",
				},
			},
			want: want{err: errors.New("get error")},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{}, errors.New("get error"))
			},
		},
		{
			name: "same guid allowed",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "SameName",
				},
			},
			want: want{},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, args.guid).
					Return(entity.Product{GUID: args.guid}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil)).
					Return([]entity.Product{
						{GUID: args.guid},
					}, nil)

				s.productRepo.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			s.productRepo.ExpectedCalls = nil
			s.productRepo.Calls = nil
			s.categoryRepo.ExpectedCalls = nil
			s.categoryRepo.Calls = nil

			tt.prepare(tt.args)

			_, err := s.svc.Update(s.ctx, tt.args.guid, tt.args.req)

			if tt.want.err != nil {
				s.ErrorIs(err, tt.want.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestDeleteProductSuit(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {
	testGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		guid    uuid.UUID
		wantErr error
		prepare func()
	}{
		{
			name:    "success",
			guid:    testGUID,
			wantErr: nil,
			prepare: func() {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, testGUID).
					Return(entity.Product{GUID: testGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(mock.Anything, testGUID).
					Return(nil).
					Once()
			},
		},
		{
			name:    "not found",
			guid:    testGUID,
			wantErr: entity.ErrNotFound,
			prepare: func() {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, testGUID).
					Return(entity.Product{}, entity.ErrNotFound).
					Once()
			},
		},
		{
			name:    "delete error",
			guid:    testGUID,
			wantErr: errors.New("db delete failed"),
			prepare: func() {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(mock.Anything, testGUID).
					Return(entity.Product{GUID: testGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(mock.Anything, testGUID).
					Return(errors.New("db delete failed")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.productRepo.ExpectedCalls = nil
			s.productRepo.Calls = nil

			tc.prepare()

			err := s.svc.Delete(s.ctx, tc.guid)

			if tc.wantErr != nil {

				s.ErrorContains(err, tc.wantErr.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *listProductSuite) TestList() {
	testCases := []struct {
		name     string
		mockResp []entity.Product
		mockErr  error
		wantErr  error
		check    func(res []entity.Product)
	}{
		{
			name: "success",
			mockResp: []entity.Product{
				{
					GUID:  uuid.Must(uuid.NewV4()),
					Name:  "iPhone",
					Price: 1000,
				},
				{
					GUID:  uuid.Must(uuid.NewV4()),
					Name:  "MacBook",
					Price: 2000,
				},
			},
			mockErr: nil,
			wantErr: nil,
			check: func(res []entity.Product) {
				s.Len(res, 2)

				s.Equal("iPhone", res[0].Name)
				s.Equal(float64(1000), res[0].Price)

				s.Equal("MacBook", res[1].Name)
				s.Equal(float64(2000), res[1].Price)
			},
		},
		{
			name:     "empty result",
			mockResp: []entity.Product{},
			mockErr:  nil,
			wantErr:  nil,
			check: func(res []entity.Product) {
				s.Len(res, 0)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			s.productRepo.ExpectedCalls = nil
			s.productRepo.Calls = nil

			s.productRepo.EXPECT().
				List(mock.Anything, (*string)(nil), (*uuid.UUID)(nil)).
				Return(tc.mockResp, tc.mockErr).
				Once()

			res, err := s.svc.List(s.ctx)

			if tc.wantErr != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				tc.check(res)
			}
		})
	}
}

func (s *updateProductSuite) TestNewService() {
	svc := NewService(s.productRepo, s.categoryRepo)
	s.NotNil(svc)
}
