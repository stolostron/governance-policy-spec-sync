module github.com/stolostron/governance-policy-spec-sync

go 1.14

require (
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/prometheus/common v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/stolostron/governance-policy-propagator v0.0.0-20220112203818-b0bc3438ecd7
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	golang.org/x/text => golang.org/x/text v0.3.3 // CVE-2020-14040
	howett.net/plist => github.com/DHowett/go-plist v0.0.0-20181124034731-591f970eefbb
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)

replace github.com/open-cluster-management/api => open-cluster-management.io/api v0.0.0-20200610161514-939cead3902c
