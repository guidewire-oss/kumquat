package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	kumquatv1beta1 "kumquat/api/v1beta1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
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

	var useExistingCluster bool = true
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.30.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
		UseExistingCluster: &useExistingCluster,
	}

	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = kumquatv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	namespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "templates"}}
	err = k8sClient.Create(context.Background(), &namespace)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := kindCluster.Destroy(context.TODO())
	Expect(err).NotTo(HaveOccurred())
	By("tearing down the test environment")
	err = testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
