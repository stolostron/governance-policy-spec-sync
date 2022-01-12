// Copyright (c) 2020 Red Hat, Inc.

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stolostron/governance-policy-propagator/test/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const case2PolicyName string = "default.case2-test-policy"
const case2PolicyYaml string = "../resources/case2_uninstall_ns/case2-test-policy.yaml"
const case2UninstallYaml string = "../resources/case2_uninstall_ns/case2-uninstall-ns.yaml"

var _ = Describe("Test uninstall ns", func() {
	BeforeEach(func() {
		By("Creating a ns on managed cluster")
		utils.Kubectl("create", "ns", "uninstall",
			"--kubeconfig=../../kubeconfig_managed")
		Eventually(func() interface{} {
			_, err := clientManaged.CoreV1().Namespaces().Get("uninstall", metav1.GetOptions{})
			return err
		}, defaultTimeoutSeconds, 1).Should(BeNil())
		By("Creating a policy on mananged cluster in ns: uninstall")
		utils.Kubectl("apply", "-f", case2PolicyYaml, "-n", "uninstall",
			"--kubeconfig=../../kubeconfig_managed")
		opt := metav1.ListOptions{}
		utils.ListWithTimeout(clientManagedDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
	})
	AfterEach(func() {
		By("Delete the job on managed cluster")
		utils.Kubectl("delete", "job", "uninstall-ns", "-n", "multicluster-endpoint",
			"--kubeconfig=../../kubeconfig_managed")
	})
	It("should remove ns on managed cluster", func() {
		By("Running uninstall ns job")
		utils.Kubectl("apply", "-f", case2UninstallYaml, "-n", "multicluster-endpoint",
			"--kubeconfig=../../kubeconfig_managed")
		By("Checking if ns uninstall has been deleted eventually")
		Eventually(func() interface{} {
			_, err := clientManaged.CoreV1().Namespaces().Get("uninstall", metav1.GetOptions{})
			return errors.IsNotFound(err)
		}, 120, 1).Should(BeTrue())
	})
})
