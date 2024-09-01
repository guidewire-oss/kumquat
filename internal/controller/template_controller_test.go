/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"                   // Import the core API group
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Import for metadata
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	kumquatv1beta1 "kumquat/api/v1beta1"
)

var _ = Describe("Template Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "generate-role"
		const namespaceName = "templates"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "templates", // TODO(user):Modify as needed
		}
		template := &kumquatv1beta1.Template{}

		BeforeEach(func() {
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

			By("creating the custom resource for the Kind Template")
			err = k8sClient.Get(ctx, typeNamespacedName, template)
			if err != nil && errors.IsNotFound(err) {
				templatePath := filepath.Join("internal/controller/resources/template_resource.yaml")
				templateData, err := os.ReadFile(templatePath)
				Expect(err).NotTo(HaveOccurred())
				var data map[string]interface{}
				err = yaml.Unmarshal(templateData, &data)
				if err != nil {
					log.Fatalf("Error unmarshalling YAML: %v", err)
				}

				// Step 3: Marshal the map into JSON
				// Step 2: Access the "query" field under "spec"
				spec := data["spec"].(map[string]interface{})
				query := spec["query"].(string)
				language := spec["template"].(map[string]interface{})["language"].(string)
				batchModeProcessing := spec["template"].(map[string]interface{})["batchModeProcessing"].(bool)
				fileName := spec["template"].(map[string]interface{})["fileName"].(string)
				templateDataa := spec["template"].(map[string]interface{})["data"].(string)

				// convert object to kumquatv1beta1.Template
				template := &kumquatv1beta1.Template{}

				template.SetName(resourceName)
				template.SetNamespace(namespaceName)
				template.APIVersion = "kumquat.guidewire.com/v1beta1"
				template.Kind = "Template"
				template.Spec = kumquatv1beta1.TemplateSpec{
					Query: query,
					TemplateDefinition: kumquatv1beta1.TemplateDefinition{
						Language:            language,
						BatchModeProcessing: batchModeProcessing,
						Data:                templateDataa,
						FileName:            fileName,
					},
				}
				fmt.Println("this is template", template)
				Expect(k8sClient.Create(ctx, template)).To(Succeed())

				obj1 := &unstructured.Unstructured{}
				// set name, namespace and kind and apiVersion
				obj1.SetName(resourceName)
				obj1.SetNamespace(namespaceName)
				obj1.SetKind("Template")
				obj1.SetAPIVersion("kumquat.guidewire.com/v1beta1")

				err = k8sClient.Get(ctx, typeNamespacedName, obj1)
				Expect(err).NotTo(HaveOccurred())

				//	fmt.Fprintf(GinkgoWriter, "Created object: %+v\n", obj)

			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &unstructured.Unstructured{}
			resource.SetName(resourceName)
			resource.SetNamespace(namespaceName)
			resource.SetKind("Template")
			resource.SetAPIVersion("kumquat.guidewire.com/v1beta1")

			err := k8sClient.Get(ctx, typeNamespacedName, resource)

			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Template")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {

			obj := &kumquatv1beta1.Template{}
			err := k8sClient.Get(ctx, typeNamespacedName, obj)
			Expect(err).NotTo(HaveOccurred())

			// fmt.Fprintf(GinkgoWriter, "Created object: %+v\n", obj)

			By("Reconciling the created resource")
			controllerReconciler := &TemplateReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
