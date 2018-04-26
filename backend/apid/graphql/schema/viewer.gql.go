// Code generated by scripts/gengraphql.go. DO NOT EDIT.

package schema

import (
	graphql1 "github.com/graphql-go/graphql"
	mapstructure "github.com/mitchellh/mapstructure"
	graphql "github.com/sensu/sensu-go/graphql"
)

// ViewerEntitiesFieldResolverArgs contains arguments provided to entities when selected
type ViewerEntitiesFieldResolverArgs struct {
	First  int     // First - self descriptive
	Last   int     // Last - self descriptive
	Before *string // Before - self descriptive
	After  *string // After - self descriptive
}

// ViewerEntitiesFieldResolverParams contains contextual info to resolve entities field
type ViewerEntitiesFieldResolverParams struct {
	graphql.ResolveParams
	Args ViewerEntitiesFieldResolverArgs
}

// ViewerEntitiesFieldResolver implement to resolve requests for the Viewer's entities field.
type ViewerEntitiesFieldResolver interface {
	// Entities implements response to request for entities field.
	Entities(p ViewerEntitiesFieldResolverParams) (interface{}, error)
}

// ViewerChecksFieldResolverArgs contains arguments provided to checks when selected
type ViewerChecksFieldResolverArgs struct {
	First  int     // First - self descriptive
	Last   int     // Last - self descriptive
	Before *string // Before - self descriptive
	After  *string // After - self descriptive
}

// ViewerChecksFieldResolverParams contains contextual info to resolve checks field
type ViewerChecksFieldResolverParams struct {
	graphql.ResolveParams
	Args ViewerChecksFieldResolverArgs
}

// ViewerChecksFieldResolver implement to resolve requests for the Viewer's checks field.
type ViewerChecksFieldResolver interface {
	// Checks implements response to request for checks field.
	Checks(p ViewerChecksFieldResolverParams) (interface{}, error)
}

// ViewerOrganizationsFieldResolver implement to resolve requests for the Viewer's organizations field.
type ViewerOrganizationsFieldResolver interface {
	// Organizations implements response to request for organizations field.
	Organizations(p graphql.ResolveParams) (interface{}, error)
}

// ViewerUserFieldResolver implement to resolve requests for the Viewer's user field.
type ViewerUserFieldResolver interface {
	// User implements response to request for user field.
	User(p graphql.ResolveParams) (interface{}, error)
}

//
// ViewerFieldResolvers represents a collection of methods whose products represent the
// response values of the 'Viewer' type.
//
// == Example SDL
//
//   """
//   Dog's are not hooman.
//   """
//   type Dog implements Pet {
//     "name of this fine beast."
//     name:  String!
//
//     "breed of this silly animal; probably shibe."
//     breed: [Breed]
//   }
//
// == Example generated interface
//
//   // DogResolver ...
//   type DogFieldResolvers interface {
//     DogNameFieldResolver
//     DogBreedFieldResolver
//
//     // IsTypeOf is used to determine if a given value is associated with the Dog type
//     IsTypeOf(interface{}, graphql.IsTypeOfParams) bool
//   }
//
// == Example implementation ...
//
//   // DogResolver implements DogFieldResolvers interface
//   type DogResolver struct {
//     logger logrus.LogEntry
//     store interface{
//       store.BreedStore
//       store.DogStore
//     }
//   }
//
//   // Name implements response to request for name field.
//   func (r *DogResolver) Name(p graphql.ResolveParams) (interface{}, error) {
//     // ... implementation details ...
//     dog := p.Source.(DogGetter)
//     return dog.GetName()
//   }
//
//   // Breed implements response to request for breed field.
//   func (r *DogResolver) Breed(p graphql.ResolveParams) (interface{}, error) {
//     // ... implementation details ...
//     dog := p.Source.(DogGetter)
//     breed := r.store.GetBreed(dog.GetBreedName())
//     return breed
//   }
//
//   // IsTypeOf is used to determine if a given value is associated with the Dog type
//   func (r *DogResolver) IsTypeOf(p graphql.IsTypeOfParams) bool {
//     // ... implementation details ...
//     _, ok := p.Value.(DogGetter)
//     return ok
//   }
//
type ViewerFieldResolvers interface {
	ViewerEntitiesFieldResolver
	ViewerChecksFieldResolver
	ViewerOrganizationsFieldResolver
	ViewerUserFieldResolver
}

// ViewerAliases implements all methods on ViewerFieldResolvers interface by using reflection to
// match name of field to a field on the given value. Intent is reduce friction
// of writing new resolvers by removing all the instances where you would simply
// have the resolvers method return a field.
//
// == Example SDL
//
//    type Dog {
//      name:   String!
//      weight: Float!
//      dob:    DateTime
//      breed:  [Breed]
//    }
//
// == Example generated aliases
//
//   type DogAliases struct {}
//   func (_ DogAliases) Name(p graphql.ResolveParams) (interface{}, error) {
//     // reflect...
//   }
//   func (_ DogAliases) Weight(p graphql.ResolveParams) (interface{}, error) {
//     // reflect...
//   }
//   func (_ DogAliases) Dob(p graphql.ResolveParams) (interface{}, error) {
//     // reflect...
//   }
//   func (_ DogAliases) Breed(p graphql.ResolveParams) (interface{}, error) {
//     // reflect...
//   }
//
// == Example Implementation
//
//   type DogResolver struct { // Implements DogResolver
//     DogAliases
//     store store.BreedStore
//   }
//
//   // NOTE:
//   // All other fields are satisified by DogAliases but since this one
//   // requires hitting the store we implement it in our resolver.
//   func (r *DogResolver) Breed(p graphql.ResolveParams) interface{} {
//     dog := v.(*Dog)
//     return r.BreedsById(dog.BreedIDs)
//   }
//
type ViewerAliases struct{}

