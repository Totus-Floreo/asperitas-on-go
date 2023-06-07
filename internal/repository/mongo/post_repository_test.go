package mongo_repository

import (
	"context"
	"testing"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllPost_Success(t *testing.T) {
	postObjectID, _ := primitive.ObjectIDFromHex("64534b74aed82e0020e916e8")
	expected := []*model.Post{
		&model.Post{
			ID:               postObjectID,
			Category:         "programming",
			Text:             "TestText",
			Title:            "TestTitle",
			Type:             "text",
			Created:          "2006-01-02T15:04:05.000Z",
			Views:            0,
			Score:            0,
			UpvotePercentage: 0,
			Comments:         []*model.Comment{},
			Votes:            []*model.Vote{},
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockIClient(ctrl)
	pool := mocks.NewMockIDBReadersPool(ctrl)

	mockDB := mocks.NewMockIMongoDB(ctrl)

	mockPostColl := mocks.NewMockICollection(ctrl)
	mockCommentColl := mocks.NewMockICollection(ctrl)
	mockCursor := mocks.NewMockICursor(ctrl)

	client.EXPECT().Database("asperitas").Return(mockDB)

	mockDB.EXPECT().Collection("posts").Return(mockPostColl)
	mockDB.EXPECT().Collection("comments").Return(mockCommentColl)

	postStorage := NewPostStorage(client, pool)

	pool.EXPECT().GetConnection().Return(client)

	client.EXPECT().Database("asperitas").Return(mockDB)
	mockDB.EXPECT().Collection("posts").Return(mockPostColl)

	mockPostColl.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockCursor, nil)

	mockCursor.EXPECT().All(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, posts *[]*model.Post) error {
		*posts = append(*posts, expected[0])
		return nil
	})

	mockCursor.EXPECT().Close(gomock.Any())

	pool.EXPECT().ReleaseConnection(client)

	posts, err := postStorage.GetAllPosts(ctx)

	require.NoError(t, err)
	require.Equal(t, expected, posts)
}
