package sproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository/mocks"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/testutil"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
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
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(mock.Anything, []uuid.UUID{args.req.CategoryGuid}).
					Return([]entity.Category{{GUID: args.req.CategoryGuid}}, nil).
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
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(mock.Anything, []uuid.UUID{args.req.CategoryGuid}).
					Return([]entity.Category{}, nil).
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
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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

	productGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		wantLen int
		prepare func(args args)
	}{
		{
			name:    "success",
			args:    args{guid: productGUID},
			wantLen: 1,
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{
						GUID: args.guid,
						Name: "Test Product",
					}}, nil).
					Once()
			},
		},
		{
			name:    "not found",
			args:    args{guid: productGUID},
			wantLen: 0,
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			results, err := s.svc.GetByGUIDs(s.ctx, []uuid.UUID{tc.args.guid})

			s.NoError(err)
			s.Len(results, tc.wantLen)
			if tc.wantLen > 0 {
				s.Equal(tc.args.guid, results[0].GUID)
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

	listErr := errors.New("list error")
	getErr := errors.New("get error")

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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil)
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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil)

				s.productRepo.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("entity.Product")).
					Return(nil)
			},
		},
		{
			name: "list error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "ErrList",
				},
			},
			want: want{err: listErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return(nil, listErr)
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
			want: want{err: getErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				s.productRepo.EXPECT().
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return(nil, getErr)
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
					GetByGUIDs(mock.Anything, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil)

				s.productRepo.EXPECT().
					List(mock.Anything, mock.AnythingOfType("*string"), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{testGUID}).
					Return([]entity.Product{{GUID: testGUID}}, nil).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{testGUID}).
					Return([]entity.Product{}, nil).
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
					GetByGUIDs(mock.Anything, []uuid.UUID{testGUID}).
					Return([]entity.Product{{GUID: testGUID}}, nil).
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
				List(mock.Anything, (*string)(nil), (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
				Return(tc.mockResp, tc.mockErr).
				Once()

			res, err := s.svc.List(s.ctx, entity.RequestProductList{})

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

type getByGUIDsProductSuite struct {
	suite.Suite

	svc          *svc
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *getByGUIDsProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.svc = &svc{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDsProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDsProductSuite))
}

func (s *getByGUIDsProductSuite) TestGetByGUIDs() {
	guid1 := uuid.Must(uuid.NewV4())
	guid2 := uuid.Must(uuid.NewV4())
	guid3 := uuid.Must(uuid.NewV4())
	dbErr := errors.New("db error")

	testCases := []struct {
		name    string
		guids   []uuid.UUID
		wantLen int
		wantErr error
		prepare func()
	}{
		{
			name:    "all found",
			guids:   []uuid.UUID{guid1, guid2},
			wantLen: 2,
			prepare: func() {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{guid1, guid2}).
					Return([]entity.Product{{GUID: guid1}, {GUID: guid2}}, nil).Once()
			},
		},
		{
			name:    "some missing",
			guids:   []uuid.UUID{guid1, guid2, guid3},
			wantLen: 2,
			prepare: func() {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{guid1, guid2, guid3}).
					Return([]entity.Product{{GUID: guid1}, {GUID: guid3}}, nil).Once()
			},
		},
		{
			name:    "db error",
			guids:   []uuid.UUID{guid1},
			wantErr: dbErr,
			prepare: func() {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{guid1}).
					Return(nil, dbErr).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare()

			result, err := s.svc.GetByGUIDs(s.ctx, tc.guids)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Len(result, tc.wantLen)
			}
		})
	}
}
