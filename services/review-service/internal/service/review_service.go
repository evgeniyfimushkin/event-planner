package service

import (
    "net/http"
    "net/url"
	grpcclient "review-service/internal/client/grpc-client"
	"review-service/internal/models"
	"review-service/internal/repository"
    "fmt"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
)

// ReviewService specializes in handling business logic for Review entities.
// It embeds GenericService for basic CRUD operations and adds additional dependencies (e.g., a verifier).
type ReviewService struct {
	*service.GenericService[models.Review]
    eventClient *grpcclient.EventClient
}

// NewReviewService creates a new instance of ReviewService using the provided verifier and repository.
// It initializes the underlying GenericService using the given repository.
func NewReviewService(repo *repository.ReviewRepository, grpcclient *grpcclient.EventClient) *ReviewService {
	return &ReviewService{
		GenericService: service.NewGenericService[models.Review](repo),
        eventClient: grpcclient,
	}
}

func (s *ReviewService) Create(claims jwt.MapClaims, entity *models.Review) (*models.Review, error) {
    userIDFloat, ok := claims["userID"].(float64)
    if !ok {
        return nil, fmt.Errorf("UserID is not a number")
    }

    entity.UserID = uint(userIDFloat) 

    existing, err := s.FindFirst(claims, "event_id = ? AND user_id = ?", entity.EventID, entity.UserID)
    if err == nil && existing != nil {
        return nil, fmt.Errorf("user is already reviewed this event")
    }
    
   
    //-------------------HTTP----------------
    baseURL := "http://event-service:8082/api/v1/events"
    params := url.Values{}
    params.Add("id", string(entity.EventID))
    fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
    resp, err := http.Get(fullURL)
    if err != nil {
        return nil, fmt.Errorf("HTTP Request Error:", err)
    }
    defer resp.Body.Close()

    updatedReview , err := s.GenericService.Create(claims, entity)
    if err != nil {
        return nil, err
    }

    return updatedReview, nil
}

