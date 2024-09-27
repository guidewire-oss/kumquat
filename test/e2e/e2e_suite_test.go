package e2e

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/test/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var kindClusterName = "kind"
var kindCluster *kind.Cluster

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	_, _ = fmt.Fprintf(GinkgoWriter, "Starting kumquat suite\n")

	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))

	By("bootstrapping test environment")
	kindCluster = kind.NewCluster(kindClusterName)
	_, err := kindCluster.SetDefaults().Create(context.TODO())
	Expect(err).NotTo(HaveOccurred())

	err = kumquatv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("kubectl", "create", "ns", "templates")

	_, err = utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred())

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := kindCluster.Destroy(context.TODO())
	Expect(err).NotTo(HaveOccurred())
})
