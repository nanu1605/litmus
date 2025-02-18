package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/generated"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/authorization"
	data_store "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/data-store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateChaosExperiment(ctx context.Context, request model.ChaosExperimentRequest, projectID string) (*model.ChaosExperimentResponse, error) {
	logFields := logrus.Fields{
		"projectId": projectID,
	}

	logrus.WithFields(logFields).Info("request received to create chaos experiment")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.CreateEnvironment],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	uiResponse, err := r.chaosExperimentHandler.CreateChaosExperiment(ctx, &request, projectID, data_store.Store)
	if err != nil {
		return nil, errors.New("could not create experiment, error: " + err.Error())
	}
	if request.RunExperiment != nil && *request.RunExperiment {

		query := bson.D{
			{"experiment_id", uiResponse.ExperimentID},
			{"is_removed", false},
		}

		experiment, err := r.chaosExperimentHandler.GetDBExperiment(query)
		if err != nil {
			return nil, errors.New("could not get experiment run, error: " + err.Error())
		}

		if experiment.CronSyntax != "" {

			if err = r.chaosExperimentRunHandler.RunCronExperiment(ctx, projectID, experiment, data_store.Store); err != nil {
				logrus.WithFields(logFields).Error(err)
				return nil, err
			}
			logrus.WithFields(logFields).WithField("workflowId", experiment.ExperimentID).Info("cron experiment created successfully")
			return uiResponse, nil
		}

		_, err = r.chaosExperimentRunHandler.RunChaosWorkFlow(ctx, projectID, experiment, data_store.Store)
		if err != nil {
			logrus.WithFields(logFields).Error(err)
			return nil, err
		}
		return uiResponse, nil

	}
	return uiResponse, nil
}

func (r *mutationResolver) SaveChaosExperiment(ctx context.Context, request model.SaveChaosExperimentRequest, projectID string) (string, error) {
	logFields := logrus.Fields{
		"experimentId": request.ID,
		"projectId":    projectID,
	}
	logrus.WithFields(logFields).Info("request received to save chaos experiment")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.CreateChaosWorkFlow],
		model.InvitationAccepted.String())
	if err != nil {
		return "", err
	}

	var uiResponse string

	uiResponse, err = r.chaosExperimentHandler.SaveChaosExperiment(ctx, request, projectID, data_store.Store)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return "", err
	}

	return uiResponse, nil
}

func (r *mutationResolver) UpdateChaosExperiment(ctx context.Context, request *model.ChaosExperimentRequest, projectID string) (*model.ChaosExperimentResponse, error) {
	logFields := logrus.Fields{
		"projectId":         projectID,
		"chaosExperimentId": request.ExperimentID,
	}

	logrus.WithFields(logFields).Info("request received to update chaos experiment")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ReRunChaosWorkFlow],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	uiResponse, err := r.chaosExperimentHandler.UpdateChaosExperiment(ctx, request, projectID, data_store.Store)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return uiResponse, err
}

func (r *mutationResolver) DeleteChaosExperiment(ctx context.Context, experimentID string, experimentRunID *string, projectID string) (bool, error) {
	logFields := logrus.Fields{
		"projectId":            projectID,
		"chaosExperimentId":    experimentID,
		"chaosExperimentRunId": experimentRunID,
	}

	logrus.WithFields(logFields).Info("request received to delete chaos experiment")

	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.DeleteChaosWorkflow],
		model.InvitationAccepted.String())
	if err != nil {
		return false, err
	}

	uiResponse, err := r.chaosExperimentHandler.DeleteChaosExperiment(ctx, projectID, experimentID, experimentRunID, data_store.Store)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return false, err
	}
	return uiResponse, err
}

func (r *mutationResolver) UpdateCronExperimentState(ctx context.Context, experimentID string, disable bool, projectID string) (bool, error) {
	logFields := logrus.Fields{
		"projectId":         projectID,
		"chaosExperimentId": experimentID,
	}

	logrus.WithFields(logFields).Info("request received to update cron chaos experiment")

	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.UpdateChaosWorkflow],
		model.InvitationAccepted.String())
	if err != nil {
		return false, err
	}

	uiResponse, err := r.chaosExperimentHandler.UpdateCronExperimentState(ctx, experimentID, disable, projectID, data_store.Store)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return false, err
	}
	return uiResponse, err
}

func (r *queryResolver) GetExperiment(ctx context.Context, projectID string, experimentID string) (*model.GetExperimentResponse, error) {
	logFields := logrus.Fields{
		"projectId":         projectID,
		"chaosExperimentId": experimentID,
	}
	logrus.WithFields(logFields).Info("request received to get chaos experiment")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ListWorkflow],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	expResponse, err := r.chaosExperimentHandler.GetExperiment(ctx, projectID, experimentID)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return expResponse, err
}

func (r *queryResolver) ListExperiment(ctx context.Context, projectID string, request model.ListExperimentRequest) (*model.ListExperimentResponse, error) {
	logFields := logrus.Fields{
		"projectId":          projectID,
		"chaosExperimentIds": request.ExperimentIDs,
	}
	logrus.WithFields(logFields).Info("request received to list chaos experiments")

	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ListWorkflow],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}
	uiResponse, err := r.chaosExperimentHandler.ListExperiment(projectID, request)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return uiResponse, err
}

func (r *queryResolver) GetExperimentStats(ctx context.Context, projectID string) (*model.GetExperimentStatsResponse, error) {
	logFields := logrus.Fields{
		"projectId": projectID,
	}
	logrus.WithFields(logFields).Info("request received to get chaos experiment stats")

	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ListWorkflow],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	uiResponse, err := r.chaosExperimentHandler.GetExperimentStats(ctx, projectID)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return uiResponse, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
