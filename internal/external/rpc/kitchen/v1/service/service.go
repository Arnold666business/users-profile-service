package service

import (
	"context"
	"time"
	kitchenpb "users-profile-service/api/generation/kitchen/v1/service"
	"users-profile-service/internal/user/delete"
)

func (g *GrpcKitchenClient) GetKitchenActiveBookingsByOwnerId(ctx context.Context, id int64) []delete.Booking {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, _ := g.client.GetKitchenActiveBookingsByOwnerId(ctx, &kitchenpb.GetBookingsRequest{
		OwnerId: id,
	})
	return bookingsMapper(response)
}

// todo
func (g *GrpcKitchenClient) UnPublishKitchenByOwerId(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	g.client.UnPublishKitchenByOwnerId(ctx, &kitchenpb.UnPublishKitchenRequest{
		OwnerId: id,
	})
	return nil
}

func bookingsMapper(req *kitchenpb.GetBookingsResponse) []delete.Booking {
	res := make([]delete.Booking, len(req.GetBookings()))
	for _, booking := range req.GetBookings() {
		res = append(res, delete.Booking{
			BookingId: booking.BookingId,
			KitchenId: booking.KitchenId,
		})
	}
	return res
}
