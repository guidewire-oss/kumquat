package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"kumquat/test/utils"
)

const namespace = "kumquat-system"

var _ = Describe("controller", Ordered, func() {
	BeforeAll(func() {
		By("installing prometheus operator")
		Expect(utils.InstallPrometheusOperator()).To(Succeed())

		By("installing the cert-manager")
		Expect(utils.InstallCertManager()).To(Succeed())

		By("creating manager namespace")
		cmd := exec.Command("kubectl", "create", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	AfterAll(func() {
		By("uninstalling the Prometheus manager bundle")
		utils.UninstallPrometheusOperator()

		By("uninstalling the cert-manager bundle")
		utils.UninstallCertManager()

		By("removing manager namespace")
		cmd := exec.Command("kubectl", "delete", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	Context("Operator", func() {
		It("should ensure the controller-manager pod is running", func() {
			var controllerPodName string
			var err error

			// projectimage stores the name of the image used in the example
			var projectimage = "guidewire.com/kumquat:v0.0.1"

			By("building the manager(Operator) image")
			cmd := exec.Command("make", "docker-build", fmt.Sprintf("IMG=%s", projectimage))
			_, err = utils.Run(cmd)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			By("loading the the manager(Operator) image on Kind")
			err = utils.LoadImageToKindClusterWithName(projectimage)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			By("installing CRDs")
			cmd = exec.Command("make", "install")
			_, err = utils.Run(cmd)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			By("deploying the controller-manager")
			cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectimage))
			_, err = utils.Run(cmd)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			By("validating that the controller-manager pod is running as expected")
			verifyControllerUp := func() error {
				// Get pod name
				cmd = exec.Command("kubectl", "get",
					"pods", "-l", "control-plane=controller-manager",
					"-o", "go-template={{ range .items }}"+
						"{{ if not .metadata.deletionTimestamp }}"+
						"{{ .metadata.name }}"+
						"{{ \"\\n\" }}{{ end }}{{ end }}",
					"-n", namespace,
				)

				podOutput, err := utils.Run(cmd)
				ExpectWithOffset(2, err).NotTo(HaveOccurred())
				podNames := utils.GetNonEmptyLines(string(podOutput))
				if len(podNames) != 1 {
					return fmt.Errorf("expect 1 controller pods running, but got %d", len(podNames))
				}
				controllerPodName = podNames[0]
				ExpectWithOffset(2, controllerPodName).Should(ContainSubstring("controller-manager"))

				// Validate pod status

				cmd = exec.Command("kubectl", "get",
					"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
					"-n", namespace,
				)
				status, err := utils.Run(cmd)
				ExpectWithOffset(2, err).NotTo(HaveOccurred())
				if string(status) != "Running" {
					return fmt.Errorf("controller pod in %s status", status)
				}
				return nil
			}
			EventuallyWithOffset(1, verifyControllerUp, time.Minute, time.Second).Should(Succeed())
		})

		It("should apply and verify examples", func() {
			exampleFolderPath := path.Join("examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())

			for _, exampleFolder := range exampleFolders {
				By(fmt.Sprintf("Running example %s", exampleFolder))
				applyExample(exampleFolder)
			}

			for _, exampleFolder := range exampleFolders {
				verifyExampleOutput(exampleFolder)
			}
		})

		It("should delete the created resources if the associated template is deleted", func() {
			exampleFolderPath := path.Join("examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())

			for _, exampleFolder := range exampleFolders {
				By(fmt.Sprintf("Running example %s", exampleFolder))
				deleteExample(exampleFolder)
			}

			for _, exampleFolder := range exampleFolders {
				verifyExampleDeletion(exampleFolder)
			}
		})
		//TODO: Add more tests here
		// It("Should ensure the resources are updated if there is an update to dependent resources")
		// It("Should ensure that the resource are updated if there is an update to the template")
		// It("should ensure that the resources are created if the number of resources required to be created is more than one")
		// It("should ensure that the resources are updated if the number of resources required to be created is more than one")

	})
})

func applyExample(exampleFolder string) {
	TemplateFilePath := path.Join("examples", exampleFolder, "input")
	filePath, err := filepath.Abs(TemplateFilePath)
	Expect(err).NotTo(HaveOccurred())
	cmd := exec.Command("kubectl", "apply", "-R", "-f", filePath)
	_, err = utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred())
}

func verifyExampleOutput(exampleFolder string) {
	ExpectedFilePath := path.Join("examples", exampleFolder, "expected", "out.yaml")
	filePath, err := filepath.Abs(ExpectedFilePath)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() error {
		cmd := exec.Command("kubectl", "get", "-f", filePath)
		_, err := utils.Run(cmd)
		return err
	}, 30*time.Second, 2*time.Second).Should(Succeed())

	expectedOutput, err := os.ReadFile(filePath)
	Expect(err).NotTo(HaveOccurred())
	var expectedData interface{}
	err = yaml.Unmarshal(expectedOutput, &expectedData)
	Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("kubectl", "get", "-f", filePath, "-o", "yaml")
	actualOutput, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred())

	var actualData map[string]interface{}
	err = yaml.Unmarshal(actualOutput, &actualData)
	Expect(err).NotTo(HaveOccurred())

	if metadata, ok := actualData["metadata"].(map[string]interface{}); ok {
		delete(metadata, "creationTimestamp")
		delete(metadata, "resourceVersion")
		delete(metadata, "uid")
	}

	Expect(actualData).To(Equal(expectedData))
}

func deleteExample(exampleFolder string) {
	TemplateFilePath := path.Join("examples", exampleFolder, "input", "templates")
	filePath, err := filepath.Abs(TemplateFilePath)
	Expect(err).NotTo(HaveOccurred())
	cmd := exec.Command("kubectl", "delete", "-R", "-f", filePath)
	_, err = utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred())
}

func verifyExampleDeletion(exampleFolder string) {
	ExpectedFilePath := path.Join("examples", exampleFolder, "expected", "out.yaml")
	filePath, err := filepath.Abs(ExpectedFilePath)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() error {
		cmd := exec.Command("kubectl", "get", "-f", filePath)
		_, err := utils.Run(cmd)
		if err == nil {
			return fmt.Errorf("resource still exists")
		}
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		return nil
	}, 10*time.Second, 250*time.Millisecond).Should(Succeed())

	templatePath := path.Join("examples", exampleFolder, "input", "templates")
	filePath, err = filepath.Abs(templatePath)
	Expect(err).NotTo(HaveOccurred())

	files, err := utils.GetFiles(filePath)
	Expect(err).NotTo(HaveOccurred())
	for _, file := range files {
		file = path.Join(filePath, file)
		log.Printf("file: %s", file)
		Eventually(func() error {
			cmd := exec.Command("kubectl", "get", "-f", file)
			_, err = utils.Run(cmd)
			if err == nil {
				return fmt.Errorf("resource still exists")
			}
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
			return nil
		}, 10*time.Second, 250*time.Millisecond).Should(Succeed())
	}
}
