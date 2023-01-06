package annotation

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/metrics"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	v1 "k8s.io/api/core/v1"
)

// ErrNotSet signals that an annotation was not set.
var ErrNotSet = errors.New("not set")

// Name defines the name of a K8S annotation.
type Name string

// AnnotateService adds the value v as an annotation with s.Name to svc.
//
// AnnotateService returns an error if converting v to a string fails.
func (s Name) AnnotateService(svc *v1.Service, v interface{}) error {
	const op = "annotation/Name.AnnotateService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	if svc.ObjectMeta.Annotations == nil {
		svc.ObjectMeta.Annotations = make(map[string]string)
	}
	k := string(s)
	switch vt := v.(type) {
	case bool:
		svc.ObjectMeta.Annotations[k] = strconv.FormatBool(vt)
	case int:
		svc.ObjectMeta.Annotations[k] = strconv.Itoa(vt)
	case string:
		svc.ObjectMeta.Annotations[k] = vt
	case []string:
		svc.ObjectMeta.Annotations[k] = strings.Join(vt, ",")
	case hv.CertificateType:
		svc.ObjectMeta.Annotations[k] = string(vt)
	case []*hv.Certificate:
		idsOrNames := make([]string, len(vt))
		for i, c := range vt {
			if c.ID == 0 && c.Name != "" {
				idsOrNames[i] = c.Name
				continue
			}
			idsOrNames[i] = strconv.Itoa(c.ID)
		}
		svc.ObjectMeta.Annotations[k] = strings.Join(idsOrNames, ",")
	case hv.NetworkZone:
		svc.ObjectMeta.Annotations[k] = string(vt)
	case hvlb.LoadBalancerAlgorithmType:
		svc.ObjectMeta.Annotations[k] = string(vt)
	case hvlb.LoadBalancerServiceProtocol:
		svc.ObjectMeta.Annotations[k] = string(vt)
	case fmt.Stringer:
		svc.ObjectMeta.Annotations[k] = vt.String()
	default:
		return fmt.Errorf("%s: %v: unsupported type: %T", op, s, v)
	}
	return nil
}

// StringFromService retrieves the value belonging to the annotation from svc.
//
// If svc has no value for the annotation the second return value is false.
func (s Name) StringFromService(svc *v1.Service) (string, bool) {
	if svc.Annotations == nil {
		return "", false
	}
	v, ok := svc.Annotations[string(s)]
	return v, ok
}

// StringsFromService retrieves the []string value belonging to the annotation
// from svc.
//
// StringsFromService returns ErrNotSet annotation was not set.
func (s Name) StringsFromService(svc *v1.Service) ([]string, error) {
	const op = "annotation/Name.StringsFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var ss []string

	err := s.applyToValue(op, svc, func(v string) error {
		ss = strings.Split(v, ",")
		return nil
	})

	return ss, err
}

// BoolFromService retrieves the boolean value belonging to the annotation from
// svc.
//
// BoolFromService returns an error if the value could not be converted to a
// boolean, or the annotation was not set. In the case of a missing value, the
// error wraps ErrNotSet.
func (s Name) BoolFromService(svc *v1.Service) (bool, error) {
	const op = "annotation/Name.BoolFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	v, ok := s.StringFromService(svc)
	if !ok {
		return false, fmt.Errorf("%s: %v: %w", op, s, ErrNotSet)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s: %v: %w", op, s, err)
	}
	return b, nil
}

// IntFromService retrieves the int value belonging to the annotation from svc.
//
// IntFromService returns an error if the value could not be converted to an
// int, or the annotation was not set. In the case of a missing value, the
// error wraps ErrNotSet.
func (s Name) IntFromService(svc *v1.Service) (int, error) {
	const op = "annotation/Name.IntFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	v, ok := s.StringFromService(svc)
	if !ok {
		return 0, fmt.Errorf("%s: %v: %w", op, s, ErrNotSet)
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %v: %w", op, s, err)
	}
	return i, nil
}

// IntsFromService retrieves the []int value belonging to the annotation from
// svc.
//
// IntsFromService returns an error if the value could not be converted to a
// []int, or the annotation was not set. In the case of a missing value, the
// error wraps ErrNotSet.
func (s Name) IntsFromService(svc *v1.Service) ([]int, error) {
	const op = "annotation/Name.IntsFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var is []int

	err := s.applyToValue(op, svc, func(v string) error {
		ss := strings.Split(v, ",")
		is = make([]int, len(ss))

		for i, s := range ss {
			iv, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			is[i] = iv
		}
		return nil
	})

	return is, err
}

// IPFromService retrieves the net.IP value belonging to the annotation from
// svc.
//
// IPFromService returns an error if the value could not be converted to a
// net.IP, or the annotation was not set. In the case of a missing value, the
// error wraps ErrNotSet.
func (s Name) IPFromService(svc *v1.Service) (net.IP, error) {
	const op = "annotation/Name.IPFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var ip net.IP

	err := s.applyToValue(op, svc, func(v string) error {
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("invalid ip address: %s", v)
		}
		return nil
	})

	return ip, err
}

// DurationFromService retrieves the time.Duration value belonging to the
// annotation from svc.
//
// DurationFromService returns an error if the value could not be converted to
// a time.Duration, or the annotation was not set. In the case of a missing
// value, the error wraps ErrNotSet.
func (s Name) DurationFromService(svc *v1.Service) (time.Duration, error) {
	const op = "annotation/Name.DurationFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var d time.Duration

	err := s.applyToValue(op, svc, func(v string) error {
		var err error

		d, err = time.ParseDuration(v)
		return err
	})

	return d, err
}

