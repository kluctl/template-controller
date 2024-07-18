package controllers

import (
	"context"
	"errors"
	"fmt"
	templatesv1alpha1 "github.com/kluctl/template-controller/api/v1alpha1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"
)

type watchesUtil struct {
	watchesCtx context.Context
	watchChan  chan event.TypedGenericEvent[client.ObjectKey]
	watches    map[client.ObjectKey]*watchesForTemplate
	mutex      sync.Mutex
}

type watchesForTemplate struct {
	watchesUtil    *watchesUtil
	templateKey    client.ObjectKey
	serviceAccount string
	client         client.WithWatch
	watches        map[templatesv1alpha1.ObjectRef]watch.Interface
	mutex          sync.Mutex
}

func (wu *watchesUtil) init(watchesCtx context.Context, controller controller.Controller) error {
	wu.watchesCtx = watchesCtx
	wu.watchChan = make(chan event.TypedGenericEvent[client.ObjectKey])
	wu.watches = make(map[client.ObjectKey]*watchesForTemplate)
	watchesSource := source.Channel(wu.watchChan, handler.TypedEnqueueRequestsFromMapFunc[client.ObjectKey](func(ctx context.Context, x client.ObjectKey) []reconcile.Request {
		return []reconcile.Request{{NamespacedName: x}}
	}))
	err := controller.Watch(watchesSource)
	if err != nil {
		return err
	}
	return nil
}

func (wu *watchesUtil) getWatchesForTemplate(templateKey client.ObjectKey) *watchesForTemplate {
	wu.mutex.Lock()
	defer wu.mutex.Unlock()

	wt := wu.watches[templateKey]
	if wt == nil {
		wt = &watchesForTemplate{
			watchesUtil: wu,
			templateKey: templateKey,
			watches:     map[templatesv1alpha1.ObjectRef]watch.Interface{},
		}
		wu.watches[templateKey] = wt
	}

	return wt
}

func (wu *watchesUtil) removeWatchesForTemplate(ctx context.Context, templateKey client.ObjectKey) {
	wu.mutex.Lock()
	defer wu.mutex.Unlock()

	wt := wu.watches[templateKey]
	if wt != nil {
		wt.removeDeletedWatches(ctx, nil)
		delete(wu.watches, templateKey)
	}
}

func (wt *watchesForTemplate) setClient(ctx context.Context, objClient client.WithWatch, serviceAccount string) {
	if wt.serviceAccount != serviceAccount {
		wt.removeDeletedWatches(ctx, nil)
		wt.serviceAccount = serviceAccount
	}
	wt.client = objClient
}

func (wt *watchesForTemplate) addWatchForObject(ctx context.Context, objectRef templatesv1alpha1.ObjectRef) error {
	logger := log.FromContext(ctx)

	wt.mutex.Lock()
	defer wt.mutex.Unlock()

	w := wt.watches[objectRef]
	if w != nil {
		return nil
	}

	logger.V(1).Info("Starting watch for object", "templateKey", wt.templateKey, "objectRef", objectRef)

	gvk, err := objectRef.GroupVersionKind()
	if err != nil {
		return err
	}

	// this is a single-object watch that does NOT require global watch permissions!
	var dummy unstructured.UnstructuredList
	dummy.SetGroupVersionKind(gvk)
	w, err = wt.client.Watch(wt.watchesUtil.watchesCtx, &dummy, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(v1.ObjectNameField, objectRef.Name),
		Namespace:     objectRef.Namespace,
	})
	if err != nil {
		logger.Info("Failed to start watch for object", "templateKey", wt.templateKey, "objectRef", objectRef, "error", err.Error())
		var err2 *errors2.StatusError
		if errors.As(err, &err2) {
			if err2.ErrStatus.Code == http.StatusForbidden {
				err = fmt.Errorf("watch for %s \"%s\" is forbidden: %w", gvk.Kind, objectRef.Name, err)
			}
		}
		return err
	}
	wt.watches[objectRef] = w

	go func() {
		for range w.ResultChan() {
			wt.watchesUtil.watchChan <- event.TypedGenericEvent[client.ObjectKey]{
				Object: wt.templateKey,
			}
		}
	}()

	return nil
}

func (wt *watchesForTemplate) removeDeletedWatches(ctx context.Context, newRefs map[templatesv1alpha1.ObjectRef]struct{}) {
	logger := log.FromContext(ctx)

	wt.mutex.Lock()
	defer wt.mutex.Unlock()

	for k, w := range wt.watches {
		if _, ok := newRefs[k]; !ok {
			logger.V(1).Info("Stopping watch for object", "templateKey", wt.templateKey, "objectRef", k)
			w.Stop()
			delete(wt.watches, k)
		}
	}
}
