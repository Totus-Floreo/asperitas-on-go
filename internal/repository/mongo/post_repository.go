package mongo_repository

import (
	"context"
	"sync"
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostStorage struct {
	PostStorage    model.ICollection
	CommentStorage model.ICollection
	mu             *sync.Mutex
	ReadersPool    model.IDBReadersPool
}

func NewPostStorage(client model.IClient, pool model.IDBReadersPool) *PostStorage {
	db := client.Database("asperitas")
	postStorage := db.Collection("posts")
	commentStorage := db.Collection("comments")
	return &PostStorage{
		PostStorage:    postStorage,
		CommentStorage: commentStorage,
		ReadersPool:    pool,
		mu:             new(sync.Mutex),
	}
}

func (s *PostStorage) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	client := s.ReadersPool.GetConnection()
	defer s.ReadersPool.ReleaseConnection(client)

	collection := client.Database("asperitas").Collection("posts")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{GetSort(), GetLookup()})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStorage) GetPostByID(ctx context.Context, postID primitive.ObjectID) (*model.Post, error) {
	client := s.ReadersPool.GetConnection()
	defer s.ReadersPool.ReleaseConnection(client)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := client.Database("asperitas").Collection("posts")

	match := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{
					Key:   "_id",
					Value: postID,
				},
			},
		},
	}

	var post *model.Post
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{match, GetLookup()})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.TryNext(ctx) {
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
	} else {
		/* Честно не понимаю, почему при отсутствии документа он не возврощает ошибку
		возможно cursor в принципе не возврощает mongo.ErrNoDocuments */
		if cursor.Err() == mongo.ErrNoDocuments {
			return nil, model.ErrPostNotFound
		}
		if cursor.Err() != nil {
			return nil, cursor.Err()
		}
		return nil, model.ErrPostNotFound
	}

	return post, nil

}

func (s *PostStorage) GetPostsByCategory(ctx context.Context, category string) ([]*model.Post, error) {
	client := s.ReadersPool.GetConnection()
	defer s.ReadersPool.ReleaseConnection(client)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := client.Database("asperitas").Collection("posts")

	match := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{
					Key:   "category",
					Value: category,
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{match, GetSort(), GetLookup()})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStorage) GetPostsByUser(ctx context.Context, userName string) ([]*model.Post, error) {
	client := s.ReadersPool.GetConnection()
	defer s.ReadersPool.ReleaseConnection(client)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := client.Database("asperitas").Collection("posts")

	match := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{
					Key:   "author.username",
					Value: userName,
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{match, GetSort(), GetLookup()})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStorage) AddPost(ctx context.Context, post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post.ID = primitive.NewObjectID()
	_, err := s.PostStorage.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStorage) DeletePost(ctx context.Context, postID primitive.ObjectID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.D{
		{
			Key:   "_id",
			Value: postID,
		},
	}

	post, err := s.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	commentObjectIDs := make([]primitive.ObjectID, 0)
	for _, comment := range post.Comments {
		commentObjectIDs = append(commentObjectIDs, comment.ID)
	}

	result := s.PostStorage.FindOneAndDelete(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.ErrPostNotFound
	}
	if result.Err() != nil {
		return result.Err()
	}

	if len(post.Comments) == 0 {
		return nil
	}

	delete := bson.D{
		{
			Key: "_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: commentObjectIDs,
				},
			},
		},
	}

	commentResult, err := s.CommentStorage.DeleteMany(ctx, delete)
	if err != nil {
		return err
	}
	if commentResult.DeletedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}

func (s *PostStorage) AddView(ctx context.Context, post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "views",
					Value: post.Views + 1},
			},
		},
	}

	postResult, err := s.PostStorage.UpdateByID(ctx, post.ID, update)
	if err != nil {
		return err
	}
	if postResult.MatchedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}

func (s *PostStorage) AddComment(ctx context.Context, post *model.Post, comment *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment.ID = primitive.NewObjectID()
	_, err := s.CommentStorage.InsertOne(ctx, comment)
	if err != nil {
		return err
	}

	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{
					Key:   "comments",
					Value: comment.ID,
				},
			},
		},
	}

	postResult, err := s.PostStorage.UpdateByID(ctx, post.ID, update)
	if err != nil {
		return err
	}
	if postResult.MatchedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}

func (s *PostStorage) DeleteComment(ctx context.Context, post *model.Post, commentID primitive.ObjectID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{
					Key: "comments",
					Value: bson.D{
						{
							Key:   "$eq",
							Value: commentID,
						},
					},
				},
			},
		},
	}

	result, err := s.PostStorage.UpdateByID(ctx, post.ID, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrPostNotFound
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: commentID,
		},
	}

	commentResult, err := s.CommentStorage.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if commentResult.DeletedCount == 0 {
		return model.ErrCommentNotFound
	}

	return nil
}

func (s *PostStorage) Vote(ctx context.Context, post *model.Post, vote *model.Vote) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filterWithArrayFilter := bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
		{
			Key: "votes",
			Value: bson.D{
				{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "user",
							Value: vote.UserID,
						},
					},
				},
			},
		},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "votes.$.vote",
					Value: vote.Score,
				},
			},
		},
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
	}

	push := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{
					Key:   "votes",
					Value: vote,
				},
			},
		},
	}

	findResult := s.PostStorage.FindOne(ctx, filterWithArrayFilter)
	if findResult.Err() == mongo.ErrNoDocuments {
		result, err := s.PostStorage.UpdateOne(ctx, filter, push)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return model.ErrPostNotFound
		}
	} else {
		resultArrayFiltred, err := s.PostStorage.UpdateOne(ctx, filterWithArrayFilter, update)
		if err != nil {
			return err
		}
		if resultArrayFiltred.MatchedCount == 0 {
			return model.ErrPostNotFound
		}
	}

	return nil
}

func (s *PostStorage) UnVote(ctx context.Context, post *model.Post, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
		{
			Key: "votes",
			Value: bson.D{
				{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "user",
							Value: userID,
						},
					},
				},
			},
		},
	}

	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{
					Key: "votes",
					Value: bson.D{
						{
							Key:   "user",
							Value: userID,
						},
					},
				},
			},
		},
	}

	postResult, err := s.PostStorage.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if postResult.MatchedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}

func (s *PostStorage) UpdateScore(ctx context.Context, post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	score := 0
	var upVotePercentage int64 = 0

	if len(post.Votes) != 0 {
		positiveVotes := 0
		for _, vote := range post.Votes {
			switch vote.Score {
			case 1:
				positiveVotes++
				score++
			case -1:
				score--
			}
		}

		upVotePercentage = int64((float64(positiveVotes) / float64(len(post.Votes))) * 100)
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: post.ID,
		},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "upvotepercentage",
					Value: upVotePercentage,
				},
				{
					Key:   "score",
					Value: score,
				},
			},
		},
	}

	postResult, err := s.PostStorage.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if postResult.MatchedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}

func GetSort() primitive.D {
	sort := bson.D{
		{
			Key: "$sort",
			Value: bson.D{
				{Key: "score", Value: -1},
				{Key: "views", Value: -1},
			},
		},
	}
	return sort
}

func GetLookup() primitive.D {
	lookup := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "comments"},
				{Key: "localField", Value: "comments"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "comments"},
			},
		},
	}

	return lookup
}
