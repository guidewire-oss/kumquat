package controller

import (
	"context"
	"path/filepath" // Alias the standard library runtime package

	// Alias the standard library runtime package
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	scheme    = runtime.NewScheme()
)

type MockedGetPreferredGVKClient struct {
	K8sClient
}

func (m *MockedGetPreferredGVKClient) GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error) {
	// Return a predefined GVK
	return schema.GroupVersionKind{Group: group, Version: "v1", Kind: kind}, nil
}

var stopMgr context.CancelFunc

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kumquatv1beta1.AddToScheme(scheme))

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// Start the manager and controller
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(err).NotTo(HaveOccurred())
	Expect(err).NotTo(HaveOccurred())
	Expect(err).NotTo(HaveOccurred())
	dynamicK8sClient, err := GetK8sClient(k8sClient, k8sManager.GetRESTMapper())
	Expect(err).ToNot(HaveOccurred())
	mockedClient := &MockedGetPreferredGVKClient{
		K8sClient: dynamicK8sClient,
	}

	err = (&TemplateReconciler{
		Client:    k8sClient,
		Scheme:    scheme,
		K8sClient: mockedClient,

		// No WatchManager since it's not used in main.go
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	var mgrCtx context.Context
	mgrCtx, stopMgr = context.WithCancel(context.Background())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(mgrCtx)
		Expect(err).ToNot(HaveOccurred())
	}()

	// Wait for the cache to sync
	Expect(k8sManager.GetCache().WaitForCacheSync(context.Background())).To(BeTrue())

})

var _ = AfterSuite(func() {
	By("Stopping the manager")
	stopMgr() // Signal the manager to stop

	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())

})
