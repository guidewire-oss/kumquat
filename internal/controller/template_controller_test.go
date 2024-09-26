package controller

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/test/utils"

	yamlv3 "gopkg.in/yaml.v3" // Ensure this is included
)

var _ = Describe("Template Controller Integration Test", func() {

	var ctx context.Context
	const namespaceName = "templates"

	BeforeEach(func() {
		ctx = context.Background()

		By("ensuring the namespace exists")
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		err := k8sClient.Create(ctx, namespace)
		if err != nil && !errors.IsAlreadyExists(err) {
			Expect(err).NotTo(HaveOccurred())
		}

	})

	AfterEach(func() {
		By("deleting all Template resources")
		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)

		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			deleteExample(path.Join(exampleFolderPath, exampleFolder, "input", "templates"), ctx)
		}

		By("deleting the namespace")
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		err = k8sClient.Delete(ctx, namespace)
		if err != nil && !errors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should reconcile the Template and create expected resources", func() {

		const (
			resourceName  = "generate-role"
			namespaceName = "templates"
			timeout       = time.Second * 10
			interval      = time.Millisecond * 250
		)
		By("applying all resources from the crds directory")
		crdsDir := "../../examples/extending-rbac/input/crds"
		applyYAMLFilesFromDirectory(ctx, crdsDir, namespaceName)

		By("creating the Template resource")
		// print the current working directory
		templatePath := filepath.Join("resources/template_resource.yaml")

		templateData, err := os.ReadFile(templatePath)
		Expect(err).NotTo(HaveOccurred())

		var data map[string]interface{}
		err = yamlv3.Unmarshal(templateData, &data)
		Expect(err).NotTo(HaveOccurred())

		// Extract the necessary fields from the YAML
		spec := data["spec"].(map[string]interface{})
		query := spec["query"].(string)
		templateSpec := spec["template"].(map[string]interface{})
		language := templateSpec["language"].(string)
		batchModeProcessing := templateSpec["batchModeProcessing"].(bool)
		fileName := templateSpec["fileName"].(string)
		templateDataField := templateSpec["data"].(string)

		// Create the Template object
		template := &kumquatv1beta1.Template{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Spec: kumquatv1beta1.TemplateSpec{
				Query: query,
				TemplateDefinition: kumquatv1beta1.TemplateDefinition{
					Language:            language,
					BatchModeProcessing: batchModeProcessing,
					Data:                templateDataField,
					FileName:            fileName,
				},
			},
		}

		Expect(k8sClient.Create(ctx, template)).To(Succeed())

		By("verifying that the Template has been reconciled")
		Eventually(func() bool {

			// Define the lookup key for the resource
			resourceLookupKey := client.ObjectKey{
				Namespace: "templates",
				Name:      "generate-role",
			}

			// Define the resource object
			resource := &kumquatv1beta1.Template{} // Replace with the actual type if different

			// Attempt to get the resource
			err := k8sClient.Get(ctx, resourceLookupKey, resource)
			if err != nil {
				return false
			}
			// print the resource

			// Resource exists
			return true
		}, 10*time.Second, 1*time.Second).Should(BeTrue())

		// another eventuallly block to check if the output.yaml file has been created
		By("verifying that the output.yaml file has been created")
		Eventually(func() error {

			outputFilePath := filepath.Join("resources/out.yaml")
			outputData, err := os.ReadFile(outputFilePath)
			Expect(err).NotTo(HaveOccurred())

			decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
			obj := &unstructured.Unstructured{}
			_, _, err = decoder.Decode(outputData, nil, obj)
			Expect(err).NotTo(HaveOccurred())

			// Use the information from the output.yaml to verify the resource in the cluster
			resourceLookupKey := client.ObjectKey{
				Namespace: obj.GetNamespace(),
				Name:      obj.GetName(),
			}

			// Fetch the resource from the cluster
			clusterResource := &unstructured.Unstructured{}
			clusterResource.SetGroupVersionKind(obj.GroupVersionKind())

			err = k8sClient.Get(ctx, resourceLookupKey, clusterResource)
			if err != nil {

				return err

			}
			return nil

			// Verify the fetched resource from the cluster matches the expected resource
			// Expect(clusterResource.Object).To(Equal(obj.Object),
			// "The resource created in the cluster should match the output.yaml file")
		}, 10*time.Second, 1*time.Second).Should(Succeed())

	})
})

// Function to apply all YAML files from a directory
func applyYAMLFilesFromDirectory(ctx context.Context, dir string, namespaceName string) {
	files, err := os.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred())

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			filePath := filepath.Join(dir, file.Name())

			content, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Split the file into individual YAML documents
			documents := strings.Split(string(content), "\n---")
			for _, doc := range documents {
				doc = strings.TrimSpace(doc)
				if len(doc) == 0 {
					continue
				}

				// Decode YAML into unstructured.Unstructured
				decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				obj := &unstructured.Unstructured{}
				_, _, err := decoder.Decode([]byte(doc), nil, obj)
				Expect(err).NotTo(HaveOccurred())

				// Remove resourceVersion if set
				obj.SetResourceVersion("")

				// Set namespace if necessary
				if obj.GetNamespace() == "" && obj.GetKind() != "Namespace" {
					obj.SetNamespace(namespaceName) // Use your test namespace
				}

				// Apply the resource to the cluster
				err = k8sClient.Create(ctx, obj)
				if err != nil {
					if errors.IsAlreadyExists(err) {
						err = k8sClient.Update(ctx, obj)
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).NotTo(HaveOccurred())
					}
				}
			}
		}
	}
}

func deleteExample(exampleFolder string, ctx context.Context) {
	files, err := os.ReadDir(exampleFolder)
	Expect(err).NotTo(HaveOccurred())

	Expect(err).NotTo(HaveOccurred())
	for _, file := range files {
		filePath := filepath.Join(exampleFolder, file.Name())
		fileData, err := os.ReadFile(filePath)
		Expect(err).NotTo(HaveOccurred())
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		obj := &unstructured.Unstructured{}
		_, _, err = decoder.Decode(fileData, nil, obj)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Delete(ctx, obj)
		if errors.IsNotFound(err) {
			continue
		}
		Expect(err).NotTo(HaveOccurred())
	}

}
