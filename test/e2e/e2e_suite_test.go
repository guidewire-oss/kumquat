package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/test/utils"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var (
	kindClusterName = "kind"
	kindCluster     *kind.Cluster
	fernClient      *client.FernApiClient
)

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

var _ = ReportAfterSuite("", func(report Report) {
	projectID := os.Getenv("E2E_FERN_PROJECT_ID")

	fernReporterBaseURL := "https://fern-reporter.int.ccs.guidewire.net/"
	if os.Getenv("FERN_REPORTER_BASE_URL") != "" {
		fernReporterBaseURL = os.Getenv("FERN_REPORTER_BASE_URL")
	}

	if os.Getenv("REPORT_TO_FERN") == "true" {
		fernClient = client.New(projectID, client.WithBaseURL(fernReporterBaseURL))
		err := fernClient.Report(report)
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "Unable to push report to Fern %v", err)
	} else {
		fmt.Println("Skipping report to Fern as REPORT_TO_FERN is set to false")
	}
})
