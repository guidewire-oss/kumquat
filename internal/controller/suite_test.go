package controller

import (
	"context"
	"path/filepath" // Alias the standard library runtime package

	// Alias the standard library runtime package
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	kumquatv1beta1 "kumquat/api/v1beta1"
	// +kubebuilder:scaffold:imports
)

var (
	cfg             *rest.Config
	k8sClient       client.Client
	testEnv         *envtest.Environment
	scheme          = runtime.NewScheme()
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface
)

var stopMgr context.CancelFunc

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	var err error

	// Set up the logger
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	// Add schemes
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kumquatv1beta1.AddToScheme(scheme))

	// Bootstrap the test environment
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	// Start the test environment and obtain the configuration
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// Initialize the k8sClient using the test environment's configuration
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// Initialize the dynamic client using the test environment's configuration
	dynamicClient, err = dynamic.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred(), "Failed to create dynamic client")

	// Initialize the discovery client using the test environment's configuration
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(cfg)
	Expect(err).NotTo(HaveOccurred(), "Failed to create discovery client")

	// Start the manager and controller
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())

	dynamicK8sClient, err := NewDynamicK8sClient(k8sManager.GetClient(), k8sManager.GetRESTMapper())
	Expect(err).NotTo(HaveOccurred())

	err = (&TemplateReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    scheme,
		K8sClient: dynamicK8sClient,
	}).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	var mgrCtx context.Context
	mgrCtx, stopMgr = context.WithCancel(context.Background())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(mgrCtx)
		Expect(err).NotTo(HaveOccurred())
	}()

	// Wait for the cache to sync
	Expect(k8sManager.GetCache().WaitForCacheSync(context.Background())).To(BeTrue())
})
