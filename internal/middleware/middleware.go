package middleware

import (
	hotelcalifornia "hotel-california-backend"
)

type Middleware func(hotelcalifornia.Service) hotelcalifornia.Service
