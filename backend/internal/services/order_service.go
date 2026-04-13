package services

import (
	"context"
	"errors"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/models"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)

type OrderService struct {
	orderRepo   repository.OrderStore
	productRepo repository.ProductStore
	userRepo    *repository.UserRepository
}

func NewOrderService(orderRepo repository.OrderStore, productRepo repository.ProductStore, userRepo *repository.UserRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, req *models.CreateOrderRequest) (*models.Order, error) {
	if s.userRepo != nil {
		if _, err := s.userRepo.FindByID(ctx, req.UserID); err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
	}
	var subtotal float64
	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		if s.productRepo != nil {
			p, err := s.productRepo.FindByID(ctx, it.ProductID)
			if err != nil {
				return nil, err
			}
			if p == nil {
				return nil, ErrProductNotFound
			}
		}
		lineTotal := it.UnitPrice * float64(it.Quantity)
		subtotal += lineTotal
		items = append(items, models.OrderItem{
			ProductID:     it.ProductID,
			ProductName:   it.ProductName,
			SelectedSize:  it.SelectedSize,
			SelectedColor: it.SelectedColor,
			Quantity:      it.Quantity,
			UnitPrice:     it.UnitPrice,
			LineTotal:     lineTotal,
		})
	}
	total := subtotal
	order := &models.Order{
		UserID:          req.UserID,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		DeliveryMethod:  req.DeliveryMethod,
		DeliveryAddress: req.DeliveryAddress,
		Comment:         req.Comment,
		Subtotal:        subtotal,
		DeliveryFee:     0,
		Total:           total,
		Items:           items,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, userID string) ([]*models.Order, error) {
	return s.orderRepo.FindByUser(ctx, userID)
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID, status string) error {
	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (s *OrderService) GetByID(ctx context.Context, orderID string) (*models.Order, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

func (s *OrderService) ListAll(ctx context.Context) ([]*models.Order, error) {
	return s.orderRepo.FindAll(ctx)
}
