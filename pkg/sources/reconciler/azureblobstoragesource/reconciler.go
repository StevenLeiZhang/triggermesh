/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azureblobstoragesource

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/reconciler"

	"github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	reconcilerv1alpha1 "github.com/triggermesh/triggermesh/pkg/client/generated/injection/reconciler/sources/v1alpha1/azureblobstoragesource"
	listersv1alpha1 "github.com/triggermesh/triggermesh/pkg/client/generated/listers/sources/v1alpha1"
	"github.com/triggermesh/triggermesh/pkg/sources/client/azure/storage"
	"github.com/triggermesh/triggermesh/pkg/sources/reconciler/common"
	"github.com/triggermesh/triggermesh/pkg/sources/reconciler/common/event"
)

// Reconciler implements controller.Reconciler for the event source type.
type Reconciler struct {
	// Getter than can obtain clients for interacting with Azure APIs
	cg storage.ClientGetter

	srcLister func(namespace string) listersv1alpha1.AzureBlobStorageSourceNamespaceLister

	// Event Hubs adapter
	base       common.GenericDeploymentReconciler
	adapterCfg *adapterConfig
}

// Check that our Reconciler implements Interface
var _ reconcilerv1alpha1.Interface = (*Reconciler)(nil)

// Check that our Reconciler implements Finalizer
var _ reconcilerv1alpha1.Finalizer = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *v1alpha1.AzureBlobStorageSource) reconciler.Event {
	// inject source into context for usage in reconciliation logic
	ctx = v1alpha1.WithSource(ctx, o)

	eventSubsCli, eventHubsCli, err := r.cg.Get(o)
	switch {
	case isNoCredentials(err):
		o.Status.MarkNotSubscribed(v1alpha1.AzureReasonNoClient, "Azure credentials missing: "+toErrMsg(err))
		return controller.NewPermanentError(reconciler.NewEvent(corev1.EventTypeWarning, ReasonFailedSubscribe,
			"Azure credentials missing: %s", toErrMsg(err)))
	case err != nil:
		o.Status.MarkNotSubscribed(v1alpha1.AzureReasonNoClient, "Error obtaining Azure clients: "+toErrMsg(err))
		// wrap any other error to fail the reconciliation
		return fmt.Errorf("%w", reconciler.NewEvent(corev1.EventTypeWarning, ReasonFailedSubscribe,
			"Error obtaining Azure clients: %s", err))
	}

	eventHubResID, err := ensureEventHub(ctx, eventHubsCli)
	if err != nil {
		return fmt.Errorf("failed to reconcile Event Hub: %w", err)
	}

	if err := r.base.ReconcileSource(ctx, r); err != nil {
		return fmt.Errorf("failed to reconcile Event Hubs event source adapter: %w", err)
	}

	return ensureEventSubscription(ctx, eventSubsCli, eventHubResID)
}

// FinalizeKind is called when the resource is deleted.
func (r *Reconciler) FinalizeKind(ctx context.Context, o *v1alpha1.AzureBlobStorageSource) reconciler.Event {
	// inject source into context for usage in finalization logic
	ctx = v1alpha1.WithSource(ctx, o)

	eventSubsCli, eventHubsCli, err := r.cg.Get(o)
	switch {
	case isNoCredentials(err):
		// the finalizer is unlikely to recover from missing
		// credentials, so we simply record a warning event and return
		event.Warn(ctx, ReasonFailedUnsubscribe, "Azure credentials missing while finalizing event source. "+
			"Ignoring: %s", err)
		return nil
	case err != nil:
		return reconciler.NewEvent(corev1.EventTypeWarning, ReasonFailedUnsubscribe,
			"Error creating Azure clients: %s", err)
	}

	// The finalizer blocks the deletion of the source object until the
	// deletion of the event subscription and Event Hub succeed to ensure
	// that we don't leave any dangling resources behind us.

	if err := ensureNoEventSubscription(ctx, eventSubsCli); err != nil {
		return fmt.Errorf("failed to finalize event subscription: %w", err)
	}

	return ensureNoEventHub(ctx, eventHubsCli)
}

// isNoCredentials returns whether the given error indicates that some required
// Azure credentials could not be obtained.
func isNoCredentials(err error) bool {
	// consider that missing Secrets indicate missing credentials in this context
	if k8sErr := apierrors.APIStatus(nil); errors.As(err, &k8sErr) {
		return k8sErr.Status().Reason == metav1.StatusReasonNotFound
	}
	if ecErr := (interface{ IsEmptyCredentials() })(nil); errors.As(err, &ecErr) {
		return true
	}
	return false
}