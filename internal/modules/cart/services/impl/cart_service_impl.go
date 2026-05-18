package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/services"
	productRepo "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type newCartServiceImpl struct {
	repo        repository.CartRepository
	productRepo productRepo.ProductRepository
}

func NewNewCartService(repo repository.CartRepository, productRepo productRepo.ProductRepository) services.CartService {
	return &newCartServiceImpl{repo: repo, productRepo: productRepo}
}

func canonicalize(input []delivery.SelectedVariantInput) models.SelectedVariants {
	out := make(models.SelectedVariants, 0, len(input))
	for _, v := range input {
		out = append(out, models.SelectedVariant{
			Value: v.Value,
			Name:  v.Name,
			Type:  v.Type,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Type != out[j].Type {
			return out[i].Type < out[j].Type
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Value < out[j].Value
	})
	return out
}

func (s *newCartServiceImpl) getOrCreateCart(ctx context.Context, userID uint) (*models.Cart, error) {
	cart, err := s.repo.FindCartByUserID(ctx, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cart = &models.Cart{UserID: userID}
	if err := s.repo.CreateCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}
	return cart, nil
}

func (s *newCartServiceImpl) GetCart(ctx context.Context, userID uint) (*delivery.CartResponse, error) {
	cart, err := s.repo.FindCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &delivery.CartResponse{
				UserID: userID,
				Items:  []delivery.CartItemResponse{},
			}, nil
		}
		return nil, err
	}

	full, err := s.repo.FindCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	return delivery.ToCartResponse(full), nil
}

func (s *newCartServiceImpl) AddItem(ctx context.Context, userID uint, req *delivery.AddItemRequest) (*delivery.CartResponse, error) {
	if _, err := s.productRepo.FindByID(ctx, req.ProductID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	variants := canonicalize(req.SelectedVariants)
	variantsJSON, err := json.Marshal(variants)
	if err != nil {
		return nil, fmt.Errorf("failed to encode selected_variants: %w", err)
	}

	var cartID uuid.UUID
	err = s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		cart, err := s.getOrCreateCart(txCtx, userID)
		if err != nil {
			return err
		}
		cartID = cart.ID

		existing, err := s.repo.FindItemByMatch(txCtx, cart.ID, req.ProductID, string(variantsJSON))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existing != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			existing.Quantity += req.Quantity
			return s.repo.UpdateItem(txCtx, existing)
		}

		item := &models.CartItem{
			CartID:           cart.ID,
			ProductID:        req.ProductID,
			Quantity:         req.Quantity,
			SelectedVariants: variants,
		}
		return s.repo.CreateItem(txCtx, item)
	})
	if err != nil {
		return nil, err
	}

	full, err := s.repo.FindCartWithItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	return delivery.ToCartResponse(full), nil
}

func (s *newCartServiceImpl) UpdateItemQuantity(ctx context.Context, userID uint, itemID uuid.UUID, req *delivery.UpdateItemQuantityRequest) (*delivery.CartResponse, error) {
	cart, err := s.repo.FindCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, err
	}

	item, err := s.repo.FindItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart item not found")
		}
		return nil, err
	}

	if item.CartID != cart.ID {
		return nil, fmt.Errorf("cart item not found")
	}

	item.Quantity = req.Quantity
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	full, err := s.repo.FindCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	return delivery.ToCartResponse(full), nil
}

func (s *newCartServiceImpl) DeleteItem(ctx context.Context, userID uint, itemID uuid.UUID) error {
	cart, err := s.repo.FindCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart not found")
		}
		return err
	}

	item, err := s.repo.FindItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart item not found")
		}
		return err
	}

	if item.CartID != cart.ID {
		return fmt.Errorf("cart item not found")
	}

	if err := s.repo.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (s *newCartServiceImpl) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.repo.FindCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if err := s.repo.DeleteItemsByCartID(ctx, cart.ID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}
