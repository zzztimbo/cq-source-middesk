package services

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/zzztimbo/cq-source-middesk/client"
	"github.com/zzztimbo/cq-source-middesk/internal/middesk"
	"golang.org/x/sync/errgroup"
)

func BusinessesTable() *schema.Table {
	return &schema.Table{
		Name:      "businesses",
		Resolver:  fetchBusinessesTable,
		Transform: transformers.TransformWithStruct(&middesk.MiddeskBusiness{}),
	}
}

func fetchBusinessesTable(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)

	pageNumber := 1
	perPage := 30 // this is the max value

	g := errgroup.Group{}
	g.SetLimit(perPage)

	for {
		businesses, err := cl.MiddeskClient.GetMiddeskBusinessIds(pageNumber, perPage)
		if err != nil {
			return err
		}
		for _, business := range businesses.Data {
			business := business
			g.Go(func() error {
				rawBusiness, err := cl.MiddeskClient.GetMiddeskBusiness(business.Id)
				if err != nil {
					cl.Logger().Error().Msgf("business id: %s", business.Id)
					return err
				}
				middeskBusiness := middesk.MiddeskBusiness{
					Id:                       rawBusiness.Id,
					ExternalId:               rawBusiness.ExternalId,
					CreatedAt:                rawBusiness.CreatedAt,
					Status:                   rawBusiness.Status,
					FormationEntityType:      rawBusiness.Formation.EntityType,
					RegistrationsEntityTypes: make([]string, len(rawBusiness.Registrations)),
				}
				for i, registration := range rawBusiness.Registrations {
					middeskBusiness.RegistrationsEntityTypes[i] = registration.EntityType
				}
				res <- middeskBusiness
				return nil
			})

		}
		if !businesses.HasMore {
			break
		}
		pageNumber++
	}

	return g.Wait()
}
