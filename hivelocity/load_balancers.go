package hivelocity

import (
	"context"
	"errors"
	"fmt"

	"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/annotation"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/metrics"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

// LoadBalancerOps defines the Load Balancer related operations required by
// the hv-cloud-controller-manager.
type LoadBalancerOps interface {
	GetByName(ctx context.Context, name string) (*hvlb.LoadBalancer, error)
	GetByID(ctx context.Context, id int) (*hvlb.LoadBalancer, error)
	GetByK8SServiceUID(ctx context.Context, svc *v1.Service) (*hvlb.LoadBalancer, error)
	Create(ctx context.Context, lbName string, service *v1.Service) (*hvlb.LoadBalancer, error)
	Delete(ctx context.Context, lb *hvlb.LoadBalancer) error
	ReconcileHCLB(ctx context.Context, lb *hvlb.LoadBalancer, svc *v1.Service) (bool, error)
	ReconcileHCLBTargets(ctx context.Context, lb *hvlb.LoadBalancer, svc *v1.Service, nodes []*v1.Node) (bool, error)
	ReconcileHCLBServices(ctx context.Context, lb *hvlb.LoadBalancer, svc *v1.Service) (bool, error)
}

type loadBalancers struct {
	lbOps                        LoadBalancerOps
	disablePrivateIngressDefault bool
	disableIPv6Default           bool
}

/*
func newLoadBalancers(lbOps LoadBalancerOps, ac hcops.HCloudActionClient, disablePrivateIngressDefault bool, disableIPv6Default bool) *loadBalancers {
	return &loadBalancers{
		lbOps:                        lbOps,
		ac:                           ac,
		disablePrivateIngressDefault: disablePrivateIngressDefault,
		disableIPv6Default:           disableIPv6Default,
	}
}
*/

func (l *loadBalancers) GetLoadBalancer(
	ctx context.Context, _ string, service *v1.Service,
) (status *v1.LoadBalancerStatus, exists bool, err error) {
	const op = "hv/loadBalancers.GetLoadBalancer"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	lb, err := l.lbOps.GetByK8SServiceUID(ctx, service)
	if err != nil {
		if errors.Is(err, hcops.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("%s: %v", op, err)
	}

	if v, ok := annotation.LBHostname.StringFromService(service); ok {
		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{{Hostname: v}},
		}, true, nil
	}

	ingresses := []v1.LoadBalancerIngress{
		{
			IP: lb.PublicNet.IPv4.IP.String(),
		},
	}

	disableIPV6, err := l.getDisableIPv6(service)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %v", op, err)
	}
	if !disableIPV6 {
		ingresses = append(ingresses, v1.LoadBalancerIngress{
			IP: lb.PublicNet.IPv6.IP.String(),
		})
	}

	return &v1.LoadBalancerStatus{Ingress: ingresses}, true, nil
}

func (l *loadBalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	if v, ok := annotation.LBName.StringFromService(service); ok {
		return v
	}
	return cloudprovider.DefaultLoadBalancerName(service)
}

func (l *loadBalancers) EnsureLoadBalancer(
	ctx context.Context, clusterName string, svc *v1.Service, nodes []*v1.Node,
) (*v1.LoadBalancerStatus, error) {
	return nil, fmt.Errorf("TODO implement EnsureLoadBalancer()")
}

func (l *loadBalancers) getDisablePrivateIngress(svc *v1.Service) (bool, error) {
	disable, err := annotation.LBDisablePrivateIngress.BoolFromService(svc)
	if err == nil {
		return disable, nil
	}
	if errors.Is(err, annotation.ErrNotSet) {
		return l.disablePrivateIngressDefault, nil
	}
	return false, err
}

func (l *loadBalancers) getDisableIPv6(svc *v1.Service) (bool, error) {
	disable, err := annotation.LBIPv6Disabled.BoolFromService(svc)
	if err == nil {
		return disable, nil
	}
	if errors.Is(err, annotation.ErrNotSet) {
		return l.disableIPv6Default, nil
	}
	return false, err
}

func (l *loadBalancers) UpdateLoadBalancer(
	ctx context.Context, clusterName string, svc *v1.Service, nodes []*v1.Node,
) error {
	const op = "hv/loadBalancers.UpdateLoadBalancer"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var (
		lb  *hvlb.LoadBalancer
		err error
	)

	nodeNames := make([]string, len(nodes))
	for i, n := range nodes {
		nodeNames[i] = n.Name
	}
	klog.InfoS("update Load Balancer", "op", op, "service", svc.Name, "nodes", nodeNames)

	lb, err = l.lbOps.GetByK8SServiceUID(ctx, svc)
	if errors.Is(err, hcops.ErrNotFound) {
		lbName := l.GetLoadBalancerName(ctx, clusterName, svc)

		lb, err = l.lbOps.GetByName(ctx, lbName)
		if errors.Is(err, hcops.ErrNotFound) {
			return nil
		}
		// further error types handled below
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = l.lbOps.ReconcileHCLB(ctx, lb, svc); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if _, err = l.lbOps.ReconcileHCLBTargets(ctx, lb, svc, nodes); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if _, err = l.lbOps.ReconcileHCLBServices(ctx, lb, svc); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (l *loadBalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	const op = "hv/loadBalancers.EnsureLoadBalancerDeleted"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	loadBalancer, err := l.lbOps.GetByK8SServiceUID(ctx, service)
	if errors.Is(err, hcops.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if loadBalancer.Protection.Delete {
		klog.InfoS("ignored: load balancer deletion protected", "op", op, "loadBalancerID", loadBalancer.ID)
		return nil
	}

	klog.InfoS("delete Load Balancer", "op", op, "loadBalancerID", loadBalancer.ID)
	err = l.lbOps.Delete(ctx, loadBalancer)
	if errors.Is(err, hcops.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
