package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	sigyaml "sigs.k8s.io/yaml"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"kumquat/test/utils"
	// Ensure this is included
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
		// Wait for the namespace to be created
		Eventually(func() error {
			err = k8sClient.Get(ctx, client.ObjectKeyFromObject(namespace), namespace)
			if err != nil {
				return err
			}
			if namespace.Status.Phase == corev1.NamespaceActive {
				fmt.Println("Namespace created")
				return nil
			}
			fmt.Printf("Namespace: %s, Phase: %s\n", namespace.Name, namespace.Status.Phase)
			return fmt.Errorf("namespace not active yet")
		}, 10*time.Second, 2*time.Second).Should(Succeed())

	})

	AfterEach(func() {
		By("deleting all applied resources")
		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			deleteExampleResources(ctx, path.Join(exampleFolderPath, exampleFolder, "input"))
		}

	})

	It("should apply and verify examples", func() {

		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			By(fmt.Sprintf("Running example %s", exampleFolder))
			applyExampleResources(ctx, path.Join(exampleFolderPath, exampleFolder))
		}

		// another eventuallly block to check if the output.yaml file has been created
		By("verifying that the output.yaml file has been created")
		for _, exampleFolder := range exampleFolders {
			exampleFolderPath := path.Join("..", "..", "examples", exampleFolder)
			verifyExampleOutput(exampleFolderPath, "out.yaml")
		}

	})
	It("repeate apply and verify examples to check some race conditions", func() {

		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			By(fmt.Sprintf("Running example %s", exampleFolder))
			applyExampleResources(ctx, path.Join(exampleFolderPath, exampleFolder))
		}

		// another eventuallly block to check if the output.yaml file has been created
		By("verifying that the output.yaml file has been created")
		for _, exampleFolder := range exampleFolders {
			exampleFolderPath := path.Join("..", "..", "examples", exampleFolder)
			verifyExampleOutput(exampleFolderPath, "out.yaml")
		}

	})
})

// Function to apply resources from an example
func applyExampleResources(ctx context.Context, examplePath string) {
	inputPath := filepath.Join(examplePath, "input")
	// get all folders in the input directory
	folders, err := utils.GetSubDirs(inputPath)
	Expect(err).NotTo(HaveOccurred())
	for _, folder := range folders {
		applyYAMLFilesFromDirectory(ctx, path.Join(inputPath, folder))
	}
}

// Function to apply all YAML files from a directory
func applyYAMLFilesFromDirectory(ctx context.Context, dir string) {
	files, err := os.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred())

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			filePath := filepath.Join(dir, file.Name())

			content, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())

			decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
			obj := &unstructured.Unstructured{}
			_, _, err = decoder.Decode(content, nil, obj)
			Expect(err).NotTo(HaveOccurred())

			// Remove resourceVersion if set
			obj.SetResourceVersion("")

			// Apply the resource to the cluster
			err = k8sClient.Create(ctx, obj)
			if errors.IsAlreadyExists(err) {
				err = k8sClient.Update(ctx, obj)
			}
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

func deleteYAMLFilesFromDirectory(ctx context.Context, dir string) {
	files, err := os.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred())

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			filePath := filepath.Join(dir, file.Name())
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
			fmt.Printf("Deleting %s/%s\n", obj.GetKind(), obj.GetName())
			// Wait for the resource to be deleted
			Eventually(func() error {
				err = k8sClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
				if errors.IsNotFound(err) {
					fmt.Printf("Resource %s/%s deleted done\n", obj.GetKind(), obj.GetName())
					return nil
				}
				if err != nil {
					return err
				}
				fmt.Printf("Resource %s/%s not deleted yet\n", obj.GetKind(), obj.GetName())
				return fmt.Errorf("resource still exists")
			}, 10*time.Second, 5*time.Millisecond).Should(Succeed())
		}
	}
}

func verifyExampleOutput(exampleFolder string, exampleFile string) {
	expectedFilePath := path.Join(exampleFolder, "expected", exampleFile)
	filePath, err := filepath.Abs(expectedFilePath)
	Expect(err).NotTo(HaveOccurred())
	expectedOutput, err := os.ReadFile(filePath)
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() error {
		// convert expected data to unsructured
		expectedData := &unstructured.Unstructured{}
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		_, _, err := decoder.Decode(expectedOutput, nil, expectedData)
		Expect(err).NotTo(HaveOccurred())

		resourceLookupKey := client.ObjectKey{
			Namespace: expectedData.GetNamespace(),
			Name:      expectedData.GetName(),
		}

		actualOutput := &unstructured.Unstructured{}
		actualOutput.SetGroupVersionKind(expectedData.GroupVersionKind())

		err = k8sClient.Get(context.Background(), resourceLookupKey, actualOutput)
		if err != nil {
			return err
		}

		yamlData, err := sigyaml.Marshal(actualOutput)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println(string(yamlData), "this issss")
		// delete metadata creationTimestamp, resourceVersion, uid, generation, selfLink
		if metadata, ok := actualOutput.Object["metadata"].(map[string]interface{}); ok {
			delete(metadata, "creationTimestamp")
			delete(metadata, "resourceVersion")
			delete(metadata, "uid")
			delete(metadata, "managedFields") // Ensure this is removed
			delete(metadata, "generation")
			delete(metadata, "selfLink")
		}

		// Convert to JSON for comparison
		actualJSON, err := json.Marshal(actualOutput)
		if err != nil {
			return err
		}

		expectedJSON, err := json.Marshal(expectedData)
		if err != nil {

			return err
		}

		var actualJSONData interface{}
		var expectedJSONData interface{}

		err = json.Unmarshal(actualJSON, &actualJSONData)
		if err != nil {
			return err
		}

		err = json.Unmarshal(expectedJSON, &expectedJSONData)
		if err != nil {
			return err
		}

		// Compare JSON objects ignoring the order of elements in arrays
		if diff := cmp.Diff(actualJSONData, expectedJSONData, cmpopts.SortSlices(func(x, y interface{}) bool {
			return fmt.Sprintf("%v", x) < fmt.Sprintf("%v", y)
		})); diff != "" {
			fmt.Printf("Actual data: %v\n", actualJSONData)
			fmt.Printf("Expected data: %v\n", expectedJSONData)
			fmt.Printf("Difference: %v\n", diff)
			return fmt.Errorf("actual data does not match expected data")
		}

		return nil
	}, 10*time.Second, 2*time.Second).Should(Succeed())

}

func deleteExampleResources(ctx context.Context, inputPath string) {
	folders, err := utils.GetSubDirs(inputPath)
	Expect(err).NotTo(HaveOccurred())
	for _, folder := range folders {
		deleteYAMLFilesFromDirectory(ctx, path.Join(inputPath, folder))
	}
}
