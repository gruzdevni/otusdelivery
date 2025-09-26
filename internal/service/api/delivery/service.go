package delivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"otusdelivery/internal/models"
	query "otusdelivery/internal/repo"
)

var ErrNoAvailableSlots = errors.New("No available slots for delivery")

type repo interface {
	CheckOrderDelivery(ctx context.Context, orderID string) (query.SlotReservation, error)
	DecreaseAvailableSlot(ctx context.Context, slotID int32) error
	GetAvailableSlot(ctx context.Context, slotID int32) (int32, error)
	IncreaseAvailableSlot(ctx context.Context, slotID int32) error
	ReserveSlotForOrder(ctx context.Context, arg query.ReserveSlotForOrderParams) error
	UnreserveSlotForOrder(ctx context.Context, orderID string) error
	WithTx(tx *sql.Tx) *query.Queries
}

type service struct {
	dbRW *sql.DB
	repo repo
}

type Service interface {
	CheckOrderDelivery(ctx context.Context, orderID string) (*models.DeliveryStatus, error)
	ReserveSlotForOrder(ctx context.Context, arg query.ReserveSlotForOrderParams) error
	UnreserveSlotForOrder(ctx context.Context, orderID string) error
}

func NewService(dbRW *sql.DB, repo repo) Service {
	return &service{
		dbRW: dbRW,
		repo: repo,
	}
}

func (s *service) CheckOrderDelivery(ctx context.Context, orderID string) (*models.DeliveryStatus, error) {
	res, err := s.repo.CheckOrderDelivery(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.DeliveryStatus{
				OrderID:        res.OrderID,
				DeliveryStatus: "cancelled",
			}, nil
		}

		return nil, fmt.Errorf("checking order delivery: %w", err)
	}

	return &models.DeliveryStatus{
		OrderID:        res.OrderID,
		DeliveryStatus: "confirmed",
	}, nil
}

func (s *service) ReserveSlotForOrder(ctx context.Context, arg query.ReserveSlotForOrderParams) error {
	tx, err := s.dbRW.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	isCommitted := false
	defer func() {
		if !isCommitted {
			if err := tx.Rollback(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to rollback the transaction")
			}
		}
	}()

	res2, err := s.repo.WithTx(tx).GetAvailableSlot(ctx, arg.SlotID)
	if err != nil {
		return fmt.Errorf("getting available slots after reserving: %w", err)
	}

	if res2 < 0 {
		return ErrNoAvailableSlots
	}

	res, err := s.repo.WithTx(tx).CheckOrderDelivery(ctx, arg.OrderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("checking order delivery: %w", err)
	}

	if res.SlotID == arg.SlotID {
		return nil
	}

	if !lo.IsEmpty(res) {
		if err = s.repo.WithTx(tx).IncreaseAvailableSlot(ctx, res.SlotID); err != nil {
			return fmt.Errorf("increasing slot quantity: %w", err)
		}
	}

	if err = s.repo.WithTx(tx).DecreaseAvailableSlot(ctx, arg.SlotID); err != nil {
		return fmt.Errorf("decreasing slot quantity: %w", err)
	}

	if err = s.repo.WithTx(tx).ReserveSlotForOrder(ctx, query.ReserveSlotForOrderParams{
		OrderID: arg.OrderID,
		SlotID:  arg.SlotID,
	}); err != nil {
		return fmt.Errorf("reserving slot for order: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	isCommitted = true

	return nil
}

func (s *service) UnreserveSlotForOrder(ctx context.Context, orderID string) error {
	tx, err := s.dbRW.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	isCommitted := false
	defer func() {
		if !isCommitted {
			if err := tx.Rollback(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to rollback the transaction")
			}
		}
	}()

	res, err := s.repo.WithTx(tx).CheckOrderDelivery(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("checking order delivery: %w", err)
	}

	if err = s.repo.WithTx(tx).IncreaseAvailableSlot(ctx, res.SlotID); err != nil {
		return fmt.Errorf("increasing slot quantity: %w", err)
	}

	if err = s.repo.WithTx(tx).UnreserveSlotForOrder(ctx, orderID); err != nil {
		return fmt.Errorf("unreserving slot for order: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	isCommitted = true

	return nil
}
