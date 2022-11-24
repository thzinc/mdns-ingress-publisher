package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/syncromatics/go-kit/v2/log"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	TTSAnnotation string = "mdns-ingress-publisher/tts"
)

type Watcher interface {
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type Zone interface {
	Publish(record string) error
	Unpublish(record string) error
}

func NewWatcher(ctx context.Context, ingresses Watcher, zone Zone, defaultTTS int) func() error {
	observedIngresses := map[types.UID][]string{}
	return func() error {
		wi, err := ingresses.Watch(ctx, metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to watch Ingresses")
		}
		defer wi.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-wi.ResultChan():
				if !ok {
					return nil
				}

				if event.Object == nil {
					continue
				}

				ingress := event.Object.(*networkingv1.Ingress)

				log.Debug("handling event",
					"type", event.Type,
					"uid", ingress.UID,
					"namespace", ingress.ObjectMeta.Namespace,
					"name", ingress.ObjectMeta.Name)

				switch event.Type {
				case watch.Added:
					fallthrough
				case watch.Modified:
					fallthrough
				case watch.Deleted:
					oldRecords, ok := observedIngresses[ingress.UID]
					if ok {
						delete(observedIngresses, ingress.UID)
						for _, record := range oldRecords {
							log.Info("unpublishing record",
								"type", event.Type,
								"uid", ingress.UID,
								"namespace", ingress.ObjectMeta.Namespace,
								"name", ingress.ObjectMeta.Name,
								"record", record)

							err = zone.Unpublish(record)
							if err != nil {
								return errors.Wrap(err, "failed to unpublish record")
							}
						}
					}
				}

				switch event.Type {
				case watch.Added:
					fallthrough
				case watch.Modified:
					addresses := []string{}
					for _, lbIngress := range ingress.Status.LoadBalancer.Ingress {
						addresses = append(addresses, lbIngress.IP)
					}

					if len(addresses) == 0 {
						log.Warn("skipping ingress because it is not ready; missing Status.LoadBalancer.Ingress[].IP",
							"type", event.Type,
							"uid", ingress.UID,
							"namespace", ingress.ObjectMeta.Namespace,
							"name", ingress.ObjectMeta.Name,
							"ingress.Status", ingress.Status)
						continue
					}

					recordData := strings.Join(addresses, " ")

					ttsAnnotation, ok := ingress.Annotations[TTSAnnotation]
					tts := defaultTTS
					if ok {
						i, err := strconv.Atoi(ttsAnnotation)
						if err == nil {
							tts = i
						} else {
							log.Warn("ignoring annotation because it is not a number",
								"type", event.Type,
								"uid", ingress.UID,
								"namespace", ingress.ObjectMeta.Namespace,
								"name", ingress.ObjectMeta.Name,
								TTSAnnotation, ttsAnnotation)
						}
					}

					records := []string{}
					for _, rule := range ingress.Spec.Rules {
						if !strings.HasSuffix(rule.Host, ".local") {
							log.Info("skipping ingress rule; host is not under .local TLD",
								"type", event.Type,
								"uid", ingress.UID,
								"namespace", ingress.ObjectMeta.Namespace,
								"name", ingress.ObjectMeta.Name,
								"rule", rule)
							continue
						}

						record := fmt.Sprintf("%s %d IN A %s", rule.Host, tts, recordData)
						records = append(records, record)
					}

					observedIngresses[ingress.UID] = records
					for _, record := range records {
						log.Info("publishing record",
							"type", event.Type,
							"uid", ingress.UID,
							"namespace", ingress.ObjectMeta.Namespace,
							"name", ingress.ObjectMeta.Name,
							"record", record)
						err = zone.Publish(record)
						if err != nil {
							return errors.Wrap(err, "failed to publish record")
						}
					}
				}
			}
		}
	}
}
