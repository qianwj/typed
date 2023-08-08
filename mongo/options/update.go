package options

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArrayFilters interface {
	Raw() options.ArrayFilters
}
