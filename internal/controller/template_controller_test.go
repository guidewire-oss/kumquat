package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

		By("ensuring the templates namespace exists")
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

		By("creating input resources")
		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			inputFolder := path.Join(exampleFolderPath, exampleFolder, "input")
			applyYAMLFilesFromDirectory(ctx, inputFolder)
		}
	})

	AfterEach(func() {
		By("deleting input resources and templates")
		exampleFolderPath := path.Join("../", "../", "examples")
		exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
		Expect(err).NotTo(HaveOccurred())
		for _, exampleFolder := range exampleFolders {
			inputFolder := path.Join(exampleFolderPath, exampleFolder, "input")
			deleteYAMLFilesFromDirectory(ctx, inputFolder)
			templateFolder := path.Join(exampleFolderPath, exampleFolder, "templates")
			deleteYAMLFilesFromDirectory(ctx, templateFolder)
		}
	})

	Context("When Template is applied", func() {
		It("creates expected managed resources", func() {
			exampleFolderPath := path.Join("../", "../", "examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())

			for _, exampleFolder := range exampleFolders {
				By(fmt.Sprintf("Running example %s", exampleFolder))

				By("applying example templates")
				templateFolder := path.Join(exampleFolderPath, exampleFolder, "templates")
				applyYAMLFilesFromDirectory(ctx, templateFolder)

				By("verifying that the output.yaml file has been created")
				verifyExampleOutput(path.Join(exampleFolderPath, exampleFolder), "out.yaml")
			}
		})

		It("creates expected managed resources again to test for race conditions", func() {
			exampleFolderPath := path.Join("../", "../", "examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())

			for _, exampleFolder := range exampleFolders {
				By(fmt.Sprintf("Running example %s", exampleFolder))

				By("applying example templates")
				templateFolder := path.Join(exampleFolderPath, exampleFolder, "templates")
				applyYAMLFilesFromDirectory(ctx, templateFolder)

				By("verifying that the output.yaml file has been created")
				verifyExampleOutput(path.Join(exampleFolderPath, exampleFolder), "out.yaml")
			}
		})
	})

	Context("When Template is deleted", func() {
		JustBeforeEach(func() {
			By("applying example templates")
			exampleFolderPath := path.Join("../", "../", "examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())
			for _, exampleFolder := range exampleFolders {
				templateFolder := path.Join(exampleFolderPath, exampleFolder, "templates")
				applyYAMLFilesFromDirectory(ctx, templateFolder)
			}
		})

		It("deletes managed resources", func() {
			exampleFolderPath := path.Join("../", "../", "examples")
			exampleFolders, err := utils.GetSubDirs(exampleFolderPath)
			Expect(err).NotTo(HaveOccurred())

			for _, exampleFolder := range exampleFolders {
				By("deleting example template")
				templateFolder := path.Join(exampleFolderPath, exampleFolder, "templates")
				deleteYAMLFilesFromDirectory(ctx, templateFolder)

				By("verifying expected managed resources got deleted")
				expectedOutputPath := path.Join(exampleFolderPath, exampleFolder, "expected", "out.yaml")
				expectedOutputResource, err := getK8sClientObject(expectedOutputPath)
				Expect(err).NotTo(HaveOccurred())
				waitForDeletion(ctx, expectedOutputResource, 10*time.Second, 5*time.Millisecond)
			}
		})
	})
})

// Get Kubernetes object defined by the YAML in the given path, filePath
func getK8sClientObject(path string) (*unstructured.Unstructured, error) {
	fileInfo, err := os.Stat(path)
	Expect(err).NotTo(HaveOccurred())

	if fileInfo.IsDir() || (filepath.Ext(fileInfo.Name()) != ".yaml" && filepath.Ext(fileInfo.Name()) != ".yml") {
		return nil, fmt.Errorf("Invalid path: %s should be a YAML file", path)
	}

	content, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err = decoder.Decode(content, nil, obj)
	Expect(err).NotTo(HaveOccurred())

	// Remove resourceVersion if set
	obj.SetResourceVersion("")

	return obj, nil
}

// Apply all YAML files defined in the given directory, dir, as well as all sub-directories
func applyYAMLFilesFromDirectory(ctx context.Context, dir string) {
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		Expect(err).NotTo(HaveOccurred())

		// Entry is a YAML file - Apply manifest
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".yaml" || filepath.Ext(entry.Name()) == ".yml") {
			// Get K8s resource from YAML file
			obj, err := getK8sClientObject(path)
			Expect(err).NotTo(HaveOccurred())

			// Apply the resource to the cluster
			err = k8sClient.Create(ctx, obj)
			if errors.IsAlreadyExists(err) {
				err = k8sClient.Update(ctx, obj)
			}
			Expect(err).NotTo(HaveOccurred())
		}

		// Entry is a directory - Continue traversal
		return nil
	})
	Expect(err).NotTo(HaveOccurred())
}

// Waits until the given K8s Client resource is deleted, or times out
func waitForDeletion(ctx context.Context, obj *unstructured.Unstructured, timeout time.Duration, interval time.Duration) {
	fmt.Printf("Waiting for deletion of resource %s/%s\n", obj.GetKind(), obj.GetName())
	Eventually(func() error {
		err := k8sClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
		if errors.IsNotFound(err) {
			fmt.Printf("Resource %s/%s deleted done\n", obj.GetKind(), obj.GetName())
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("Resource %s/%s not deleted yet\n", obj.GetKind(), obj.GetName())
		return fmt.Errorf("Resource %s/%s still exists", obj.GetKind(), obj.GetName())
	}, timeout, interval).Should(Succeed())
}

// Delete all YAML files defined in the given directory, dir, as well as all sub-directories
func deleteYAMLFilesFromDirectory(ctx context.Context, dir string) {
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		Expect(err).NotTo(HaveOccurred())

		// Entry is a YAML file - Delete manifest
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".yaml" || filepath.Ext(entry.Name()) == ".yml") {
			// Get K8s resource from YAML file
			obj, err := getK8sClientObject(path)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("Deleting %s/%s\n", obj.GetKind(), obj.GetName())
			err = k8sClient.Delete(ctx, obj)
			if errors.IsNotFound(err) {
				return nil
			}
			Expect(err).NotTo(HaveOccurred())

			// Wait for the resource to be deleted
			waitForDeletion(ctx, obj, 10*time.Second, 5*time.Millisecond)
		}

		// Entry is a directory - Continue traversal
		return nil
	})
	Expect(err).NotTo(HaveOccurred())
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