// Entities implements response to request for 'entities' field.
func (_ ViewerAliases) Entities(p ViewerEntitiesFieldResolverParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// Checks implements response to request for 'checks' field.
func (_ ViewerAliases) Checks(p ViewerChecksFieldResolverParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// Organizations implements response to request for 'organizations' field.
func (_ ViewerAliases) Organizations(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// User implements response to request for 'user' field.
func (_ ViewerAliases) User(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// ViewerType Describes a viewer of the system; generally an authenticated user.
var ViewerType = graphql.NewType("Viewer", graphql.ObjectKind)

// RegisterViewer registers Viewer object type with given service.
func RegisterViewer(svc *graphql.Service, impl ViewerFieldResolvers) {
	svc.RegisterObject(_ObjectTypeViewerDesc, impl)
}
func _ObjTypeViewerEntitiesHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(ViewerEntitiesFieldResolver)
	return func(p graphql1.ResolveParams) (interface{}, error) {
		frp := ViewerEntitiesFieldResolverParams{ResolveParams: p}
		err := mapstructure.Decode(p.Args, &frp.Args)
		if err != nil {
			return nil, err
		}

		return resolver.Entities(frp)
	}
}

func _ObjTypeViewerChecksHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(ViewerChecksFieldResolver)
	return func(p graphql1.ResolveParams) (interface{}, error) {
		frp := ViewerChecksFieldResolverParams{ResolveParams: p}
		err := mapstructure.Decode(p.Args, &frp.Args)
		if err != nil {
			return nil, err
		}

		return resolver.Checks(frp)
	}
}

func _ObjTypeViewerOrganizationsHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(ViewerOrganizationsFieldResolver)
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Organizations(frp)
	}
}

func _ObjTypeViewerUserHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(ViewerUserFieldResolver)
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.User(frp)
	}
}

func _ObjectTypeViewerConfigFn() graphql1.ObjectConfig {
	return graphql1.ObjectConfig{
		Description: "Describes a viewer of the system; generally an authenticated user.",
		Fields: graphql1.Fields{
			"checks": &graphql1.Field{
				Args: graphql1.FieldConfigArgument{
					"after": &graphql1.ArgumentConfig{
						Description: "self descriptive",
						Type:        graphql1.String,
					},
					"before": &graphql1.ArgumentConfig{
						Description: "self descriptive",
						Type:        graphql1.String,
					},
					"first": &graphql1.ArgumentConfig{
						DefaultValue: 10,
						Description:  "self descriptive",
						Type:         graphql1.Int,
					},
					"last": &graphql1.ArgumentConfig{
						DefaultValue: 10,
						Description:  "self descriptive",
						Type:         graphql1.Int,
					},
				},
				DeprecationReason: "",
				Description:       "All check configurations the viewer has access to view.",
				Name:              "checks",
				Type:              graphql.OutputType("CheckConfigConnection"),
			},
			"entities": &graphql1.Field{
				Args: graphql1.FieldConfigArgument{
					"after": &graphql1.ArgumentConfig{
						Description: "self descriptive",
						Type:        graphql1.String,
					},
					"before": &graphql1.ArgumentConfig{
						Description: "self descriptive",
						Type:        graphql1.String,
					},
					"first": &graphql1.ArgumentConfig{
						DefaultValue: 10,
						Description:  "self descriptive",
						Type:         graphql1.Int,
					},
					"last": &graphql1.ArgumentConfig{
						DefaultValue: 10,
						Description:  "self descriptive",
						Type:         graphql1.Int,
					},
				},
				DeprecationReason: "",
				Description:       "All entities the viewer has access to view.",
				Name:              "entities",
				Type:              graphql.OutputType("EntityConnection"),
			},
			"organizations": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "All organizations the viewer has access to view.",
				Name:              "organizations",
				Type:              graphql1.NewNonNull(graphql1.NewList(graphql1.NewNonNull(graphql.OutputType("Organization")))),
			},
			"user": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "User account associated with the viewer.",
				Name:              "user",
				Type:              graphql.OutputType("User"),
			},
		},
		Interfaces: []*graphql1.Interface{},
		IsTypeOf: func(_ graphql1.IsTypeOfParams) bool {
			// NOTE:
			// Panic by default. Intent is that when Service is invoked, values of
			// these fields are updated with instantiated resolvers. If these
			// defaults are called it is most certainly programmer err.
			// If you're see this comment then: 'Whoops! Sorry, my bad.'
			panic("Unimplemented; see ViewerFieldResolvers.")
		},
		Name: "Viewer",
	}
}

// describe Viewer's configuration; kept private to avoid unintentional tampering of configuration at runtime.
var _ObjectTypeViewerDesc = graphql.ObjectDesc{
	Config: _ObjectTypeViewerConfigFn,
	FieldHandlers: map[string]graphql.FieldHandler{
		"checks":        _ObjTypeViewerChecksHandler,
		"entities":      _ObjTypeViewerEntitiesHandler,
		"organizations": _ObjTypeViewerOrganizationsHandler,
		"user":          _ObjTypeViewerUserHandler,
	},
}