// LBSvcProtocolFromService retrieves the hvlb.LoadBalancerServiceProtocol
// value belonging to the annotation from svc.
//
// LBSvcProtocolFromService returns an error if the value could not be
// converted to a hvlb.LoadBalancerServiceProtocol, or the annotation was not
// set. In the case of a missing value, the error wraps ErrNotSet.
func (s Name) LBSvcProtocolFromService(svc *v1.Service) (hvlb.LoadBalancerServiceProtocol, error) {
	const op = "annotation/Name.LBSvcProtocolFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var p hvlb.LoadBalancerServiceProtocol

	err := s.applyToValue(op, svc, func(v string) error {
		var err error

		p, err = validateServiceProtocol(v)
		return err
	})

	return p, err
}

// LBAlgorithmTypeFromService retrieves the hvlb.LoadBalancerAlgorithmType
// value belonging to the annotation from svc.
//
// LBAlgorithmTypeFromService returns an error if the value could not be
// converted to a hvlb.LoadBalancerAlgorithmType, or the annotation was not
// set. In the case of a missing value, the error wraps ErrNotSet.
func (s Name) LBAlgorithmTypeFromService(svc *v1.Service) (hvlb.LoadBalancerAlgorithmType, error) {
	const op = "annotation/Name.LBAlgorithmTypeFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var alg hvlb.LoadBalancerAlgorithmType

	err := s.applyToValue(op, svc, func(v string) error {
		var err error

		alg, err = validateAlgorithmType(v)
		return err
	})

	return alg, err
}

// NetworkZoneFromService retrieves the hv.NetworkZone value belonging to
// the annotation from svc.
//
// NetworkZoneFromService returns ErrNotSet if the annotation was not set.
func (s Name) NetworkZoneFromService(svc *v1.Service) (hv.NetworkZone, error) {
	const op = "annotation/Name.NetworkZoneFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var nz hv.NetworkZone

	err := s.applyToValue(op, svc, func(v string) error {
		nz = hv.NetworkZone(v)
		return nil
	})

	return nz, err
}

// CertificatesFromService retrieves the []*hv.Certificate value belonging
// to the annotation from svc.
//
// CertificatesFromService returns an error if the value could not be converted
// to a []*hv.Certificate, or the annotation was not set. In the case of a
// missing value, the error wraps ErrNotSet.
func (s Name) CertificatesFromService(svc *v1.Service) ([]*hv.Certificate, error) {
	const op = "annotation/Name.CertificatesFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var cs []*hv.Certificate

	err := s.applyToValue(op, svc, func(v string) error {
		ss := strings.Split(v, ",")
		cs = make([]*hv.Certificate, len(ss))

		for i, s := range ss {
			id, err := strconv.Atoi(s)
			if err != nil {
				// If we could not parse the string as an integer we assume it
				// is a name not an id.
				cs[i] = &hv.Certificate{Name: s}
				continue
			}
			cs[i] = &hv.Certificate{ID: id}
		}

		return nil
	})

	return cs, err
}

// CertificateTypeFromService retrieves the hv.CertificateType value
// belonging to the annotation from svc.
//
// CertificateTypeFromService returns an error if the value could not be
// converted to a hv.CertificateType. In the case of a missing value, the
// error wraps ErrNotSet.
func (s Name) CertificateTypeFromService(svc *v1.Service) (hv.CertificateType, error) {
	const op = "annotation/Name.CertificateTypeFromService"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	var ct hv.CertificateType

	err := s.applyToValue(op, svc, func(v string) error {
		switch strings.ToLower(v) {
		case string(hv.CertificateTypeUploaded):
			ct = hv.CertificateTypeUploaded
		case string(hv.CertificateTypeManaged):
			ct = hv.CertificateTypeManaged
		default:
			return fmt.Errorf("%s: unsupported certificate type: %s", op, v)
		}
		return nil
	})

	return ct, err
}

func (s Name) applyToValue(op string, svc *v1.Service, f func(string) error) error {
	v, ok := s.StringFromService(svc)
	if !ok {
		return fmt.Errorf("%s: %v: %w", op, s, ErrNotSet)
	}
	if err := f(v); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func validateAlgorithmType(algorithmType string) (hvlb.LoadBalancerAlgorithmType, error) {
	const op = "annotation/validateAlgorithmType"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	algorithmType = strings.ToLower(algorithmType) // Lowercase because all our protocols are lowercase
	hvAlgorithmType := hvlb.LoadBalancerAlgorithmType(algorithmType)

	switch hvAlgorithmType {
	case hvlb.LoadBalancerAlgorithmTypeLeastConnections:
	case hvlb.LoadBalancerAlgorithmTypeRoundRobin:
	default:
		return "", fmt.Errorf("%s: invalid: %s", op, algorithmType)
	}

	return hvAlgorithmType, nil
}

func validateServiceProtocol(protocol string) (hvlb.LoadBalancerServiceProtocol, error) {
	const op = "annotation/validateServiceProtocol"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	protocol = strings.ToLower(protocol) // Lowercase because all our protocols are lowercase
	hvProtocol := hvlb.LoadBalancerServiceProtocol(protocol)
	switch hvProtocol {
	case hvlb.LoadBalancerServiceProtocolTCP:
	case hvlb.LoadBalancerServiceProtocolHTTPS:
	case hvlb.LoadBalancerServiceProtocolHTTP:
		// Valid
		break
	default:
		return "", fmt.Errorf("%s: invalid: %s", op, protocol)
	}
	return hvProtocol, nil
}

type serviceAnnotator struct {
	Svc *v1.Service
	Err error
}

func (sa *serviceAnnotator) Annotate(n Name, v interface{}) {
	if sa.Err != nil {
		return
	}
	sa.Err = n.AnnotateService(sa.Svc, v)
}
